package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/httpclient"

	compliance "github.com/opengovern/opensecurity/services/compliance/api"
)

type ComplianceServiceClient interface {
	ListAssignmentsByBenchmark(ctx *httpclient.Context, benchmarkID string) (*compliance.BenchmarkAssignedEntities, error)
	GetBenchmark(ctx *httpclient.Context, benchmarkID string) (*compliance.Benchmark, error)
	GetBenchmarkSummary(ctx *httpclient.Context, benchmarkID string, connectionId []string, timeAt *time.Time) (*compliance.BenchmarkEvaluationSummary, error)
	GetBenchmarkControls(ctx *httpclient.Context, benchmarkID string, connectionId []string, timeAt *time.Time) (*compliance.BenchmarkControlSummary, error)
	GetControl(ctx *httpclient.Context, controlID string) (*compliance.Control, error)
	ListBenchmarks(ctx *httpclient.Context, frameworkIDs []string, tags map[string][]string) ([]compliance.Benchmark, error)
	ListAllBenchmarks(ctx *httpclient.Context, isBare bool) ([]compliance.Benchmark, error)
	ListQueries(ctx *httpclient.Context) ([]compliance.Policy, error)
	ListControl(ctx *httpclient.Context, controlIDs []string, tags map[string][]string) ([]compliance.Control, error)
	GetControlDetails(ctx *httpclient.Context, controlID string) (*compliance.GetControlDetailsResponse, error)
	ListBenchmarksNestedForBenchmark(ctx *httpclient.Context, benchmarkId string) (*compliance.NestedBenchmark, error)
	PurgeSampleData(ctx *httpclient.Context) error
}

type complianceClient struct {
	baseURL string
}

func NewComplianceClient(baseURL string) ComplianceServiceClient {
	return &complianceClient{baseURL: baseURL}
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

func (s *complianceClient) ListQueries(ctx *httpclient.Context) ([]compliance.Policy, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks/queries", s.baseURL)

	var response []compliance.Policy
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return response, nil
}

func (s *complianceClient) ListBenchmarks(ctx *httpclient.Context, frameworkIDs []string, tags map[string][]string) ([]compliance.Benchmark, error) {
	url := fmt.Sprintf("%s/api/v1/benchmarks", s.baseURL)

	isFirstParamAttached := false
	if len(frameworkIDs) > 0 {
		for _, controlID := range frameworkIDs {
			if !isFirstParamAttached {
				url += "?"
				isFirstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("framework_id=%s", controlID)
		}
	}
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

func (s *complianceClient) ListBenchmarksNestedForBenchmark(ctx *httpclient.Context, benchmarkId string) (*compliance.NestedBenchmark, error) {
	url := fmt.Sprintf("%s/api/v3/benchmarks/%s/nested", s.baseURL, benchmarkId)

	var response compliance.NestedBenchmark
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *complianceClient) PurgeSampleData(ctx *httpclient.Context) error {
	url := fmt.Sprintf("%s/api/v3/sample/purge", s.baseURL)

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}
	return nil
}
