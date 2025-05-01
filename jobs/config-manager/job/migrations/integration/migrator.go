package integration

import (
	"context"
	"fmt"

	"github.com/opengovern/opensecurity/jobs/config-manager/config"
	"github.com/opengovern/opensecurity/jobs/config-manager/db"
	integrationModels "github.com/opengovern/opensecurity/services/integration/models"
	"gorm.io/gorm/clause"

	"github.com/opengovern/og-util/pkg/postgres"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Migration struct {
}

func (m Migration) IsGitBased() bool {
	return true
}
func (m Migration) AttachmentFolderPath() string {
	return config.IntegrationsGitPath
}

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "integration",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{ORM: orm}

	parser := GitParser{}
	err = parser.ExtractConnectionGroups(m.AttachmentFolderPath())
	if err != nil {
		return err
	}

	err = dbm.ORM.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&integrationModels.IntegrationGroup{}).Where("1 = 1").Unscoped().Delete(&integrationModels.IntegrationGroup{}).Error
		if err != nil {
			logger.Error("failed to delete integration groups", zap.Error(err))
			return err
		}

		for _, integrationGroup := range parser.integrationGroups {
			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&integrationGroup).Error
			if err != nil {
				logger.Error("failed to create integration group", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failure in integration group transaction: %w", err)
	}

	return nil
}
