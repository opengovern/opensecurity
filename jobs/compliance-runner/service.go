package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	authApi "github.com/opengovern/og-util/pkg/api"
	cloudql_init_job "github.com/opengovern/opensecurity/jobs/cloudql-init"
	coreApi "github.com/opengovern/opensecurity/services/core/api"
	"github.com/opengovern/opensecurity/services/integration/client"
	schedulerApi "github.com/opengovern/opensecurity/services/scheduler/api"
	"github.com/opengovern/opensecurity/services/scheduler/db/model"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/config"
	esSinkClient "github.com/opengovern/og-util/pkg/es/ingest/client"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/steampipe"
	complianceApi "github.com/opengovern/opensecurity/services/compliance/api"
	complianceClient "github.com/opengovern/opensecurity/services/compliance/client"
	coreClient "github.com/opengovern/opensecurity/services/core/client"
	schedulerClient "github.com/opengovern/opensecurity/services/scheduler/client"
	"go.uber.org/zap"
)

type Config struct {
	ElasticSearch         config.ElasticSearch
	NATS                  config.NATS
	Compliance            config.OpenGovernanceService
	Scheduler             config.OpenGovernanceService
	Integration           config.OpenGovernanceService
	Inventory             config.OpenGovernanceService
	Core                  config.OpenGovernanceService
	EsSink                config.OpenGovernanceService
	Steampipe             config.Postgres
	PostgresPlugin        config.Postgres
	PrometheusPushAddress string
}

type Worker struct {
	config        Config
	logger        *zap.Logger
	steampipeConn *steampipe.Database
	esClient      opengovernance.Client
	jq            *jq.JobQueue
	//regoEngine        *regoService.RegoEngine
	complianceClient  complianceClient.ComplianceServiceClient
	integrationClient client.IntegrationServiceClient
	schedulerClient   schedulerClient.SchedulerServiceClient

	coreClient coreClient.CoreServiceClient
	sinkClient esSinkClient.EsSinkServiceClient

	benchmarkCache           map[string]complianceApi.Benchmark
	canceledComplianceJobs   map[uint]bool
	canceledComplianceJobsMu sync.RWMutex
	queryParameters          []coreApi.QueryParameter
	queryParamsMu            sync.RWMutex
}

var (
	ManualTrigger = os.Getenv("MANUAL_TRIGGER")
)

func NewWorker(
	config Config,
	logger *zap.Logger,
	ctx context.Context,
) (*Worker, error) {
	integrationClient := client.NewIntegrationServiceClient(config.Integration.BaseURL)

	pluginJob := cloudql_init_job.NewJob(logger, cloudql_init_job.Config{
		Postgres:      config.PostgresPlugin,
		ElasticSearch: config.ElasticSearch,
		Steampipe:     config.Steampipe,
		Integration:   config.Integration,
	}, integrationClient)
	logger.Info("running plugin job to initialize integrations in cloudql")
	steampipeConn, err := pluginJob.Run(ctx)
	if err != nil {
		logger.Error("failed to run plugin job", zap.Error(err))
		return nil, err
	}

	logger.Info("steampipe service started")
	logger.Sync()

	esClient, err := opengovernance.NewClient(opengovernance.ClientConfig{
		Addresses:     []string{config.ElasticSearch.Address},
		Username:      &config.ElasticSearch.Username,
		Password:      &config.ElasticSearch.Password,
		IsOnAks:       &config.ElasticSearch.IsOnAks,
		IsOpenSearch:  &config.ElasticSearch.IsOpenSearch,
		AwsRegion:     &config.ElasticSearch.AwsRegion,
		AssumeRoleArn: &config.ElasticSearch.AssumeRoleArn,
	})
	if err != nil {
		logger.Error("failed to create elasticsearch client", zap.Error(err))
		logger.Sync()
		return nil, err
	}
	logger.Info("elasticsearch client created")
	logger.Sync()

	jq, err := jq.New(config.NATS.URL, logger)
	if err != nil {
		logger.Error("failed to create job queue", zap.Error(err))
		logger.Sync()
		return nil, err
	}
	logger.Info("job queue connection created")
	logger.Sync()

	queueTopic := JobQueueTopic
	if ManualTrigger == "true" {
		queueTopic = JobQueueTopicManuals
	}

	logger.Info("creating stream", zap.String("stream", StreamName), zap.String("topic", queueTopic), zap.String("resultTopic", ResultQueueTopic))
	logger.Sync()
	if err := jq.Stream(ctx, StreamName, "compliance runner job queue", []string{queueTopic, ResultQueueTopic}, 1000000); err != nil {
		logger.Error("failed to create stream", zap.Error(err), zap.String("stream", StreamName), zap.String("topic", queueTopic), zap.String("resultTopic", ResultQueueTopic))
		return nil, err
	}
	logger.Info("stream created", zap.String("stream", StreamName), zap.String("topic", queueTopic), zap.String("resultTopic", ResultQueueTopic))
	logger.Sync()

	//logger.Info("initializing rego engine")
	//logger.Sync()
	//regoEngine, err := regoService.NewRegoEngine(ctx, logger, steampipeConn)
	//if err != nil {
	//	logger.Error("failed to create rego engine", zap.Error(err))
	//	logger.Sync()
	//	return nil, err
	//}

	w := &Worker{
		config:        config,
		logger:        logger,
		steampipeConn: steampipeConn,
		esClient:      esClient,
		jq:            jq,
		//regoEngine:        regoEngine,
		complianceClient:  complianceClient.NewComplianceClient(config.Compliance.BaseURL),
		schedulerClient:   schedulerClient.NewSchedulerServiceClient(config.Scheduler.BaseURL),
		integrationClient: integrationClient,

		coreClient:               coreClient.NewCoreServiceClient(config.Core.BaseURL),
		sinkClient:               esSinkClient.NewEsSinkServiceClient(logger, config.EsSink.BaseURL),
		benchmarkCache:           make(map[string]complianceApi.Benchmark),
		canceledComplianceJobs:   make(map[uint]bool),
		canceledComplianceJobsMu: sync.RWMutex{},
		queryParamsMu:            sync.RWMutex{},
	}
	ctx2 := &httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}
	benchmarks, err := w.complianceClient.ListAllBenchmarks(ctx2, true)
	if err != nil {
		logger.Error("failed to get benchmarks", zap.Error(err))
		logger.Sync()
		return nil, err
	}
	for _, benchmark := range benchmarks {
		w.benchmarkCache[benchmark.ID] = benchmark
	}

	return w, nil
}

