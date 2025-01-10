package integration_types

import (
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opencomply/pkg/utils"
	"github.com/opengovern/opencomply/services/integration/api/models"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"go.uber.org/zap"
)

type API struct {
	logger      *zap.Logger
	typeManager *integration_type.IntegrationTypeManager
}

func New(logger *zap.Logger, typeManager *integration_type.IntegrationTypeManager) *API {
	return &API{
		logger:      logger.Named("integration_types"),
		typeManager: typeManager,
	}
}

func (a *API) Register(e *echo.Group) {
	e.GET("", httpserver.AuthorizeHandler(a.List, api.ViewerRole))
	e.GET("/:integration_type/resource-type/table/:table_name", httpserver.AuthorizeHandler(a.GetResourceTypeFromTableName, api.ViewerRole))
	e.POST("/:integration_type/resource-type/label", httpserver.AuthorizeHandler(a.GetResourceTypesByLabels, api.ViewerRole))
	e.GET("/:integration_type/configuration", httpserver.AuthorizeHandler(a.GetConfiguration, api.ViewerRole))
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
		resourceType := value.GetResourceTypeFromTableName(tableName)
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
		return c.JSON(200, value.GetConfiguration())
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
			res.ResourceTypes[k] = utils.GetPointer(v.ToAPI())
		}
		return c.JSON(200, res)
	} else {
		return echo.NewHTTPError(404, "integration type not found")
	}
}
