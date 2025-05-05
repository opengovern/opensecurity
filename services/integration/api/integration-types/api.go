package integration_types

import (
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/elasticsearch"
	coreClient "github.com/opengovern/opensecurity/services/core/client"
	demo_import "github.com/opengovern/opensecurity/services/integration/demo-import"
	utils2 "github.com/opengovern/opensecurity/services/integration/utils"
	hczap "github.com/zaffka/zap-to-hclog"
	"golang.org/x/net/context"
	"io"
	"io/fs"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/env"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	"github.com/opengovern/opensecurity/pkg/utils"
	"github.com/opengovern/opensecurity/services/integration/api/models"
	"github.com/opengovern/opensecurity/services/integration/db"
	integration_type "github.com/opengovern/opensecurity/services/integration/integration-type"
	models2 "github.com/opengovern/opensecurity/services/integration/models"
	"go.uber.org/zap"
	"sort"
)

const OneGB = 1024 * 1024 * 1024

type API struct {
	logger        *zap.Logger
	typeManager   *integration_type.IntegrationTypeManager
	database      db.Database
	elastic       opengovernance.Client
	coreClient    coreClient.CoreServiceClient
	elasticConfig config.ElasticSearch
}

func New(typeManager *integration_type.IntegrationTypeManager, database db.Database, logger *zap.Logger, elastic opengovernance.Client, coreClient coreClient.CoreServiceClient, elasticConfig config.ElasticSearch) *API {
	return &API{
		logger:        logger.Named("integration_types"),
		typeManager:   typeManager,
		database:      database,
		elastic:       elastic,
		coreClient:    coreClient,
		elasticConfig: elasticConfig,
	}
}

func (a *API) Register(e *echo.Group) {
	e.GET("", httpserver.AuthorizeHandler(a.List, api.ViewerRole))
	e.GET("/:integration_type/resource-type/table/:table_name", httpserver.AuthorizeHandler(a.GetResourceTypeFromTableName, api.ViewerRole))
	e.GET("/:integration_type/table", httpserver.AuthorizeHandler(a.ListTables, api.ViewerRole))
	e.POST("/:integration_type/resource-type/label", httpserver.AuthorizeHandler(a.GetResourceTypesByLabels, api.ViewerRole))
	e.GET("/:integration_type/configuration", httpserver.AuthorizeHandler(a.GetConfiguration, api.ViewerRole))

	plugin := e.Group("/plugin")
	plugin.GET("/load-v1", httpserver.AuthorizeHandler(a.AddPluginSpec, api.ViewerRole))
	plugin.GET("/:id/setup", httpserver.AuthorizeHandler(a.GetSetup, api.ViewerRole))
	plugin.GET("/:id/manifest", httpserver.AuthorizeHandler(a.GetManifest, api.ViewerRole))
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
	plugin.GET("/tables", httpserver.AuthorizeHandler(a.GetPluginsTables, api.ViewerRole))
	plugin.PUT("/:id/demo/load", httpserver.AuthorizeHandler(a.LoadPluginDemoData, api.EditorRole))
	plugin.PUT("/:id/demo/remove", httpserver.AuthorizeHandler(a.RemovePluginDemoData, api.EditorRole))
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

	plugin, err := a.database.GetPluginByID(integrationType)
	if err != nil {
		a.logger.Error("failed to get integration plugin", zap.String("integrationType", string(integrationType)), zap.Error(err))
		return echo.NewHTTPError(500, "failed to get integration plugin")
	}

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		conf, err := value.GetConfiguration()
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}

		return c.JSON(200, models.IntegrationTypeConfiguration{
			IntegrationConfiguration: conf,
			DiscoveryType:            string(plugin.DiscoveryType),
		})
	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}

// GetSetup godoc
// @Summary			Get integration setup
// @Description		Get integration setup
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			integration_type	path	string	true	"Integration type"
// @Success			200	{object} interfaces.IntegrationConfiguration
// @Router			/integration/api/v1/integration-types/{id}/setup [get]
func (a *API) GetSetup(c echo.Context) error {
	integrationType := c.Param("id")

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		conf, err := value.GetConfiguration()
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		var setup string
		// convert byte to string
		setup = string(conf.SetupMD)
		return c.String(200, setup)

	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}

