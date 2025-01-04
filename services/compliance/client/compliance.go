package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/httpclient"

	"github.com/opengovern/og-util/pkg/source"
	compliance "github.com/opengovern/opencomply/services/compliance/api"
)

type ComplianceServiceClient interface {
	ListAssignmentsByBenchmark(ctx *httpclient.Context, benchmarkID string) (*compliance.BenchmarkAssignedEntities, error)
	GetBenchmark(ctx *httpclient.Context, benchmarkID string) (*compliance.Benchmark, error)
	GetBenchmarkSummary(ctx *httpclient.Context, benchmarkID string, connectionId []string, timeAt *time.Time) (*compliance.BenchmarkEvaluationSummary, error)
	GetBenchmarkControls(ctx *httpclient.Context, benchmarkID string, connectionId []string, timeAt *time.Time) (*compliance.BenchmarkControlSummary, error)
	GetControl(ctx *httpclient.Context, controlID string) (*compliance.Control, error)
	GetQuery(ctx *httpclient.Context, queryID string) (*compliance.Query, error)
	GetComplianceResults(ctx *httpclient.Context, req compliance.GetComplianceResultsRequest) (compliance.GetComplianceResultsResponse, error)
	ListBenchmarks(ctx *httpclient.Context, tags map[string][]string) ([]compliance.Benchmark, error)
	ListAllBenchmarks(ctx *httpclient.Context, isBare bool) ([]compliance.Benchmark, error)
	GetAccountsComplianceResultsSummary(ctx *httpclient.Context, benchmarkId string, connectionId []string, connector []source.Type) (compliance.GetAccountsComplianceResultsSummaryResponse, error)
	CreateBenchmarkAssignment(ctx *httpclient.Context, benchmarkID, connectionId string) ([]compliance.BenchmarkAssignment, error)
	ListQueries(ctx *httpclient.Context) ([]compliance.Query, error)
	ListControl(ctx *httpclient.Context, controlIDs []string, tags map[string][]string) ([]compliance.Control, error)
	GetControlDetails(ctx *httpclient.Context, controlID string) (*compliance.GetControlDetailsResponse, error)
	SyncQueries(ctx *httpclient.Context) error
}

type complianceClient struct {
	baseURL string
}

func NewComplianceClient(baseURL string) ComplianceServiceClient {
	return &complianceClient{baseURL: baseURL}
}

func (s *complianceClient) SyncQueries(ctx *httpclient.Context) error {
	url := fmt.Sprintf("%s/api/v1/queries/sync", s.baseURL)

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}
	return nil
}

func (s *complianceClient) ListAssignmentsByBenchmark(ctx *httpclient.Context, benchmarkID string) (*compliance.BenchmarkAssignedEntities, error) {
	url := fmt.Sprintf("%s/api/v1/assignments/benchmark/%s", s.baseURL, benchmarkID)

	var response compliance.BenchmarkAssignedEntities
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *complianceClient) GetBenchmark(ctx *httpclient.Context, benchmarkID string) (*compliance.Benchmark, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/%s", s.baseURL, benchmarkID)

	var response compliance.Benchmark
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *complianceClient) GetControlDetails(ctx *httpclient.Context, controlID string) (*compliance.GetControlDetailsResponse, error) {
	url := fmt.Sprintf("%s/api/v3/control/%s", s.baseURL, controlID)

	var response compliance.GetControlDetailsResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if statusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	return &response, nil
}

func (s *complianceClient) GetBenchmarkSummary(ctx *httpclient.Context, benchmarkID string, connectionId []string, timeAt *time.Time) (*compliance.BenchmarkEvaluationSummary, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/%s/summary", s.baseURL, benchmarkID)

	firstParamAttached := false
	if len(connectionId) > 0 {
		for _, connection := range connectionId {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("connectionId=%s", connection)
		}
	}
	if timeAt != nil {
		if !firstParamAttached {
			url += "?"
			firstParamAttached = true
		} else {
			url += "&"
		}
		url += fmt.Sprintf("timeAt=%d", timeAt.Unix())
	}

	var response compliance.BenchmarkEvaluationSummary
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}



