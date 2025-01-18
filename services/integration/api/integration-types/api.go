package integration_types

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-getter"
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/pkg/utils"
	"github.com/opengovern/opencomply/services/integration/api/models"
	"github.com/opengovern/opencomply/services/integration/db"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	models2 "github.com/opengovern/opencomply/services/integration/models"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type API struct {
	logger      *zap.Logger
	typeManager *integration_type.IntegrationTypeManager
	database    db.Database
}

func New(typeManager *integration_type.IntegrationTypeManager, database db.Database, logger *zap.Logger) *API {
	return &API{
		logger:      logger.Named("integration_types"),
		typeManager: typeManager,
		database:    database,
	}
}

func (a *API) Register(e *echo.Group) {
	e.GET("", httpserver.AuthorizeHandler(a.List, api.ViewerRole))
	e.GET("/:integration_type/resource-type/table/:table_name", httpserver.AuthorizeHandler(a.GetResourceTypeFromTableName, api.ViewerRole))
	e.GET("/:integration_type/table", httpserver.AuthorizeHandler(a.ListTables, api.ViewerRole))
	e.POST("/:integration_type/resource-type/label", httpserver.AuthorizeHandler(a.GetResourceTypesByLabels, api.ViewerRole))
	e.GET("/:integration_type/configuration", httpserver.AuthorizeHandler(a.GetConfiguration, api.ViewerRole))

	plugin := e.Group("/plugin")
	plugin.POST("/load/id/:id", httpserver.AuthorizeHandler(a.LoadPluginWithID, api.EditorRole))
	plugin.POST("/load/url/:http_url", httpserver.AuthorizeHandler(a.LoadPluginWithURL, api.EditorRole))
	plugin.DELETE("/uninstall/id/:id", httpserver.AuthorizeHandler(a.UninstallPlugin, api.EditorRole))
	plugin.DELETE("/plugin/:id/enable", httpserver.AuthorizeHandler(a.EnablePlugin, api.EditorRole))
	plugin.DELETE("/plugin/:id/disable", httpserver.AuthorizeHandler(a.DisablePlugin, api.EditorRole))
	plugin.GET("", httpserver.AuthorizeHandler(a.ListPlugins, api.ViewerRole))
	plugin.GET("/:id", httpserver.AuthorizeHandler(a.GetPlugin, api.ViewerRole))
	plugin.GET("/:id/integrations", httpserver.AuthorizeHandler(a.ListPluginIntegrations, api.ViewerRole))
	plugin.GET("/:id/credentials", httpserver.AuthorizeHandler(a.ListPluginCredentials, api.ViewerRole))
	plugin.POST("/:id/healthcheck", httpserver.AuthorizeHandler(a.HealthCheck, api.ViewerRole))
}

// List godoc
//
// @Summary			List integration types
// @Description		List integration types
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Success			200	{object} []string
// @Router			/integration/api/v1/integration_types [get]
func (a *API) List(c echo.Context) error {
	types := a.typeManager.GetIntegrationTypes()

	typesApi := make([]string, 0, len(types))
	for _, t := range types {
		typesApi = append(typesApi, t.String())
	}

	return c.JSON(200, typesApi)
}

// GetResourceTypeFromTableName godoc
//
// @Summary			Get resource type from table name
// @Description		Get resource type from table name
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			integration_type	path	string	true	"Integration type"
// @Param			table_name		path	string	true	"Table name"
// @Success			200	{object} models.GetResourceTypeFromTableNameResponse
// @Router			/integration/api/v1/integration_types/{integration_type}/resource-type/table/{table_name} [get]
func (a *API) GetResourceTypeFromTableName(c echo.Context) error {
	integrationType := c.Param("integration_type")
	tableName := c.Param("table_name")

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		resourceType, err := value.GetResourceTypeFromTableName(tableName)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		if resourceType != "" {
			res := models.GetResourceTypeFromTableNameResponse{
				ResourceType: resourceType,
			}
			return c.JSON(200, res)
		} else {
			return echo.NewHTTPError(404, "resource type not found")
		}
	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}