// GetManifest godoc
// @Summary			Get integration manifest
// @Description		Get integration manifest
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			integration_type	path	string	true	"Integration type"
// @Success			200	{object} interfaces.IntegrationConfiguration
// @Router			/integration/api/v1/integration-types/{id}/manifest [get]
func (a *API) GetManifest(c echo.Context) error {
	integrationType := c.Param("id")

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		conf, err := value.GetConfiguration()
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
		var manifest models2.Manifest
		if err := yaml.Unmarshal(conf.Manifest, &manifest); err != nil {
			return echo.NewHTTPError(500, err.Error())
		}

		if plugin, err := a.database.GetPluginByID(integrationType); err == nil {
			if plugin.SupportedPlatformVersion != "" {
				manifest.SupportedPlatformVersion = plugin.SupportedPlatformVersion
			}
		}

		return c.JSON(200, manifest)

	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}

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

	validResourcetypes := make(map[string]bool)
	if req.IntegrationID != nil {
		resourceTypes, err := a.database.GetIntegrationResourcetypes(*req.IntegrationID)
		if err != nil {
			a.logger.Error("could not get integration resource types", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "could not get integration resource types")
		}
		if resourceTypes != nil {
			for _, rt := range resourceTypes.ResourceTypes {
				validResourcetypes[rt] = true
			}
		}
	}

	rtMap := a.typeManager.GetIntegrationTypeMap()
	if value, ok := rtMap[a.typeManager.ParseType(integrationType)]; ok {
		rts, err := value.GetResourceTypesByLabels(req.Labels)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}

		var finalRts []interfaces.ResourceTypeConfiguration

		if len(validResourcetypes) > 0 {
			for _, rt := range rts {
				if _, ok2 := validResourcetypes[rt.Name]; ok2 {
					finalRts = append(finalRts, rt)
				}
			}
		} else {
			finalRts = rts
		}
		res := models.GetResourceTypesByLabelsResponse{
			ResourceTypes: finalRts,
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

	installingPlugins, err := a.database.ListInstallingPlugins()
	if err != nil {
		a.logger.Error("failed to list installing plugins", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list installing plugins")
	}

	if len(installingPlugins) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "other plugin install is in progress")
	}

	plugin, err := a.database.GetPluginByID(pluginID)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	err = a.CheckEnoughMemory()
	if err != nil {
		a.logger.Error("checking enough memory failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	plugin.InstallState = models2.IntegrationTypeInstallStateInstalling

	err = a.database.UpdatePlugin(plugin)
	if err != nil {
		a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", pluginID))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
	}

	go func() {
		err = a.InstallOrUpdatePlugin(context.Background(), plugin)
		if err != nil {
			a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", pluginID))
		}
	}()

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

	installingPlugins, err := a.database.ListInstallingPlugins()
	if err != nil {
		a.logger.Error("failed to list installing plugins", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list installing plugins")
	}

	if len(installingPlugins) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "other plugin install is in progress")
	}

	if pluginURL == "" {
		return echo.NewHTTPError(http.StatusNotFound, "plugin url is empty")
	}
	url := pluginURL

	err = a.CheckEnoughMemory()
	if err != nil {
		a.logger.Error("checking enough memory failed", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	plugin, err := a.database.GetPluginByURL(url)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}

	if plugin == nil {
		baseDir := "/integration-types"

		// create tmp directory if not exists
		if _, err = os.Stat(baseDir); os.IsNotExist(err) {
			if err = os.Mkdir(baseDir, os.ModePerm); err != nil {
				a.logger.Error("failed to create tmp directory", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to create tmp directory")
			}
		}

		// remove existing files
		if err = os.RemoveAll(baseDir + "/integration_type"); err != nil {
			a.logger.Error("failed to remove existing files", zap.Error(err), zap.String("url", url), zap.String("path", baseDir+"/integration_type"))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to remove existing files")
		}

		downloader := getter.Client{
			Src:  url,
			Dst:  baseDir + "/integration_type",
			Mode: getter.ClientModeDir,
		}
		err = downloader.Get()
		if err != nil {
			a.logger.Error("failed to get integration binaries", zap.Error(err), zap.String("url", url))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get integration binaries")
		}

		// make scope to delete integrationPlugin and cloudqlPlugin after usage
		{
			// read integration-plugin file
			integrationPlugin, err := os.ReadFile(baseDir + "/integration_type/integration-plugin")
			if err != nil {
				a.logger.Error("failed to open integration-plugin file", zap.Error(err), zap.String("url", url))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to open integration-plugin file")
			}
			cloudqlPlugin, err := os.ReadFile(baseDir + "/integration_type/cloudql-plugin")
			if err != nil {
				a.logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("url", url))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to open cloudql-plugin file")
			}
			//// read manifest file
			manifestFile, err := os.ReadFile(baseDir + "/integration_type/manifest.yaml")
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

			plugin = &models2.IntegrationPlugin{
				PluginID:        m.IntegrationType.String(),
				IntegrationType: m.IntegrationType,
				DescriberURL:    m.DescriberURL,
				DescriberTag:    m.DescriberTag,
				DemoDataURL:     m.DemoDataURL,
				InstallState:    models2.IntegrationTypeInstallStateInstalling,
				URL:             url,
			}
			err = a.database.CreatePlugin(plugin)
			if err != nil {
				a.logger.Error("failed to create plugin", zap.Error(err), zap.String("id", m.IntegrationType.String()))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to create plugin")
			}

			pluginBinary := &models2.IntegrationPluginBinary{
				PluginID: m.IntegrationType.String(),

				IntegrationPlugin: integrationPlugin,
				CloudQlPlugin:     cloudqlPlugin,
			}
			err = a.database.CreatePluginBinary(pluginBinary)
			if err != nil {
				a.logger.Error("failed to create plugin binary", zap.Error(err), zap.String("id", m.IntegrationType.String()))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to create plugin binary")
			}

			go func() {
				err = a.LoadPlugin(c.Request().Context(), plugin, pluginBinary)
				if err != nil {
					a.logger.Error("failed to load plugin", zap.Error(err), zap.String("id", plugin.PluginID))

					plugin.InstallState = models2.IntegrationTypeInstallStateNotInstalled
					plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusFailed

					err = a.database.UpdatePlugin(plugin)
					if err != nil {
						a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.IntegrationType.String()))
						return
					}
					return
				}

				plugin.InstallState = models2.IntegrationTypeInstallStateInstalled
				plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusEnabled

				err = a.database.UpdatePlugin(plugin)
				if err != nil {
					a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.IntegrationType.String()))
					return
				}
			}()
		}
	} else {
		plugin.InstallState = models2.IntegrationTypeInstallStateInstalling

		err = a.database.UpdatePlugin(plugin)
		if err != nil {
			a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.PluginID))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to update plugin")
		}

		go func() {
			err = a.InstallOrUpdatePlugin(context.Background(), plugin)
			if err != nil {
				a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.PluginID))
			}
		}()
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

	err = a.database.UpdatePlugin(plugin)
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

	err = a.database.UpdatePlugin(plugin)
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

	err = a.database.UpdatePlugin(plugin)
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

		operationalStatusUpdatesStr, err := plugin.GetStringOperationalStatusUpdates()
		if err != nil {
			a.logger.Error("failed to get operational status updates", zap.Error(err))
		}

		var operationalStatusUpdates []models.OperationalStatusUpdate
		for _, operationalStatusUpdateStr := range operationalStatusUpdatesStr {
			var operationalStatusUpdate models2.OperationalStatusUpdate
			err = json.Unmarshal([]byte(operationalStatusUpdateStr), &operationalStatusUpdate)
			if err != nil {
				a.logger.Error("failed to unmarshal operational status updates", zap.Error(err))
			}
			operationalStatusUpdates = append(operationalStatusUpdates, models.OperationalStatusUpdate{
				Time:      operationalStatusUpdate.Time,
				OldStatus: string(operationalStatusUpdate.OldStatus),
				NewStatus: string(operationalStatusUpdate.NewStatus),
				Reason:    operationalStatusUpdate.Reason,
			})
		}

		items = append(items, models.IntegrationPlugin{
			ID:                       plugin.ID,
			PluginID:                 plugin.PluginID,
			IntegrationType:          plugin.IntegrationType.String(),
			InstallState:             string(plugin.InstallState),
			OperationalStatus:        string(plugin.OperationalStatus),
			OperationalStatusUpdates: operationalStatusUpdates,
			URL:                      plugin.URL,
			Tier:                     plugin.Tier,
			Description:              plugin.Description,
			Version:                  plugin.Version,
			Icon:                     plugin.Icon,
			Availability:             plugin.Availability,
			SourceCode:               plugin.SourceCode,
			PackageType:              plugin.PackageType,
			DescriberURL:             plugin.DescriberURL,
			DiscoveryType:            string(plugin.DiscoveryType),
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

	operationalStatusUpdatesStr, err := plugin.GetStringOperationalStatusUpdates()
	if err != nil {
		a.logger.Error("failed to get operational status updates", zap.Error(err))
	}

	var operationalStatusUpdates []models.OperationalStatusUpdate
	for _, operationalStatusUpdateStr := range operationalStatusUpdatesStr {
		var operationalStatusUpdate models2.OperationalStatusUpdate
		err = json.Unmarshal([]byte(operationalStatusUpdateStr), &operationalStatusUpdate)
		if err != nil {
			a.logger.Error("failed to unmarshal operational status updates", zap.Error(err))
		}
		operationalStatusUpdates = append(operationalStatusUpdates, models.OperationalStatusUpdate{
			Time:      operationalStatusUpdate.Time,
			OldStatus: string(operationalStatusUpdate.OldStatus),
			NewStatus: string(operationalStatusUpdate.NewStatus),
			Reason:    operationalStatusUpdate.Reason,
		})
	}

	return c.JSON(http.StatusOK, models.IntegrationPlugin{
		ID:                       plugin.ID,
		PluginID:                 plugin.PluginID,
		IntegrationType:          plugin.IntegrationType.String(),
		InstallState:             string(plugin.InstallState),
		OperationalStatus:        string(plugin.OperationalStatus),
		OperationalStatusUpdates: operationalStatusUpdates,
		URL:                      plugin.URL,
		Tier:                     plugin.Tier,
		Description:              plugin.Description,
		Version:                  plugin.Version,
		Icon:                     plugin.Icon,
		Availability:             plugin.Availability,
		SourceCode:               plugin.SourceCode,
		PackageType:              plugin.PackageType,
		DescriberURL:             plugin.DescriberURL,
		DiscoveryType:            string(plugin.DiscoveryType),
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

// GetPluginsTables godoc
//
// @Summary			Get plugins tables
// @Description		Get plugins tables
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Success			200 {object}  []models.PluginTables
// @Router			/integration/api/v1/integration-types/plugin/tables [get]
func (a *API) GetPluginsTables(c echo.Context) error {
	var plugins []models.PluginTables

	for it, plugin := range a.typeManager.GetIntegrationTypeMap() {
		tablesMap, err := plugin.ListAllTables()
		if err != nil {
			a.logger.Error("failed to list all tables", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to list all tables")
		}
		var tables []string
		for t := range tablesMap {
			tables = append(tables, t)
		}
		plugins = append(plugins, models.PluginTables{
			PluginID: it.String(),
			Tables:   tables,
		})
	}

	return c.JSON(http.StatusOK, plugins)
}

func (a *API) InstallOrUpdatePlugin(ctx context.Context, plugin *models2.IntegrationPlugin) (err error) {
	defer func() {
		if err != nil {
			plugin.InstallState = models2.IntegrationTypeInstallStateNotInstalled
			plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusFailed
		} else {
			plugin.InstallState = models2.IntegrationTypeInstallStateInstalled
			plugin.OperationalStatus = models2.IntegrationPluginOperationalStatusEnabled
		}
		err = a.database.UpdatePlugin(plugin)
		if err != nil {
			a.logger.Error("failed to update plugin", zap.Error(err), zap.String("id", plugin.PluginID))
			return
		}
	}()

	baseDir := "/integration-types"

	// create tmp directory if not exists
	if _, err = os.Stat(baseDir); os.IsNotExist(err) {
		if err = os.Mkdir(baseDir, os.ModePerm); err != nil {
			a.logger.Error("failed to create tmp directory", zap.Error(err))
			return err
		}
	}

	// download files from urls

	if plugin.URL == "" {
		return err
	}
	url := plugin.URL
	// remove existing files
	if err = os.RemoveAll(baseDir + "/integration_type"); err != nil {
		a.logger.Error("failed to remove existing files", zap.Error(err), zap.String("id", plugin.PluginID), zap.String("path", baseDir+"/integration_type"))
		return err
	}

	downloader := getter.Client{
		Src:  url,
		Dst:  baseDir + "/integration_type",
		Mode: getter.ClientModeDir,
	}
	err = downloader.Get()
	if err != nil {
		a.logger.Error("failed to get integration binaries", zap.Error(err), zap.String("id", plugin.PluginID))
		return err
	}

	// read integration-plugin file
	integrationPlugin, err := os.ReadFile(baseDir + "/integration_type/integration-plugin")
	if err != nil {
		a.logger.Error("failed to open integration-plugin file", zap.Error(err), zap.String("id", plugin.PluginID))
		return err
	}
	cloudqlPlugin, err := os.ReadFile(baseDir + "/integration_type/cloudql-plugin")
	if err != nil {
		a.logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("id", plugin.PluginID))
		return err
	}

	//// read manifest file
	manifestFile, err := os.ReadFile(baseDir + "/integration_type/manifest.yaml")
	if err != nil {
		a.logger.Error("failed to open manifest file", zap.Error(err))
		return err
	}
	a.logger.Info("manifestFile", zap.String("file", string(manifestFile)))

	var m models2.Manifest
	// decode yaml
	if err = yaml.Unmarshal(manifestFile, &m); err != nil {
		a.logger.Error("failed to decode manifest", zap.Error(err), zap.String("url", url))
		return err
	}

	// Opensearch templates
	a.logger.Info("checking for index-templates", zap.String("id", plugin.PluginID))
	var files []string
	if stats, err := os.Stat(filepath.Join(baseDir, "integration_type", "index-templates")); err == nil && stats.IsDir() {
		a.logger.Info("found index-templates directory", zap.String("id", plugin.PluginID))
		err = filepath.Walk(filepath.Join(baseDir, "integration_type", "index-templates"), func(path string, info fs.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".json") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			a.logger.Error("failed to get files", zap.Error(err))
			return err
		}

		// We need to create component templates first hence we are iterating over the files twice
		for _, fp := range files {
			if strings.Contains(fp, "_component_template") {
				a.logger.Info("creating component template", zap.String("filepath", fp))
				err = elasticsearch.CreateTemplate(ctx, a.elastic, a.logger, fp)
				if err != nil {
					a.logger.Error("failed to create component template", zap.Error(err), zap.String("filepath", fp), zap.String("id", plugin.PluginID))
				}
			}
		}

		for _, fp := range files {
			if !strings.Contains(fp, "_component_template") {
				a.logger.Info("creating template", zap.String("filepath", fp))
				err = elasticsearch.CreateTemplate(ctx, a.elastic, a.logger, fp)
				if err != nil {
					a.logger.Error("failed to create template", zap.Error(err), zap.String("filepath", fp), zap.String("id", plugin.PluginID))
				}
			}
		}
	}

	a.logger.Info("done reading files", zap.String("id", plugin.PluginID), zap.String("url", url), zap.String("integrationType", plugin.IntegrationType.String()), zap.Int("integrationPluginSize", len(integrationPlugin)), zap.Int("cloudqlPluginSize", len(cloudqlPlugin)))

	plugin.DescriberURL = m.DescriberURL
	plugin.DescriberTag = m.DescriberTag
	plugin.DemoDataURL = m.DemoDataURL

	pluginBinary := &models2.IntegrationPluginBinary{
		PluginID: m.IntegrationType.String(),

		IntegrationPlugin: integrationPlugin,
		CloudQlPlugin:     cloudqlPlugin,
	}
	err = a.database.CreatePluginBinary(pluginBinary)
	if err != nil {
		a.logger.Error("failed to create plugin binary", zap.Error(err), zap.String("id", m.IntegrationType.String()))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create plugin binary")
	}
	err = a.LoadPlugin(context.Background(), plugin, pluginBinary)
	if err != nil {
		a.logger.Error("failed to load plugin", zap.Error(err), zap.String("id", plugin.PluginID))
		return err
	}

	return nil
}

func (a *API) LoadPlugin(ctx context.Context, plugin *models2.IntegrationPlugin, pluginBinary *models2.IntegrationPluginBinary) error {
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

	if v, ok := a.typeManager.Clients[plugin.IntegrationType]; ok {
		v.Kill()
		delete(a.typeManager.Clients, plugin.IntegrationType)
	}

	pluginPath := filepath.Join(baseDir, plugin.IntegrationType.String()+".so")

	if _, err := os.Stat(pluginPath); err == nil {
		if err := os.Remove(pluginPath); err != nil {
			a.logger.Error("failed to delete existing plugin file", zap.Error(err), zap.String("pluginPath", pluginPath))
		}
	}

	err := os.WriteFile(pluginPath, pluginBinary.IntegrationPlugin, 0755)
	if err != nil {
		a.logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", pluginName))
		return err
	}
	hcLogger := hczap.Wrap(a.logger)

	client := plugin2.NewClient(&plugin2.ClientConfig{
		HandshakeConfig: interfaces.HandshakeConfig,
		Plugins:         map[string]plugin2.Plugin{pluginName: &interfaces.IntegrationTypePlugin{}},
		Cmd:             exec.Command(pluginPath),
		Logger:          hcLogger,
		Managed:         true,
	})
	a.typeManager.Clients[plugin.IntegrationType] = client

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

	itInterface, ok := raw.(interfaces.IntegrationType)
	if !ok {
		a.logger.Error("failed to cast plugin to integration type", zap.String("plugin", pluginName), zap.String("path", pluginPath))
		client.Kill()
		return err
	}
	a.typeManager.IntegrationTypes[plugin.IntegrationType] = itInterface

	a.typeManager.PingLocks[plugin.IntegrationType] = &sync.Mutex{}

	if plugin.DiscoveryType == models2.IntegrationPluginDiscoveryTypeClassic {
		err = a.typeManager.EnableIntegrationTypeHelper(ctx, plugin)
		if err != nil {
			a.logger.Error("failed to enable integration type describer", zap.Error(err))
			return err
		}
	}

	err = a.coreClient.ReloadPluginSteampipeConfig(&httpclient.Context{UserRole: api.AdminRole}, plugin.PluginID)
	if err != nil {
		a.logger.Error("failed to reload plugin config", zap.Error(err), zap.String("id", plugin.PluginID))
		return err
	}

	// Restart cloudql enabled services so that they can use the new plugin
	err = a.typeManager.RestartCloudQLEnabledServices(ctx)
	if err != nil {
		a.logger.Error("failed to restart cloudql enabled services", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to restart cloudql enabled services", err)
	}

	return nil
}

func (a *API) UnLoadPlugin(ctx context.Context, plugin models2.IntegrationPlugin) error {
	err := a.typeManager.DisableIntegrationTypeHelper(ctx, plugin.IntegrationType.String())
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

	err = a.coreClient.RemovePluginSteampipeConfig(&httpclient.Context{UserRole: api.AdminRole}, plugin.PluginID)
	if err != nil {
		a.logger.Error("failed to reload plugin config", zap.Error(err), zap.String("id", plugin.PluginID))
		return err
	}

	return nil
}

func (a *API) CheckEnoughMemory() error {
	if a.typeManager.MetricsClient == nil {
		a.logger.Warn("Metrics API is not available")
		return nil
	}
	a.logger.Info("start checking memory")
	namespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
	if !ok {
		a.logger.Error("current namespace lookup failed")
		return nil
	}

	labelSelector := "app=integration-service"

	// Get the pod associated with the deployment
	pods, err := a.typeManager.KubeClientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		a.logger.Error("failed to get pods")
		return nil
	}

	for _, pod := range pods.Items {
		var memoryLimit *resource.Quantity
		for _, container := range pod.Spec.Containers {
			limit := container.Resources.Limits.Memory()
			memoryLimit = limit
		}
		if memoryLimit == nil {
			a.logger.Error("no memory limit found")
			return nil
		}
		memoryLimitBytes := memoryLimit.Value()

		podMetrics, err := a.typeManager.MetricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		if err != nil {
			a.logger.Error("failed to get metrics")
			return nil
		}
		var memoryUsage *resource.Quantity
		for _, container := range podMetrics.Containers {
			usage := container.Usage.Memory()
			memoryUsage = usage
		}
		if memoryUsage == nil {
			a.logger.Error("no memory usage found")
			return nil
		}
		memoryUsageBytes := memoryUsage.Value()

		availableMemory := big.NewInt(memoryLimitBytes).Sub(big.NewInt(memoryLimitBytes), big.NewInt(memoryUsageBytes)).Int64()

		a.logger.Info("memory usage and available", zap.Int64("usage", memoryLimitBytes), zap.Int64("limit", memoryLimitBytes),
			zap.Int64("available", availableMemory))
		if availableMemory >= OneGB {
			return nil
		} else {
			return fmt.Errorf("not enough available memory: %v", availableMemory)
		}
	}

	return nil
}

// LoadPluginDemoData godoc
//
// @Summary			Load demo data for plugin
// @Description		Load demo data for plugin by the given url
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id}/demo/load [put]
func (a *API) LoadPluginDemoData(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	if plugin.DemoDataLoaded {
		return echo.NewHTTPError(http.StatusConflict, "plugin demo data already loaded")
	}

	if plugin.DemoDataURL == "" {
		return echo.NewHTTPError(http.StatusNotFound, "plugin does not contain demo data")
	}

	demoImportConfig := demo_import.Config{
		DemoDataURL:       plugin.DemoDataURL,
		OpenSSLPassword:   env.GetString("OPENSSL_PASSWORD", ""),
		ElasticsearchUser: a.elasticConfig.Username,
		ElasticsearchPass: a.elasticConfig.Password,
		ElasticsearchAddr: a.elasticConfig.Address,
	}
	err = demo_import.LoadDemoData(demoImportConfig, a.logger, a.database, plugin)
	if err != nil {
		a.logger.Error("failed to load demo data", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load demo data")
	}

	err = a.database.PluginDemoDataLoad(plugin.PluginID)
	if err != nil {
		a.logger.Error("failed to set load plugin demo data flag", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to set load plugin demo data flag")
	}

	return c.NoContent(http.StatusOK)
}

// RemovePluginDemoData godoc
//
// @Summary			Remove demo data for plugin
// @Description		Remove demo data for plugin by the given url
// @Security		BearerToken
// @Tags			integration_types
// @Produce			json
// @Param			id	path	string	true	"plugin id"
// @Success			200
// @Router			/integration/api/v1/plugin/{id}/demo/remove [put]
func (a *API) RemovePluginDemoData(c echo.Context) error {
	id := c.Param("id")

	plugin, err := a.database.GetPluginByID(id)
	if err != nil {
		a.logger.Error("failed to get plugin", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get plugin")
	}
	if plugin == nil {
		return echo.NewHTTPError(http.StatusNotFound, "plugin not found")
	}

	if !plugin.DemoDataLoaded {
		return echo.NewHTTPError(http.StatusConflict, "plugin demo data not loaded")
	}

	if err = a.database.DeletePluginSampleIntegrations(plugin.IntegrationType); err != nil {
		a.logger.Error("failed to delete plugin sample integration", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete plugin sample integration")
	}

	err = a.database.PluginDemoDataUnLoad(plugin.PluginID)
	if err != nil {
		a.logger.Error("failed to unset load plugin demo data flag", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to unset load plugin demo data flag")
	}

	return c.NoContent(http.StatusOK)
}

// AddPluginSpec godoc
//
//	@Summary	Load Task
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		cursor		query	int	false	"cursor"
//	@Param		per_page	query	int	false	"per page"
//	@Produce	json
//	@Success	200	{object}	api.ListTaskRunsResponse
//	@Router		/tasks/api/v1/plugin/load-v1 [post]
func (a *API) AddPluginSpec(ctx echo.Context) error {
	bodyBytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		a.logger.Error("failed to read request body", zap.Error(err))
		return ctx.JSON(http.StatusBadRequest, "failed to read body")
	}

	err = utils2.ValidateAndLoadPlugin(a.typeManager.IntegrationTypeDb, a.logger, bodyBytes)
	if err != nil {
		a.logger.Error("failed to add plugin spec", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, "failed to add plugin spec")
	}

	return ctx.NoContent(http.StatusOK)
}