func (s *complianceClient) GetBenchmarkControls(ctx *httpclient.Context, benchmarkID string, connectionId []string, timeAt *time.Time) (*compliance.BenchmarkControlSummary, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/%s/controls", s.baseURL, benchmarkID)

	firstParamAttached := false
	if len(connectionId) > 0 {
		for _, connection := range connectionId {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("connectionId=%s", connection)
		}
	}
	if timeAt != nil {
		if !firstParamAttached {
			url += "?"
			firstParamAttached = true
		} else {
			url += "&"
		}
		url += fmt.Sprintf("timeAt=%d", timeAt.Unix())
	}

	var response compliance.BenchmarkControlSummary
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *complianceClient) GetControl(ctx *httpclient.Context, controlID string) (*compliance.Control, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/controls/%s", s.baseURL, controlID)

	var response compliance.Control
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *complianceClient) ListControl(ctx *httpclient.Context, controlIDs []string, tags map[string][]string) ([]compliance.Control, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/controls", s.baseURL)

	firstParamAttached := false
	if len(controlIDs) > 0 {
		for _, controlID := range controlIDs {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("control_id=%s", controlID)
		}
	}
	for tagKey, tagValues := range tags {
		for _, tagValue := range tagValues {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("tag=%s=%s", tagKey, tagValue)
		}
		if len(tagValues) == 0 {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("tag=%s=", tagKey)
		}
	}

	var response []compliance.Control
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (s *complianceClient) ListQueries(ctx *httpclient.Context) ([]compliance.Query, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/queries", s.baseURL)

	var response []compliance.Query
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (s *complianceClient) GetQuery(ctx *httpclient.Context, queryID string) (*compliance.Query, error) {
	url := fmt.Sprintf("%s/api/v1/queries/%s", s.baseURL, queryID)

	var response compliance.Query
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *complianceClient) GetComplianceResults(ctx *httpclient.Context, req compliance.GetComplianceResultsRequest) (compliance.GetComplianceResultsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/compliance_result", s.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return compliance.GetComplianceResultsResponse{}, err
	}

	var response compliance.GetComplianceResultsResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), payload, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return compliance.GetComplianceResultsResponse{}, echo.NewHTTPError(statusCode, err.Error())
		}
		return compliance.GetComplianceResultsResponse{}, err
	}

	return response, nil
}



func (s *complianceClient) ListBenchmarks(ctx *httpclient.Context, tags map[string][]string) ([]compliance.Benchmark, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks", s.baseURL)

	isFirstParamAttached := false
	for tagKey, tagValues := range tags {
		for _, tagValue := range tagValues {
			if !isFirstParamAttached {
				url += "?"
				isFirstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("tag=%s=%s", tagKey, tagValue)
		}
		if len(tagValues) == 0 {
			if !isFirstParamAttached {
				url += "?"
				isFirstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("tag=%s=", tagKey)
		}
	}

	var benchmarks []compliance.Benchmark
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &benchmarks); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return benchmarks, nil
}

func (s *complianceClient) ListAllBenchmarks(ctx *httpclient.Context, isBare bool) ([]compliance.Benchmark, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/all", s.baseURL)

	isFirstParamAttached := false
	if !isBare {
		if isFirstParamAttached {
			url += "&"
		} else {
			url += "?"
			isFirstParamAttached = true
		}
		url += fmt.Sprintf("bare=%v", isBare)
	}

	var benchmarks []compliance.Benchmark
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &benchmarks); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return benchmarks, nil
}

func (s *complianceClient) GetAccountsComplianceResultsSummary(ctx *httpclient.Context, benchmarkId string, connectionIds []string, connector []source.Type) (compliance.GetAccountsComplianceResultsSummaryResponse, error) {
	url := fmt.Sprintf("%s/api/v1/compliance_result/%s/accounts", s.baseURL, benchmarkId)

	var firstParamAttached bool
	firstParamAttached = false

	if len(connectionIds) > 0 {
		for _, connectionId := range connectionIds {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("connectionId=%v", &connectionId)
		}
	}

	if connector != nil {
		if !firstParamAttached {
			url += "?"
			firstParamAttached = true
		} else {
			url += "&"
		}
		url += fmt.Sprintf("connector=%v", &connector)
	}

	var res compliance.GetAccountsComplianceResultsSummaryResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &res); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return compliance.GetAccountsComplianceResultsSummaryResponse{}, echo.NewHTTPError(statusCode, err.Error())
		}
		return compliance.GetAccountsComplianceResultsSummaryResponse{}, err
	}
	return res, nil
}

func (s *complianceClient) CreateBenchmarkAssignment(ctx *httpclient.Context, benchmarkID, connectionId string) ([]compliance.BenchmarkAssignment, error) {
	url := fmt.Sprintf("%s/api/v1/assignments/%s/connection?connectionId=%s", s.baseURL, benchmarkID, connectionId)

	var assignments []compliance.BenchmarkAssignment
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, ctx.ToHeaders(), nil, &assignments); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return assignments, nil
}
