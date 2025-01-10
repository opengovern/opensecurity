package compliance_quick_run_job

import (
	"context"
	"encoding/json"
	"fmt"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opencomply/services/integration/client"
	"os"
	"time"

	"github.com/opengovern/opencomply/services/scheduler/db/model"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/og-util/pkg/config"
	esSinkClient "github.com/opengovern/og-util/pkg/es/ingest/client"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/steampipe"
	complianceClient "github.com/opengovern/opencomply/services/compliance/client"
	coreClient "github.com/opengovern/opencomply/services/core/client"
	"go.uber.org/zap"
)

type Config struct {
	ElasticSearch config.ElasticSearch
	NATS          config.NATS
	Compliance    config.OpenGovernanceService
	Core          config.OpenGovernanceService
	Integration   config.OpenGovernanceService
	EsSink        config.OpenGovernanceService
	Steampipe     config.Postgres
}

type Worker struct {
	config            Config
	logger            *zap.Logger
	steampipeConn     *steampipe.Database
	esClient          opengovernance.Client
	jq                *jq.JobQueue
	complianceClient  complianceClient.ComplianceServiceClient
	coreClient        coreClient.CoreServiceClient
	integrationClient client.IntegrationServiceClient
	sinkClient        esSinkClient.EsSinkServiceClient
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

	httpCtx := httpclient.Context{Ctx: ctx, UserRole: authApi.ViewerRole}

	integrationTypes, err := integrationClient.ListIntegrationTypes(&httpCtx)
	if err != nil {
		logger.Error("failed to list integration types", zap.Error(err))
		return nil, err
	}

	for _, integrationType := range integrationTypes {
		describerConfig, err := integrationClient.GetIntegrationConfiguration(&httpCtx, integrationType)
		if err != nil {
			logger.Error("failed to get integration configuration", zap.Error(err))
			return nil, err
		}
		err = steampipe.PopulateSteampipeConfig(config.ElasticSearch, describerConfig.SteampipePluginName)
		if err != nil {
			return nil, err
		}
	}
	if err := steampipe.PopulateOpenGovernancePluginSteampipeConfig(config.ElasticSearch, config.Steampipe); err != nil {
		return nil, err
	}

	time.Sleep(2 * time.Minute)

	steampipeConn, err := steampipe.StartSteampipeServiceAndGetConnection(logger)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	jq, err := jq.New(config.NATS.URL, logger)
	if err != nil {
		return nil, err
	}

	if err := jq.Stream(ctx, StreamName, "audit job queue", []string{JobQueueTopic, ResultQueueTopic}, 1000); err != nil {
		logger.Error("failed to create stream", zap.Error(err))
		return nil, err
	}

	return &Worker{
		config:            config,
		logger:            logger,
		steampipeConn:     steampipeConn,
		esClient:          esClient,
		jq:                jq,
		complianceClient:  complianceClient.NewComplianceClient(config.Compliance.BaseURL),
		coreClient:        coreClient.NewCoreServiceClient(config.Core.BaseURL),
		sinkClient:        esSinkClient.NewEsSinkServiceClient(logger, config.EsSink.BaseURL),
		integrationClient: integrationClient,
	}, nil
}

// Run is a blocking function so you may decide to call it in another goroutine.
// It runs a NATS consumer and it will close it when the given context is closed.
func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("starting to consume")

	consumeCtx, err := w.jq.ConsumeWithConfig(ctx, ConsumerGroup, StreamName, []string{JobQueueTopic},
		jetstream.ConsumerConfig{
			DeliverPolicy:     jetstream.DeliverAllPolicy,
			AckPolicy:         jetstream.AckExplicitPolicy,
			AckWait:           time.Hour,
			MaxDeliver:        1,
			InactiveThreshold: time.Hour,
			Replicas:          1,
			MemoryStorage:     false,
		}, nil,
		func(msg jetstream.Msg) {
			w.logger.Info("received a new job")
			w.logger.Info("committing")
			if err := msg.InProgress(); err != nil {
				w.logger.Error("failed to send the initial in progress message", zap.Error(err), zap.Any("msg", msg))
			}
			ticker := time.NewTicker(15 * time.Second)
			go func() {
				for range ticker.C {
					if err := msg.InProgress(); err != nil {
						w.logger.Error("failed to send an in progress message", zap.Error(err), zap.Any("msg", msg))
					}
				}
			}()

			err := w.ProcessMessage(ctx, msg)
			if err != nil {
				w.logger.Error("failed to process message", zap.Error(err))
			}
			ticker.Stop()

			if err := msg.Ack(); err != nil {
				w.logger.Error("failed to send the ack message", zap.Error(err), zap.Any("msg", msg))
			}

			w.logger.Info("processing a job completed")
		})
	if err != nil {
		return err
	}

	w.logger.Info("consuming")

	<-ctx.Done()
	consumeCtx.Drain()
	consumeCtx.Stop()

	return nil
}

func (w *Worker) ProcessMessage(ctx context.Context, msg jetstream.Msg) (err error) {
	var job AuditJob

	if err := json.Unmarshal(msg.Data(), &job); err != nil {
		return err
	}

	result := JobResult{
		JobID:          job.JobID,
		Status:         model.ComplianceJobRunnersInProgress,
		FailureMessage: "",
	}

	defer func() {
		if err != nil {
			result.FailureMessage = err.Error()
			result.Status = model.ComplianceJobFailed
		} else {
			result.Status = model.ComplianceJobSucceeded
		}

		resultJson, err := json.Marshal(result)
		if err != nil {
			w.logger.Error("failed to create job result json", zap.Error(err))
			return
		}

		if _, err := w.jq.Produce(ctx, ResultQueueTopic, resultJson, fmt.Sprintf("audit-job-result-%d", job.JobID)); err != nil {
			w.logger.Error("failed to publish job result", zap.String("jobResult", string(resultJson)), zap.Error(err))
		}
	}()

	resultJson, err := json.Marshal(result)
	if err != nil {
		w.logger.Error("failed to create job in progress json", zap.Error(err))
		return err
	}

	if _, err := w.jq.Produce(ctx, ResultQueueTopic, resultJson, fmt.Sprintf("audit-job-inprogress-%d", job.JobID)); err != nil {
		w.logger.Error("failed to publish job in progress", zap.String("jobInProgress", string(resultJson)), zap.Error(err))
	}

	w.logger.Info("running job", zap.ByteString("job", msg.Data()))

	err = w.RunJob(ctx, &job)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Stop() error {
	w.steampipeConn.Conn().Close()
	steampipe.StopSteampipeService(w.logger)
	return nil
}
