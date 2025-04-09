package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/integration"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/opensecurity/services/core/api"
	"github.com/opengovern/opensecurity/services/core/db/models"
)

type CoreServiceClient interface {
	RunQuery(ctx *httpclient.Context, req api.RunQueryRequest) (*api.RunQueryResponse, error)
	GetQuery(ctx *httpclient.Context, id string) (*api.NamedQueryItemV2, error)
	ListQueriesV2(ctx *httpclient.Context, req *api.ListQueryV2Request) (*api.ListQueriesV2Response, error)
	CountResources(ctx *httpclient.Context) (int64, error)
	ListIntegrationsData(ctx *httpclient.Context, integrationIds []string, resourceCollections []string, startTime, endTime *time.Time, metricIDs []string, needCost, needResourceCount bool) (map[string]api.ConnectionData, error)
	ListResourceTypesMetadata(ctx *httpclient.Context, integrationTypes []integration.Type, services []string, resourceTypes []string, summarized bool, tags map[string]string, pageSize, pageNumber int) (*api.ListResourceTypeMetadataResponse, error)
	ListResourceCollections(ctx *httpclient.Context) ([]api.ResourceCollection, error)
	GetResourceCollectionMetadata(ctx *httpclient.Context, id string) (*api.ResourceCollection, error)
	ListResourceCollectionsMetadata(ctx *httpclient.Context, ids []string) ([]api.ResourceCollection, error)
	GetTablesResourceCategories(ctx *httpclient.Context, tables []string) ([]api.CategoriesTables, error)
	GetResourceCategories(ctx *httpclient.Context, tables []string, categories []string) (*api.GetResourceCategoriesResponse, error)
	RunQueryByID(ctx *httpclient.Context, req api.RunQueryByIDRequest) (*api.RunQueryResponse, error)
	GetConfigMetadata(ctx *httpclient.Context, key models.MetadataKey) (models.IConfigMetadata, error)
	SetConfigMetadata(ctx *httpclient.Context, key models.MetadataKey, value any) error
	ListQueryParameters(ctx *httpclient.Context, request api.ListQueryParametersRequest) (*api.ListQueryParametersResponse, error)
	SetQueryParameter(ctx *httpclient.Context, request api.SetQueryParameterRequest) error
	VaultConfigured(ctx *httpclient.Context) (*string, error)
	GetViewsCheckpoint(ctx *httpclient.Context) (*api.GetViewsCheckpointResponse, error)
	ReloadViews(ctx *httpclient.Context) error
	GetAbout(ctx *httpclient.Context) (*api.About, error)
	ReloadPluginSteampipeConfig(ctx *httpclient.Context, pluginId string) error
	RemovePluginSteampipeConfig(ctx *httpclient.Context, pluginId string) error
	ListCacheEnabledQueries(ctx *httpclient.Context) ([]models.NamedQueryWithCacheStatus, error)
}

type coreClient struct {
	baseURL string
}

var ErrConfigNotFound = errors.New("config not found")

func NewCoreServiceClient(baseURL string) CoreServiceClient {
	return &coreClient{baseURL: baseURL}
}

// inventory

func (s *coreClient) RunQuery(ctx *httpclient.Context, req api.RunQueryRequest) (*api.RunQueryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/query/run", s.baseURL)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp api.RunQueryResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), reqBytes, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &resp, nil
}

func (s *coreClient) CountResources(ctx *httpclient.Context) (int64, error) {
	url := fmt.Sprintf("%s/api/v2/resources/count", s.baseURL)

	var count int64
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &count); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return 0, echo.NewHTTPError(statusCode, err.Error())
		}
		return 0, err
	}
	return count, nil
}

func (s *coreClient) GetQuery(ctx *httpclient.Context, id string) (*api.NamedQueryItemV2, error) {
	url := fmt.Sprintf("%s/api/v3/queries/%s", s.baseURL, id)

	var namedQuery api.NamedQueryItemV2
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &namedQuery); err != nil {
		if statusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	return &namedQuery, nil
}

func (s *coreClient) ListQueriesV2(ctx *httpclient.Context, req *api.ListQueryV2Request) (*api.ListQueriesV2Response, error) {
	url := fmt.Sprintf("%s/api/v3/queries", s.baseURL)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var namedQuery api.ListQueriesV2Response
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), reqBytes, &namedQuery); err != nil {
		if statusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	return &namedQuery, nil
}

func (s *coreClient) GetTablesResourceCategories(ctx *httpclient.Context, tables []string) ([]api.CategoriesTables, error) {
	url := fmt.Sprintf("%s/api/v3/tables/categories", s.baseURL)

	firstParamAttached := false
	if len(tables) > 0 {
		for _, t := range tables {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("tables=%s", t)
		}
	}

	var resp []api.CategoriesTables
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return resp, nil
}

