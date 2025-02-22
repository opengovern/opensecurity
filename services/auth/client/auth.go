package client

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opensecurity/services/auth/api"
	"net/http"
)

type AuthServiceClient interface {
	ListUsers(ctx *httpclient.Context) ([]api.GetUsersResponse, error)
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
