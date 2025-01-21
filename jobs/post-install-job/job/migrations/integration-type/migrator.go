package integration_type

import (
	"context"
	"fmt"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/jobs/post-install-job/db"
	"github.com/opengovern/opencomply/services/integration/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Migration struct {
}

func (m Migration) IsGitBased() bool {
	return true
}
func (m Migration) AttachmentFolderPath() string {
	return config.ConfigzGitPath
}

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "integration_types",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{ORM: orm}

	err = dbm.ORM.AutoMigrate(&models.IntegrationPlugin{})
	if err != nil {
		logger.Error("failed to auto migrate integration binaries", zap.Error(err))
		return err
	}

	parser := GitParser{}
	err = parser.ExtractIntegrations(logger)
	if err != nil {
		return err
	}

	for _, iPlugin := range parser.Integrations.Plugins {
		plugin, pluginBinary, err := parser.ExtractIntegrationBinaries(logger, iPlugin)
		if err != nil {
			return err
		}
		if plugin == nil {
			continue
		}

		err = dbm.ORM.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "plugin_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"id", "integration_type", "name", "tier", "description", "icon",
				"availability", "source_code", "package_type", "url", "tags"}),
		}).Create(plugin).Error
		if err != nil {
			logger.Error("failed to create integration binary", zap.Error(err))
			return err
		}

		err = dbm.ORM.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "plugin_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"integration_plugin": gorm.Expr("CASE WHEN ? <> '' THEN ? ELSE integration_plugin_binaries.integration_plugin END", pluginBinary.IntegrationPlugin, pluginBinary.IntegrationPlugin),
				"cloud_ql_plugin":    gorm.Expr("CASE WHEN ? <> '' THEN ? ELSE integration_plugin_binaries.cloud_ql_plugin END", pluginBinary.CloudQlPlugin, pluginBinary.CloudQlPlugin),
			}),
		}).Create(pluginBinary).Error
		if err != nil {
			logger.Error("failed to create integration binary", zap.Error(err))
			return err
		}
	}

	return nil
}
