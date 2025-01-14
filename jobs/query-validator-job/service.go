package query_validator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	cloudql_init_job "github.com/opengovern/opencomply/jobs/cloudql-init-job"
	"github.com/opengovern/opencomply/services/integration/client"
	"strconv"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/steampipe"
	"go.uber.org/zap"
)

type Config struct {
	ElasticSearch         config.ElasticSearch
	NATS                  config.NATS
	Compliance            config.OpenGovernanceService
	Integration           config.OpenGovernanceService
	Inventory             config.OpenGovernanceService
	Core                  config.OpenGovernanceService
	EsSink                config.OpenGovernanceService
	Steampipe             config.Postgres
	PostgresPlugin        config.Postgres
	PrometheusPushAddress string
}

type Worker struct {
	config            Config
	logger            *zap.Logger
	steampipeConn     *steampipe.Database
	esClient          opengovernance.Client
	jq                *jq.JobQueue
	integrationClient client.IntegrationServiceClient
}

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
	}, integrationClient)
	err := pluginJob.Run(ctx)
	if err != nil {
		logger.Error("failed to run plugin job", zap.Error(err))
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

	if err := jq.Stream(ctx, StreamName, "compliance runner job queue", []string{JobQueueTopic, JobResultQueueTopic}, 1000); err != nil {
		return nil, err
	}

	w := &Worker{
		config:            config,
		logger:            logger,
		steampipeConn:     steampipeConn,
		esClient:          esClient,
		jq:                jq,
		integrationClient: integrationClient,
	}

	return w, nil
}

// Run is a blocking function so you may decide to call it in another goroutine.
// It runs a NATS consumer and it will close it when the given context is closed.
func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("starting to consume")

	queueTopic := JobQueueTopic
	consumer := ConsumerGroup

	consumeCtx, err := w.jq.ConsumeWithConfig(ctx, consumer, StreamName, []string{queueTopic},
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
	var job Job

	if err = json.Unmarshal(msg.Data(), &job); err != nil {
		return err
	}

	w.logger.Info("job message delivered", zap.String("jobID", strconv.Itoa(int(job.ID))))

	result := JobResult{
		ID:             job.ID,
		QueryType:      job.QueryType,
		QueryId:        job.QueryId,
		ControlId:      job.ControlId,
		Status:         QueryValidatorInProgress,
		FailureMessage: "",
	}

	defer func() {
		if err != nil {
			result.FailureMessage = err.Error()
			result.Status = QueryValidatorFailed
		} else {
			result.Status = QueryValidatorSucceeded
		}

		w.logger.Info("job is finished with status", zap.String("ID", strconv.Itoa(int(job.ID))), zap.String("status", string(result.Status)))

		resultJson, err := json.Marshal(result)
		if err != nil {
			w.logger.Error("failed to create job result json", zap.Error(err))
			return
		}

		if _, err := w.jq.Produce(ctx, JobResultQueueTopic, resultJson, fmt.Sprintf("query-runner-result-%d", job.ID)); err != nil {
			w.logger.Error("failed to publish job result", zap.String("jobResult", string(resultJson)), zap.Error(err))
		}
	}()

	resultJson, err := json.Marshal(result)
	if err != nil {
		w.logger.Error("failed to create job in progress json", zap.Error(err))
		return err
	}

	if _, err := w.jq.Produce(ctx, JobResultQueueTopic, resultJson, fmt.Sprintf("query-validator-inprogress-%d", job.ID)); err != nil {
		w.logger.Error("failed to publish job in progress", zap.String("jobInProgress", string(resultJson)), zap.Error(err))
	}

	w.logger.Info("running job", zap.ByteString("job", msg.Data()), zap.Any("job object", job))

	err = w.RunJob(ctx, job)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Stop() error {
	w.steampipeConn.Conn().Close()
	err := steampipe.StopSteampipeService(w.logger)
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) Initialize(ctx context.Context, j Job) error {
	providerAccountID := "all"

	err := w.steampipeConn.SetConfigTableValue(ctx, steampipe.OpenGovernanceConfigKeyIntegrationID, providerAccountID)
	if err != nil {
		w.logger.Error("failed to set account id", zap.Error(err))
		return err
	}
	err = w.steampipeConn.SetConfigTableValue(ctx, steampipe.OpenGovernanceConfigKeyClientType, "compliance")
	if err != nil {
		w.logger.Error("failed to set client type", zap.Error(err))
		return err
	}

	return nil
}
