package compliance_quick_run

import (
	"context"
	"time"

	"github.com/opengovern/og-util/pkg/jq"
	coreClient "github.com/opengovern/opencomply/services/core/client"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/ticker"
	"github.com/opengovern/opencomply/pkg/utils"
	complianceClient "github.com/opengovern/opencomply/services/compliance/client"
	"github.com/opengovern/opencomply/services/scheduler/config"
	"github.com/opengovern/opencomply/services/scheduler/db"

	"go.uber.org/zap"
)

const JobSchedulingInterval = 10 * time.Second

type JobScheduler struct {
	runSetupNatsStreams func(context.Context) error
	conf                config.SchedulerConfig
	logger              *zap.Logger
	db                  db.Database
	jq                  *jq.JobQueue
	esClient            opengovernance.Client

	complianceClient    complianceClient.ComplianceServiceClient
	coreClient      coreClient.CoreServiceClient
}

func New(
	runSetupNatsStreams func(context.Context) error,
	conf config.SchedulerConfig,
	logger *zap.Logger,
	db db.Database,
	jq *jq.JobQueue,
	esClient opengovernance.Client,

	complianceClient complianceClient.ComplianceServiceClient,
	coreClient coreClient.CoreServiceClient,
) *JobScheduler {
	return &JobScheduler{
		runSetupNatsStreams: runSetupNatsStreams,
		conf:                conf,
		logger:              logger,
		db:                  db,
		jq:                  jq,
		esClient:            esClient,
	
		complianceClient:    complianceClient,
		coreClient:      coreClient,
	}
}

func (s *JobScheduler) Run(ctx context.Context) {
	utils.EnsureRunGoroutine(func() {
		s.RunPublisher(ctx)
	})
	utils.EnsureRunGoroutine(func() {
		s.logger.Fatal("RunAuditJobResultsConsumer exited", zap.Error(s.RunAuditJobResultsConsumer(ctx)))
	})
}

func (s *JobScheduler) RunPublisher(ctx context.Context) {
	s.logger.Info("Scheduling publisher on a timer")

	t := ticker.NewTicker(JobSchedulingInterval, time.Second*10)
	defer t.Stop()

	for ; ; <-t.C {
		if err := s.runPublisher(ctx); err != nil {
			s.logger.Error("failed to run compliance publisher", zap.Error(err))
			continue
		}
	}
}
