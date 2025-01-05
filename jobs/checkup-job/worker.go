package checkup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/opencomply/jobs/checkup-job/config"
	authClient "github.com/opengovern/opencomply/services/auth/client"
	coreClient "github.com/opengovern/opencomply/services/core/client"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/opencomply/services/integration/client"
	"go.uber.org/zap"
)

type Worker struct {
	id                string
	jq                *jq.JobQueue
	logger            *zap.Logger
	config            config.WorkerConfig
	integrationClient client.IntegrationServiceClient
	authClient        authClient.AuthServiceClient
	coreClient    coreClient.CoreServiceClient
}

func NewWorker(
	id string,
	natsURL string,
	logger *zap.Logger,
	integrationBaseURL string,
	authBaseURL string,
	coreBaseURL string,
	config config.WorkerConfig,
	ctx context.Context,
) (w *Worker, err error) {
	if id == "" {
		return nil, fmt.Errorf("'id' must be set to a non empty string")
	}

	w = &Worker{id: id}
	defer func() {
		if err != nil && w != nil {
			w.Stop()
		}
	}()

	jq, err := jq.New(natsURL, logger)
	if err != nil {
		return nil, err
	}

	if err := jq.Stream(ctx, StreamName, "checkup job queue", []string{JobsQueueName, ResultsQueueName}, 1000); err != nil {
		return nil, err
	}

	w.jq = jq

	w.logger = logger

	w.integrationClient = client.NewIntegrationServiceClient(integrationBaseURL)
	w.authClient = authClient.NewAuthClient(authBaseURL)
	w.coreClient = coreClient.NewCoreServiceClient(coreBaseURL)
	w.config = config
	return w, nil
}

func (w *Worker) Run(ctx context.Context) error {
	consumeCtx, err := w.jq.Consume(
		ctx,
		"checkup-service",
		StreamName,
		[]string{JobsQueueName},
		"checkup-service",
		func(msg jetstream.Msg) {
			var job Job
			if err := json.Unmarshal(msg.Data(), &job); err != nil {
				w.logger.Error("Failed to unmarshal task", zap.Error(err))

				// sending ack for message because we cannot do anything
				// more by repeating the process.
				if err = msg.Ack(); err != nil {
					w.logger.Error("Failed to ack the message", zap.Error(err))
				}

				return
			}

			w.logger.Info("Processing job", zap.Int("jobID", int(job.JobID)))

			result := job.Do(w.integrationClient, w.authClient, w.coreClient, w.logger, w.config)

			bytes, err := json.Marshal(result)
			if err != nil {
				return
			}

			w.logger.Info("Publishing job result", zap.Int("jobID", int(job.JobID)))

			if _, err := w.jq.Produce(ctx, ResultsQueueName, bytes, fmt.Sprintf("job-result-%d", result.JobID)); err != nil {
				w.logger.Error("Failed to send results to queue: %s", zap.Error(err))
			}

			if err := msg.Ack(); err != nil {
				w.logger.Error("Failed to ack the message", zap.Error(err))
			}
		},
	)
	if err != nil {
		return err
	}

	w.logger.Error("Waiting indefinitely for messages. To exit press CTRL+C")
	<-ctx.Done()
	consumeCtx.Drain()
	consumeCtx.Stop()

	return nil
}

func (w *Worker) Stop() {
}
