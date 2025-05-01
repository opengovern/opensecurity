package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/services/compliance/db"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/og-util/pkg/config"
	esSinkClient "github.com/opengovern/og-util/pkg/es/ingest/client"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/jobs/compliance-summarizer/types"
	coreClient "github.com/opengovern/opensecurity/services/core/client"
	integrationClient "github.com/opengovern/opensecurity/services/integration/client"
	"go.uber.org/zap"
)

type Config struct {
	PostgreSQL            config.Postgres
	ElasticSearch         config.ElasticSearch
	NATS                  config.NATS
	PrometheusPushAddress string
	Core                  config.OpenGovernanceService
	Integration           config.OpenGovernanceService
	EsSink                config.OpenGovernanceService
}

type Worker struct {
	db db.Database

	config   Config
	logger   *zap.Logger
	esClient opengovernance.Client
	jq       *jq.JobQueue

	coreClient        coreClient.CoreServiceClient
	integrationClient integrationClient.IntegrationServiceClient
	esSinkClient      esSinkClient.EsSinkServiceClient
}

var (
	ManualTrigger = os.Getenv("MANUAL_TRIGGER")
)

func NewWorker(
	config Config,
	logger *zap.Logger,
	ctx context.Context,
) (*Worker, error) {
	cfg := postgres.Config{
		Host:    config.PostgreSQL.Host,
		Port:    config.PostgreSQL.Port,
		User:    config.PostgreSQL.Username,
		Passwd:  config.PostgreSQL.Password,
		DB:      config.PostgreSQL.DB,
		SSLMode: config.PostgreSQL.SSLMode,
	}
	orm, err := postgres.NewClient(&cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("new postgres client: %w", err)
	}

	database := db.Database{Orm: orm}
	fmt.Println("Connected to the postgres database: ", config.PostgreSQL.DB)

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

	queueTopic := JobQueueTopic
	if ManualTrigger == "true" {
		queueTopic = JobQueueTopicManuals
	}

	if err := jq.Stream(ctx, StreamName, "compliance summarizer job runner queue", []string{queueTopic, ResultQueueTopic}, 1000); err != nil {
		return nil, err
	}

	w := &Worker{
		config:            config,
		logger:            logger,
		esClient:          esClient,
		db:                database,
		jq:                jq,
		coreClient:        coreClient.NewCoreServiceClient(config.Core.BaseURL),
		integrationClient: integrationClient.NewIntegrationServiceClient(config.Integration.BaseURL),
		esSinkClient:      esSinkClient.NewEsSinkServiceClient(logger, config.EsSink.BaseURL),
	}

	return w, nil
}

// Run is a blocking function so you may decide to call it in another goroutine.
// It runs a NATS consumer and it will close it when the given context is closed.
func (w *Worker) Run(ctx context.Context) error {
	w.logger.Info("starting to consume")

	queueTopic := JobQueueTopic
	service := "compliance-summarizer"
	consumer := ConsumerGroup
	if ManualTrigger == "true" {
		queueTopic = JobQueueTopicManuals
		consumer = ConsumerGroupManuals
		service = "compliance-summarizer-manuals"
	}

	consumeCtx, err := w.jq.Consume(ctx, service, StreamName, []string{queueTopic}, consumer, func(msg jetstream.Msg) {
		w.logger.Info("received a new job")

		if err := w.ProcessMessage(ctx, msg); err != nil {
			w.logger.Error("failed to process message", zap.Error(err))
		}
		err := msg.Ack()
		if err != nil {
			w.logger.Error("failed to ack message", zap.Error(err))
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

func (w *Worker) ProcessMessage(ctx context.Context, msg jetstream.Msg) error {
	startTime := time.Now()

	var job types.Job
	err := json.Unmarshal(msg.Data(), &job)
	if err != nil {
		return err
	}

	defer func() {
		result := JobResult{
			Job:       job,
			StartedAt: startTime,
			Status:    ComplianceSummarizerSucceeded,
			Error:     "",
		}

		if err != nil {
			result.Error = err.Error()
			result.Status = ComplianceSummarizerFailed
		}

		resultJson, err := json.Marshal(result)
		if err != nil {
			w.logger.Error("failed to create job result json", zap.Error(err))
			return
		}

		if _, err := w.jq.Produce(ctx, ResultQueueTopic, resultJson, fmt.Sprintf("job-result-%d-%d", job.ID, job.RetryCount)); err != nil {
			w.logger.Error("failed to publish job result", zap.String("jobResult", string(resultJson)), zap.Error(err))
		}
	}()

	runtime.GC()

	w.logger.Info("running job", zap.ByteString("job", msg.Data()))

	err = w.RunJob(ctx, job)
	if err != nil {
		w.logger.Info("failure while running job", zap.Error(err))
		return err
	}

	return nil
}

func (w *Worker) Stop() error {
	return nil
}