func (s *coreClient) GetResourceCategories(ctx *httpclient.Context, tables []string, categories []string) (*api.GetResourceCategoriesResponse, error) {
	url := fmt.Sprintf("%s/api/v3/resources/categories", s.baseURL)

	firstParamAttached := false
	if len(tables) > 0 {
		for _, t := range tables {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("tables=%s", t)
		}
	}
	if len(categories) > 0 {
		for _, t := range categories {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("categories=%s", t)
		}
	}

	var resp api.GetResourceCategoriesResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &resp, nil
}

func (s *coreClient) ListIntegrationsData(
	ctx *httpclient.Context,
	integrationIds, resourceCollections []string,
	startTime, endTime *time.Time, metricIDs []string,
	needCost, needResourceCount bool,
) (map[string]api.ConnectionData, error) {
	params := url.Values{}

	url := fmt.Sprintf("%s/api/v2/integrations/data", s.baseURL)

	if len(integrationIds) > 0 {
		for _, integrationId := range integrationIds {
			params.Add("integrationId", integrationId)
		}
	}
	if len(resourceCollections) > 0 {
		for _, resourceCollection := range resourceCollections {
			params.Add("resourceCollection", resourceCollection)
		}
	}
	if len(metricIDs) > 0 {
		for _, metricID := range metricIDs {
			params.Add("metricId", metricID)
		}

	}
	if startTime != nil {
		params.Set("startTime", strconv.FormatInt(startTime.Unix(), 10))
	}
	if endTime != nil {
		params.Set("endTime", strconv.FormatInt(endTime.Unix(), 10))
	}
	if !needCost {
		params.Set("needCost", "false")
	}
	if !needResourceCount {
		params.Set("needResourceCount", "false")
	}
	if len(params) > 0 {
		url += "?" + params.Encode()
	}
	var response map[string]api.ConnectionData
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (s *coreClient) ListResourceTypesMetadata(ctx *httpclient.Context, integrationTypes []integration.Type, services []string, resourceTypes []string, summarized bool, tags map[string]string, pageSize, pageNumber int) (*api.ListResourceTypeMetadataResponse, error) {
	url := fmt.Sprintf("%s/api/v2/metadata/resourcetype", s.baseURL)
	firstParamAttached := false
	if len(integrationTypes) > 0 {
		for _, integrationType := range integrationTypes {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("integrationType=%s", integrationType)
		}
	}
	if len(services) > 0 {
		for _, service := range services {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("service=%s", service)
		}
	}
	if len(resourceTypes) > 0 {
		for _, resourceType := range resourceTypes {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("resourceType=%s", resourceType)
		}
	}
	if summarized {
		if !firstParamAttached {
			url += "?"
			firstParamAttached = true
		} else {
			url += "&"
		}
		url += "summarized=true"
	}
	if len(tags) > 0 {
		for key, value := range tags {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("tags=%s=%s", key, value)
		}
	}
	if pageSize > 0 {
		if !firstParamAttached {
			url += "?"
			firstParamAttached = true
		} else {
			url += "&"
		}
		url += fmt.Sprintf("pageSize=%d", pageSize)
	}
	if pageNumber > 0 {
		if !firstParamAttached {
			url += "?"
			firstParamAttached = true
		} else {
			url += "&"
		}
		url += fmt.Sprintf("pageNumber=%d", pageNumber)
	}

	var response api.ListResourceTypeMetadataResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *coreClient) GetResourceCollectionMetadata(ctx *httpclient.Context, id string) (*api.ResourceCollection, error) {
	url := fmt.Sprintf("%s/api/v2/metadata/resource-collection/%s", s.baseURL, id)

	var response api.ResourceCollection
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *coreClient) ListResourceCollectionsMetadata(ctx *httpclient.Context, ids []string) ([]api.ResourceCollection, error) {
	url := fmt.Sprintf("%s/api/v2/metadata/resource-collection", s.baseURL)

	firstParamAttached := false
	for _, id := range ids {
		if !firstParamAttached {
			url += "?"
			firstParamAttached = true
		} else {
			url += "&"
		}
		url += fmt.Sprintf("id=%s", id)
	}

	var response []api.ResourceCollection
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (s *coreClient) ListResourceCollections(ctx *httpclient.Context) ([]api.ResourceCollection, error) {
	url := fmt.Sprintf("%s/api/v2/metadata/resource-collection", s.baseURL)

	var response []api.ResourceCollection
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (s *coreClient) RunQueryByID(ctx *httpclient.Context, req api.RunQueryByIDRequest) (*api.RunQueryResponse, error) {
	url := fmt.Sprintf("%s/api/v3/query/run", s.baseURL)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var resp api.RunQueryResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), reqBytes, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &resp, nil
}

// metadata

func (s *coreClient) GetConfigMetadata(ctx *httpclient.Context, key models.MetadataKey) (models.IConfigMetadata, error) {
	url := fmt.Sprintf("%s/api/v1/metadata/%s", s.baseURL, string(key))
	var cnf models.ConfigMetadata
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &cnf); err != nil {
		if statusCode == 404 {
			return nil, ErrConfigNotFound
		}
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}

	switch cnf.Type {
	case models.ConfigMetadataTypeString:
		return &models.StringConfigMetadata{
			ConfigMetadata: cnf,
		}, nil
	case models.ConfigMetadataTypeInt:
		intValue, err := strconv.ParseInt(cnf.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert string to int: %w", err)
		}
		return &models.IntConfigMetadata{
			ConfigMetadata: cnf,
			Value:          int(intValue),
		}, nil
	case models.ConfigMetadataTypeBool:
		boolValue, err := strconv.ParseBool(cnf.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert string to bool: %w", err)
		}
		return &models.BoolConfigMetadata{
			ConfigMetadata: cnf,
			Value:          boolValue,
		}, nil
	case models.ConfigMetadataTypeJSON:
		return &models.JSONConfigMetadata{
			ConfigMetadata: cnf,
			Value:          cnf.Value,
		}, nil
	}

	return nil, fmt.Errorf("unknown config metadata type: %s", cnf.Type)
}

func (s *coreClient) SetConfigMetadata(ctx *httpclient.Context, key models.MetadataKey, value any) error {
	url := fmt.Sprintf("%s/api/v1/metadata", s.baseURL)

	req := api.SetConfigMetadataRequest{
		Key:   string(key),
		Value: value,
	}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	var cnf models.ConfigMetadata
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), jsonReq, &cnf); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}

	return nil
}

