package errors

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	ErrMetadataKeyNotFound              = echo.NewHTTPError(http.StatusNotFound, errors.New("metadata key not found").Error())
	ErrorMetadataValueTypeMismatch      = echo.NewHTTPError(http.StatusBadRequest, errors.New("metadata value type mismatch").Error())
	ErrorConfigMetadataTypeNotSupported = echo.NewHTTPError(http.StatusBadRequest, errors.New("config metadata type not supported").Error())
	ErrQueryParameterKeyNotFound        = echo.NewHTTPError(http.StatusNotFound, errors.New("query parameter key is not valid").Error())
)
