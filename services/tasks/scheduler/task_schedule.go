package scheduler

import (
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/ticker"
	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opensecurity/pkg/utils"
	"github.com/opengovern/opensecurity/services/tasks/config"
	"github.com/opengovern/opensecurity/services/tasks/db"
	"github.com/opengovern/opensecurity/services/tasks/db/models"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"time"
)

type NatsConfig struct {
	Stream         string `json:"stream"`
	Topic          string `json:"topic"`
	ResultTopic    string `json:"result_topic"`
	Consumer       string `json:"consumer"`
	ResultConsumer string `json:"result_consumer"`
}

type TaskScheduler struct {
	runSetupNatsStreams func(context.Context) error
	jq                  *jq.JobQueue
	db                  db.Database
	logger              *zap.Logger

	cfg config.Config

	TaskID           string
	Timeout          float64
	NatsConfig       NatsConfig
	vault            vault.VaultSourceConfig
	TaskRunSchedules []models.TaskRunSchedule
}

func NewTaskScheduler(
	runSetupNatsStreams func(context.Context) error,
	logger *zap.Logger,
	db db.Database,
	jq *jq.JobQueue,

	cfg config.Config,
	vault vault.VaultSourceConfig,
	taskID string, natsConfig NatsConfig,
	taskRunSchedules []models.TaskRunSchedule) *TaskScheduler {
	return &TaskScheduler{
		runSetupNatsStreams: runSetupNatsStreams,
		logger:              logger,
		db:                  db,
		jq:                  jq,

		cfg:   cfg,
		vault: vault,

		TaskID:           taskID,
		NatsConfig:       natsConfig,
		TaskRunSchedules: taskRunSchedules,
	}
}

func (s *TaskScheduler) Run(ctx context.Context) {
	s.logger.Info("Run task scheduler started", zap.String("task", s.TaskID),
		zap.Any("nats config", s.NatsConfig))
	utils.EnsureRunGoroutine(func() {
		s.RunPublisher(ctx)
	})

	for _, taskRunSchedule := range s.TaskRunSchedules {
		utils.EnsureRunGoroutine(func() {
			s.CreateTask(ctx, taskRunSchedule)
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

func (s *TaskScheduler) CreateTask(ctx context.Context, runSchedule models.TaskRunSchedule) {
	s.logger.Info("Scheduling publisher on a timer")

	interval := time.Duration(runSchedule.Frequency) * time.Minute
	t := ticker.NewTicker(interval, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		if err := s.createTask(ctx, runSchedule); err != nil {
			s.logger.Error("failed to run compliance publisher", zap.Error(err))
			continue
		}
	}
}
