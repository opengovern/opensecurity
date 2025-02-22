package query_validator

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/steampipe"
	coreApi "github.com/opengovern/opensecurity/services/core/api"
	"go.uber.org/zap"
)

func (w *Worker) RunSQLNamedQuery(ctx context.Context, query string) (*steampipe.Result, error) {
	var err error

	direction := coreApi.DirectionType("")

	for i := 0; i < 10; i++ {
		err = w.steampipeConn.Conn().Ping(ctx)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		w.logger.Error("failed to ping steampipe", zap.Error(err))
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	w.logger.Info("executing named query", zap.String("query", query))
	res, err := w.steampipeConn.Query(ctx, query, nil, nil, "", steampipe.DirectionType(direction))
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return res, nil
}