// GetConfiguration godoc
//
// @Summary			Get integration configuration
// @Description		Get integration configuration
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			integration_type	path	string	true	"Integration type"
// @Success			200	{object} interfaces.IntegrationConfiguration
// @Router			/integration/api/v1/integration_types/{integration_type}/configuration [get]
func (a *API) GetConfiguration(c echo.Context) error {
	integrationType := c.Param("integration_type")

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		conf, err := value.GetConfiguration()
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}

		return c.JSON(200, conf)
	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}

// GetResourceTypesByLabelsRequest in body

// GetResourceTypesByLabels godoc
//
// @Summary			Get resource types by labels
// @Description		Get resource types by labels
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			integration_type	path	string	true	"Integration type"
// @Param			request	body	models.GetResourceTypesByLabelsRequest	true	"Request"
// @Success			200	{object} models.GetResourceTypesByLabelsResponse
// @Router			/integration/api/v1/integration_types/{integration_type}/resource-type/label [post]
func (a *API) GetResourceTypesByLabels(c echo.Context) error {
	integrationType := c.Param("integration_type")

	req := new(models.GetResourceTypesByLabelsRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(400, "invalid request")
	}

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		rts, err := value.GetResourceTypesByLabels(req.Labels)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		res := models.GetResourceTypesByLabelsResponse{
			ResourceTypes: make(map[string]*models.ResourceTypeConfiguration),
		}
		for k, v := range rts {
			res.ResourceTypes[k] = utils.GetPointer(models.ApiResourceTypeConfiguration(v))
		}
		return c.JSON(200, res)
	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}

// ListTables godoc
//
// @Summary			List tables
// @Description		List tables
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			integration_type	path	string	true	"Integration type"
// @Success			200	{object} models.ListTablesResponse
// @Router			/integration/api/v1/integration_types/{integration_type}/table [get]
func (a *API) ListTables(c echo.Context) error {
	integrationType := c.Param("integration_type")

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		tables, err := value.ListAllTables()
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		return c.JSON(200, models.ListTablesResponse{Tables: tables})
	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}