func (s *coreClient) ListQueryParameters(ctx *httpclient.Context, request api.ListQueryParametersRequest) (*api.ListQueryParametersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/query_parameter", s.baseURL)
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	var resp api.ListQueryParametersResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), jsonReq, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return &resp, echo.NewHTTPError(statusCode, err.Error())
		}
		return &resp, err
	}
	return &resp, nil
}

func (s *coreClient) SetQueryParameter(ctx *httpclient.Context, request api.SetQueryParameterRequest) error {
	url := fmt.Sprintf("%s/api/v1/query_parameter/set", s.baseURL)
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return err
	}

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), jsonReq, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}

	return nil
}

func (s *coreClient) VaultConfigured(ctx *httpclient.Context) (*string, error) {
	url := fmt.Sprintf("%s/api/v3/vault/configured", s.baseURL)
	var status string
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &status); err != nil {
		if statusCode == 404 {
			return nil, ErrConfigNotFound
		}
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}

	return &status, nil
}

func (s *coreClient) ReloadViews(ctx *httpclient.Context) error {
	url := fmt.Sprintf("%s/api/v3/views/reload", s.baseURL)

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}
	return nil
}

func (s *coreClient) GetViewsCheckpoint(ctx *httpclient.Context) (*api.GetViewsCheckpointResponse, error) {
	url := fmt.Sprintf("%s/api/v3/views/checkpoint", s.baseURL)
	var resp api.GetViewsCheckpointResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &resp, nil
}

func (s *coreClient) GetAbout(ctx *httpclient.Context) (*api.About, error) {
	url := fmt.Sprintf("%s/api/v3/about", s.baseURL)

	var about api.About
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &about); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}

	return &about, nil
}

func (s *coreClient) ReloadPluginSteampipeConfig(ctx *httpclient.Context, pluginId string) error {
	url := fmt.Sprintf("%s/api/v3/plugins/%s/reload", s.baseURL, pluginId)

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}

	return nil
}

func (s *coreClient) RemovePluginSteampipeConfig(ctx *httpclient.Context, pluginId string) error {
	url := fmt.Sprintf("%s/api/v3/plugins/%s/remove", s.baseURL, pluginId)

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}

	return nil
}

func (s *coreClient) ListCacheEnabledQueries(ctx *httpclient.Context) ([]models.NamedQueryWithCacheStatus, error) {
	url := fmt.Sprintf("%s/api/v3/queries/cache-enabled", s.baseURL)

	var resp []models.NamedQueryWithCacheStatus
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}

	return resp, nil
}
