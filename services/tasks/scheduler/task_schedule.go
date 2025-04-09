package scheduler

import (
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/ticker"
	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opensecurity/pkg/utils"
	"github.com/opengovern/opensecurity/services/tasks/config"
	"github.com/opengovern/opensecurity/services/tasks/db"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

type NatsConfig struct {
	Stream         string
	Topic          string
	ResultTopic    string
	Consumer       string
	ResultConsumer string
}

type TaskScheduler struct {
	runSetupNatsStreams func(context.Context) error
	jq                  *jq.JobQueue
	db                  db.Database
	logger              *zap.Logger

	cfg config.Config

	TaskID     string
	ResultType string
	NatsConfig NatsConfig
	Interval   uint64
	Timeout    uint64
	vault      vault.VaultSourceConfig
}

func NewTaskScheduler(
	runSetupNatsStreams func(context.Context) error,
	logger *zap.Logger,
	db db.Database,
	jq *jq.JobQueue,

	cfg config.Config,
	vault vault.VaultSourceConfig,
	taskID, ResultType string, natsConfig NatsConfig, interval uint64, timeout uint64) *TaskScheduler {
	return &TaskScheduler{
		runSetupNatsStreams: runSetupNatsStreams,
		logger:              logger,
		db:                  db,
		jq:                  jq,

		cfg:   cfg,
		vault: vault,

		TaskID:     taskID,
		ResultType: ResultType,
		NatsConfig: natsConfig,
		Interval:   interval,
		Timeout:    timeout,
	}
}

func (s *TaskScheduler) Run(ctx context.Context) {
	s.logger.Info("Run task scheduler started", zap.String("task", s.TaskID),
		zap.Any("nats config", s.NatsConfig), zap.Uint64("interval", s.Interval))
	utils.EnsureRunGoroutine(func() {
		s.RunPublisher(ctx)
	})
	if s.Interval > 0 {
		utils.EnsureRunGoroutine(func() {
			s.CreateTask(ctx)
		})
	}
	utils.EnsureRunGoroutine(func() {
		s.logger.Fatal("RunTaskResponseConsumer exited", zap.Error(s.RunTaskResponseConsumer(ctx)))
	})
}

func (s *TaskScheduler) RunPublisher(ctx context.Context) {
	s.logger.Info("Scheduling publisher on a timer")

	t := ticker.NewTicker(time.Second*10, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		if err := s.runPublisher(ctx); err != nil {
			s.logger.Error("failed to run compliance publisher", zap.Error(err))
			continue
		}
	}
}

func (s *TaskScheduler) CreateTask(ctx context.Context) {
	s.logger.Info("Scheduling publisher on a timer")

	interval := time.Duration(s.Interval) * time.Minute
	t := ticker.NewTicker(interval, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		if err := s.createTask(ctx); err != nil {
			s.logger.Error("failed to run compliance publisher", zap.Error(err))
			continue
		}
	}
}