// Run is a blocking function so you may decide to call it in another goroutine.
// It runs a NATS consumer and it will close it when the given context is closed.
func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("starting to consume")
	w.logger.Sync()

	w.logger.Info("fetching parameters values")
	queryParams, err := w.coreClient.ListQueryParameters(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole}, coreApi.ListQueryParametersRequest{})
	if err != nil {
		w.logger.Error("failed to get query parameters", zap.Error(err))
	} else {
		w.queryParamsMu.Lock()
		w.queryParameters = queryParams.Items
		w.queryParamsMu.Unlock()
	}
	go w.fetchParameters(ctx)

	queueTopic := JobQueueTopic
	consumer := ConsumerGroup
	if ManualTrigger == "true" {
		queueTopic = JobQueueTopicManuals
		consumer = ConsumerGroupManuals
	}

	consumeCtx, err := w.jq.ConsumeWithConfig(ctx, consumer, StreamName, []string{queueTopic},
		jetstream.ConsumerConfig{
			DeliverPolicy:     jetstream.DeliverAllPolicy,
			AckPolicy:         jetstream.AckExplicitPolicy,
			AckWait:           time.Hour,
			MaxDeliver:        1,
			InactiveThreshold: time.Hour,
			Replicas:          1,
			MemoryStorage:     false,
			Durable:           fmt.Sprintf("%s-service", consumer),
		}, nil,
		func(msg jetstream.Msg) {
			if err := msg.Ack(); err != nil {
				w.logger.Error("failed to send the ack message", zap.Error(err), zap.Any("msg", msg))
				w.logger.Sync()
			}

			w.logger.Info("received a new job")
			w.logger.Sync()

			jobCtx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go w.pollAPI(jobCtx, cancel, msg)

			_, _, err := w.ProcessMessage(jobCtx, msg)
			if err != nil {
				w.logger.Error("failed to process message", zap.Error(err))
				w.logger.Sync()

			}

			w.logger.Info("processing a job completed")
			w.logger.Sync()

		})
	if err != nil {
		return err
	}

	w.logger.Info("consuming")
	w.logger.Sync()

	<-ctx.Done()
	consumeCtx.Drain()
	consumeCtx.Stop()

	return nil
}

