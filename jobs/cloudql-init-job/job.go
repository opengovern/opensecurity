package cloudql_init_job

import (
	"context"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/integration-type/models"
	"github.com/opengovern/opencomply/services/integration/client"
	"go.uber.org/zap"
	"os"
	"os/exec"
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

func (j *Job) Run(ctx context.Context) error {
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
		return err
	}

	var integrations []models.IntegrationTypeBinaries
	err = db.Find(&integrations).Error
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
			pluginPath := dirPath + "/steampipe-plugin-" + describerConfig.SteampipePluginName + ".plugin"
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

	// execute command to start steampipe service
	_, err = steampipe.StartSteampipeServiceAndGetConnection(j.logger)
	if err != nil {
		j.logger.Error("failed to start steampipe service", zap.Error(err))
		return err
	}

	cmd := exec.Command("steampipe", "service", "stop", "--force")
	err = cmd.Run()
	if err != nil {
		j.logger.Error("first stop failed", zap.Error(err))
		return err
	}
	//NOTE: stop must be called twice. it's not a mistake
	cmd = exec.Command("steampipe", "service", "stop", "--force")
	err = cmd.Run()
	if err != nil {
		j.logger.Error("second stop failed", zap.Error(err))
		return err
	}

	return nil
}
