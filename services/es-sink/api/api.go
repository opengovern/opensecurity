package api

import (
	"github.com/labstack/echo/v4"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opensecurity/services/es-sink/api/ingest"
	"github.com/opengovern/opensecurity/services/es-sink/service"
	"go.uber.org/zap"
)

type API struct {
	logger    *zap.Logger
	ingestApi *ingest.API
}

func New(logger *zap.Logger, esSinkService *service.EsSinkService) *API {
	logger = logger.Named("api-es-sink")
	ingestApi := ingest.New(logger, esSinkService)
	return &API{
		logger:    logger,
		ingestApi: ingestApi,
	}
}

func (api *API) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	v1.POST("/ingest", httpserver.AuthorizeHandler(api.ingestApi.Ingest, authApi.AdminRole))
}