// LoadPluginWithID godoc
//
// @Summary			Load plugin with the given plugin ID
// @Description		Load plugin with the given plugin ID
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/load/id/{id} [post]
func (a *API) LoadPluginWithID(c echo.Context) error {
	pluginID := c.Param("id")

	plugin, err := a.database.GetPluginByID(pluginID)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}
	baseDir := "/integration-types"

	// create tmp directory if not exists
	if _, err = os.Stat(baseDir); os.IsNotExist(err) {
		if err = os.Mkdir(baseDir, os.ModePerm); err != nil {
			a.logger.Error("failed to create tmp directory", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create tmp directory")
		}
	}

	// download files from urls

	if plugin.URL == "" {
		return echo.NewHTTPError(http.StatusNotFound, "plugin url is empty")
	}
	url := plugin.URL
	// remove existing files
	if err = os.RemoveAll(baseDir + "/integarion_type"); err != nil {
		a.logger.Error("failed to remove existing files", zap.Error(err), zap.String("id", pluginID), zap.String("path", baseDir+"/integarion_type"))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to remove existing files")
	}

	downloader := getter.Client{
		Src:  url,
		Dst:  baseDir + "/integarion_type",
		Mode: getter.ClientModeDir,
	}
	err = downloader.Get()
	if err != nil {
		a.logger.Error("failed to get integration binaries", zap.Error(err), zap.String("id", pluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get integration binaries")
	}

	// read integration-plugin file
	integrationPlugin, err := os.ReadFile(baseDir + "/integarion_type/integration-plugin")
	if err != nil {
		a.logger.Error("failed to open integration-plugin file", zap.Error(err), zap.String("id", pluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open integration-plugin file")
	}
	cloudqlPlugin, err := os.ReadFile(baseDir + "/integarion_type/cloudql-plugin")
	if err != nil {
		a.logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("id", pluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open cloudql-plugin file")
	}

	a.logger.Info("done reading files", zap.String("id", pluginID), zap.String("url", url), zap.String("integrationType", plugin.IntegrationType.String()), zap.Int("integrationPluginSize", len(integrationPlugin)), zap.Int("cloudqlPluginSize", len(cloudqlPlugin)))

	err = a.database.UpdatePlugin(models2.IntegrationPlugin{
		PluginID:        pluginID,
		IntegrationType: plugin.IntegrationType,
		InstallState:    models2.IntegrationTypeInstallStateInstalled,
		URL:             plugin.URL,

		IntegrationPlugin: integrationPlugin,
		CloudQlPlugin:     cloudqlPlugin,
	})
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", pluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
	}
	return c.NoContent(http.StatusOK)
}

// LoadPluginWithURL godoc
//
// @Summary			Load plugin with the given plugin URL
// @Description		Load plugin with the given plugin URL
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			http_url 	path	string	true	"plugin url"
// @Success			200
// @Router			/integration/api/v1/plugin/load/url/{http_url} [post]
func (a *API) LoadPluginWithURL(c echo.Context) error {
	pluginURL := c.Param("http_url")

	var err error

	baseDir := "/integration-types"

	// create tmp directory if not exists
	if _, err = os.Stat(baseDir); os.IsNotExist(err) {
		if err = os.Mkdir(baseDir, os.ModePerm); err != nil {
			a.logger.Error("failed to create tmp directory", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create tmp directory")
		}
	}

	// download files from urls

	if pluginURL == "" {
		return echo.NewHTTPError(http.StatusNotFound, "plugin url is empty")
	}
	url := pluginURL
	// remove existing files
	if err = os.RemoveAll(baseDir + "/integarion_type"); err != nil {
		a.logger.Error("failed to remove existing files", zap.Error(err), zap.String("url", url), zap.String("path", baseDir+"/integarion_type"))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to remove existing files")
	}

	downloader := getter.Client{
		Src:  url,
		Dst:  baseDir + "/integarion_type",
		Mode: getter.ClientModeDir,
	}
	err = downloader.Get()
	if err != nil {
		a.logger.Error("failed to get integration binaries", zap.Error(err), zap.String("url", url))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get integration binaries")
	}

	// read integration-plugin file
	integrationPlugin, err := os.ReadFile(baseDir + "/integarion_type/integration-plugin")
	if err != nil {
		a.logger.Error("failed to open integration-plugin file", zap.Error(err), zap.String("url", url))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open integration-plugin file")
	}
	cloudqlPlugin, err := os.ReadFile(baseDir + "/integarion_type/cloudql-plugin")
	if err != nil {
		a.logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("url", url))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open cloudql-plugin file")
	}
	// read manifest file
	manifestFile, err := os.Open(baseDir + "/integarion_type/manifest.yaml")
	if err != nil {
		a.logger.Error("failed to open manifest file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open manifest file")
	}
	defer manifestFile.Close()
	var m models2.Manifest
	// decode yaml
	if err = yaml.NewDecoder(manifestFile).Decode(&m); err != nil {
		a.logger.Error("failed to decode manifest", zap.Error(err), zap.String("url", url))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to decode manifest file")
	}

	a.logger.Info("done reading files", zap.String("url", url), zap.String("url", url), zap.String("integrationType", m.IntegrationType.String()), zap.Int("integrationPluginSize", len(integrationPlugin)), zap.Int("cloudqlPluginSize", len(cloudqlPlugin)))

	err = a.database.UpdatePlugin(models2.IntegrationPlugin{
		PluginID:          m.PluginID,
		IntegrationType:   m.IntegrationType,
		InstallState:      models2.IntegrationTypeInstallStateInstalled,
		OperationalStatus: models2.IntegrationPluginOperationalStatusEnabled,
		URL:               url,

		IntegrationPlugin: integrationPlugin,
		CloudQlPlugin:     cloudqlPlugin,
	})
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", m.PluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
	}
	return c.NoContent(http.StatusOK)
}

// UninstallPlugin godoc
//
// @Summary			Load plugin with the given plugin URL
// @Description		Load plugin with the given plugin URL
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/uninstall/id/{id} [delete]
func (a *API) UninstallPlugin(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	integrations, err := a.database.ListIntegration([]integration.Type{plugin.IntegrationType})
	if err != nil {
		a.logger.Error("failed to list integrations", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list integrations")
	}
	if len(integrations) > 0 {
		return echo.NewHTTPError(http.StatusNotFound, "integration type has integrations")
	}

	credentials, err := a.database.ListCredentialsFiltered(nil, []string{plugin.IntegrationType.String()})
	if err != nil {
		a.logger.Error("failed to list credentials", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list credentials")
	}
	if len(credentials) > 0 {
		return echo.NewHTTPError(http.StatusNotFound, "integration type has credentials")
	}

	plugin.InstallState = models2.IntegrationTypeInstallStateNotInstalled

	err = a.database.UpdatePlugin(*plugin)
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.PluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
	}

	return c.NoContent(http.StatusOK)
}

// EnablePlugin godoc
//
// @Summary			Enable plugin with the given plugin id
// @Description		Enable plugin with the given plugin id
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id}/enable [put]
func (a *API) EnablePlugin(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusEnabled

	err = a.database.UpdatePlugin(*plugin)
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.PluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
	}

	return c.NoContent(http.StatusOK)
}

// DisablePlugin godoc
//
// @Summary			Disable plugin with the given plugin id
// @Description		Disable plugin with the given plugin id
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id}/disable [put]
func (a *API) DisablePlugin(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusDisabled

	err = a.database.UpdatePlugin(*plugin)
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.PluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
	}

	return c.NoContent(http.StatusOK)
}

