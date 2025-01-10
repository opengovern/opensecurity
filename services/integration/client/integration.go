package client

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opencomply/services/integration/api/models"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"net/http"
)

type IntegrationServiceClient interface {
	ListIntegrationTypes(ctx *httpclient.Context) ([]string, error)
	GetResourceTypeFromTableName(ctx *httpclient.Context, integrationType string, tableName string) (string, error)
	GetResourceTypesByLabels(ctx *httpclient.Context, integrationType string, labels map[string]string) (map[string]*models.ResourceTypeConfiguration, error)
	GetIntegrationConfiguration(ctx *httpclient.Context, integrationType string) (interfaces.IntegrationConfiguration, error)
	GetIntegration(ctx *httpclient.Context, integrationID string) (*models.Integration, error)
	ListIntegrations(ctx *httpclient.Context, integrationTypes []string) (*models.ListIntegrationsResponse, error)
	ListIntegrationsByFilters(ctx *httpclient.Context, req models.ListIntegrationsRequest) (*models.ListIntegrationsResponse, error)
	IntegrationHealthcheck(ctx *httpclient.Context, integrationID string) (*models.Integration, error)
	GetCredential(ctx *httpclient.Context, credentialID string) (*models.Credential, error)
	ListCredentials(ctx *httpclient.Context) (*models.ListCredentialsResponse, error)
	GetIntegrationGroup(ctx *httpclient.Context, integrationGroupName string) (*models.IntegrationGroup, error)
	ListIntegrationGroups(ctx *httpclient.Context) ([]models.IntegrationGroup, error)
	PurgeSampleData(ctx *httpclient.Context) ([]string, error)
}

type integrationClient struct {
	baseURL string
}

func NewIntegrationServiceClient(baseURL string) IntegrationServiceClient {
	return &integrationClient{baseURL: baseURL}
}

func (c *integrationClient) ListIntegrationTypes(ctx *httpclient.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/integration-types", c.baseURL)
	var response []string

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (c *integrationClient) GetResourceTypeFromTableName(ctx *httpclient.Context, integrationType string, tableName string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/integration-types/%s/resource-type/table/%s", c.baseURL, integrationType, tableName)
	var response models.GetResourceTypeFromTableNameResponse

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return "", echo.NewHTTPError(statusCode, err.Error())
		}
		return "", err
	}
	return response.ResourceType, nil
}

func (c *integrationClient) GetResourceTypesByLabels(ctx *httpclient.Context, integrationType string, labels map[string]string) (map[string]*models.ResourceTypeConfiguration, error) {
	url := fmt.Sprintf("%s/api/v1/integration-types/%s/resource-type/label", c.baseURL, integrationType)

	req := models.GetResourceTypesByLabelsRequest{
		Labels: labels,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var response models.GetResourceTypesByLabelsResponse

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), payload, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response.ResourceTypes, nil
}

func (c *integrationClient) GetIntegrationConfiguration(ctx *httpclient.Context, integrationType string) (interfaces.IntegrationConfiguration, error) {
	url := fmt.Sprintf("%s/api/v1/integration-types/%s/configuration", c.baseURL, integrationType)
	var response interfaces.IntegrationConfiguration

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return response, echo.NewHTTPError(statusCode, err.Error())
		}
		return response, err
	}
	return response, nil
}

func (c *integrationClient) GetIntegration(ctx *httpclient.Context, integrationID string) (*models.Integration, error) {
	url := fmt.Sprintf("%s/api/v1/integrations/%s", c.baseURL, integrationID)
	var response *models.Integration

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (c *integrationClient) ListIntegrations(ctx *httpclient.Context, integrationTypes []string) (*models.ListIntegrationsResponse, error) {
	ctx.UserRole = authApi.AdminRole
	url := fmt.Sprintf("%s/api/v1/integrations", c.baseURL)
	for i, v := range integrationTypes {
		if i == 0 {
			url += "?"
		} else {
			url += "&"
		}
		url += "integration_type=" + v
	}

	var response models.ListIntegrationsResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (c *integrationClient) ListIntegrationsByFilters(ctx *httpclient.Context, req models.ListIntegrationsRequest) (*models.ListIntegrationsResponse, error) {
	ctx.UserRole = authApi.AdminRole
	url := fmt.Sprintf("%s/api/v1/integrations/list", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var response models.ListIntegrationsResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), payload, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (c *integrationClient) GetCredential(ctx *httpclient.Context, credentialID string) (*models.Credential, error) {
	url := fmt.Sprintf("%s/api/v1/credentials/%s", c.baseURL, credentialID)
	var response *models.Credential

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (c *integrationClient) ListCredentials(ctx *httpclient.Context) (*models.ListCredentialsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/credentials", c.baseURL)
	var response models.ListCredentialsResponse

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (c *integrationClient) IntegrationHealthcheck(ctx *httpclient.Context, integrationID string) (*models.Integration, error) {
	url := fmt.Sprintf("%s/api/v1/integrations/%s/healthcheck", c.baseURL, integrationID)
	var response *models.Integration

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (c *integrationClient) GetIntegrationGroup(ctx *httpclient.Context, integrationGroupName string) (*models.IntegrationGroup, error) {
	url := fmt.Sprintf("%s/api/v1/integrations/integration-groups/%s", c.baseURL, integrationGroupName)

	var integrationGroup models.IntegrationGroup
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &integrationGroup); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}

	return &integrationGroup, nil
}

func (c *integrationClient) ListIntegrationGroups(ctx *httpclient.Context) ([]models.IntegrationGroup, error) {
	url := fmt.Sprintf("%s/api/v1/integrations/integration-groups", c.baseURL)

	var integrationGroup []models.IntegrationGroup
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &integrationGroup); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}

	return integrationGroup, nil
}

func (c *integrationClient) PurgeSampleData(ctx *httpclient.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/integrations/sample/purge", c.baseURL)

	var resp struct {
		Integrations []string `json:"integrations"`
	}

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}

	return resp.Integrations, nil
}
