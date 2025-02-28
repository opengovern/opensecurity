package utils

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strconv"
	"time"
)

func TimeFromQueryParam(ctx echo.Context, paramName string, defaultValue time.Time) (time.Time, error) {
	timeStr := ctx.QueryParam(paramName)
	theTime := defaultValue
	if timeStr != "" {
		endTimeVal, err := strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			return theTime, errors.New("invalid time")
		}
		theTime = time.Unix(endTimeVal, 0)
	}
	return theTime, nil
}