// ListPlugins godoc
//
// @Summary			List plugins
// @Description		List plugins
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Success			200
// @Router			/integration/api/v1/plugin [get]
func (a *API) ListPlugins(c echo.Context) error {
	plugins, err := a.database.ListPlugins()
	if err != nil {
		a.logger.Error("failed to list plugins", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list plugins")
	}

	return c.JSON(http.StatusOK, plugins)
}

// GetPlugin godoc
//
// @Summary			Get plugin
// @Description		Get plugin
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id} [get]
func (a *API) GetPlugin(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	return c.JSON(http.StatusOK, plugin)
}

// ListPluginIntegrations godoc
//
// @Summary			List plugin integrations
// @Description		List plugin integrations
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id}/integrations [get]
func (a *API) ListPluginIntegrations(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	integrations, err := a.database.ListIntegration([]integration.Type{plugin.IntegrationType})
	if err != nil {
		a.logger.Error("failed to list integrations", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list integrations")
	}

	return c.JSON(http.StatusOK, integrations)
}

// ListPluginCredentials godoc
//
// @Summary			List plugin credentials
// @Description		List plugin credentials
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id}/credentials [get]
func (a *API) ListPluginCredentials(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	credentials, err := a.database.ListCredentialsFiltered(nil, []string{plugin.IntegrationType.String()})
	if err != nil {
		a.logger.Error("failed to list credentials", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list credentials")
	}

	return c.JSON(http.StatusOK, credentials)
}

// HealthCheck godoc
//
// @Summary			Health check
// @Description		Health check
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id}/healthcheck [post]
func (a *API) HealthCheck(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[plugin.IntegrationType]; ok {
		err := value.Ping()
		if err == nil {
			return c.JSON(http.StatusOK, "plugin is healthy")
		}

		err = a.typeManager.RetryRebootIntegrationType(plugin)
		if err != nil {
			return echo.NewHTTPError(400, fmt.Sprintf("plugin was found unhealthy and failed to reboot with error: %v", err))
		} else {
			return c.JSON(http.StatusOK, "plugin was found unhealthy and got successfully rebooted")
		}
	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}
