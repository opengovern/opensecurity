package integration_types

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	hczap "github.com/zaffka/zap-to-hclog"
	"golang.org/x/net/context"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-getter"
	plugin2 "github.com/hashicorp/go-plugin"
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	"github.com/opengovern/opencomply/pkg/utils"
	"github.com/opengovern/opencomply/services/integration/api/models"
	"github.com/opengovern/opencomply/services/integration/db"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	models2 "github.com/opengovern/opencomply/services/integration/models"
	"go.uber.org/zap"
	"sort"
)

const (
	TemplateDeploymentPath          string = "/integrations/deployment-template.yaml"
	TemplateManualsDeploymentPath   string = "/integrations/deployment-template-manuals.yaml"
	TemplateScaledObjectPath        string = "/integrations/scaled-object-template.yaml"
	TemplateManualsScaledObjectPath string = "/integrations/scaled-object-template-manuals.yaml"
)

type API struct {
	logger      *zap.Logger
	typeManager *integration_type.IntegrationTypeManager
	database    db.Database
	kubeClient  client.Client
}

func New(typeManager *integration_type.IntegrationTypeManager, database db.Database, logger *zap.Logger, kubeClient client.Client) *API {
	return &API{
		logger:      logger.Named("integration_types"),
		typeManager: typeManager,
		database:    database,
		kubeClient:  kubeClient,
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
	plugin.POST("/:id/enable", httpserver.AuthorizeHandler(a.EnablePlugin, api.EditorRole))
	plugin.POST("/:id/disable", httpserver.AuthorizeHandler(a.DisablePlugin, api.EditorRole))
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

	//// read manifest file
	manifestFile, err := os.ReadFile(baseDir + "/integarion_type/manifest.yaml")
	if err != nil {
		a.logger.Error("failed to open manifest file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open manifest file")
	}
	a.logger.Info("manifestFile", zap.String("file", string(manifestFile)))

	var m models2.Manifest
	// decode yaml
	if err := yaml.Unmarshal(manifestFile, &m); err != nil {
		a.logger.Error("failed to decode manifest", zap.Error(err), zap.String("url", url))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to decode manifest")
	}

	a.logger.Info("done reading files", zap.String("id", pluginID), zap.String("url", url), zap.String("integrationType", plugin.IntegrationType.String()), zap.Int("integrationPluginSize", len(integrationPlugin)), zap.Int("cloudqlPluginSize", len(cloudqlPlugin)))

	plugin.IntegrationPlugin = integrationPlugin
	plugin.CloudQlPlugin = cloudqlPlugin
	plugin.DescriberURL = m.DescriberURL
	plugin.DescriberTag = m.DescriberTag
	plugin.InstallState = models2.IntegrationTypeInstallStateInstalled
	plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusEnabled

	err = a.database.UpdatePlugin(*plugin)
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", pluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
	}

	err = a.LoadPlugin(c.Request().Context(), *plugin)
	if err != nil {
		a.logger.Error("failed to load plugin", zap.Error(err), zap.String("id", pluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load plugin")
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
	//// read manifest file
	manifestFile, err := os.ReadFile(baseDir + "/integarion_type/manifest.yaml")
	if err != nil {
		a.logger.Error("failed to open manifest file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open manifest file")
	}
	a.logger.Info("manifestFile", zap.String("file", string(manifestFile)))

	var m models2.Manifest
	// decode yaml
	if err := yaml.Unmarshal(manifestFile, &m); err != nil {
		a.logger.Error("failed to decode manifest", zap.Error(err), zap.String("url", url))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to decode manifest file")
	}

	a.logger.Info("done reading files", zap.String("url", url), zap.String("url", url), zap.String("integrationType", m.IntegrationType.String()), zap.Int("integrationPluginSize", len(integrationPlugin)), zap.Int("cloudqlPluginSize", len(cloudqlPlugin)))

	plugin, err := a.database.GetPluginByURL(url)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}

	if plugin == nil {
		plugin = &models2.IntegrationPlugin{
			PluginID:          m.IntegrationType.String(),
			IntegrationType:   m.IntegrationType,
			DescriberURL:      m.DescriberURL,
			DescriberTag:      m.DescriberTag,
			InstallState:      models2.IntegrationTypeInstallStateInstalled,
			OperationalStatus: models2.IntegrationPluginOperationalStatusEnabled,
			URL:               url,

			IntegrationPlugin: integrationPlugin,
			CloudQlPlugin:     cloudqlPlugin,
		}
		err = a.database.CreatePlugin(*plugin)
		if err != nil {
			a.logger.Error("failed to create plugin", zap.Error(err), zap.String("id", m.IntegrationType.String()))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create plugin")
		}
	} else {
		plugin.PluginID = m.IntegrationType.String()
		plugin.IntegrationType = m.IntegrationType
		plugin.DescriberURL = m.DescriberURL
		plugin.DescriberTag = m.DescriberTag
		plugin.InstallState = models2.IntegrationTypeInstallStateInstalled
		plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusEnabled
		plugin.URL = url
		plugin.IntegrationPlugin = integrationPlugin
		plugin.CloudQlPlugin = cloudqlPlugin
		err = a.database.UpdatePlugin(*plugin)
		if err != nil {
			a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", m.IntegrationType.String()))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
		}
	}
	err = a.LoadPlugin(c.Request().Context(), *plugin)
	if err != nil {
		a.logger.Error("failed to load plugin", zap.Error(err), zap.String("id", plugin.PluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load plugin")
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
	plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusDisabled

	err = a.UnLoadPlugin(c.Request().Context(), *plugin)
	if err != nil {
		a.logger.Error("failed to unload plugin", zap.Error(err), zap.String("id", plugin.PluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unload plugin")
	}

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

	err = a.database.InactiveIntegrationType(plugin.IntegrationType)
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.PluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
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
	perPageStr := c.QueryParam("per_page")
	cursorStr := c.QueryParam("cursor")
	filteredEnabled := c.QueryParam("enabled")
	hasIntegration := c.QueryParam("has_integration")
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")
	var perPage, cursor int64
	if perPageStr != "" {
		perPage, _ = strconv.ParseInt(perPageStr, 10, 64)
	}
	if cursorStr != "" {
		cursor, _ = strconv.ParseInt(cursorStr, 10, 64)
	}

	plugins, err := a.database.ListPlugins()
	if err != nil {
		a.logger.Error("failed to list integration types", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list integration types")
	}

	var items []models.IntegrationPlugin
	for _, plugin := range plugins {
		integrations, err := a.database.ListIntegrationsByFilters(nil, []string{plugin.IntegrationType.String()}, nil, nil)
		if err != nil {
			a.logger.Error("failed to list integrations", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to list integrations")
		}
		if hasIntegration == "true" {
			if len(integrations) == 0 {
				continue
			}
		}
		count := models.IntegrationTypeIntegrationCount{}
		for _, i := range integrations {
			count.Total += 1
			if i.State == integration.IntegrationStateActive {
				count.Active += 1
			}
			if i.State == integration.IntegrationStateInactive {
				count.Inactive += 1
			}
			if i.State == integration.IntegrationStateArchived {
				count.Archived += 1
			}
			if i.State == integration.IntegrationStateSample {
				count.Demo += 1
			}
		}
		if plugin.OperationalStatus == models2.IntegrationPluginOperationalStatusDisabled {
			if filteredEnabled == "true" {
				continue
			}
		}
		items = append(items, models.IntegrationPlugin{
			ID:                       plugin.ID,
			PluginID:                 plugin.PluginID,
			IntegrationType:          plugin.IntegrationType.String(),
			InstallState:             string(plugin.InstallState),
			OperationalStatus:        string(plugin.OperationalStatus),
			OperationalStatusUpdates: plugin.OperationalStatusUpdates,
			URL:                      plugin.URL,
			Tier:                     plugin.Tier,
			Description:              plugin.Description,
			Icon:                     plugin.Icon,
			Availability:             plugin.Availability,
			SourceCode:               plugin.SourceCode,
			PackageType:              plugin.PackageType,
			DescriberURL:             plugin.DescriberURL,
			Name:                     plugin.Name,
			Count:                    count,
		})
	}

	totalCount := len(items)
	if sortOrder == "desc" {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ID > items[j].ID
		})
		if sortBy == "count" {
			sort.Slice(items, func(i, j int) bool {
				if items[i].Count.Total == items[j].Count.Total {
					return items[i].ID < items[j].ID
				}
				return items[i].Count.Total > items[j].Count.Total
			})
		}
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].ID < items[j].ID
		})
		if sortBy == "count" {
			sort.Slice(items, func(i, j int) bool {
				if items[i].Count.Total == items[j].Count.Total {
					return items[i].ID < items[j].ID
				}
				return items[i].Count.Total < items[j].Count.Total
			})
		}
	}

	if perPage != 0 {
		if cursor == 0 {
			items = utils.Paginate(1, perPage, items)
		} else {
			items = utils.Paginate(cursor, perPage, items)
		}
	}

	return c.JSON(http.StatusOK, models.IntegrationPluginListResponse{
		Items:      items,
		TotalCount: totalCount,
	})
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

	return c.JSON(http.StatusOK, models.IntegrationPlugin{
		ID:                       plugin.ID,
		PluginID:                 plugin.PluginID,
		IntegrationType:          plugin.IntegrationType.String(),
		InstallState:             string(plugin.InstallState),
		OperationalStatus:        string(plugin.OperationalStatus),
		OperationalStatusUpdates: plugin.OperationalStatusUpdates,
		URL:                      plugin.URL,
		Tier:                     plugin.Tier,
		Description:              plugin.Description,
		Icon:                     plugin.Icon,
		Availability:             plugin.Availability,
		SourceCode:               plugin.SourceCode,
		PackageType:              plugin.PackageType,
		DescriberURL:             plugin.DescriberURL,
		Name:                     plugin.Name,
	})
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

func (a *API) LoadPlugin(ctx context.Context, plugin models2.IntegrationPlugin) error {
	// create directory for plugins if not exists
	baseDir := "/plugins"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.Mkdir(baseDir, os.ModePerm)
		if err != nil {
			a.logger.Error("failed to create plugins directory", zap.Error(err))
			return nil
		}
	}
	pluginName := plugin.IntegrationType.String()

	// write the plugin to the file system
	pluginPath := filepath.Join(baseDir, plugin.IntegrationType.String()+".so")
	err := os.WriteFile(pluginPath, plugin.IntegrationPlugin, 0755)
	if err != nil {
		a.logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", pluginName))
		return err
	}
	hcLogger := hczap.Wrap(a.logger)

	var client *plugin2.Client
	if v, ok := a.typeManager.Clients[plugin.IntegrationType]; ok {
		client = v
	} else {
		client = plugin2.NewClient(&plugin2.ClientConfig{
			HandshakeConfig: interfaces.HandshakeConfig,
			Plugins:         map[string]plugin2.Plugin{pluginName: &interfaces.IntegrationTypePlugin{}},
			Cmd:             exec.Command(pluginPath),
			Logger:          hcLogger,
			Managed:         true,
		})
		a.typeManager.Clients[plugin.IntegrationType] = client
	}

	rpcClient, err := client.Client()
	if err != nil {
		a.logger.Error("failed to create plugin client", zap.Error(err), zap.String("plugin", pluginName), zap.String("path", pluginPath))
		client.Kill()
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginName)
	if err != nil {
		a.logger.Error("failed to dispense plugin", zap.Error(err), zap.String("plugin", pluginName), zap.String("path", pluginPath))
		client.Kill()
		return err
	}

	var itInterface interfaces.IntegrationType
	var ok2 bool
	if v, ok := a.typeManager.IntegrationTypes[plugin.IntegrationType]; ok {
		itInterface = v
	} else {
		itInterface, ok2 = raw.(interfaces.IntegrationType)
		if !ok2 {
			a.logger.Error("failed to cast plugin to integration type", zap.String("plugin", pluginName), zap.String("path", pluginPath))
			client.Kill()
			return err
		}
		a.typeManager.IntegrationTypes[plugin.IntegrationType] = itInterface
	}

	a.typeManager.PingLocks[plugin.IntegrationType] = &sync.Mutex{}

	err = a.EnableIntegrationTypeHelper(ctx, plugin.IntegrationType.String())
	if err != nil {
		a.logger.Error("failed to enable integration type describer", zap.Error(err))
		return err
	}

	return nil
}

func (a *API) UnLoadPlugin(ctx context.Context, plugin models2.IntegrationPlugin) error {
	err := a.DisableIntegrationTypeHelper(ctx, plugin.IntegrationType.String())
	if err != nil {
		a.logger.Error("failed to disable integration type describer", zap.Error(err))
		return err
	}

	if _, ok := a.typeManager.Clients[plugin.IntegrationType]; ok {
		a.typeManager.Clients[plugin.IntegrationType].Kill()
		delete(a.typeManager.Clients, plugin.IntegrationType)
	}
	if _, ok := a.typeManager.IntegrationTypes[plugin.IntegrationType]; ok {
		delete(a.typeManager.IntegrationTypes, plugin.IntegrationType)
	}
	if _, ok := a.typeManager.PingLocks[plugin.IntegrationType]; ok {
		delete(a.typeManager.PingLocks, plugin.IntegrationType)
	}

	return nil
}

func (a *API) DisableIntegrationTypeHelper(ctx context.Context, integrationTypeName string) error {
	plugin, err := a.database.GetPluginByIntegrationType(integrationTypeName)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	var integrationTypes []integration.Type
	integrationTypes = append(integrationTypes, integration.Type(integrationTypeName))

	integrations, err := a.database.ListIntegration(integrationTypes)
	if err != nil {
		a.logger.Error("failed to list credentials", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list credential")
	}
	if len(integrations) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "integration type contains integrations, you can not disable it")
	}

	currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "current namespace lookup failed")
	}
	integrationType, ok := a.typeManager.GetIntegrationTypeMap()[integration.Type(integrationTypeName)]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "invalid integration type")
	}
	cnf, err := integrationType.GetConfiguration()
	if err != nil {
		a.logger.Error("failed to get configuration", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get configuration"+err.Error())
	}

	// Scheduled deployment
	var describerDeployment appsv1.Deployment
	err = a.kubeClient.Get(ctx, client.ObjectKey{
		Namespace: currentNamespace,
		Name:      cnf.DescriberDeploymentName,
	}, &describerDeployment)
	if err != nil {
		a.logger.Error("failed to get manual deployment", zap.Error(err))
	} else {
		err = a.kubeClient.Delete(ctx, &describerDeployment)
		if err != nil {
			a.logger.Error("failed to delete deployment", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete deployment")
		}
	}

	// Manual deployment
	var describerDeploymentManuals appsv1.Deployment
	err = a.kubeClient.Get(ctx, client.ObjectKey{
		Namespace: currentNamespace,
		Name:      cnf.DescriberDeploymentName + "-manuals",
	}, &describerDeploymentManuals)
	if err != nil {
		a.logger.Error("failed to get manual deployment", zap.Error(err))
	} else {
		err = a.kubeClient.Delete(ctx, &describerDeploymentManuals)
		if err != nil {
			a.logger.Error("failed to delete manual deployment", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete manual deployment")
		}
	}

	kedaEnabled, ok := os.LookupEnv("KEDA_ENABLED")
	if !ok {
		kedaEnabled = "false"
	}
	if strings.ToLower(kedaEnabled) == "true" {
		// Scheduled ScaledObject
		var describerScaledObject kedav1alpha1.ScaledObject
		err = a.kubeClient.Get(ctx, client.ObjectKey{
			Namespace: currentNamespace,
			Name:      cnf.DescriberDeploymentName + "-scaled-object",
		}, &describerScaledObject)
		if err != nil {
			a.logger.Error("failed to get scaled object", zap.Error(err))
		} else {
			err = a.kubeClient.Delete(ctx, &describerScaledObject)
			if err != nil {
				a.logger.Error("failed to delete scaled object", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete scaled object")
			}
		}

		// Manual ScaledObject
		var describerScaledObjectManuals kedav1alpha1.ScaledObject
		err = a.kubeClient.Get(ctx, client.ObjectKey{
			Namespace: currentNamespace,
			Name:      cnf.DescriberDeploymentName + "-manuals-scaled-object",
		}, &describerScaledObjectManuals)
		if err != nil {
			a.logger.Error("failed to get manual scaled object", zap.Error(err))
		} else {
			err = a.kubeClient.Delete(ctx, &describerScaledObjectManuals)
			if err != nil {
				a.logger.Error("failed to delete manual scaled object", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete manual scaled object")
			}
		}
	}
	return nil
}

func (a *API) EnableIntegrationTypeHelper(ctx context.Context, integrationTypeName string) error {
	currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "current namespace lookup failed")
	}

	plugin, err := a.database.GetPluginByIntegrationType(integrationTypeName)
	if err != nil {
		a.logger.Error("failed to get integration type", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get integration type")
	}
	kedaEnabled, ok := os.LookupEnv("KEDA_ENABLED")
	if !ok {
		kedaEnabled = "false"
	}

	// Scheduled deployment
	var describerDeployment appsv1.Deployment
	templateDeploymentFile, err := os.Open(TemplateDeploymentPath)
	if err != nil {
		a.logger.Error("failed to open template deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template deployment file")
	}
	defer templateDeploymentFile.Close()

	data, err := ioutil.ReadAll(templateDeploymentFile)
	if err != nil {
		a.logger.Error("failed to read template deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template deployment file")
	}

	err = yaml.Unmarshal(data, &describerDeployment)
	if err != nil {
		a.logger.Error("failed to unmarshal template deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template deployment file")
	}

	integrationType, ok := a.typeManager.GetIntegrationTypeMap()[integration.Type(integrationTypeName)]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "invalid integration type")
	}
	cnf, err := integrationType.GetConfiguration()
	if err != nil {
		a.logger.Error("failed to get integration type configuration", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get integration type configuration")
	}

	describerDeployment.ObjectMeta.Name = cnf.DescriberDeploymentName
	describerDeployment.ObjectMeta.Namespace = currentNamespace
	if kedaEnabled == "true" {
		describerDeployment.Spec.Replicas = aws.Int32(0)
	} else {
		describerDeployment.Spec.Replicas = aws.Int32(5)
	}
	describerDeployment.Spec.Selector.MatchLabels["app"] = cnf.DescriberDeploymentName
	describerDeployment.Spec.Template.ObjectMeta.Labels["app"] = cnf.DescriberDeploymentName
	describerDeployment.Spec.Template.Spec.ServiceAccountName = "og-describer"

	container := describerDeployment.Spec.Template.Spec.Containers[0]
	container.Name = cnf.DescriberDeploymentName
	container.Image = fmt.Sprintf("%s:%s", plugin.DescriberURL, plugin.DescriberTag)
	container.Command = []string{cnf.DescriberRunCommand}
	natsUrl, ok := os.LookupEnv("NATS_URL")
	if ok {
		container.Env = append(container.Env, v1.EnvVar{
			Name:  "NATS_URL",
			Value: natsUrl,
		})
	}
	describerDeployment.Spec.Template.Spec.Containers[0] = container

	newDeployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cnf.DescriberDeploymentName,
			Namespace: currentNamespace,
			Labels: map[string]string{
				"app": cnf.DescriberDeploymentName,
			},
		},
		Spec: describerDeployment.Spec,
	}

	err = a.kubeClient.Create(ctx, &newDeployment)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}

	// Manual deployment
	var describerDeploymentManuals appsv1.Deployment
	templateManualsDeploymentFile, err := os.Open(TemplateManualsDeploymentPath)
	if err != nil {
		a.logger.Error("failed to open template manuals deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template manuals deployment file")
	}
	defer templateManualsDeploymentFile.Close()

	data, err = ioutil.ReadAll(templateManualsDeploymentFile)
	if err != nil {
		a.logger.Error("failed to read template manuals deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template manuals deployment file")
	}

	err = yaml.Unmarshal(data, &describerDeploymentManuals)
	if err != nil {
		a.logger.Error("failed to unmarshal template manuals deployment file", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template manuals deployment file")
	}

	describerDeploymentManuals.ObjectMeta.Name = cnf.DescriberDeploymentName + "-manuals"
	describerDeploymentManuals.ObjectMeta.Namespace = currentNamespace
	if kedaEnabled == "true" {
		describerDeploymentManuals.Spec.Replicas = aws.Int32(0)
	} else {
		describerDeploymentManuals.Spec.Replicas = aws.Int32(2)
	}
	describerDeploymentManuals.Spec.Selector.MatchLabels["app"] = cnf.DescriberDeploymentName + "-manuals"
	describerDeploymentManuals.Spec.Template.ObjectMeta.Labels["app"] = cnf.DescriberDeploymentName + "-manuals"
	describerDeploymentManuals.Spec.Template.Spec.ServiceAccountName = "og-describer"

	containerManuals := describerDeploymentManuals.Spec.Template.Spec.Containers[0]
	containerManuals.Name = cnf.DescriberDeploymentName
	containerManuals.Image = fmt.Sprintf("%s:%s", plugin.DescriberURL, plugin.DescriberTag)
	containerManuals.Command = []string{cnf.DescriberRunCommand}
	natsUrl, ok = os.LookupEnv("NATS_URL")
	if ok {
		containerManuals.Env = append(containerManuals.Env, v1.EnvVar{
			Name:  "NATS_URL",
			Value: natsUrl,
		})
	}
	describerDeploymentManuals.Spec.Template.Spec.Containers[0] = containerManuals

	newDeploymentManuals := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cnf.DescriberDeploymentName + "-manuals",
			Namespace: currentNamespace,
			Labels: map[string]string{
				"app": cnf.DescriberDeploymentName + "-manuals",
			},
		},
		Spec: describerDeploymentManuals.Spec,
	}

	err = a.kubeClient.Create(ctx, &newDeploymentManuals)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}

	if strings.ToLower(kedaEnabled) == "true" {
		// Scheduled ScaledObject
		var describerScaledObject kedav1alpha1.ScaledObject
		templateScaledObjectFile, err := os.Open(TemplateScaledObjectPath)
		if err != nil {
			a.logger.Error("failed to open template scaledobject file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template scaledobject file")
		}
		defer templateScaledObjectFile.Close()

		data, err = ioutil.ReadAll(templateScaledObjectFile)
		if err != nil {
			a.logger.Error("failed to read template manuals deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template scaledobject file")
		}

		err = yaml.Unmarshal(data, &describerScaledObject)
		if err != nil {
			a.logger.Error("failed to unmarshal template deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template deployment file")
		}

		describerScaledObject.Spec.ScaleTargetRef.Name = cnf.DescriberDeploymentName

		trigger := describerScaledObject.Spec.Triggers[0]
		trigger.Metadata["stream"] = cnf.NatsStreamName
		soNatsUrl, _ := os.LookupEnv("SCALED_OBJECT_NATS_URL")
		trigger.Metadata["natsServerMonitoringEndpoint"] = soNatsUrl
		trigger.Metadata["consumer"] = cnf.NatsConsumerGroup + "-service"
		describerScaledObject.Spec.Triggers[0] = trigger

		newScaledObject := kedav1alpha1.ScaledObject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cnf.DescriberDeploymentName + "-scaled-object",
				Namespace: currentNamespace,
			},
			Spec: describerScaledObject.Spec,
		}

		err = a.kubeClient.Create(ctx, &newScaledObject)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}

		// Manual ScaledObject
		var describerScaledObjectManuals kedav1alpha1.ScaledObject
		templateManualsScaledObjectFile, err := os.Open(TemplateManualsScaledObjectPath)
		if err != nil {
			a.logger.Error("failed to open template manuals scaledobject file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to open template manuals scaledobject file")
		}
		defer templateManualsScaledObjectFile.Close()

		data, err = ioutil.ReadAll(templateManualsScaledObjectFile)
		if err != nil {
			a.logger.Error("failed to read template manuals deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read template manuals scaledobject file")
		}

		err = yaml.Unmarshal(data, &describerScaledObjectManuals)
		if err != nil {
			a.logger.Error("failed to unmarshal template deployment file", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to unmarshal template deployment file")
		}

		describerScaledObjectManuals.Spec.ScaleTargetRef.Name = cnf.DescriberDeploymentName + "-manuals"

		triggerManuals := describerScaledObjectManuals.Spec.Triggers[0]
		triggerManuals.Metadata["stream"] = cnf.NatsStreamName
		triggerManuals.Metadata["natsServerMonitoringEndpoint"] = soNatsUrl
		triggerManuals.Metadata["consumer"] = cnf.NatsConsumerGroupManuals + "-service"
		describerScaledObjectManuals.Spec.Triggers[0] = triggerManuals

		newScaledObjectManuals := kedav1alpha1.ScaledObject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cnf.DescriberDeploymentName + "-manuals-scaled-object",
				Namespace: currentNamespace,
			},
			Spec: describerScaledObjectManuals.Spec,
		}

		err = a.kubeClient.Create(ctx, &newScaledObjectManuals)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
	}

	return nil
}
