package plugins

import (
	"context"
	"fmt"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/jobs/post-install-job/config"
	"github.com/opengovern/opensecurity/services/integration/models"
	"github.com/opengovern/opensecurity/services/integration/utils"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Migration struct {
}

func (m Migration) IsGitBased() bool {
	return true
}
func (m Migration) AttachmentFolderPath() string {
	return config.PluginsGitPath
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

	err = orm.AutoMigrate(&models.IntegrationPlugin{}, &models.IntegrationPluginBinary{})
	if err != nil {
		logger.Error("failed to auto migrate integration binaries", zap.Error(err))
		return err
	}

	err = filepath.Walk(config.PluginsGitPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = utils.ValidateAndLoadPlugin(orm, logger, content)
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("migrated plugin: %s", path))
		}

		return nil
	})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("plugins migrated"))

	return nil
}
