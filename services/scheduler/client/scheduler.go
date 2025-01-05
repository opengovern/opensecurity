package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/opencomply/services/scheduler/db/model"

	"github.com/labstack/echo/v4"

	"github.com/opengovern/opencomply/services/scheduler/api"
)

type TimeRangeFilter struct {
	From int64 // from epoch millisecond
	To   int64 // from epoch millisecond
}
type SchedulerServiceClient interface {
	GetDescribeStatus(ctx *httpclient.Context, resourceType string) ([]api.DescribeStatus, error)
	GetConnectionDescribeStatus(ctx *httpclient.Context, connectionID string) ([]api.IntegrationDescribeStatus, error)
	ListPendingConnections(ctx *httpclient.Context) ([]string, error)
	GetLatestComplianceJobForBenchmark(ctx *httpclient.Context, benchmarkID string) (*api.ComplianceJob, error)
	GetDescribeAllJobsStatus(ctx *httpclient.Context) (*api.DescribeAllJobsStatus, error)
	CountJobsByDate(ctx *httpclient.Context, includeCost *bool, jobType api.JobType, startDate, endDate time.Time) (int64, error)
	GetAsyncQueryRunJobStatus(ctx *httpclient.Context, jobID string) (*api.GetAsyncQueryRunJobStatusResponse, error)
	RunQuery(ctx *httpclient.Context, queryID string) (*model.QueryRunnerJob, error)
	PurgeSampleData(ctx *httpclient.Context, integrations []string) error
	RunDiscovery(ctx *httpclient.Context, userId string, request api.RunDiscoveryRequest) (*api.RunDiscoveryResponse, error)
	ListComplianceJobsHistory(ctx *httpclient.Context, interval, triggerType, createdBy string, cursor, perPage int) (*api.ListComplianceJobsHistoryResponse, error)
	GetSummaryJobs(ctx *httpclient.Context, jobIDs []string) ([]string, error)
	GetIntegrationLastDiscoveryJob(ctx *httpclient.Context, request api.GetIntegrationLastDiscoveryJobRequest) (*model.DescribeIntegrationJob, error)
	GetComplianceQuickSequence(ctx *httpclient.Context, jobID string) (*api.QuickScanSequence, error)
	GetComplianceJobStatus(ctx *httpclient.Context, jobId string) (*api.GetComplianceJobStatusResponse, error)
}

type schedulerClient struct {
	baseURL string
}

func NewSchedulerServiceClient(baseURL string) SchedulerServiceClient {
	return &schedulerClient{baseURL: baseURL}
}

func (s *schedulerClient) RunDiscovery(ctx *httpclient.Context, userId string, request api.RunDiscoveryRequest) (*api.RunDiscoveryResponse, error) {
	url := fmt.Sprintf("%s/api/v3/discovery/run", s.baseURL)

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		httpserver.XPlatformUserIDHeader:   userId,
		httpserver.XPlatformUserRoleHeader: string(ctx.UserRole),
	}

	var response api.RunDiscoveryResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPost, url, headers, payload, &response); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &response, nil
}

func (s *schedulerClient) GetSummaryJobs(ctx *httpclient.Context, jobIDs []string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v3/jobs/compliance/summary/jobs", s.baseURL)
	firstParamAttached := false
	if len(jobIDs) > 0 {
		for _, connection := range jobIDs {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("job_ids=%s", connection)
		}
	}
	var jobs []string
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &jobs); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return jobs, nil
}

func (s *schedulerClient) GetDescribeAllJobsStatus(ctx *httpclient.Context) (*api.DescribeAllJobsStatus, error) {
	url := fmt.Sprintf("%s/api/v1/describe/all/jobs/state", s.baseURL)

	var status api.DescribeAllJobsStatus
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &status); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &status, nil
}

func (s *schedulerClient) GetAsyncQueryRunJobStatus(ctx *httpclient.Context, jobID string) (*api.GetAsyncQueryRunJobStatusResponse, error) {
	url := fmt.Sprintf("%s/api/v3/job/query/%s", s.baseURL, jobID)

	var job api.GetAsyncQueryRunJobStatusResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &job); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &job, nil
}

func (s *schedulerClient) RunQuery(ctx *httpclient.Context, queryID string) (*model.QueryRunnerJob, error) {
	url := fmt.Sprintf("%s/api/v3/query/%s/run", s.baseURL, queryID)

	var job model.QueryRunnerJob
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, &job); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &job, nil
}

