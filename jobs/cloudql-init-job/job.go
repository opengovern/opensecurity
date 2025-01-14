package query_runner

import (
	"context"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/integration-type/models"
	"github.com/opengovern/opencomply/services/integration/client"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

type Job struct {
	logger            *zap.Logger
	db                *gorm.DB
	cfg               Config
	integrationClient client.IntegrationServiceClient
}

func NewJob(logger *zap.Logger, cfg Config, db *gorm.DB, integrationClient client.IntegrationServiceClient) *Job {
	return &Job{
		logger:            logger,
		db:                db,
		cfg:               cfg,
		integrationClient: integrationClient,
	}
}

func (j *Job) Run(ctx context.Context) error {
	var integrations []models.IntegrationTypeBinaries
	err := j.db.Find(&integrations).Error
	if err != nil {
		j.logger.Error("failed to get integration binaries", zap.Error(err))
		return err
	}
	var integrationMap = make(map[string]*models.IntegrationTypeBinaries)
	for _, integration := range integrations {
		integration := integration
		integrationMap[integration.IntegrationType.String()] = &integration
	}

	httpCtx := httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}

	integrationTypes, err := j.integrationClient.ListIntegrationTypes(&httpCtx)
	if err != nil {
		j.logger.Error("failed to list integration types", zap.Error(err))
		return err
	}

	basePath := "/home/steampipe/.steampipe/plugins/hub.steampipe.io/plugins/turbot"
	for _, integrationType := range integrationTypes {
		describerConfig, err := j.integrationClient.GetIntegrationConfiguration(&httpCtx, integrationType)
		if err != nil {
			j.logger.Error("failed to get integration configuration", zap.Error(err))
			return err
		}
		err = steampipe.PopulateSteampipeConfig(j.cfg.ElasticSearch, describerConfig.SteampipePluginName)
		if err != nil {
			return err
		}

		if integrationBin, ok := integrationMap[integrationType]; ok {
			dirPath := basePath + "/" + describerConfig.SteampipePluginName + "@latest"
			// create directory if not exists
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				err := os.MkdirAll(dirPath, os.ModePerm)
				if err != nil {
					j.logger.Error("failed to create directory", zap.Error(err), zap.String("path", dirPath))
					return err
				}
			}

			// write the plugin to the file system
			pluginPath := dirPath + "/" + describerConfig.SteampipePluginName + ".plugin"
			err := os.WriteFile(pluginPath, integrationBin.CloudQlPlugin, 0777)
			if err != nil {
				j.logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", describerConfig.SteampipePluginName))
				return err
			}
		}
	}
	if err := steampipe.PopulateOpenGovernancePluginSteampipeConfig(j.cfg.ElasticSearch, j.cfg.Steampipe); err != nil {
		return err
	}

	return nil
}
