package rego

import (
	"github.com/labstack/echo/v4"
	"github.com/opengovern/opencomply/services/rego/api/models"
	"github.com/opengovern/opencomply/services/rego/service"
	"go.uber.org/zap"
)

type API struct {
	logger  *zap.Logger
	Service *service.RegoEngine
}

func New(logger *zap.Logger, service *service.RegoEngine) *API {
	return &API{
		logger:  logger.Named("evaluate"),
		Service: service,
	}
}

func (r API) Register(g *echo.Group) {
	g.POST("/evaluate", r.EvaluateEndpoint)
}

// EvaluateEndpoint godoc
//	@Summary		Evaluate a rego policy
//	@Description	Evaluate a rego policy
//	@Tags			rego
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegoEvaluateRequest	true	"Rego Evaluate Request"
//	@Success		200		{object}	RegoEvaluateResponse
//	@Router			/evaluate [post]
func (r *API) EvaluateEndpoint(c echo.Context) error {
	req := new(models.RegoEvaluateRequest)
	if err := c.Bind(req); err != nil {
		r.logger.Error("Unable to bind request", zap.Error(err))
		r.logger.Sync()
		return err
	}

	results, err := r.Service.Evaluate(c.Request().Context(), req.Policies, req.Query)
	if err != nil {
		r.logger.Error("Unable to evaluate rego", zap.Error(err))
		r.logger.Sync()
		return err
	}

	return c.JSON(200, models.RegoEvaluateResponse{Results: results})
}
