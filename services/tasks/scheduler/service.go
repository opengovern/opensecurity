package scheduler

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opensecurity/services/tasks/config"
	"github.com/opengovern/opensecurity/services/tasks/db"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

type MainScheduler struct {
	jq     *jq.JobQueue
	db     db.Database
	logger *zap.Logger

	cfg   config.Config
	vault vault.VaultSourceConfig

	Tasks []TaskScheduler
}

var RunningTasks = make(map[string]bool)

func NewMainScheduler(cfg config.Config, logger *zap.Logger, db db.Database, vault vault.VaultSourceConfig) (*MainScheduler, error) {
	jq, err := jq.New(cfg.NATS.URL, logger)
	if err != nil {
		logger.Error("Failed to create job queue", zap.Error(err))
		return nil, err
	}

	return &MainScheduler{
		jq:     jq,
		db:     db,
		logger: logger,
		cfg:    cfg,
		vault:  vault,
	}, nil
}

func (s *MainScheduler) Start(ctx context.Context) error {
	tasks, err := s.db.GetTaskList()
	if err != nil {
		s.logger.Error("failed to get task list", zap.Error(err))
		return err
	}

	for _, task := range tasks {
		if _, ok := RunningTasks[task.ID]; ok {
			continue
		}
		var natsConfig NatsConfig
		if task.NatsConfig.Status != pgtype.Present {
			return fmt.Errorf("JSONB data is not present")
		}
		if err := json.Unmarshal(task.NatsConfig.Bytes, &natsConfig); err != nil {
			return fmt.Errorf("failed to unmarshal JSONB: %w", err)
		}

		err = s.SetupNats(ctx, task.ID, natsConfig)
		if err != nil {
			s.logger.Error("Failed to setup nats streams", zap.Error(err))
			return err
		}

		taskScheduler := NewTaskScheduler(
			func(ctx context.Context) error {
				return s.SetupNats(ctx, task.ID, natsConfig)
			},
			s.logger,
			s.db,
			s.jq,
			s.cfg,
			s.vault,
			task.ID,
			natsConfig)
		taskScheduler.Run(ctx)
		RunningTasks[task.ID] = true
	}
	return nil
}

func (s *MainScheduler) SetupNats(ctx context.Context, taskID string, natsConfig NatsConfig) error {
	s.logger.Info("Subscribing to stream", zap.String("task", taskID), zap.String("stream", natsConfig.Stream),
		zap.Strings("topics", []string{natsConfig.Topic, natsConfig.ResultTopic}))
	if err := s.jq.Stream(ctx, natsConfig.Stream, "task job queue", []string{natsConfig.Topic, natsConfig.ResultTopic}, 100); err != nil {
		s.logger.Error("Failed to stream to task queue", zap.String("task", taskID), zap.Error(err))
		return err
	}

	s.logger.Info("Creating or Updating Consumer", zap.String("task", taskID), zap.String("stream", natsConfig.Stream),
		zap.Strings("topics", []string{natsConfig.Topic}))
	if err := s.jq.CreateOrUpdateConsumer(ctx, natsConfig.Consumer, natsConfig.Stream,
		[]string{natsConfig.Topic}, jetstream.ConsumerConfig{
			Replicas:          1,
			AckPolicy:         jetstream.AckExplicitPolicy,
			DeliverPolicy:     jetstream.DeliverAllPolicy,
			MaxAckPending:     -1,
			AckWait:           time.Minute * 30,
			InactiveThreshold: time.Hour,
		}); err != nil {
		s.logger.Error("Failed to create consumer for task", zap.String("task", taskID), zap.Error(err))
		return err
	}
	return nil
}
