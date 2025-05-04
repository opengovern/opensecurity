package db

import (
	"github.com/opengovern/opensecurity/services/scheduler/db/model"
	"gorm.io/gorm"
)

type Database struct {
	ORM *gorm.DB
}

func (db Database) Initialize() error {
	return db.ORM.AutoMigrate(
		&model.ComplianceJob{}, &model.ComplianceSummarizer{}, &model.ComplianceRunner{}, &model.IntegrationHealthCheckJob{},
		&model.DescribeIntegrationJob{}, &model.IntegrationDiscovery{},
		&model.JobSequencer{}, &model.QueryRunnerJob{},
		&model.QuickScanSequence{}, &model.FrameworkValidation{}, &model.ManualDiscoverySchedule{},
		&model.ResourceTypeDescribedCount{},
	)
}
