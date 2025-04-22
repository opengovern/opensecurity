package cloudql_init_job

import (
	"context"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/opensecurity/services/integration/client"
	"github.com/opengovern/opensecurity/services/integration/models"
	taskModels "github.com/opengovern/opensecurity/services/tasks/db/models"
	"go.uber.org/zap"
	"os"
	"strings"
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

	tasksDb, err := postgres.NewClient(&postgres.Config{
		Host:    j.cfg.Postgres.Host,
		Port:    j.cfg.Postgres.Port,
		User:    j.cfg.Postgres.Username,
		Passwd:  j.cfg.Postgres.Password,
		DB:      "task",
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

	// Integrations
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

	// Tasks
	var tasks []taskModels.Task
	if err := tasksDb.Model(&taskModels.Task{}).Where("is_enabled = ?", true).Find(&tasks).Error; err != nil {
		j.logger.Error("failed to get tasks", zap.Error(err))
		return nil, err
	}
	for _, task := range tasks {
		if task.SteampipePluginName == "" {
			j.logger.Debug("steampipe plugin name not found - skipping task", zap.String("id", task.ID))
			continue
		}

		steampipePluginName := strings.ReplaceAll(task.SteampipePluginName, "steampipe-plugin-", "")
		var cloudqlBinary string
		err = db.Raw("SELECT cloud_ql_plugin FROM task_binaries WHERE task_id = ?", task.ID).Scan(&cloudqlBinary).Error
		if err != nil {
			j.logger.Error("failed to get plugin binary", zap.Error(err), zap.String("task_id", task.ID))
			return nil, err
		}

		if cloudqlBinary == "" {
			j.logger.Debug("steampipe plugin binary not found - skipping task", zap.String("id", task.ID))
			continue
		}
		err = steampipe.PopulateSteampipeConfig(j.cfg.ElasticSearch, steampipePluginName)
		if err != nil {
			return nil, err
		}

		dirPath := basePath + "/" + steampipePluginName + "@latest"
		// create directory if not exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				j.logger.Error("failed to create directory", zap.Error(err), zap.String("path", dirPath))
				return nil, err
			}
		}

		// write the plugin to the file system
		pluginPath := dirPath + "/" + task.SteampipePluginName + ".plugin"
		err := os.WriteFile(pluginPath, []byte(cloudqlBinary), 0777)
		if err != nil {
			j.logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", steampipePluginName))
			return nil, err
		}

		cloudqlBinary = ""
	}

	if j.cfg.Postgres.DB == "" {
		j.cfg.Postgres.DB = "integration"
	}
	if err := steampipe.PopulateOpenGovernancePluginSteampipeConfig(j.cfg.ElasticSearch, j.cfg.Postgres, j.cfg.Integration.BaseURL); err != nil {
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

func (j *Job) ReloadSinglePlugin(ctx context.Context, pluginID string) error {
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

	var plugin models.IntegrationPlugin
	err = db.Model(models.IntegrationPlugin{}).Where("plugin_id = ?", pluginID).First(&plugin).Error
	if err != nil {
		j.logger.Error("failed to get integration binaries", zap.Error(err))
		return err
	}

	httpCtx := httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}

	basePath := "/home/steampipe/.steampipe/plugins/hub.steampipe.io/plugins/turbot"
	describerConfig, err := j.integrationClient.GetIntegrationConfiguration(&httpCtx, plugin.IntegrationType.String())
	if err != nil {
		j.logger.Error("failed to get integration configuration", zap.Error(err))
		return err
	}
	err = steampipe.PopulateSteampipeConfig(j.cfg.ElasticSearch, describerConfig.SteampipePluginName)
	if err != nil {
		return err
	}

	dirPath := basePath + "/" + describerConfig.SteampipePluginName + "@latest"
	// create directory if not exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			j.logger.Error("failed to create directory", zap.Error(err), zap.String("path", dirPath))
			return err
		}
	}

	var cloudqlBinary string
	err = db.Raw("SELECT cloud_ql_plugin FROM integration_plugin_binaries WHERE plugin_id = ?", plugin.PluginID).Scan(&cloudqlBinary).Error
	if err != nil {
		j.logger.Error("failed to get plugin binary", zap.Error(err), zap.String("plugin_id", plugin.PluginID))
		return err
	}

	// write the plugin to the file system
	pluginPath := dirPath + "/steampipe-plugin-" + describerConfig.SteampipePluginName + ".plugin"
	err = os.WriteFile(pluginPath, []byte(cloudqlBinary), 0777)
	if err != nil {
		j.logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", describerConfig.SteampipePluginName))
		return err
	}

	cloudqlBinary = ""

	return nil
}

func (j *Job) RemoveSinglePlugin(ctx context.Context, pluginID string) error {
	httpCtx := httpclient.Context{Ctx: ctx, UserRole: api.AdminRole}

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
	var plugin models.IntegrationPlugin
	err = db.Model(models.IntegrationPlugin{}).Where("plugin_id = ?", pluginID).First(&plugin).Error
	if err != nil {
		j.logger.Error("failed to get integration binaries", zap.Error(err))
		return err
	}

	basePath := "/home/steampipe/.steampipe/plugins/hub.steampipe.io/plugins/turbot"

	describerConfig, err := j.integrationClient.GetIntegrationConfiguration(&httpCtx, plugin.IntegrationType.String())
	if err != nil {
		j.logger.Error("failed to get integration configuration", zap.Error(err))
		return err
	}

	err = steampipe.RemoveSpecFile(describerConfig.SteampipePluginName)
	if err != nil {
		return err
	}

	dirPath := basePath + "/" + describerConfig.SteampipePluginName + "@latest"

	pluginPath := dirPath + "/steampipe-plugin-" + describerConfig.SteampipePluginName + ".plugin"

	err = os.Remove(pluginPath)
	if err != nil {
		j.logger.Error("failed to remove plugin file", zap.Error(err))
		return err
	}

	return nil
}
