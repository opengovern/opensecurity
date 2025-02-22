package api

import (
	"github.com/labstack/echo/v4"
	"github.com/opengovern/opensecurity/services/rego/api/rego"
	"github.com/opengovern/opensecurity/services/rego/service"
	"go.uber.org/zap"
)

type API struct {
	logger  *zap.Logger
	Service *service.RegoEngine
}

func New(logger *zap.Logger, service *service.RegoEngine) *API {
	return &API{
		logger:  logger.Named("api"),
		Service: service,
	}
}

func (api *API) Register(e *echo.Echo) {
	evaluateApi := rego.New(api.logger, api.Service)
	evaluateApi.Register(e.Group("/api/v1/rego"))
}