func (s *schedulerClient) GetDescribeStatus(ctx *httpclient.Context, resourceType string) ([]api.DescribeStatus, error) {
	url := fmt.Sprintf("%s/api/v1/describe/status/%s", s.baseURL, resourceType)

	var res []api.DescribeStatus
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &res); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *schedulerClient) GetLatestComplianceJobForBenchmark(ctx *httpclient.Context, benchmarkId string) (*api.ComplianceJob, error) {
	url := fmt.Sprintf("%s/api/v1/compliance/status/%s", s.baseURL, benchmarkId)

	var res *api.ComplianceJob
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &res); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *schedulerClient) GetConnectionDescribeStatus(ctx *httpclient.Context, connectionID string) ([]api.IntegrationDescribeStatus, error) {
	url := fmt.Sprintf("%s/api/v1/describe/connection/status?connection_id=%s", s.baseURL, connectionID)

	var res []api.IntegrationDescribeStatus
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &res); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *schedulerClient) ListPendingConnections(ctx *httpclient.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/describe/pending/connections", s.baseURL)

	var res []string
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &res); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return res, nil
}

func (s *schedulerClient) CountJobsByDate(ctx *httpclient.Context, includeCost *bool, jobType api.JobType, startDate, endDate time.Time) (int64, error) {
	url := fmt.Sprintf("%s/api/v1/jobs/bydate?startDate=%d&endDate=%d&jobType=%s", s.baseURL, startDate.UnixMilli(), endDate.UnixMilli(), string(jobType))
	if includeCost != nil {
		url += fmt.Sprintf("&include_cost=%v", *includeCost)
	}

	var resp int64
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return 0, echo.NewHTTPError(statusCode, err.Error())
		}
		return 0, err
	}
	return resp, nil
}

func (s *schedulerClient) ListComplianceJobsHistory(ctx *httpclient.Context, interval, triggerType, createdBy string, cursor, perPage int) (*api.ListComplianceJobsHistoryResponse, error) {
	url := fmt.Sprintf("%s/api/v3/jobs/history/compliance?interval=%s&trigger_type=%s&created_by=%s&cursor=%d&per_page=%d",
		s.baseURL, interval, triggerType, createdBy, cursor, perPage)

	var resp api.ListComplianceJobsHistoryResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &resp); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &resp, nil
}

func (s *schedulerClient) GetIntegrationLastDiscoveryJob(ctx *httpclient.Context, request api.GetIntegrationLastDiscoveryJobRequest) (*model.DescribeIntegrationJob, error) {
	url := fmt.Sprintf("%s/api/v3/integration/discovery/last-job", s.baseURL)
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	var job model.DescribeIntegrationJob
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), payload, &job); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &job, nil
}

func (s *schedulerClient) GetComplianceQuickRun(ctx *httpclient.Context, jobID string) (*api.ComplianceJob, error) {
	url := fmt.Sprintf("%s/api/v3/compliance/quick/%s", s.baseURL, jobID)

	var job api.ComplianceJob
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &job); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &job, nil
}

func (s *schedulerClient) GetComplianceQuickSequence(ctx *httpclient.Context, jobID string) (*api.QuickScanSequence, error) {
	url := fmt.Sprintf("%s/api/v3/compliance/quick/sequence/%s", s.baseURL, jobID)

	var job api.QuickScanSequence
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &job); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &job, nil
}

func (s *schedulerClient) GetComplianceJobStatus(ctx *httpclient.Context, jobId string) (*api.GetComplianceJobStatusResponse, error) {
	url := fmt.Sprintf("%s/api/v3/job/compliance/%s", s.baseURL, jobId)

	var job api.GetComplianceJobStatusResponse
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &job); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, err.Error())
		}
		return nil, err
	}
	return &job, nil
}

func (s *schedulerClient) PurgeSampleData(ctx *httpclient.Context, integrations []string) error {
	url := fmt.Sprintf("%s/api/v3/sample/purge", s.baseURL)
	firstParamAttached := false
	if len(integrations) > 0 {
		for _, connection := range integrations {
			if !firstParamAttached {
				url += "?"
				firstParamAttached = true
			} else {
				url += "&"
			}
			url += fmt.Sprintf("integrations=%s", connection)
		}
	}

	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodPut, url, ctx.ToHeaders(), nil, nil); err != nil {
		if 400 <= statusCode && statusCode < 500 {
			return echo.NewHTTPError(statusCode, err.Error())
		}
		return err
	}
	return nil
}