func (w *Worker) ProcessMessage(ctx context.Context, msg jetstream.Msg) (commit bool, requeue bool, err error) {
	var job Job

	if err := json.Unmarshal(msg.Data(), &job); err != nil {
		return true, false, err
	}

	result := JobResult{
		Job:                        job,
		StartedAt:                  time.Now(),
		Status:                     model.ComplianceRunnerInProgress,
		PodName:                    os.Getenv("POD_NAME"),
		Error:                      "",
		TotalComplianceResultCount: nil,
	}

	w.canceledComplianceJobsMu.RLock()
	if _, ok := w.canceledComplianceJobs[job.ParentJobID]; ok {
		return true, false, nil
	}
	w.canceledComplianceJobsMu.RUnlock()

	defer func() {
		if err != nil {
			result.Error = err.Error()
			result.Status = model.ComplianceRunnerFailed
		} else {
			result.Status = model.ComplianceRunnerSucceeded
		}

		w.canceledComplianceJobsMu.RLock()
		if _, ok := w.canceledComplianceJobs[job.ParentJobID]; ok {
			result.Status = model.ComplianceRunnerCanceled
		}
		w.canceledComplianceJobsMu.RUnlock()

		resultJson, err := json.Marshal(result)
		if err != nil {
			w.logger.Error("failed to create job result json", zap.Error(err))
			w.logger.Sync()

			return
		}

		if _, err := w.jq.Produce(ctx, ResultQueueTopic, resultJson, fmt.Sprintf("compliance-runner-result-%d-%d", job.ID, job.RetryCount)); err != nil {
			w.logger.Error("failed to publish job result", zap.String("jobResult", string(resultJson)), zap.Error(err))
			w.logger.Sync()

		}
	}()

	resultJson, err := json.Marshal(result)
	if err != nil {
		w.logger.Error("failed to create job in progress json", zap.Error(err))
		w.logger.Sync()

		return true, false, err
	}

	if _, err := w.jq.Produce(ctx, ResultQueueTopic, resultJson, fmt.Sprintf("compliance-runner-inprogress-%d-%d", job.ID, job.RetryCount)); err != nil {
		w.logger.Error("failed to publish job in progress", zap.String("jobInProgress", string(resultJson)), zap.Error(err))
		w.logger.Sync()

	}

	w.logger.Info("running job", zap.ByteString("job", msg.Data()))
	w.logger.Sync()

	totalComplianceResultCount, err := w.RunJob(ctx, job)
	if err != nil {
		return true, false, err
	}

	result.TotalComplianceResultCount = &totalComplianceResultCount
	return true, false, nil
}

// **pollAPI runs every 15 seconds and cancels the process if needed**
func (w *Worker) fetchParameters(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.logger.Info("fetching parameters values")
			queryParams, err := w.coreClient.ListQueryParameters(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole}, coreApi.ListQueryParametersRequest{})
			if err != nil {
				w.logger.Error("failed to get query parameters", zap.Error(err))
			} else {
				w.queryParamsMu.Lock()
				w.queryParameters = queryParams.Items
				w.queryParamsMu.Unlock()
			}
		}
	}
}

// **pollAPI runs every 15 seconds and cancels the process if needed**
func (w *Worker) pollAPI(ctx context.Context, cancelFunc context.CancelFunc, msg jetstream.Msg) {
	var job Job

	if err := json.Unmarshal(msg.Data(), &job); err != nil {
		w.logger.Error("failed to unmarshal job msg")
		w.logger.Sync()

		return
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	w.logger.Info("Polling API...")
	stop, err := w.checkAPIResponse(strconv.Itoa(int(job.ParentJobID)))
	if err != nil {
		w.logger.Error("Failed to check compliance job status", zap.Uint("compliance-job-id", job.ParentJobID),
			zap.Uint("runner-job-id", job.ID), zap.Error(err))
	}
	if stop { // If API returns a special response
		w.canceledComplianceJobsMu.Lock()
		w.canceledComplianceJobs[job.ParentJobID] = true
		w.canceledComplianceJobsMu.Unlock()
		w.logger.Warn("Received stop signal from API! Cancelling job.")
		cancelFunc()
		return
	}

	for {
		select {
		case <-ticker.C:
			w.logger.Info("Polling API...")
			stop, err = w.checkAPIResponse(strconv.Itoa(int(job.ParentJobID)))
			if err != nil {
				w.logger.Error("Failed to check compliance job status", zap.Uint("compliance-job-id", job.ParentJobID),
					zap.Uint("runner-job-id", job.ID), zap.Error(err))
			}
			if stop { // If API returns a special response
				w.canceledComplianceJobsMu.Lock()
				w.canceledComplianceJobs[job.ParentJobID] = true
				w.canceledComplianceJobsMu.Unlock()
				w.logger.Warn("Received stop signal from API! Cancelling job.")
				cancelFunc()
				return
			}

		case <-ctx.Done(): // Stop if context is canceled
			w.logger.Info("Stopping API polling.")
			return
		}
	}
}

// **checkAPIResponse simulates an API request**
func (w *Worker) checkAPIResponse(jobId string) (bool, error) {
	clientCtx := &httpclient.Context{UserRole: authApi.AdminRole}

	status, err := w.schedulerClient.GetComplianceJobStatus(clientCtx, jobId)
	if err != nil {
		return false, err
	}
	if status.JobStatus == schedulerApi.ComplianceJobCanceled {
		return true, nil
	}
	return false, nil
}

func (w *Worker) Stop() error {
	w.steampipeConn.Conn().Close()
	steampipe.StopSteampipeService(w.logger)
	return nil
}
