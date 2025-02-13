package cloudql_init_job

import (
	"context"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/opencomply/services/integration/client"
	"github.com/opengovern/opencomply/services/integration/models"
	"go.uber.org/zap"
	"os"
)

type Job struct {
	logger            *zap.Logger
	cfg               Config
	integrationClient client.IntegrationServiceClient
}

func NewJob(logger *zap.Logger, cfg Config, integrationClient client.IntegrationServiceClient) *Job {
	return &Job{
		logger:            logger,
		cfg:               cfg,
		integrationClient: integrationClient,
	}
}

func (j *Job) Run(ctx context.Context) (*steampipe.Database, error) {
	db, err := postgres.NewClient(&postgres.Config{
		Host:    j.cfg.Postgres.Host,
		Port:    j.cfg.Postgres.Port,
		User:    j.cfg.Postgres.Username,
		Passwd:  j.cfg.Postgres.Password,
		DB:      "integration_types",
		SSLMode: j.cfg.Postgres.SSLMode,
	}, j.logger.Named("postgres"))
	if err != nil {
		j.logger.Error("failed to create postgres client", zap.Error(err))
		return nil, err
	}

	var integrations []models.IntegrationPlugin
	err = db.Find(&integrations).Error
	if err != nil {
		j.logger.Error("failed to get integration binaries", zap.Error(err))
		return nil, err
	}
	var integrationMap = make(map[string]*models.IntegrationPlugin)
	for _, integration := range integrations {
		integration := integration
		integrationMap[integration.IntegrationType.String()] = &integration
	}

	httpCtx := httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}

	integrationTypes, err := j.integrationClient.ListIntegrationTypes(&httpCtx)
	if err != nil {
		j.logger.Error("failed to list integration types", zap.Error(err))
		return nil, err
	}

	basePath := "/home/steampipe/.steampipe/plugins/hub.steampipe.io/plugins/turbot"
	for _, integrationType := range integrationTypes {
		describerConfig, err := j.integrationClient.GetIntegrationConfiguration(&httpCtx, integrationType)
		if err != nil {
			j.logger.Error("failed to get integration configuration", zap.Error(err))
			return nil, err
		}
		err = steampipe.PopulateSteampipeConfig(j.cfg.ElasticSearch, describerConfig.SteampipePluginName)
		if err != nil {
			return nil, err
		}

		if plugin, ok := integrationMap[integrationType]; ok {
			dirPath := basePath + "/" + describerConfig.SteampipePluginName + "@latest"
			// create directory if not exists
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				err := os.MkdirAll(dirPath, os.ModePerm)
				if err != nil {
					j.logger.Error("failed to create directory", zap.Error(err), zap.String("path", dirPath))
					return nil, err
				}
			}

			var cloudqlBinary string
			err = db.Raw("SELECT cloud_ql_plugin FROM integration_plugin_binaries WHERE plugin_id = ?", plugin.PluginID).Scan(&cloudqlBinary).Error
			if err != nil {
				j.logger.Error("failed to get plugin binary", zap.Error(err), zap.String("plugin_id", plugin.PluginID))
				return nil, err
			}

			// write the plugin to the file system
			pluginPath := dirPath + "/steampipe-plugin-" + describerConfig.SteampipePluginName + ".plugin"
			err := os.WriteFile(pluginPath, []byte(cloudqlBinary), 0777)
			if err != nil {
				j.logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", describerConfig.SteampipePluginName))
				return nil, err
			}

			cloudqlBinary = ""
		}
	}
	if j.cfg.Postgres.DB == "" {
		j.cfg.Postgres.DB = "integration"
	}
	if err := steampipe.PopulateOpenGovernancePluginSteampipeConfig(j.cfg.ElasticSearch, j.cfg.Postgres); err != nil {
		return nil, err
	}

	// execute command to start steampipe service
	cloudqlConn, err := steampipe.StartSteampipeServiceAndGetConnection(j.logger)
	if err != nil {
		j.logger.Error("failed to start steampipe service", zap.Error(err))
		return nil, err
	}

	return cloudqlConn, nil
}
