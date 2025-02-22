package discovery

import (
	"context"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/pkg/utils"
	"github.com/opengovern/opensecurity/services/compliance/client"
	config2 "github.com/opengovern/opensecurity/services/scheduler/config"
	"github.com/opengovern/opensecurity/services/scheduler/db"
	"go.uber.org/zap"
)

type Scheduler struct {
	conf             config2.SchedulerConfig
	logger           *zap.Logger
	complianceClient client.ComplianceServiceClient
	db               db.Database
	esClient         opengovernance.Client
}

func New(conf config2.SchedulerConfig, logger *zap.Logger, complianceClient client.ComplianceServiceClient, db db.Database, esClient opengovernance.Client) *Scheduler {
	return &Scheduler{
		conf:             conf,
		logger:           logger,
		complianceClient: complianceClient,
		db:               db,
		esClient:         esClient,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	utils.EnsureRunGoroutine(func() {
		s.OldResourceDeleter(ctx)
	})
}
