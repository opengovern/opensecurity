package client

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opensecurity/services/auth/api"
)

type AuthServiceClient interface {
	ListUsers(ctx *httpclient.Context) ([]api.GetUsersResponse, error)
	ListApiKeys(ctx *httpclient.Context) ([]api.APIKeyResponse, error)
	GetConnectors(ctx *httpclient.Context) ([]api.GetConnectorsResponse, error)
}

type authClient struct {
	baseURL string
}

func NewAuthClient(baseURL string) AuthServiceClient {
	return &authClient{baseURL: baseURL}
}

func (s *authClient) ListUsers(ctx *httpclient.Context) ([]api.GetUsersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users", s.baseURL)

	var users []api.GetUsersResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &users); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return users, nil
}

func (s *authClient) ListApiKeys(ctx *httpclient.Context) ([]api.APIKeyResponse, error) {
	url := fmt.Sprintf("%s/api/v1/keys", s.baseURL)

	var keys []api.APIKeyResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &keys); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return keys, nil
}

func (s *authClient) GetConnectors(ctx *httpclient.Context) ([]api.GetConnectorsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/connectors", s.baseURL)

	var connectors []api.GetConnectorsResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &connectors); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return connectors, nil
}
