package describe

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/opengovern/opencomply/pkg/types"
	"golang.org/x/net/context"

	"github.com/jackc/pgtype"
	"github.com/labstack/echo/v4"
	apiAuth "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	runner2 "github.com/opengovern/opencomply/jobs/compliance-runner-job"
	queryrunner "github.com/opengovern/opencomply/jobs/query-runner-job"
	"github.com/opengovern/opencomply/pkg/utils"
	integrationapi "github.com/opengovern/opencomply/services/integration/api/models"
	"github.com/sony/sonyflake"

	complianceapi "github.com/opengovern/opencomply/services/compliance/api"
	"github.com/opengovern/opencomply/services/scheduler/api"
	"github.com/opengovern/opencomply/services/scheduler/db"
	model2 "github.com/opengovern/opencomply/services/scheduler/db/model"
	"github.com/opengovern/opencomply/services/scheduler/es"
	"go.uber.org/zap"
	"gorm.io/gorm"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type HttpServer struct {
	Address    string
	DB         db.Database
	Scheduler  *Scheduler
	kubeClient k8sclient.Client
}

func NewHTTPServer(
	address string,
	db db.Database,
	s *Scheduler,
) *HttpServer {
	return &HttpServer{
		Address:   address,
		DB:        db,
		Scheduler: s,
	}
}

func (h HttpServer) Register(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	v1.PUT("/describe/trigger/:connection_id", httpserver.AuthorizeHandler(h.TriggerPerConnectionDescribeJob, apiAuth.AdminRole))
	v1.PUT("/describe/trigger", httpserver.AuthorizeHandler(h.TriggerDescribeJob, apiAuth.AdminRole))
	v1.PUT("/compliance/trigger", httpserver.AuthorizeHandler(h.TriggerConnectionsComplianceJobs, apiAuth.AdminRole))
	v1.PUT("/compliance/trigger/:benchmark_id", httpserver.AuthorizeHandler(h.TriggerConnectionsComplianceJob, apiAuth.AdminRole))
	v1.PUT("/compliance/trigger/:benchmark_id/summary", httpserver.AuthorizeHandler(h.TriggerConnectionsComplianceJobSummary, apiAuth.AdminRole))
	v1.GET("/compliance/re-evaluate/:benchmark_id", httpserver.AuthorizeHandler(h.CheckReEvaluateComplianceJob, apiAuth.AdminRole))
	v1.PUT("/compliance/re-evaluate/:benchmark_id", httpserver.AuthorizeHandler(h.ReEvaluateComplianceJob, apiAuth.AdminRole))
	v1.GET("/compliance/status/:benchmark_id", httpserver.AuthorizeHandler(h.GetComplianceBenchmarkStatus, apiAuth.ViewerRole))
	v1.GET("/describe/status/:resource_type", httpserver.AuthorizeHandler(h.GetDescribeStatus, apiAuth.ViewerRole))
	v1.GET("/describe/connection/status", httpserver.AuthorizeHandler(h.GetConnectionDescribeStatus, apiAuth.ViewerRole))
	v1.GET("/describe/pending/connections", httpserver.AuthorizeHandler(h.ListAllPendingConnection, apiAuth.ViewerRole))
	v1.GET("/describe/all/jobs/state", httpserver.AuthorizeHandler(h.GetDescribeAllJobsStatus, apiAuth.ViewerRole))

	v1.POST("/jobs", httpserver.AuthorizeHandler(h.ListJobs, apiAuth.ViewerRole))
	v1.GET("/jobs/bydate", httpserver.AuthorizeHandler(h.CountJobsByDate, apiAuth.ViewerRole))

	v3 := e.Group("/api/v3")
	v3.POST("/jobs/discovery/connections/:connection_id", httpserver.AuthorizeHandler(h.GetDescribeJobsHistory, apiAuth.ViewerRole))
	v3.POST("/jobs/compliance/connections/:connection_id", httpserver.AuthorizeHandler(h.GetComplianceJobsHistory, apiAuth.ViewerRole))
	v3.POST("/jobs/discovery/connections", httpserver.AuthorizeHandler(h.GetDescribeJobsHistoryByIntegration, apiAuth.ViewerRole))
	v3.POST("/jobs/compliance/connections", httpserver.AuthorizeHandler(h.GetComplianceJobsHistoryByIntegration, apiAuth.ViewerRole))

	v3.POST("/compliance/benchmark/:benchmark_id/run", httpserver.AuthorizeHandler(h.RunBenchmarkById, apiAuth.AdminRole))
	v3.POST("/compliance/run", httpserver.AuthorizeHandler(h.RunBenchmark, apiAuth.AdminRole))
	v3.POST("/discovery/run", httpserver.AuthorizeHandler(h.RunDiscovery, apiAuth.AdminRole))
	v3.POST("/discovery/status", httpserver.AuthorizeHandler(h.GetIntegrationDiscoveryProgress, apiAuth.ViewerRole))

	v3.PUT("/query/:query_id/run", httpserver.AuthorizeHandler(h.RunQuery, apiAuth.AdminRole))
	v3.GET("/job/discovery/:job_id", httpserver.AuthorizeHandler(h.GetDescribeJobStatus, apiAuth.ViewerRole))
	v3.GET("/job/compliance/:job_id", httpserver.AuthorizeHandler(h.GetComplianceJobStatus, apiAuth.ViewerRole))
	v3.GET("/job/query/:job_id", httpserver.AuthorizeHandler(h.GetAsyncQueryRunJobStatus, apiAuth.ViewerRole))
	v3.POST("/jobs/discovery", httpserver.AuthorizeHandler(h.ListDescribeJobs, apiAuth.ViewerRole))
	v3.POST("/jobs/compliance", httpserver.AuthorizeHandler(h.ListComplianceJobs, apiAuth.ViewerRole))
	v3.POST("/benchmark/:benchmark_id/run-history", httpserver.AuthorizeHandler(h.BenchmarkAuditHistory, apiAuth.ViewerRole))
	v3.GET("/benchmark/run-history/integrations", httpserver.AuthorizeHandler(h.BenchmarkAuditHistoryIntegrations, apiAuth.ViewerRole))
	v3.PUT("/jobs/cancel/byid", httpserver.AuthorizeHandler(h.CancelJobById, apiAuth.AdminRole))
	v3.POST("/jobs/cancel", httpserver.AuthorizeHandler(h.CancelJob, apiAuth.AdminRole))
	v3.POST("/jobs", httpserver.AuthorizeHandler(h.ListJobsByType, apiAuth.ViewerRole))
	v3.GET("/jobs/interval", httpserver.AuthorizeHandler(h.ListJobsInterval, apiAuth.ViewerRole))
	v3.GET("/jobs/compliance/summary/jobs", httpserver.AuthorizeHandler(h.GetSummaryJobs, apiAuth.ViewerRole))
	v3.GET("/jobs/history/compliance", httpserver.AuthorizeHandler(h.ListComplianceJobsHistory, apiAuth.ViewerRole))

	v3.PUT("/sample/purge", httpserver.AuthorizeHandler(h.PurgeSampleData, apiAuth.AdminRole))

	v3.GET("/integration/discovery/last-job", httpserver.AuthorizeHandler(h.GetIntegrationLastDiscoveryJob, apiAuth.ViewerRole))

	v3.POST("/compliance/quick/sequence", httpserver.AuthorizeHandler(h.CreateComplianceQuickSequence, apiAuth.EditorRole))
	v3.GET("/compliance/quick/sequence/:run_id", httpserver.AuthorizeHandler(h.GetComplianceQuickSequence, apiAuth.ViewerRole))
}

// ListJobs godoc
//
//	@Summary	Lists all jobs
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.ListJobsRequest	true	"List jobs request"
//	@Produce	json
//	@Success	200	{object}	api.ListJobsResponse
//	@Router		/schedule/api/v1/jobs [post]
func (h HttpServer) ListJobs(ctx echo.Context) error {
	var request api.ListJobsRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var jobs []api.Job

	srcs, err := h.Scheduler.integrationClient.ListIntegrations(&httpclient.Context{UserRole: apiAuth.AdminRole}, nil)
	if err != nil {
		return err
	}

	benchmarks, err := h.Scheduler.complianceClient.ListBenchmarks(&httpclient.Context{UserRole: apiAuth.AdminRole}, nil)
	if err != nil {
		return err
	}

	sortBy := "id"
	switch request.SortBy {
	case api.JobSort_ByIntegrationID, api.JobSort_ByJobID, api.JobSort_ByJobType, api.JobSort_ByStatus:
		sortBy = string(request.SortBy)
	}

	sortOrder := "DESC"
	if request.SortOrder == api.JobSortOrder_ASC {
		sortOrder = "ASC"
	}

	var startTime, endTime *time.Time
	if request.Interval != nil {
		tmpInterval, err := convertInterval(*request.Interval)
		if err != nil {
			h.Scheduler.logger.Error("failed to parse interval", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to parse interval")
		}
		request.Interval = &tmpInterval
	} else if request.To != nil && request.From != nil {
		endTimeTmp := time.Unix(*request.To, 0)
		endTime = &endTimeTmp

		startTimeTmp := time.Unix(*request.From, 0)
		startTime = &startTimeTmp
	}

	describeJobs, err := h.DB.ListAllJobs(request.PageStart, request.PageEnd, request.Interval, startTime, endTime, request.TypeFilters,
		request.StatusFilter, sortBy, sortOrder)
	if err != nil {
		return err
	}
	for _, job := range describeJobs {
		var jobSRC integrationapi.Integration
		for _, src := range srcs.Integrations {
			if src.IntegrationID == job.ConnectionID {
				jobSRC = src
			}
		}

		if job.JobType == "compliance" {
			for _, benchmark := range benchmarks {
				if fmt.Sprintf("%v", benchmark.ID) == job.Title {
					job.Title = benchmark.Title
				}
			}
		}

		jobs = append(jobs, api.Job{
			ID:                     job.ID,
			CreatedAt:              job.CreatedAt,
			UpdatedAt:              job.UpdatedAt,
			Type:                   api.JobType(job.JobType),
			ConnectionID:           job.ConnectionID,
			ConnectionProviderID:   jobSRC.ProviderID,
			ConnectionProviderName: jobSRC.Name,
			Title:                  job.Title,
			Status:                 job.Status,
			FailureReason:          job.FailureMessage,
		})
	}

	var jobSummaries []api.JobSummary
	summaries, err := h.DB.GetAllJobSummary(request.Interval, startTime, endTime, request.TypeFilters, request.StatusFilter)
	if err != nil {
		return err
	}
	for _, summary := range summaries {
		jobSummaries = append(jobSummaries, api.JobSummary{
			Type:   api.JobType(summary.JobType),
			Status: summary.Status,
			Count:  summary.Count,
		})
	}

	switch strings.ToLower(sortOrder) {
	case "asc":
		switch strings.ToLower(string(request.SortBy)) {
		case "id":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].ID < jobs[j].ID
			})
		case "connection_id":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].ConnectionID < jobs[j].ConnectionID
			})
		case "job_type":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].Type < jobs[j].Type
			})
		case "status":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].Status < jobs[j].Status
			})
		case "created_at", "createdat":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].CreatedAt.Before(jobs[j].CreatedAt)
			})
		case "updated_at", "updatedat":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].UpdatedAt.Before(jobs[j].UpdatedAt)
			})
		default:
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].UpdatedAt.Before(jobs[j].UpdatedAt)
			})
		}
	case "desc":
		switch strings.ToLower(string(request.SortBy)) {
		case "id":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].ID > jobs[j].ID
			})
		case "connection_id":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].ConnectionID > jobs[j].ConnectionID
			})
		case "job_type":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].Type > jobs[j].Type
			})
		case "status":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].Status > jobs[j].Status
			})
		case "created_at", "createdat":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
			})
		case "updated_at", "updatedat":
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].UpdatedAt.After(jobs[j].UpdatedAt)
			})
		default:
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].UpdatedAt.After(jobs[j].UpdatedAt)
			})
		}
	}

	return ctx.JSON(http.StatusOK, api.ListJobsResponse{
		Jobs:      jobs,
		Summaries: jobSummaries,
	})
}

func UniqueArray[T any](arr []T) []T {
	m := make(map[string]T)
	for _, item := range arr {
		// hash the item
		hash := sha1.New()
		hash.Write([]byte(fmt.Sprintf("%v", item)))
		hashResult := hash.Sum(nil)
		m[fmt.Sprintf("%x", hashResult)] = item
	}
	var resp []T
	for _, v := range m {
		resp = append(resp, v)
	}
	return resp
}

func (h HttpServer) CountJobsByDate(ctx echo.Context) error {
	startDate, err := strconv.ParseInt(ctx.QueryParam("startDate"), 10, 64)
	if err != nil {
		return err
	}
	endDate, err := strconv.ParseInt(ctx.QueryParam("endDate"), 10, 64)
	if err != nil {
		return err
	}
	includeCostStr := ctx.QueryParam("include_cost")

	var count int64
	switch api.JobType(ctx.QueryParam("jobType")) {
	case api.JobType_Discovery:
		var includeCost *bool
		if len(includeCostStr) > 0 {
			v, err := strconv.ParseBool(includeCostStr)
			if err != nil {
				return err
			}

			includeCost = &v
		}
		count, err = h.DB.CountDescribeJobsByDate(includeCost, time.UnixMilli(startDate), time.UnixMilli(endDate))

	case api.JobType_Compliance:
		count, err = h.DB.CountComplianceJobsByDate(true, time.UnixMilli(startDate), time.UnixMilli(endDate))
	default:
		return errors.New("invalid job type")
	}
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, count)
}

// TriggerPerConnectionDescribeJob godoc
//
//	@Summary		Triggers describer
//	@Description	Triggers a describe job to run immediately for the given connection
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Success		200
//	@Param			connection_id	path	string		true	"Connection ID"
//	@Param			force_full		query	bool		false	"Force full discovery"
//	@Param			resource_type	query	[]string	false	"Resource Type"
//	@Param			cost_discovery	query	bool		false	"Cost discovery"
//	@Router			/schedule/api/v1/describe/trigger/{connection_id} [put]
func (h HttpServer) TriggerPerConnectionDescribeJob(ctx echo.Context) error {
	connectionID := ctx.Param("connection_id")
	costFullDiscovery := ctx.QueryParam("cost_full_discovery") == "true"
	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}

	ctx2 := &httpclient.Context{UserRole: apiAuth.AdminRole}
	ctx2.Ctx = ctx.Request().Context()

	var srcs []integrationapi.Integration
	if connectionID == "all" {
		var err error
		integrationsResp, err := h.Scheduler.integrationClient.ListIntegrations(ctx2, nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		srcs = integrationsResp.Integrations
	} else {
		src, err := h.Scheduler.integrationClient.GetIntegration(ctx2, connectionID)
		if err != nil || src == nil {
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid connection id")
			}
		}
		srcs = []integrationapi.Integration{*src}
	}

	dependencyIDs := make([]int64, 0)
	var err error

	integrationTypes, err := h.Scheduler.integrationClient.ListIntegrationTypes(ctx2)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	integrationTypesMap := make(map[string]bool)
	for _, integrationType := range integrationTypes {
		integrationTypesMap[integrationType] = true
	}

	for _, src := range srcs {
		_, ok := integrationTypesMap[src.IntegrationType.String()]
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "unknown integration type")
		}

		resourceTypes := ctx.QueryParams()["resource_type"]

		if resourceTypes == nil {
			resourceTypesMap, err := h.Scheduler.integrationClient.GetResourceTypesByLabels(ctx2, src.IntegrationType.String(), src.Labels)
			if err != nil {
				h.Scheduler.logger.Error("failed to get resource types by labels", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get resource types by labels")
			}
			for rt, _ := range resourceTypesMap {
				resourceTypes = append(resourceTypes, rt)
			}
		}

		for _, resourceType := range resourceTypes {
			daj, err := h.Scheduler.describe(src, resourceType, false, costFullDiscovery, false, nil, userID, nil)
			if err != nil && errors.Is(err, ErrJobInProgress) {
				return echo.NewHTTPError(http.StatusConflict, err.Error())
			}
			if err != nil {
				return err
			}
			dependencyIDs = append(dependencyIDs, int64(daj.ID))
		}
	}

	err = h.DB.CreateJobSequencer(&model2.JobSequencer{
		DependencyList:   dependencyIDs,
		DependencySource: model2.JobSequencerJobTypeDescribe,
		NextJob:          model2.JobSequencerJobTypeBenchmark,
		Status:           model2.JobSequencerWaitingForDependencies,
	})
	if err != nil {
		return fmt.Errorf("failed to create job sequencer: %v", err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h HttpServer) TriggerDescribeJob(ctx echo.Context) error {
	resourceTypes := httpserver.QueryArrayParam(ctx, "resource_type")
	connectors := httpserver.QueryArrayParam(ctx, "connector")
	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}

	integrations, err := h.Scheduler.integrationClient.ListIntegrations(&httpclient.Context{UserRole: apiAuth.AdminRole}, connectors)
	if err != nil {
		h.Scheduler.logger.Error("failed to get list of sources", zap.String("spot", "ListSources"), zap.Error(err))
		DescribeJobsCount.WithLabelValues("failure").Inc()
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
	}

	ctx2 := &httpclient.Context{UserRole: apiAuth.AdminRole}
	ctx2.Ctx = ctx.Request().Context()

	integrationTypes, err := h.Scheduler.integrationClient.ListIntegrationTypes(ctx2)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	integrationTypesMap := make(map[string]bool)
	for _, integrationType := range integrationTypes {
		integrationTypesMap[integrationType] = true
	}

	for _, integration := range integrations.Integrations {
		if integration.State != integrationapi.IntegrationStateActive {
			continue
		}
		_, ok := integrationTypesMap[integration.IntegrationType.String()]
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "unknown integration type")
		}
		rtToDescribe := resourceTypes

		if len(rtToDescribe) == 0 {
			resourceTypesMap, err := h.Scheduler.integrationClient.GetResourceTypesByLabels(ctx2, integration.IntegrationType.String(), integration.Labels)
			if err != nil {
				h.Scheduler.logger.Error("failed to get resource types by labels", zap.Error(err))
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to get resource types by labels")
			}
			for rt, _ := range resourceTypesMap {
				rtToDescribe = append(rtToDescribe, rt)
			}
		}

		for _, resourceType := range rtToDescribe {
			_, err = h.Scheduler.describe(integration, resourceType, false, false, false, nil, userID, nil)
			if err != nil {
				h.Scheduler.logger.Error("failed to describe connection", zap.String("integration_id", integration.IntegrationID), zap.Error(err))
			}
		}
	}
	return ctx.JSON(http.StatusOK, "")
}

// TriggerConnectionsComplianceJob godoc
//
//	@Summary		Triggers compliance job
//	@Description	Triggers a compliance job to run immediately for the given benchmark
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Success		200
//	@Param			benchmark_id	path	string		true	"Benchmark ID"
//	@Param			connection_id	query	[]string	false	"Connection ID"
//	@Router			/schedule/api/v1/compliance/trigger/{benchmark_id} [put]
func (h HttpServer) TriggerConnectionsComplianceJob(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}

	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}
	benchmarkID := ctx.Param("benchmark_id")
	benchmark, err := h.Scheduler.complianceClient.GetBenchmark(clientCtx, benchmarkID)
	if err != nil {
		return fmt.Errorf("error while getting benchmarks: %v", err)
	}

	if benchmark == nil {
		return echo.NewHTTPError(http.StatusNotFound, "benchmark not found")
	}

	connectionIDs := httpserver.QueryArrayParam(ctx, "connection_id")

	lastJob, err := h.Scheduler.db.GetLastComplianceJob(true, benchmark.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if lastJob != nil && (lastJob.Status == model2.ComplianceJobRunnersInProgress ||
		lastJob.Status == model2.ComplianceJobSummarizerInProgress ||
		lastJob.Status == model2.ComplianceJobCreated) {
		return echo.NewHTTPError(http.StatusConflict, "compliance job is already running")
	}

	_, err = h.Scheduler.complianceScheduler.CreateComplianceReportJobs(true, benchmarkID, lastJob, connectionIDs, true, userID, nil)

	return ctx.JSON(http.StatusOK, "")
}

// TriggerConnectionsComplianceJobs godoc
//
//	@Summary		Triggers compliance job
//	@Description	Triggers a compliance job to run immediately for the given benchmark
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Success		200
//	@Param			benchmark_id	query	[]string	true	"Benchmark IDs leave empty for everything"
//	@Param			connection_id	query	[]string	false	"Connection IDs leave empty for default (enabled connections)"
//	@Router			/schedule/api/v1/compliance/trigger [put]
func (h HttpServer) TriggerConnectionsComplianceJobs(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}

	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}
	benchmarkIDs := httpserver.QueryArrayParam(ctx, "benchmark_id")

	connectionIDs := httpserver.QueryArrayParam(ctx, "connection_id")

	var benchmarks []complianceapi.Benchmark
	var err error
	if len(benchmarkIDs) == 0 {
		benchmarks, err = h.Scheduler.complianceClient.ListBenchmarks(clientCtx, nil)
		if err != nil {
			return fmt.Errorf("error while getting benchmarks: %v", err)
		}
	} else {
		for _, benchmarkID := range benchmarkIDs {
			benchmark, err := h.Scheduler.complianceClient.GetBenchmark(clientCtx, benchmarkID)
			if err != nil {
				return fmt.Errorf("error while getting benchmarks: %v", err)
			}
			if benchmark == nil {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("benchmark %s not found", benchmarkID))
			}
			benchmarks = append(benchmarks, *benchmark)
		}
	}

	for _, benchmark := range benchmarks {
		lastJob, err := h.Scheduler.db.GetLastComplianceJob(true, benchmark.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if lastJob != nil && (lastJob.Status == model2.ComplianceJobRunnersInProgress ||
			lastJob.Status == model2.ComplianceJobSummarizerInProgress ||
			lastJob.Status == model2.ComplianceJobCreated) {
			return echo.NewHTTPError(http.StatusConflict, "compliance job is already running")
		}

		_, err = h.Scheduler.complianceScheduler.CreateComplianceReportJobs(true, benchmark.ID, lastJob, connectionIDs, true, userID, nil)
		if err != nil {
			return fmt.Errorf("error while creating compliance job: %v", err)
		}
	}
	return ctx.JSON(http.StatusOK, "")
}

// TriggerConnectionsComplianceJobSummary godoc
//
//	@Summary		Triggers compliance job
//	@Description	Triggers a compliance job to run immediately for the given benchmark
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Success		200
//	@Param			benchmark_id	path	string	true	"Benchmark ID use 'all' for everything"
//	@Router			/schedule/api/v1/compliance/trigger/{benchmark_id}/summary [put]
func (h HttpServer) TriggerConnectionsComplianceJobSummary(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}
	benchmarkID := ctx.Param("benchmark_id")

	var benchmarks []complianceapi.Benchmark
	var err error
	if benchmarkID == "all" {
		benchmarks, err = h.Scheduler.complianceClient.ListBenchmarks(clientCtx, nil)
		if err != nil {
			return fmt.Errorf("error while getting benchmarks: %v", err)
		}
	} else {
		benchmark, err := h.Scheduler.complianceClient.GetBenchmark(clientCtx, benchmarkID)
		if err != nil {
			return fmt.Errorf("error while getting benchmarks: %v", err)
		}
		if benchmark == nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("benchmark %s not found", benchmarkID))
		}
		benchmarks = append(benchmarks, *benchmark)
	}

	for _, benchmark := range benchmarks {
		err = h.Scheduler.complianceScheduler.CreateSummarizer(benchmark.ID, nil, nil, model2.ComplianceTriggerTypeManual)
		if err != nil {
			return fmt.Errorf("error while creating compliance job summarizer: %v", err)
		}
	}
	return ctx.JSON(http.StatusOK, "")
}

type ReEvaluateDescribeJob struct {
	Integration  integrationapi.Integration
	ResourceType string
}

func (h HttpServer) getReEvaluateParams(benchmarkID string, connectionIDs, controlIDs []string) (*model2.JobSequencerJobTypeBenchmarkRunnerParameters, []ReEvaluateDescribeJob, error) {
	integrations, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(&httpclient.Context{UserRole: apiAuth.AdminRole}, integrationapi.ListIntegrationsRequest{
		IntegrationID: connectionIDs,
	})
	if err != nil {
		h.Scheduler.logger.Error("failed to get connections", zap.Error(err))
		return nil, nil, err
	}
	var describeJobs []ReEvaluateDescribeJob

	ctx2 := &httpclient.Context{UserRole: apiAuth.AdminRole}
	ctx2.Ctx = context.Background()
	integrationTypes, err := h.Scheduler.integrationClient.ListIntegrationTypes(ctx2)
	if err != nil {
		return nil, nil, err
	}
	integrationTypesMap := make(map[string]bool)
	for _, integrationType := range integrationTypes {
		integrationTypesMap[integrationType] = true
	}

	// TODO: filter needed resource types for tables for controls queries
	for _, integration := range integrations.Integrations {
		if integration.State != integrationapi.IntegrationStateActive {
			continue
		}
		_, ok := integrationTypesMap[integration.IntegrationType.String()]
		if !ok {
			return nil, nil, fmt.Errorf("unknown integration type")
		}

		possibleRt, err := h.Scheduler.integrationClient.GetResourceTypesByLabels(ctx2, integration.IntegrationType.String(), integration.Labels)
		if err != nil {
			h.Scheduler.logger.Error("failed to get resource types by labels", zap.Error(err))
			return nil, nil, fmt.Errorf("failed to get resource types by labels")
		}

		for resourceType, _ := range possibleRt {
			describeJobs = append(describeJobs, ReEvaluateDescribeJob{
				Integration:  integration,
				ResourceType: resourceType,
			})
		}
	}

	return &model2.JobSequencerJobTypeBenchmarkRunnerParameters{
		BenchmarkID:   benchmarkID,
		ControlIDs:    controlIDs,
		ConnectionIDs: connectionIDs,
	}, describeJobs, nil
}

func (h HttpServer) getBenchmarkChildrenControls(benchmarkID string) ([]string, error) {
	benchmark, err := h.Scheduler.complianceClient.GetBenchmark(&httpclient.Context{UserRole: apiAuth.AdminRole}, benchmarkID)
	if err != nil {
		h.Scheduler.logger.Error("failed to get benchmark", zap.Error(err))
		return nil, err
	}
	if benchmark == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "benchmark not found")
	}

	var controlIDs []string
	for _, control := range benchmark.Controls {
		controlIDs = append(controlIDs, control)
	}
	for _, childBenchmarkID := range benchmark.Children {
		childControlIDs, err := h.getBenchmarkChildrenControls(childBenchmarkID)
		if err != nil {
			return nil, err
		}
		controlIDs = append(controlIDs, childControlIDs...)
	}
	return controlIDs, nil
}

// ReEvaluateComplianceJob godoc
//
//	@Summary		Re-evaluates compliance job
//	@Description	Triggers a discovery job to run immediately for the given connection then triggers compliance job
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Success		200
//	@Param			benchmark_id	path	string		true	"Benchmark ID"
//	@Param			integrationID	query	[]string	true	"Connection ID"
//	@Param			control_id		query	[]string	false	"Control ID"
//	@Router			/schedule/api/v1/compliance/re-evaluate/{benchmark_id} [put]
func (h HttpServer) ReEvaluateComplianceJob(ctx echo.Context) error {
	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}
	benchmarkID := ctx.Param("benchmark_id")
	integrationID := httpserver.QueryArrayParam(ctx, "integrationID")
	if len(integrationID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "connection_id is required")
	}
	controlIDs := httpserver.QueryArrayParam(ctx, "control_id")

	jobParameters, describeJobs, err := h.getReEvaluateParams(benchmarkID, integrationID, controlIDs)
	if err != nil {
		return err
	}

	h.Scheduler.logger.Info("re-evaluating compliance job", zap.Any("job_parameters", jobParameters))

	jobParametersJSON, err := json.Marshal(jobParameters)
	if err != nil {
		h.Scheduler.logger.Error("failed to marshal job parameters", zap.Error(err))
		return err
	}

	jp := pgtype.JSONB{}
	err = jp.Set(jobParametersJSON)
	if err != nil {
		h.Scheduler.logger.Error("failed to set job parameters", zap.Error(err))
		return err
	}

	var dependencyIDs []int64
	for _, describeJob := range describeJobs {
		daj, err := h.Scheduler.describe(describeJob.Integration, describeJob.ResourceType, false, false, false, nil, userID, nil)
		if err != nil {
			h.Scheduler.logger.Error("failed to describe connection", zap.String("integration_id", describeJob.Integration.IntegrationID), zap.Error(err))
			continue
		}
		dependencyIDs = append(dependencyIDs, int64(daj.ID))
	}

	err = h.DB.CreateJobSequencer(&model2.JobSequencer{
		DependencyList:    dependencyIDs,
		DependencySource:  model2.JobSequencerJobTypeDescribe,
		NextJob:           model2.JobSequencerJobTypeBenchmarkRunner,
		NextJobParameters: &jp,
		Status:            model2.JobSequencerWaitingForDependencies,
	})
	if err != nil {
		h.Scheduler.logger.Error("failed to create job sequencer", zap.Error(err))
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

// CheckReEvaluateComplianceJob godoc
//
//	@Summary		Get re-evaluates compliance job
//	@Description	Get re-evaluate job for the given connection and control
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Param			benchmark_id	path		string		true	"Benchmark ID"
//	@Param			integrationID	query		[]string	true	"Connection ID"
//	@Param			control_id		query		[]string	false	"Control ID"
//	@Success		200				{object}	api.JobSeqCheckResponse
//	@Router			/schedule/api/v1/compliance/re-evaluate/{benchmark_id} [get]
func (h HttpServer) CheckReEvaluateComplianceJob(ctx echo.Context) error {
	benchmarkID := ctx.Param("benchmark_id")
	integrationID := httpserver.QueryArrayParam(ctx, "integrationID")
	if len(integrationID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "integrationID is required")
	}
	controlIDs := httpserver.QueryArrayParam(ctx, "control_id")

	jobParameters, describeJobs, err := h.getReEvaluateParams(benchmarkID, integrationID, controlIDs)
	if err != nil {
		return err
	}

	var dependencyIDs []int64
	for _, describeJob := range describeJobs {
		daj, err := h.Scheduler.db.GetLastDescribeIntegrationJob(describeJob.Integration.IntegrationID, describeJob.ResourceType)
		if err != nil {
			h.Scheduler.logger.Error("failed to describe connection", zap.String("integration_id", describeJob.Integration.IntegrationID), zap.Error(err))
			continue
		}
		dependencyIDs = append(dependencyIDs, int64(daj.ID))
	}

	jobs, err := h.Scheduler.db.ListJobSequencersOfTypeOfToday(model2.JobSequencerJobTypeDescribe, model2.JobSequencerJobTypeBenchmarkRunner)
	if err != nil {
		return err
	}

	var theJob *model2.JobSequencer
	for _, job := range jobs {
		var params model2.JobSequencerJobTypeBenchmarkRunnerParameters
		err := json.Unmarshal(job.NextJobParameters.Bytes, &params)
		if err != nil {
			h.Scheduler.logger.Error("failed to unmarshal job parameters", zap.Error(err))
			return err
		}

		fmt.Println(">>>", job)
		fmt.Println("<<<", params, dependencyIDs)
		fmt.Println("----", params.BenchmarkID, jobParameters.BenchmarkID)
		fmt.Println("----", params.ConnectionIDs, jobParameters.ConnectionIDs)
		fmt.Println("----", params.ControlIDs, jobParameters.ControlIDs)
		fmt.Println("----", job.DependencyList, dependencyIDs)

		if params.BenchmarkID == jobParameters.BenchmarkID &&
			utils.IncludesAll(params.ConnectionIDs, jobParameters.ConnectionIDs) &&
			utils.IncludesAll(params.ControlIDs, jobParameters.ControlIDs) &&
			utils.IncludesAll(job.DependencyList, dependencyIDs) {
			theJob = &job
			break
		}
	}

	if theJob == nil || theJob.Status == model2.JobSequencerFailed {
		fmt.Println("job not found/failed", theJob)
		return ctx.JSON(http.StatusOK, api.JobSeqCheckResponse{
			IsRunning: false,
		})
	}

	if theJob.Status == model2.JobSequencerWaitingForDependencies {
		fmt.Println("job waiting", theJob)
		return ctx.JSON(http.StatusOK, api.JobSeqCheckResponse{
			IsRunning: true,
		})
	}

	var nid []int64
	for _, m := range strings.Split(theJob.NextJobIDs, ",") {
		i, _ := strconv.ParseInt(m, 10, 64)
		nid = append(nid, i)
	}
	runnerJobs, err := h.Scheduler.db.ListRunnersWithID(nid)
	if err != nil {
		return err
	}
	for _, runner := range runnerJobs {
		if runner.Status != runner2.ComplianceRunnerSucceeded &&
			runner.Status != runner2.ComplianceRunnerFailed &&
			runner.Status != runner2.ComplianceRunnerTimeOut {
			fmt.Println("+++ job status", runner.Status)

			return ctx.JSON(http.StatusOK, api.JobSeqCheckResponse{
				IsRunning: true,
			})
		}
	}

	fmt.Println("job finished", theJob)
	return ctx.JSON(http.StatusOK, api.JobSeqCheckResponse{
		IsRunning: false,
	})
}

func (h HttpServer) GetComplianceBenchmarkStatus(ctx echo.Context) error {
	benchmarkId := ctx.Param("benchmark_id")
	lastComplianceJob, err := h.Scheduler.db.GetLastComplianceJob(true, benchmarkId)
	if err != nil {
		h.Scheduler.logger.Error("failed to get compliance job", zap.String("benchmark_id", benchmarkId), zap.Error(err))
		return err
	}
	if lastComplianceJob == nil {
		return ctx.JSON(http.StatusOK, nil)
	}
	return ctx.JSON(http.StatusOK, lastComplianceJob.ToApi())
}

func (h HttpServer) GetDescribeStatus(ctx echo.Context) error {
	resourceType := ctx.Param("resource_type")

	status, err := h.DB.GetDescribeStatus(resourceType)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, status)
}

// GetConnectionDescribeStatus godoc
//
//	@Summary	Get connection describe status
//	@Security	BearerToken
//	@Tags		describe
//	@Produce	json
//	@Success	200
//	@Param		connection_id	query	string	true	"Connection ID"
//	@Router		/schedule/api/v1/describe/connection/status [put]
func (h HttpServer) GetConnectionDescribeStatus(ctx echo.Context) error {
	connectionID := ctx.QueryParam("connection_id")

	status, err := h.DB.GetIntegrationDescribeStatus(connectionID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, status)
}

func (h HttpServer) ListAllPendingConnection(ctx echo.Context) error {
	status, err := h.DB.ListAllPendingIntegration()
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, status)
}

func (h HttpServer) GetDescribeAllJobsStatus(ctx echo.Context) error {
	count, _, err := h.DB.CountJobsAndResources()
	if err != nil {
		return err
	}

	if count == nil || *count == 0 {
		return ctx.JSON(http.StatusOK, api.DescribeAllJobsStatusNoJobToRun)
	}

	pendingDiscoveryTypes, err := h.DB.ListAllFirstTryPendingIntegration()
	if err != nil {
		return err
	}

	for _, dt := range pendingDiscoveryTypes {
		if dt == string(model2.DiscoveryType_Cost) || dt == string(model2.DiscoveryType_Fast) {
			return ctx.JSON(http.StatusOK, api.DescribeAllJobsStatusJobsRunning)
		}
	}

	succeededJobs, err := h.DB.ListAllSuccessfulDescribeJobs()
	if err != nil {
		return err
	}

	publishedJobs := 0
	totalJobs := 0
	for _, job := range succeededJobs {
		totalJobs++

		if job.DescribedResourceCount > 0 {
			resourceCount, err := es.GetInventoryCountResponse(ctx.Request().Context(), h.Scheduler.es, strings.ToLower(job.ResourceType))
			if err != nil {
				return err
			}

			if resourceCount > 0 {
				publishedJobs++
			}
		} else {
			publishedJobs++
		}
	}

	h.Scheduler.logger.Info("job count",
		zap.Int("publishedJobs", publishedJobs),
		zap.Int("totalJobs", totalJobs),
	)
	if publishedJobs == totalJobs {
		return ctx.JSON(http.StatusOK, api.DescribeAllJobsStatusResourcesPublished)
	}

	job, err := h.DB.GetLastSuccessfulDescribeJob()
	if err != nil {
		return err
	}

	if job != nil &&
		job.UpdatedAt.Before(time.Now().Add(-5*time.Minute)) {
		return ctx.JSON(http.StatusOK, api.DescribeAllJobsStatusResourcesPublished)
	}

	return ctx.JSON(http.StatusOK, api.DescribeAllJobsStatusJobsFinished)
}

type MigratorResponse struct {
	Hits  MigratorHits `json:"hits"`
	PitID string
}
type MigratorHits struct {
	Total opengovernance.SearchTotal `json:"total"`
	Hits  []MigratorHit              `json:"hits"`
}
type MigratorHit struct {
	ID      string        `json:"_id"`
	Score   float64       `json:"_score"`
	Index   string        `json:"_index"`
	Type    string        `json:"_type"`
	Version int64         `json:"_version,omitempty"`
	Source  MigrateSource `json:"_source"`
	Sort    []any         `json:"sort"`
}

type MigrateSource map[string]any

func (m MigrateSource) KeysAndIndex() ([]string, string) {
	return nil, ""
}

func bindValidate(ctx echo.Context, i interface{}) error {
	if err := ctx.Bind(i); err != nil {
		return err
	}

	if err := ctx.Validate(i); err != nil {
		return err
	}

	return nil
}

// function to remove duplicate values
func removeDuplicates(s []string) []string {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}

// GetDescribeJobsHistory godoc
//
//	@Summary	Get describe jobs history for give connection
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request			body	api.GetDescribeJobsHistoryRequest	true	"List jobs request"
//	@Param		connection_id	path	string								true	"Connection ID"
//	@Produce	json
//	@Success	200	{object}	[]api.GetDescribeJobsHistoryResponse
//	@Router		/schedule/api/v3/jobs/discovery/connections/{connection_id} [post]
func (h HttpServer) GetDescribeJobsHistory(ctx echo.Context) error {
	connectionId := ctx.Param("connection_id")

	var request api.GetDescribeJobsHistoryRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var jobsResults []api.GetDescribeJobsHistoryResponse

	jobs, err := h.DB.ListDescribeJobsByFilters(nil, []string{connectionId}, request.ResourceType,
		request.DiscoveryType, request.JobStatus, &request.StartTime, request.EndTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, j := range jobs {
		jobsResults = append(jobsResults, api.GetDescribeJobsHistoryResponse{
			JobId:          j.ID,
			ResourceType:   j.ResourceType,
			JobStatus:      j.Status,
			UpdatedAt:      j.UpdatedAt,
			CreatedAt:      j.CreatedAt,
			FailureMessage: j.FailureMessage,
			Title:          j.ResourceType,
		})
	}
	if request.SortBy != nil {
		switch strings.ToLower(*request.SortBy) {
		case "id":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		case "datetime", "updated_at", "updatedat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].UpdatedAt.After(jobsResults[j].UpdatedAt)
			})
		case "created_at", "createdat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].CreatedAt.After(jobsResults[j].CreatedAt)
			})
		case "resourcetype":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].ResourceType < jobsResults[j].ResourceType
			})
		case "jobstatus":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobStatus < jobsResults[j].JobStatus
			})
		default:
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		}
	} else {
		sort.Slice(jobsResults, func(i, j int) bool {
			return jobsResults[i].JobId < jobsResults[j].JobId
		})
	}
	if request.PerPage != nil {
		if request.Cursor == nil {
			jobsResults = utils.Paginate(1, *request.PerPage, jobsResults)
		} else {
			jobsResults = utils.Paginate(*request.Cursor, *request.PerPage, jobsResults)
		}
	}

	return ctx.JSON(http.StatusOK, jobsResults)
}

// GetComplianceJobsHistory godoc
//
//	@Summary	Get compliance jobs history for give connection
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request			body	api.GetComplianceJobsHistoryRequest	true	"List jobs request"
//	@Param		connection_id	path	string								true	"Connection ID"
//	@Produce	json
//	@Success	200	{object}	[]api.GetComplianceJobsHistoryResponse
//	@Router		/schedule/api/v3/jobs/compliance/connections/{connection_id} [post]
func (h HttpServer) GetComplianceJobsHistory(ctx echo.Context) error {
	var request api.GetComplianceJobsHistoryRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	connectionId := ctx.Param("connection_id")

	jobs, err := h.DB.ListComplianceJobsByFilters(nil, []string{connectionId}, request.BenchmarkId, request.JobStatus, &request.StartTime, request.EndTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var jobsResults []api.GetComplianceJobsHistoryResponse
	for _, j := range jobs {
		jobsResults = append(jobsResults, api.GetComplianceJobsHistoryResponse{
			JobId:         j.ID,
			WithIncidents: j.WithIncidents,
			BenchmarkId:   j.FrameworkID,
			JobStatus:     j.Status.ToApi(),
			UpdatedAt:     j.UpdatedAt,
			CreatedAt:     j.CreatedAt,
		})
	}
	if request.SortBy != nil {
		switch strings.ToLower(*request.SortBy) {
		case "id":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		case "datetime", "updated_at", "updatedat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].UpdatedAt.After(jobsResults[j].UpdatedAt)
			})
		case "created_at", "createdat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].CreatedAt.After(jobsResults[j].CreatedAt)
			})
		case "benchmarkid":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].BenchmarkId < jobsResults[j].BenchmarkId
			})
		case "jobstatus":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobStatus < jobsResults[j].JobStatus
			})
		case "with_incidents", "withincidents":
			sort.Slice(jobsResults, func(i, j int) bool {
				return !jobsResults[i].WithIncidents && jobsResults[j].WithIncidents
			})
		default:
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		}
	} else {
		sort.Slice(jobsResults, func(i, j int) bool {
			return jobsResults[i].JobId < jobsResults[j].JobId
		})
	}
	if request.PerPage != nil {
		if request.Cursor == nil {
			jobsResults = utils.Paginate(1, *request.PerPage, jobsResults)
		} else {
			jobsResults = utils.Paginate(*request.Cursor, *request.PerPage, jobsResults)
		}
	}

	return ctx.JSON(http.StatusOK, jobsResults)
}

// RunBenchmarkById godoc
//
//	@Summary		Triggers compliance job by benchmark id
//	@Description	Triggers a compliance job to run immediately for the given benchmark
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Param			benchmark_id	path		string						true	"Benchmark ID"
//	@Param			request			body		api.RunBenchmarkByIdRequest	true	"Integrations filter"
//	@Success		200				{object}	api.RunBenchmarkResponse
//	@Router			/schedule/api/v3/compliance/benchmark/{benchmark_id}/run [post]
func (h HttpServer) RunBenchmarkById(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}
	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}
	benchmarkID := strings.ToLower(ctx.Param("benchmark_id"))

	var request api.RunBenchmarkByIdRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	benchmark, err := h.Scheduler.complianceClient.GetBenchmark(&httpclient.Context{UserRole: apiAuth.AdminRole}, benchmarkID)
	if err != nil {
		return fmt.Errorf("error while getting benchmarks: %v", err)
	}
	if benchmark == nil {
		return echo.NewHTTPError(http.StatusNotFound, "benchmark not found")
	}

	validIntegrationTypes := make(map[string]bool)
	for _, it := range benchmark.IntegrationTypes {
		validIntegrationTypes[it] = true
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				h.Scheduler.logger.Error("failed to get integration", zap.Any("integration", info), zap.Error(err))
				return echo.NewHTTPError(http.StatusBadRequest, "failed to get integration")
			}
			if integration != nil {
				integrations = append(integrations, *integration)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			h.Scheduler.logger.Error("failed to get source", zap.Any("source", info), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to get source")
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}

	connectionInfo := make(map[string]api.IntegrationInfo)
	var connectionIDs []string
	for _, c := range integrations {
		if _, ok := validIntegrationTypes[c.IntegrationType.String()]; !ok {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid integration type for this framework")
		}
		connectionInfo[c.IntegrationID] = api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		}
		connectionIDs = append(connectionIDs, c.IntegrationID)
	}

	var apiJobs []api.RunBenchmarkItem
	if request.WithIncidents {
		lastJob, err := h.Scheduler.db.GetLastComplianceJob(true, benchmark.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		jobs, err := h.Scheduler.complianceScheduler.CreateComplianceReportJobs(true, benchmarkID, lastJob, connectionIDs, true, userID, nil)
		if err != nil {
			return fmt.Errorf("error while creating compliance job: %v", err)
		}
		for _, j := range jobs {
			job := api.RunBenchmarkItem{
				JobId:        j.ID,
				WithIncident: request.WithIncidents,
				BenchmarkId:  benchmark.ID,
			}
			for _, integration := range j.IntegrationIDs {
				if v, ok := connectionInfo[integration]; ok {
					job.IntegrationInfo = append(job.IntegrationInfo, v)
				}
			}
			apiJobs = append(apiJobs, job)
		}
	} else {
		jobs, err := h.Scheduler.complianceScheduler.CreateComplianceReportJobs(false, benchmarkID, nil, connectionIDs, true, userID, nil)
		if err != nil {
			return fmt.Errorf("error while creating compliance job: %v", err)
		}
		for _, j := range jobs {
			job := api.RunBenchmarkItem{
				JobId:        j.ID,
				WithIncident: request.WithIncidents,
				BenchmarkId:  benchmark.ID,
			}
			for _, integration := range j.IntegrationIDs {
				if v, ok := connectionInfo[integration]; ok {
					job.IntegrationInfo = append(job.IntegrationInfo, v)
				}
			}
			apiJobs = append(apiJobs, job)
		}
	}

	return ctx.JSON(http.StatusOK, api.RunBenchmarkResponse{
		Jobs: apiJobs,
	})
}

// RunBenchmark godoc
//
//	@Summary		Triggers compliance job
//	@Description	Triggers a compliance job to run immediately for the given benchmark
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Success		200		{object}	api.RunBenchmarkResponse
//	@Param			request	body		api.RunBenchmarkRequest	true	"Requst Body"
//	@Router			/schedule/api/v3/compliance/run [post]
func (h HttpServer) RunBenchmark(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}

	var request api.RunBenchmarkRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if len(request.IntegrationInfo) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please provide at least one connection info")
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			if integration != nil {
				integrations = append(integrations, *integration)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			h.Scheduler.logger.Error("failed to get source", zap.Any("source", info), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to get source")
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}

	connectionInfo := make(map[string]api.IntegrationInfo)
	var connectionIDs []string
	for _, c := range integrations {
		connectionIDs = append(connectionIDs, c.IntegrationID)
	}
	connections2, err := h.Scheduler.integrationClient.ListIntegrations(clientCtx, nil)
	if err != nil {
		h.Scheduler.logger.Error("failed to list connections", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	for _, c := range connections2.Integrations {
		connectionInfo[c.IntegrationID] = api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		}
	}

	var benchmarks []complianceapi.Benchmark
	if len(request.BenchmarkIds) == 0 {
		benchmarks, err = h.Scheduler.complianceClient.ListBenchmarks(clientCtx, nil)
		if err != nil {
			return fmt.Errorf("error while getting benchmarks: %v", err)
		}
	} else {
		for _, benchmarkID := range request.BenchmarkIds {
			benchmark, err := h.Scheduler.complianceClient.GetBenchmark(clientCtx, benchmarkID)
			if err != nil {
				return fmt.Errorf("error while getting benchmarks: %v", err)
			}
			if benchmark == nil {
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("benchmark %s not found", benchmarkID))
			}
			benchmarks = append(benchmarks, *benchmark)
		}
	}

	var jobs []api.RunBenchmarkItem
	for _, benchmark := range benchmarks {
		lastJob, err := h.Scheduler.db.GetLastComplianceJob(true, benchmark.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		benchmarkJobs, err := h.Scheduler.complianceScheduler.CreateComplianceReportJobs(true, benchmark.ID, lastJob, connectionIDs, true, userID, nil)
		if err != nil {
			return fmt.Errorf("error while creating compliance job: %v", err)
		}

		for _, j := range benchmarkJobs {
			job := api.RunBenchmarkItem{
				JobId:       j.ID,
				BenchmarkId: benchmark.ID,
			}
			for _, integration := range j.IntegrationIDs {
				if v, ok := connectionInfo[integration]; ok {
					job.IntegrationInfo = append(job.IntegrationInfo, v)
				}
			}
			jobs = append(jobs, job)
		}

	}

	return ctx.JSON(http.StatusOK, api.RunBenchmarkResponse{
		Jobs: jobs,
	})
}

// RunDiscovery godoc
//
//	@Summary		Run Discovery job
//	@Description	Triggers a discovery job to run immediately for the given resource types and Integrations
//	@Security		BearerToken
//	@Tags			describe
//	@Produce		json
//	@Success		200		{object}	api.RunDiscoveryResponse
//	@Param			request	body		api.RunDiscoveryRequest	true	"Request Body"
//	@Router			/schedule/api/v3/discovery/run [post]
func (h HttpServer) RunDiscovery(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole, Ctx: ctx.Request().Context()}
	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}

	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	triggerId, err := sf.NextID()
	if err != nil {
		return err
	}

	var request api.RunDiscoveryRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if len(request.IntegrationInfo) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please provide at least one connection info")
	}

	integrationTypes, err := h.Scheduler.integrationClient.ListIntegrationTypes(clientCtx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	integrationTypesMap := make(map[string]bool)
	for _, integrationType := range integrationTypes {
		integrationTypesMap[integrationType] = true
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				h.Scheduler.logger.Error("failed to get source", zap.String("source id", *info.IntegrationID), zap.Error(err))
				return echo.NewHTTPError(http.StatusBadRequest, "failed to get source")
			}
			if integration != nil {
				integrations = append(integrations, *integration)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			h.Scheduler.logger.Error("failed to get source", zap.Any("source", info), zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, "failed to get source")
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}

	var jobs []api.RunDiscoveryJob
	for _, integration := range integrations {
		if integration.State != integrationapi.IntegrationStateActive {
			continue
		}
		rtToDescribe := request.ResourceTypes

		var resourceTypesNames []string
		for _, rt := range rtToDescribe {
			resourceTypesNames = append(resourceTypesNames, rt.ResourceType)
		}

		discoveryType := model2.DiscoveryType_Fast
		integrationDiscovery := &model2.IntegrationDiscovery{
			TriggerID:     uint(triggerId),
			ConnectionID:  integration.IntegrationID,
			AccountID:     integration.ProviderID,
			TriggerType:   enums.DescribeTriggerTypeManual,
			TriggeredBy:   userID,
			DiscoveryType: discoveryType,
			ResourceTypes: resourceTypesNames,
		}
		err = h.DB.CreateIntegrationDiscovery(integrationDiscovery)
		if err != nil {
			h.Scheduler.logger.Error("failed to create integration discovery", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to create integration discovery")
		}

		_, ok := integrationTypesMap[integration.IntegrationType.String()]
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "unknown integration type")
		}

		possibleRtMap, err := h.Scheduler.integrationClient.GetResourceTypesByLabels(clientCtx, integration.IntegrationType.String(), integration.Labels)
		if err != nil {
			h.Scheduler.logger.Error("failed to get resource types by labels", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get resource types by labels")
		}
		var possibleRt []api.ResourceTypeRunDiscoveryRequest
		for rt, _ := range possibleRtMap {
			possibleRt = append(possibleRt, api.ResourceTypeRunDiscoveryRequest{
				ResourceType: rt,
			})
		}
		if len(rtToDescribe) == 0 {
			rtToDescribe = possibleRt
		}

		for _, resourceType := range rtToDescribe {
			isOK := false
			for _, rt := range possibleRt {
				if rt.ResourceType == resourceType.ResourceType {
					isOK = true
				}
			}
			if !isOK {
				continue
			}
			var status, failureReason string
			job, err := h.Scheduler.describe(integration, resourceType.ResourceType, false, false, false, &integrationDiscovery.ID, userID, resourceType.Parameters)
			if err != nil {
				if err.Error() == "job already in progress" {
					tmpJob, err := h.Scheduler.db.GetLastDescribeIntegrationJob(integration.IntegrationID, resourceType.ResourceType)
					if err != nil {
						h.Scheduler.logger.Error("failed to get last describe job", zap.String("resource_type", resourceType.ResourceType), zap.String("connection_id", integration.IntegrationID), zap.Error(err))
					}
					h.Scheduler.logger.Error("failed to describe connection", zap.String("integration_id", integration.IntegrationID), zap.Error(err))
					status = "FAILED"
					failureReason = fmt.Sprintf("job already in progress: %v", tmpJob.ID)
				} else {
					failureReason = err.Error()
				}
			}

			var jobId uint
			if job == nil {
				status = "FAILED"
				if failureReason == "" && err != nil {
					failureReason = err.Error()
				}
			} else {
				jobId = job.ID
				status = string(job.Status)
			}
			jobs = append(jobs, api.RunDiscoveryJob{
				JobId:         jobId,
				ResourceType:  resourceType.ResourceType,
				Status:        status,
				FailureReason: failureReason,
				IntegrationInfo: api.IntegrationInfo{
					IntegrationID:   integration.IntegrationID,
					IntegrationType: string(integration.IntegrationType),
					ProviderID:      integration.ProviderID,
					Name:            integration.Name,
				},
			})
		}
	}
	return ctx.JSON(http.StatusOK, api.RunDiscoveryResponse{
		Jobs:      jobs,
		TriggerID: uint(triggerId),
	})
}

// GetDescribeJobStatus godoc
//
//	@Summary	Get describe job status by job id
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		job_id	path	string	true	"Job ID"
//	@Produce	json
//	@Success	200	{object}	api.GetDescribeJobStatusResponse
//	@Router		/schedule/api/v3/jobs/discovery/{job_id} [get]
func (h HttpServer) GetDescribeJobStatus(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	jobId := ctx.Param("job_id")

	j, err := h.DB.GetDescribeJobById(jobId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	connection, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, j.IntegrationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	jobsResult := api.GetDescribeJobStatusResponse{
		JobId: j.ID,
		IntegrationInfo: api.IntegrationInfo{
			IntegrationID:   connection.IntegrationID,
			IntegrationType: string(connection.IntegrationType),
			ProviderID:      connection.ProviderID,
			Name:            connection.Name,
		},
		ResourceType: j.ResourceType,
		JobStatus:    string(j.Status),
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, jobsResult)
}

// GetComplianceJobStatus godoc
//
//	@Summary	Get compliance job status by job id
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		job_id	path	string	true	"Job ID"
//	@Produce	json
//	@Success	200	{object}	api.GetComplianceJobStatusResponse
//	@Router		/schedule/api/v3/job/compliance/{job_id} [get]
func (h HttpServer) GetComplianceJobStatus(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	jobIdString := ctx.Param("job_id")
	jobId, err := strconv.ParseUint(jobIdString, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job id")
	}

	j, err := h.DB.GetComplianceJobByID(uint(jobId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var summaryJobId *uint
	summaryJobs, err := h.DB.ListSummaryJobs([]string{jobIdString})
	if err != nil {
		h.Scheduler.logger.Error("failed to get summary job", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get summary job")
	}
	if len(summaryJobs) > 0 {
		summaryJobId = &summaryJobs[0].ID
	}

	var integrations []api.IntegrationInfo
	for _, i := range j.IntegrationIDs {
		integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, i)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrationInfo := api.IntegrationInfo{
			IntegrationID:   integration.IntegrationID,
			IntegrationType: string(integration.IntegrationType),
			ProviderID:      integration.ProviderID,
			Name:            integration.Name,
		}
		integrations = append(integrations, integrationInfo)
	}

	jobsResult := api.GetComplianceJobStatusResponse{
		JobId:           j.ID,
		WithIncidents:   j.WithIncidents,
		SummaryJobId:    summaryJobId,
		IntegrationInfo: integrations,
		FrameworkId:     j.FrameworkID,
		JobStatus:       string(j.Status),
		CreatedAt:       j.CreatedAt,
		UpdatedAt:       j.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, jobsResult)
}

// GetAsyncQueryRunJobStatus godoc
//
//	@Summary	Get async query run job status by job id
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		job_id	path	string	true	"Job ID"
//	@Produce	json
//	@Success	200	{object}	api.GetAsyncQueryRunJobStatusResponse
//	@Router		/schedule/api/v3/job/query/{job_id} [get]
func (h HttpServer) GetAsyncQueryRunJobStatus(ctx echo.Context) error {

	jobIdString := ctx.Param("job_id")
	jobId, err := strconv.ParseUint(jobIdString, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job id")
	}

	j, err := h.DB.GetQueryRunnerJob(uint(jobId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	jobsResult := api.GetAsyncQueryRunJobStatusResponse{
		JobId:          j.ID,
		QueryId:        j.QueryId,
		CreatedAt:      j.CreatedAt,
		UpdatedAt:      j.UpdatedAt,
		CreatedBy:      j.CreatedBy,
		JobStatus:      j.Status,
		FailureMessage: j.FailureMessage,
	}

	return ctx.JSON(http.StatusOK, jobsResult)
}

// ListDescribeJobs godoc
//
//	@Summary	Get describe jobs history for give connection
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.ListDescribeJobsRequest	true	"List jobs request"
//	@Produce	json
//	@Success	200	{object}	[]api.GetDescribeJobsHistoryResponse
//	@Router		/schedule/api/v3/jobs/discovery [post]
func (h HttpServer) ListDescribeJobs(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.ListDescribeJobsRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			if integration != nil {
				integrations = append(integrations, *integration)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}

	connectionInfo := make(map[string]api.IntegrationInfo)
	var connectionIDs []string
	for _, c := range integrations {
		connectionInfo[c.IntegrationID] = api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		}
		connectionIDs = append(connectionIDs, c.IntegrationID)
	}

	var jobsResults []api.GetDescribeJobsHistoryResponse
	var startTime, endTime *time.Time
	if request.Interval != nil {
		startTime, endTime, _ = parseTimeInterval(*request.Interval)
	} else {
		startTime = &request.StartTime
		endTime = request.EndTime
	}

	jobs, err := h.DB.ListDescribeJobsByFilters(nil, connectionIDs, request.ResourceType,
		request.DiscoveryType, request.JobStatus, startTime, endTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, j := range jobs {
		jobResult := api.GetDescribeJobsHistoryResponse{
			JobId:          j.ID,
			ResourceType:   j.ResourceType,
			JobStatus:      j.Status,
			UpdatedAt:      j.UpdatedAt,
			CreatedAt:      j.CreatedAt,
			FailureMessage: j.FailureMessage,
			Title:          j.ResourceType,
		}
		if info, ok := connectionInfo[j.IntegrationID]; ok {
			jobResult.IntegrationInfo = &info
		}
		jobsResults = append(jobsResults, jobResult)
	}
	if request.SortBy != nil {
		switch strings.ToLower(*request.SortBy) {
		case "id":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		case "datetime", "updated_at", "updatedat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].UpdatedAt.After(jobsResults[j].UpdatedAt)
			})
		case "created_at", "createdat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].CreatedAt.After(jobsResults[j].CreatedAt)
			})
		case "resourcetype":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].ResourceType < jobsResults[j].ResourceType
			})
		case "jobstatus":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobStatus < jobsResults[j].JobStatus
			})
		default:
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		}
	} else {
		sort.Slice(jobsResults, func(i, j int) bool {
			return jobsResults[i].JobId < jobsResults[j].JobId
		})
	}
	TotalCount := len(jobsResults)
	if request.PerPage != nil {
		if request.Cursor == nil {
			jobsResults = utils.Paginate(1, *request.PerPage, jobsResults)
		} else {
			jobsResults = utils.Paginate(*request.Cursor, *request.PerPage, jobsResults)
		}
	}

	return ctx.JSON(http.StatusOK, api.GetDescribeJobsHistoryFinalResponse{
		TotalCount: TotalCount,
		Items:      jobsResults,
	})
}

// ListComplianceJobs godoc
//
//	@Summary	Get compliance jobs history for give connection
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.ListComplianceJobsRequest	true	"List jobs request"
//	@Produce	json
//	@Success	200	{object}	[]api.GetComplianceJobsHistoryResponse
//	@Router		/schedule/api/v3/jobs/compliance [post]
func (h HttpServer) ListComplianceJobs(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.ListComplianceJobsRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			if integration != nil {
				integrations = append(integrations, *integration)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}

	connectionInfo := make(map[string]api.IntegrationInfo)
	var connectionIDs []string
	for _, c := range integrations {
		connectionInfo[c.IntegrationID] = api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		}
		connectionIDs = append(connectionIDs, c.IntegrationID)
	}
	var startTime, endTime *time.Time
	if request.Interval != nil {
		startTime, endTime, _ = parseTimeInterval(*request.Interval)

	} else {
		startTime = &request.StartTime
		endTime = request.EndTime
	}

	jobs, err := h.DB.ListComplianceJobsByFilters(nil, connectionIDs, request.BenchmarkId, request.JobStatus, startTime, endTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	benchmarks, err := h.Scheduler.complianceClient.ListBenchmarks(&httpclient.Context{UserRole: apiAuth.AdminRole}, nil)
	if err != nil {
		return err
	}

	var jobsResults []api.GetComplianceJobsHistoryResponse
	for _, j := range jobs {
		Title := ""
		for _, benchmark := range benchmarks {
			if fmt.Sprintf("%v", benchmark.ID) == j.FrameworkID {
				Title = benchmark.Title
			}
		}
		jobResult := api.GetComplianceJobsHistoryResponse{
			JobId:          j.ID,
			WithIncidents:  j.WithIncidents,
			BenchmarkId:    j.FrameworkID,
			JobStatus:      j.Status.ToApi(),
			UpdatedAt:      j.UpdatedAt,
			CreatedAt:      j.CreatedAt,
			FailureMessage: j.FailureMessage,
			Title:          Title,
		}
		for _, i := range j.IntegrationIDs {
			if info, ok := connectionInfo[i]; ok {
				jobResult.IntegrationInfo = []api.IntegrationInfo{info}
			}
		}

		jobsResults = append(jobsResults, jobResult)
	}
	if request.SortBy != nil {
		switch strings.ToLower(*request.SortBy) {
		case "id":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		case "datetime", "updated_at", "updatedat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].UpdatedAt.After(jobsResults[j].UpdatedAt)
			})
		case "created_at", "createdat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].CreatedAt.After(jobsResults[j].CreatedAt)
			})
		case "benchmarkid", "benchmark_id":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].BenchmarkId < jobsResults[j].BenchmarkId
			})
		case "jobstatus", "job_status":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobStatus < jobsResults[j].JobStatus
			})
		case "with_incidents", "withincidents":
			sort.Slice(jobsResults, func(i, j int) bool {
				return !jobsResults[i].WithIncidents && jobsResults[j].WithIncidents
			})
		default:
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		}
	} else {
		sort.Slice(jobsResults, func(i, j int) bool {
			return jobsResults[i].JobId < jobsResults[j].JobId
		})
	}
	totalCount := len(jobsResults)
	if request.PerPage != nil {
		if request.Cursor == nil {
			jobsResults = utils.Paginate(1, *request.PerPage, jobsResults)
		} else {
			jobsResults = utils.Paginate(*request.Cursor, *request.PerPage, jobsResults)
		}
	}

	return ctx.JSON(http.StatusOK, api.GetComplianceJobsHistoryFinalResponse{
		TotalCount: totalCount,
		Items:      jobsResults,
	})
}

// BenchmarkAuditHistory godoc
//
//	@Summary	Get compliance jobs history for give connection
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request			body	api.BenchmarkAuditHistoryRequest	true	"List jobs request"
//	@Param		benchmark_id	query	string								true	"Benchmark ID to get the run history for"
//	@Produce	json
//	@Success	200	{object}	api.BenchmarkAuditHistoryResponse
//	@Router		/schedule/api/v3/benchmark/:benchmark_id/run-history [post]
func (h HttpServer) BenchmarkAuditHistory(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	benchmarkID := ctx.Param("benchmark_id")

	var request api.BenchmarkAuditHistoryRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			connection, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			if connection != nil {
				integrations = append(integrations, *connection)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}

	integrationInfo := make(map[string]api.IntegrationInfo)
	var connectionIDs []string
	for _, c := range integrations {
		connectionIDs = append(connectionIDs, c.IntegrationID)
	}
	connections2, err := h.Scheduler.integrationClient.ListIntegrations(clientCtx, nil)
	if err != nil {
		h.Scheduler.logger.Error("failed to list connections", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	for _, c := range connections2.Integrations {
		integrationInfo[c.IntegrationID] = api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		}
	}

	var startTime, endTime *time.Time
	if request.Interval != nil {
		startTime, endTime, err = parseTimeInterval(*request.Interval)
	} else {
		startTime = &request.StartTime
		endTime = request.EndTime
	}
	var items []api.BenchmarkAuditHistoryItem

	// With Incidents
	jobs, err := h.DB.ListComplianceJobsByFilters(request.WithIncidents, connectionIDs, []string{benchmarkID}, request.JobStatus, startTime, endTime)
	if err != nil {
		h.Scheduler.logger.Error("failed to get list of compliance jobs", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get list of compliance jobs")
	}

	for _, j := range jobs {
		item := api.BenchmarkAuditHistoryItem{
			JobId:         j.ID,
			WithIncidents: j.WithIncidents,
			BenchmarkId:   j.FrameworkID,
			JobStatus:     j.Status.ToApi(),
			CreatedAt:     j.CreatedAt,
			UpdatedAt:     j.UpdatedAt,
		}
		for _, i := range j.IntegrationIDs {
			if info, ok := integrationInfo[i]; ok {
				item.IntegrationInfo = append(item.IntegrationInfo, info)
				item.NumberOfIntegrations += 1
			}
		}

		items = append(items, item)
	}

	if request.SortBy != nil {
		switch strings.ToLower(*request.SortBy) {
		case "id":
			sort.Slice(items, func(i, j int) bool {
				return items[i].JobId < items[j].JobId
			})
		case "updated_at", "updatedat":
			sort.Slice(items, func(i, j int) bool {
				return items[i].UpdatedAt.After(items[j].UpdatedAt)
			})
		case "created_at", "createdat":
			sort.Slice(items, func(i, j int) bool {
				return items[i].CreatedAt.After(items[j].CreatedAt)
			})
		case "benchmarkid":
			sort.Slice(items, func(i, j int) bool {
				return items[i].BenchmarkId < items[j].BenchmarkId
			})
		case "jobstatus":
			sort.Slice(items, func(i, j int) bool {
				return items[i].JobStatus < items[j].JobStatus
			})
		case "with_incidents", "withincidents":
			sort.Slice(items, func(i, j int) bool {
				return !items[i].WithIncidents && items[j].WithIncidents
			})
		default:
			sort.Slice(items, func(i, j int) bool {
				return items[i].UpdatedAt.After(items[j].UpdatedAt)
			})
		}
	} else {
		sort.Slice(items, func(i, j int) bool {
			return items[i].UpdatedAt.After(items[j].UpdatedAt)
		})
	}

	totalCount := len(items)
	if request.PerPage != nil {
		if request.Cursor == nil {
			items = utils.Paginate(1, *request.PerPage, items)
		} else {
			items = utils.Paginate(*request.Cursor, *request.PerPage, items)
		}
	}

	return ctx.JSON(http.StatusOK, api.BenchmarkAuditHistoryResponse{
		Items:      items,
		TotalCount: totalCount,
	})
}

// BenchmarkAuditHistoryIntegrations godoc
//
//	@Summary	Get compliance jobs history for give connection possible integrations
//	@Security	BearerToken
//	@Tags		scheduler
//	@Produce	json
//	@Success	200	{object}	[]api.IntegrationInfo
//	@Router		/schedule/api/v3/benchmark/run-history/integrations [get]
func (h HttpServer) BenchmarkAuditHistoryIntegrations(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	connectionIDs, err := h.DB.GetComplianceJobsIntegrations()
	if err != nil {
		h.Scheduler.logger.Error("failed to get compliance jobs integrations", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get compliance jobs integrations")
	}

	connections, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx, integrationapi.ListIntegrationsRequest{
		IntegrationID: connectionIDs,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var integrations []api.IntegrationInfo

	for _, c := range connections.Integrations {
		integrations = append(integrations, api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		})
	}

	return ctx.JSON(http.StatusOK, integrations)
}

// GetIntegrationLastDiscoveryJob godoc
//
//	@Summary	Get Last dicovery job for integration
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.GetIntegrationLastDiscoveryJobRequest	true	"List jobs request"
//	@Produce	json
//	@Success	200
//	@Router		/schedule/api/v3/integration/discovery/last-job [post]
func (h HttpServer) GetIntegrationLastDiscoveryJob(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.GetIntegrationLastDiscoveryJobRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var integration integrationapi.Integration
	if request.IntegrationInfo.IntegrationID != nil {
		connectionTmp, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *request.IntegrationInfo.IntegrationID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if connectionTmp != nil {
			integration = *connectionTmp
		}
	} else {
		var integrationTypes []string
		if request.IntegrationInfo.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *request.IntegrationInfo.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       request.IntegrationInfo.Name,
				ProviderIDRegex: request.IntegrationInfo.ProviderID,
			})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integration = connectionsTmp.Integrations[0]
	}

	job, err := h.DB.ListDescribeJobs(integration.IntegrationID)
	if err != nil {
		h.Scheduler.logger.Error("failed to get describe jobs", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get describe jobs")
	}

	return ctx.JSON(http.StatusOK, job)
}

// GetDescribeJobsHistoryByIntegration godoc
//
//	@Summary	Get describe jobs history for give connection
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.GetDescribeJobsHistoryByIntegrationRequest	true	"List jobs request"
//	@Produce	json
//	@Success	200	{object}	[]api.GetDescribeJobsHistoryResponse
//	@Router		/schedule/api/v3/jobs/discovery/connections [post]
func (h HttpServer) GetDescribeJobsHistoryByIntegration(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.GetDescribeJobsHistoryByIntegrationRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			if integration != nil {
				integrations = append(integrations, *integration)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}

	connectionInfo := make(map[string]api.IntegrationInfo)
	for _, c := range integrations {
		connectionInfo[c.IntegrationID] = api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		}
	}

	var jobsResults []api.GetDescribeJobsHistoryResponse

	for _, c := range connectionInfo {
		jobs, err := h.DB.ListDescribeJobsByFilters(nil, []string{c.IntegrationID}, request.ResourceType,
			request.DiscoveryType, request.JobStatus, &request.StartTime, request.EndTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		for _, j := range jobs {
			jobsResults = append(jobsResults, api.GetDescribeJobsHistoryResponse{
				JobId:           j.ID,
				ResourceType:    j.ResourceType,
				JobStatus:       j.Status,
				UpdatedAt:       j.UpdatedAt,
				CreatedAt:       j.CreatedAt,
				FailureMessage:  j.FailureMessage,
				Title:           j.ResourceType,
				IntegrationInfo: &c,
			})
		}
	}

	if request.SortBy != nil {
		switch strings.ToLower(*request.SortBy) {
		case "id":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		case "datetime", "updated_at", "updatedat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].UpdatedAt.After(jobsResults[j].UpdatedAt)
			})
		case "created_at", "createdat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].CreatedAt.After(jobsResults[j].CreatedAt)
			})
		case "resourcetype":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].ResourceType < jobsResults[j].ResourceType
			})
		case "jobstatus":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobStatus < jobsResults[j].JobStatus
			})
		default:
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		}
	} else {
		sort.Slice(jobsResults, func(i, j int) bool {
			return jobsResults[i].JobId < jobsResults[j].JobId
		})
	}
	if request.PerPage != nil {
		if request.Cursor == nil {
			jobsResults = utils.Paginate(1, *request.PerPage, jobsResults)
		} else {
			jobsResults = utils.Paginate(*request.Cursor, *request.PerPage, jobsResults)
		}
	}

	return ctx.JSON(http.StatusOK, jobsResults)
}

// GetComplianceJobsHistoryByIntegration godoc
//
//	@Summary	Get compliance jobs history for give connection
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.GetComplianceJobsHistoryByIntegrationRequest	true	"List jobs request"
//	@Produce	json
//	@Success	200	{object}	[]api.GetComplianceJobsHistoryResponse
//	@Router		/schedule/api/v3/jobs/compliance/connections [post]
func (h HttpServer) GetComplianceJobsHistoryByIntegration(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.GetComplianceJobsHistoryByIntegrationRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var integrations integrationapi.Integration
	if request.IntegrationInfo.IntegrationID != nil {
		connection, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *request.IntegrationInfo.IntegrationID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if connection != nil {
			integrations = *connection
		}
	} else {
		var integrationTypes []string
		if request.IntegrationInfo.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *request.IntegrationInfo.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       request.IntegrationInfo.Name,
				ProviderIDRegex: request.IntegrationInfo.ProviderID,
			})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrations = connectionsTmp.Integrations[0]
	}

	connectionInfo := make(map[string]api.IntegrationInfo)
	connectionInfo[integrations.IntegrationID] = api.IntegrationInfo{
		IntegrationID:   integrations.IntegrationID,
		IntegrationType: string(integrations.IntegrationType),
		Name:            integrations.Name,
		ProviderID:      integrations.ProviderID,
	}

	var jobsResults []api.GetComplianceJobsHistoryResponse
	for _, c := range connectionInfo {
		jobs, err := h.DB.ListComplianceJobsByFilters(request.WithIncidents, []string{c.IntegrationID}, request.BenchmarkId, request.JobStatus, &request.StartTime, request.EndTime)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		for _, j := range jobs {
			var jobIntegrations []api.IntegrationInfo
			for _, i := range j.IntegrationIDs {
				if info, ok := connectionInfo[i]; ok {
					jobIntegrations = append(jobIntegrations, info)
				} else {
					integration, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, i)
					if err != nil {
						return echo.NewHTTPError(http.StatusBadRequest, err.Error())
					}
					if integration != nil {
						info = api.IntegrationInfo{
							IntegrationID:   integration.IntegrationID,
							IntegrationType: string(integration.IntegrationType),
							Name:            integration.Name,
							ProviderID:      integration.ProviderID,
						}
						connectionInfo[i] = info
						jobIntegrations = append(jobIntegrations, info)
					}
				}
			}

			jobsResults = append(jobsResults, api.GetComplianceJobsHistoryResponse{
				JobId:           j.ID,
				WithIncidents:   j.WithIncidents,
				BenchmarkId:     j.FrameworkID,
				JobStatus:       j.Status.ToApi(),
				UpdatedAt:       j.UpdatedAt,
				CreatedAt:       j.CreatedAt,
				IntegrationInfo: jobIntegrations,
			})
		}
	}

	if request.SortBy != nil {
		switch strings.ToLower(*request.SortBy) {
		case "id":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		case "datetime", "updated_at", "updatedat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].UpdatedAt.After(jobsResults[j].UpdatedAt)
			})
		case "created_at", "createdat":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].CreatedAt.After(jobsResults[j].CreatedAt)
			})
		case "benchmarkid":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].BenchmarkId < jobsResults[j].BenchmarkId
			})
		case "jobstatus":
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobStatus < jobsResults[j].JobStatus
			})
		case "with_incidents", "withincidents":
			sort.Slice(jobsResults, func(i, j int) bool {
				return !jobsResults[i].WithIncidents && jobsResults[j].WithIncidents
			})
		default:
			sort.Slice(jobsResults, func(i, j int) bool {
				return jobsResults[i].JobId < jobsResults[j].JobId
			})
		}
	} else {
		sort.Slice(jobsResults, func(i, j int) bool {
			return jobsResults[i].JobId < jobsResults[j].JobId
		})
	}
	if request.PerPage != nil {
		if request.Cursor == nil {
			jobsResults = utils.Paginate(1, *request.PerPage, jobsResults)
		} else {
			jobsResults = utils.Paginate(*request.Cursor, *request.PerPage, jobsResults)
		}
	}

	return ctx.JSON(http.StatusOK, jobsResults)
}

// CancelJobById godoc
//
//	@Summary	Cancel job by given job type and job ID
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		job_id		query	string	true	"Job ID to Cancel"
//	@Param		job_type	query	string	true	"Job Type"
//	@Produce	json
//	@Success	200
//	@Router		/schedule/api/v3/jobs/cancel/byid [put]
func (h HttpServer) CancelJobById(ctx echo.Context) error {
	jobIdStr := ctx.QueryParam("job_id")
	jobType := strings.ToLower(ctx.QueryParam("job_type"))

	switch jobType {
	case "compliance":
		jobId, err := strconv.ParseUint(jobIdStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid job id")
		}
		complianceJob, err := h.DB.GetComplianceJobByID(uint(jobId))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if complianceJob == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "job not found")
		}
		if complianceJob.Status == model2.ComplianceJobCreated {
			err = h.DB.UpdateComplianceJob(uint(jobId), model2.ComplianceJobCanceled, "")
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return ctx.NoContent(http.StatusOK)
		} else if complianceJob.Status == model2.ComplianceJobSucceeded || complianceJob.Status == model2.ComplianceJobFailed ||
			complianceJob.Status == model2.ComplianceJobTimeOut || complianceJob.Status == model2.ComplianceJobCanceled {
			return echo.NewHTTPError(http.StatusOK, "job is already finished")
		} else if complianceJob.Status == model2.ComplianceJobSummarizerInProgress || complianceJob.Status == model2.ComplianceJobSinkInProgress {
			return echo.NewHTTPError(http.StatusOK, "job is already in progress, unable to cancel")
		}
		runners, err := h.DB.ListComplianceJobRunnersWithID(uint(jobId))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if len(runners) == 0 {
			err = h.DB.UpdateComplianceJob(uint(jobId), model2.ComplianceJobCanceled, "")
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return ctx.NoContent(http.StatusOK)
		} else {
			allInProgress := true
			for _, r := range runners {
				if r.Status == runner2.ComplianceRunnerCreated {
					allInProgress = false
					err = h.DB.UpdateRunnerJob(r.ID, runner2.ComplianceRunnerCanceled, r.StartedAt, nil, "")
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
				} else if r.Status == runner2.ComplianceRunnerQueued {
					allInProgress = false
					err = h.Scheduler.jq.DeleteMessage(ctx.Request().Context(), runner2.StreamName, r.NatsSequenceNumber)
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
					err = h.DB.UpdateRunnerJob(r.ID, runner2.ComplianceRunnerCanceled, r.StartedAt, nil, "")
					if err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
					}
				}
			}
			if allInProgress {
				return echo.NewHTTPError(http.StatusOK, "job is already in progress, unable to cancel")
			} else {
				err = h.DB.UpdateComplianceJob(uint(jobId), model2.ComplianceJobCanceled, "")
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				return ctx.NoContent(http.StatusOK)
			}
		}
	case "discovery":
		job, err := h.DB.GetDescribeJobById(jobIdStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid job id")
		}
		if job == nil {
			return echo.NewHTTPError(http.StatusBadRequest, "job not found")
		}
		if job.Status == api.DescribeResourceJobCreated {
			err = h.DB.UpdateDescribeIntegrationJobStatus(job.ID, api.DescribeResourceJobCanceled, "", "", 0, 0)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return ctx.NoContent(http.StatusOK)
		} else if job.Status == api.DescribeResourceJobCanceled || job.Status == api.DescribeResourceJobFailed ||
			job.Status == api.DescribeResourceJobSucceeded || job.Status == api.DescribeResourceJobTimeout {
			return echo.NewHTTPError(http.StatusOK, "job is already finished")
		} else if job.Status == api.DescribeResourceJobQueued {
			ctx2 := &httpclient.Context{UserRole: apiAuth.AdminRole, Ctx: ctx.Request().Context()}

			integrationTypes, err := h.Scheduler.integrationClient.ListIntegrationTypes(ctx2)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			integrationTypesMap := make(map[string]bool)
			for _, integrationType := range integrationTypes {
				integrationTypesMap[integrationType] = true
			}

			_, ok := integrationTypesMap[job.IntegrationType.String()]
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "unknown integration type")
			}

			integrationConfig, err := h.Scheduler.integrationClient.GetIntegrationConfiguration(ctx2, job.IntegrationType.String())
			err = h.Scheduler.jq.DeleteMessage(ctx.Request().Context(), integrationConfig.NatsStreamName, job.NatsSequenceNumber)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			err = h.DB.UpdateDescribeIntegrationJobStatus(job.ID, api.DescribeResourceJobCanceled, "", "", 0, 0)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return ctx.NoContent(http.StatusOK)
		} else {
			return echo.NewHTTPError(http.StatusOK, "job is already in progress, unable to cancel")
		}

	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job type")
	}
	return echo.NewHTTPError(http.StatusOK, "nothing done")
}

// CancelJob godoc
//
//	@Summary	Cancel job by given job type and job ID
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.CancelJobRequest	true	"Request Body"
//	@Produce	json
//	@Success	200	{object}	[]api.CancelJobResponse
//	@Router		/schedule/api/v3/jobs/cancel [post]
func (h HttpServer) CancelJob(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.CancelJobRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if strings.ToLower(request.JobType) != "compliance" && strings.ToLower(request.JobType) != "discovery" &&
		strings.ToLower(request.JobType) != "query" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job type")
	}

	if strings.ToLower(request.JobType) == "query" && strings.ToLower(request.Selector) != "job_id" {
		return echo.NewHTTPError(http.StatusBadRequest, "only jobId is acceptable for query run")
	}

	var jobIDs []string
	var results []api.CancelJobResponse

	switch strings.ToLower(request.Selector) {
	case "job_id":
		jobIDs = request.JobId
	case "integration_info":
		var integrations []integrationapi.Integration
		for _, info := range request.IntegrationInfo {
			if info.IntegrationID != nil {
				connection, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
				if connection != nil {
					integrations = append(integrations, *connection)
				}
				continue
			}
			var integrationTypes []string
			if info.IntegrationType != nil {
				integrationTypes = append(integrationTypes, *info.IntegrationType)
			}
			connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
				integrationapi.ListIntegrationsRequest{
					IntegrationType: integrationTypes,
					NameRegex:       info.Name,
					ProviderIDRegex: info.ProviderID,
				})
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			integrations = append(integrations, connectionsTmp.Integrations[0])

		}

		IntegrationIDsMap := make(map[string]bool)
		for _, c := range integrations {
			IntegrationIDsMap[c.IntegrationID] = true
		}
		var integrationIDs []string
		switch strings.ToLower(request.JobType) {
		case "compliance":
			jobs, err := h.DB.ListPendingComplianceJobsByIntegrationID(aws.Bool(true), integrationIDs)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			for _, j := range jobs {
				jobIDs = append(jobIDs, strconv.Itoa(int(j.ID)))
			}
		case "discovery":
			jobs, err := h.DB.ListPendingDescribeJobsByFilters(integrationIDs, nil, nil, nil, nil)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			for _, j := range jobs {
				jobIDs = append(jobIDs, strconv.Itoa(int(j.ID)))
			}

		}
	case "status":
		for _, status := range request.Status {
			switch strings.ToLower(request.JobType) {
			case "compliance":
				if model2.ComplianceJobStatus(strings.ToUpper(status)) != model2.ComplianceJobCreated &&
					model2.ComplianceJobStatus(strings.ToUpper(status)) != model2.ComplianceJobRunnersInProgress {
					continue
				}
				jobs, err := h.DB.ListComplianceJobsByStatus(aws.Bool(true), model2.ComplianceJobStatus(strings.ToUpper(status)))
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				for _, j := range jobs {
					jobIDs = append(jobIDs, strconv.Itoa(int(j.ID)))
				}
			case "discovery":
				if api.DescribeResourceJobStatus(strings.ToUpper(status)) != api.DescribeResourceJobCreated &&
					api.DescribeResourceJobStatus(strings.ToUpper(status)) != api.DescribeResourceJobQueued {
					continue
				}
				jobs, err := h.DB.ListDescribeJobsByStatus(api.DescribeResourceJobStatus(strings.ToUpper(status)))
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				for _, j := range jobs {
					jobIDs = append(jobIDs, strconv.Itoa(int(j.ID)))
				}

			}
		}
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid selector. valid values: job_id, integration, status")
	}

	for _, jobIdStr := range jobIDs {
		var failureReason string
		var canceled bool
		switch strings.ToLower(request.JobType) {
		case "compliance":
			jobId, err := strconv.ParseUint(jobIdStr, 10, 64)
			if err != nil {
				failureReason = "invalid job id"
				break
			}
			complianceJob, err := h.DB.GetComplianceJobByID(uint(jobId))
			if err != nil {
				failureReason = err.Error()
				break
			}
			if complianceJob == nil {
				failureReason = "job not found"
				break
			}
			if complianceJob.Status == model2.ComplianceJobCreated {
				err = h.DB.UpdateComplianceJob(uint(jobId), model2.ComplianceJobCanceled, "")
				if err != nil {
					failureReason = err.Error()
					break
				}
				canceled = true
				break
			} else if complianceJob.Status == model2.ComplianceJobSucceeded || complianceJob.Status == model2.ComplianceJobFailed ||
				complianceJob.Status == model2.ComplianceJobTimeOut || complianceJob.Status == model2.ComplianceJobCanceled {
				failureReason = "job is already finished"
				break
			} else if complianceJob.Status == model2.ComplianceJobSummarizerInProgress || complianceJob.Status == model2.ComplianceJobSinkInProgress {
				failureReason = "job is already in progress, unable to cancel"
				break
			}
			runners, err := h.DB.ListComplianceJobRunnersWithID(uint(jobId))
			if err != nil {
				failureReason = err.Error()
				break
			}
			if len(runners) == 0 {
				err = h.DB.UpdateComplianceJob(uint(jobId), model2.ComplianceJobCanceled, "")
				if err != nil {
					failureReason = err.Error()
					break
				}
				canceled = true
				break
			} else {
				allInProgress := true
				for _, r := range runners {
					if r.Status == runner2.ComplianceRunnerCreated {
						allInProgress = false
						err = h.DB.UpdateRunnerJob(r.ID, runner2.ComplianceRunnerCanceled, r.StartedAt, nil, "")
						if err != nil {
							failureReason = err.Error()
							break
						}
					} else if r.Status == runner2.ComplianceRunnerQueued {
						allInProgress = false
						err = h.Scheduler.jq.DeleteMessage(ctx.Request().Context(), runner2.StreamName, r.NatsSequenceNumber)
						if err != nil {
							failureReason = err.Error()
							break
						}
						err = h.DB.UpdateRunnerJob(r.ID, runner2.ComplianceRunnerCanceled, r.StartedAt, nil, "")
						if err != nil {
							failureReason = err.Error()
							break
						}
					}
				}
				if allInProgress {
					failureReason = "job is already in progress, unable to cancel"
					break
				} else {
					err = h.DB.UpdateComplianceJob(uint(jobId), model2.ComplianceJobCanceled, "")
					if err != nil {
						failureReason = err.Error()
						break
					}
					canceled = true
					break
				}
			}
		case "discovery":
			job, err := h.DB.GetDescribeJobById(jobIdStr)
			if err != nil {
				failureReason = "invalid job id"
				break
			}
			if job == nil {
				failureReason = "job not found"
				break
			}
			if job.Status == api.DescribeResourceJobCreated {
				err = h.DB.UpdateDescribeIntegrationJobStatus(job.ID, api.DescribeResourceJobCanceled, "", "", 0, 0)
				if err != nil {
					failureReason = err.Error()
					break
				}
				canceled = true
				break
			} else if job.Status == api.DescribeResourceJobCanceled || job.Status == api.DescribeResourceJobFailed ||
				job.Status == api.DescribeResourceJobSucceeded || job.Status == api.DescribeResourceJobTimeout {
				failureReason = "job is already finished"
				break
			} else if job.Status == api.DescribeResourceJobQueued {
				integrationConfig, err := h.Scheduler.integrationClient.GetIntegrationConfiguration(clientCtx, job.IntegrationType.String())
				err = h.Scheduler.jq.DeleteMessage(ctx.Request().Context(), integrationConfig.NatsStreamName, job.NatsSequenceNumber)
				if err != nil {
					failureReason = err.Error()
					break
				}
				err = h.DB.UpdateDescribeIntegrationJobStatus(job.ID, api.DescribeResourceJobCanceled, "", "", 0, 0)
				if err != nil {
					failureReason = err.Error()
					break
				}
				canceled = true
				break
			} else {
				failureReason = "job is already in progress, unable to cancel"
				break
			}

		case "query":
			jobId, err := strconv.ParseUint(jobIdStr, 10, 64)
			if err != nil {
				failureReason = "invalid job id"
				break
			}
			job, err := h.DB.GetQueryRunnerJob(uint(jobId))
			if err != nil {
				failureReason = err.Error()
				break
			}
			if job == nil {
				failureReason = "job not found"
				break
			}
			if job.Status == queryrunner.QueryRunnerCreated {
				err = h.DB.UpdateQueryRunnerJobStatus(job.ID, queryrunner.QueryRunnerCanceled, "")
				if err != nil {
					failureReason = err.Error()
					break
				}
			} else if job.Status == queryrunner.QueryRunnerInProgress {
				failureReason = "job is already in progress, unable to cancel"
				break
			} else if job.Status == queryrunner.QueryRunnerQueued {
				err = h.Scheduler.jq.DeleteMessage(ctx.Request().Context(), queryrunner.StreamName, job.NatsSequenceNumber)
				if err != nil {
					failureReason = err.Error()
					break
				}
				err = h.DB.UpdateQueryRunnerJobStatus(job.ID, queryrunner.QueryRunnerCanceled, "")
				if err != nil {
					failureReason = err.Error()
					break
				}
			} else {
				failureReason = "job is already finished"
				break
			}

		default:
			failureReason = "invalid job type"
			break
		}
		results = append(results, api.CancelJobResponse{
			JobId:    jobIdStr,
			JobType:  strings.ToLower(request.JobType),
			Canceled: canceled,
			Reason:   failureReason,
		})
	}

	return ctx.JSON(http.StatusOK, results)
}

// ListJobsByType godoc
//
//	@Summary	List jobs by job type and filters
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.ListJobsByTypeRequest	true	"Request Body"
//	@Produce	json
//	@Success	200	{object}	[]api.ListJobsByTypeResponse
//	@Router		/schedule/api/v3/jobs [post]
func (h HttpServer) ListJobsByType(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.ListJobsByTypeRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	if strings.ToLower(request.JobType) != "compliance" && strings.ToLower(request.JobType) != "discovery" &&
		strings.ToLower(request.JobType) != "query_run" && strings.ToLower(request.JobType) != "queryrun" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job type")
	}

	sortBy := "id"
	switch request.SortBy {
	case api.JobSort_ByIntegrationID, api.JobSort_ByJobID, api.JobSort_ByJobType, api.JobSort_ByStatus, api.JobSort_ByCreatedAt,
		api.JobSort_ByUpdatedAt:
		sortBy = string(request.SortBy)
	}

	if strings.ToLower(request.JobType) == "query" && strings.ToLower(request.Selector) != "job_id" {
		return echo.NewHTTPError(http.StatusBadRequest, "only jobId is acceptable for query run")
	}

	var items []api.ListJobsByTypeItem

	var err error
	switch strings.ToLower(request.Selector) {
	case "job_id":
		switch strings.ToLower(request.JobType) {
		case "query":
			var jobs []model2.QueryRunnerJob
			if len(request.JobId) == 0 {
				jobs, err = h.DB.ListQueryRunnerJobs()
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			} else {
				jobs, err = h.DB.ListQueryRunnerJobsById(request.JobId)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			}

			for _, j := range jobs {
				items = append(items, api.ListJobsByTypeItem{
					JobId:     strconv.Itoa(int(j.ID)),
					JobType:   strings.ToLower(request.JobType),
					JobStatus: string(j.Status),
					CreatedAt: j.CreatedAt,
					UpdatedAt: j.UpdatedAt,
				})
			}
		case "compliance":
			jobs, err := h.DB.ListComplianceJobsByIds(aws.Bool(true), request.JobId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			for _, j := range jobs {
				items = append(items, api.ListJobsByTypeItem{
					JobId:     strconv.Itoa(int(j.ID)),
					JobType:   strings.ToLower(request.JobType),
					JobStatus: string(j.Status),
					CreatedAt: j.CreatedAt,
					UpdatedAt: j.UpdatedAt,
				})
			}
		case "discovery":
			jobs, err := h.DB.ListDescribeJobsByIds(request.JobId)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			for _, j := range jobs {
				items = append(items, api.ListJobsByTypeItem{
					JobId:     strconv.Itoa(int(j.ID)),
					JobType:   strings.ToLower(request.JobType),
					JobStatus: string(j.Status),
					CreatedAt: j.CreatedAt,
					UpdatedAt: j.UpdatedAt,
				})
			}

		}
	case "integration_info":
		var integrations []integrationapi.Integration
		for _, info := range request.IntegrationInfo {
			if info.IntegrationID != nil {
				connection, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err.Error())
				}
				if connection != nil {
					integrations = append(integrations, *connection)
				}
				continue
			}

			var integrationTypes []string
			if info.IntegrationType != nil {
				integrationTypes = append(integrationTypes, *info.IntegrationType)
			}
			connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
				integrationapi.ListIntegrationsRequest{
					IntegrationType: integrationTypes,
					NameRegex:       info.Name,
					ProviderIDRegex: info.ProviderID,
				})
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			integrations = append(integrations, connectionsTmp.Integrations...)
		}

		connectionIDsMap := make(map[string]bool)
		for _, c := range integrations {
			connectionIDsMap[c.IntegrationID] = true
		}
		var connectionIDs []string
		switch strings.ToLower(request.JobType) {
		case "compliance":
			jobs, err := h.DB.ListComplianceJobsByIntegrationID(aws.Bool(true), connectionIDs)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			for _, j := range jobs {
				items = append(items, api.ListJobsByTypeItem{
					JobId:     strconv.Itoa(int(j.ID)),
					JobType:   strings.ToLower(request.JobType),
					JobStatus: string(j.Status),
					CreatedAt: j.CreatedAt,
					UpdatedAt: j.UpdatedAt,
				})
			}
		case "discovery":
			jobs, err := h.DB.ListDescribeJobsByFilters(nil, connectionIDs, nil, nil, nil, nil, nil)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			for _, j := range jobs {
				items = append(items, api.ListJobsByTypeItem{
					JobId:     strconv.Itoa(int(j.ID)),
					JobType:   strings.ToLower(request.JobType),
					JobStatus: string(j.Status),
					CreatedAt: j.CreatedAt,
					UpdatedAt: j.UpdatedAt,
				})
			}

		}
	case "status":
		for _, status := range request.Status {
			switch strings.ToLower(request.JobType) {
			case "compliance":
				jobs, err := h.DB.ListComplianceJobsByStatus(aws.Bool(true), model2.ComplianceJobStatus(strings.ToUpper(status)))
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				for _, j := range jobs {
					items = append(items, api.ListJobsByTypeItem{
						JobId:     strconv.Itoa(int(j.ID)),
						JobType:   strings.ToLower(request.JobType),
						JobStatus: string(j.Status),
						CreatedAt: j.CreatedAt,
						UpdatedAt: j.UpdatedAt,
					})
				}
			case "discovery":
				jobs, err := h.DB.ListDescribeJobsByStatus(api.DescribeResourceJobStatus(strings.ToUpper(status)))
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				for _, j := range jobs {
					items = append(items, api.ListJobsByTypeItem{
						JobId:     strconv.Itoa(int(j.ID)),
						JobType:   strings.ToLower(request.JobType),
						JobStatus: string(j.Status),
						CreatedAt: j.CreatedAt,
						UpdatedAt: j.UpdatedAt,
					})
				}

			}
		}
	case "benchmark":
		jobs, err := h.DB.ListComplianceJobsByFrameworkID(aws.Bool(true), request.Benchmark)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		for _, j := range jobs {
			items = append(items, api.ListJobsByTypeItem{
				JobId:     strconv.Itoa(int(j.ID)),
				JobType:   strings.ToLower(request.JobType),
				JobStatus: string(j.Status),
				CreatedAt: j.CreatedAt,
				UpdatedAt: j.UpdatedAt,
			})
		}
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid selector. valid values: job_id, integration_info, status, benchmark")
	}

	switch sortBy {
	case api.JobSort_ByJobID:
		sort.Slice(items, func(i, j int) bool {
			return items[i].JobId < items[j].JobId
		})
	case api.JobSort_ByCreatedAt:
		sort.Slice(items, func(i, j int) bool {
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		})
	case api.JobSort_ByUpdatedAt:
		sort.Slice(items, func(i, j int) bool {
			return items[i].UpdatedAt.Before(items[j].UpdatedAt)
		})
	case api.JobSort_ByStatus:
		sort.Slice(items, func(i, j int) bool {
			return items[i].JobStatus < items[j].JobStatus
		})
	}

	if request.SortOrder == api.JobSortOrder_DESC {
		sort.Slice(items, func(i, j int) bool {
			return i > j
		})
	}

	totalCount := len(items)

	if request.PerPage != nil {
		if request.Cursor == nil {
			items = utils.Paginate(1, *request.PerPage, items)
		} else {
			items = utils.Paginate(*request.Cursor, *request.PerPage, items)
		}
	}

	response := api.ListJobsByTypeResponse{
		Items:      items,
		TotalCount: totalCount,
	}

	return ctx.JSON(http.StatusOK, response)
}

// ListJobsInterval godoc
//
//	@Summary	List jobs by job type and filters
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		job_type		query	string	true	"Job Type"
//	@Param		interval		query	string	true	"Time Interval to filter by"
//	@Param		trigger_type	query	string	true	"Trigger Type: (all(default), manual, system)"
//	@Param		created_by		query	string	true	"Created By User ID"
//	@Param		cursor			query	int		true	"cursor"
//	@Param		per_page		query	int		true	"per page"
//	@Produce	json
//	@Success	200	{object}	[]api.ListJobsByTypeItem
//	@Router		/schedule/api/v3/jobs/interval [get]
func (h HttpServer) ListJobsInterval(ctx echo.Context) error {
	jobType := ctx.QueryParam("job_type")
	interval := ctx.QueryParam("interval")
	triggerType := ctx.QueryParam("trigger_type")
	createdBy := ctx.QueryParam("created_by")

	var cursor, perPage int64
	var err error
	cursorStr := ctx.QueryParam("cursor")
	if cursorStr != "" {
		cursor, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			return err
		}
	}
	perPageStr := ctx.QueryParam("per_page")
	if perPageStr != "" {
		perPage, err = strconv.ParseInt(perPageStr, 10, 64)
		if err != nil {
			return err
		}
	}

	convertedInterval, err := convertInterval(interval)
	if err != nil {
		h.Scheduler.logger.Error("invalid interval", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid interval")
	}

	if strings.ToLower(jobType) != "compliance" && strings.ToLower(jobType) != "discovery" &&
		strings.ToLower(jobType) != "query" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job type")
	}

	var items []api.ListJobsByTypeItem

	switch strings.ToLower(jobType) {
	case "compliance":
		jobs, err := h.DB.ListComplianceJobsForInterval(aws.Bool(true), convertedInterval, triggerType, createdBy)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		for _, j := range jobs {
			items = append(items, api.ListJobsByTypeItem{
				JobId:     strconv.Itoa(int(j.ID)),
				JobType:   strings.ToLower(jobType),
				JobStatus: string(j.Status),
				CreatedAt: j.CreatedAt,
				UpdatedAt: j.UpdatedAt,
			})
		}
	case "discovery":
		jobs, err := h.DB.ListDescribeJobsForInterval(convertedInterval, triggerType, createdBy)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		for _, j := range jobs {
			items = append(items, api.ListJobsByTypeItem{
				JobId:     strconv.Itoa(int(j.ID)),
				JobType:   strings.ToLower(jobType),
				JobStatus: string(j.Status),
				CreatedAt: j.CreatedAt,
				UpdatedAt: j.UpdatedAt,
			})
		}

	case "query":
		jobs, err := h.DB.ListQueryRunnerJobForInterval(convertedInterval, triggerType, createdBy)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		for _, j := range jobs {
			items = append(items, api.ListJobsByTypeItem{
				JobId:     strconv.Itoa(int(j.ID)),
				JobType:   strings.ToLower(jobType),
				JobStatus: string(j.Status),
				CreatedAt: j.CreatedAt,
				UpdatedAt: j.UpdatedAt,
			})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].JobId > items[j].JobId
	})

	totalCount := len(items)
	if perPage != 0 {
		if cursor == 0 {
			items = utils.Paginate(1, perPage, items)
		} else {
			items = utils.Paginate(cursor, perPage, items)
		}
	}

	return ctx.JSON(http.StatusOK, api.ListJobsIntervalResponse{
		TotalCount: totalCount,
		Items:      items,
	})
}

// GetSummaryJobs godoc
//
//	@Summary	List jobs by job type and filters
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		job_ids	query	[]string	true	"Compliance Job ID"
//	@Produce	json
//	@Success	200	{object}	[]string
//	@Router		/schedule/api/v3/jobs/compliance/summary/jobs [get]
func (h HttpServer) GetSummaryJobs(ctx echo.Context) error {
	jobIds := httpserver.QueryArrayParam(ctx, "job_ids")

	jobs, err := h.DB.ListSummaryJobs(jobIds)
	if err != nil {
		h.Scheduler.logger.Error("could not get summary jobs", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "could not get summary jobs")
	}

	var summaryJobIDs []string
	for _, j := range jobs {
		summaryJobIDs = append(summaryJobIDs, strconv.Itoa(int(j.ID)))
	}

	return ctx.JSON(http.StatusOK, summaryJobIDs)
}

func convertInterval(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	// Define regex to match shorthand formats like 90m, 50s, 1h
	re := regexp.MustCompile(`^(\d+)([smhd])$`)

	// Check if the input matches the shorthand format
	matches := re.FindStringSubmatch(input)
	if len(matches) == 3 {
		number := matches[1]
		unit := matches[2]

		// Map shorthand units to full words
		unitMap := map[string]string{
			"s": "seconds",
			"m": "minutes",
			"h": "hours",
			"d": "days",
		}

		// Convert the shorthand unit to full word
		if fullUnit, ok := unitMap[unit]; ok {
			return fmt.Sprintf("%s %s", number, fullUnit), nil
		}
	}

	// If the input doesn't match shorthand, assume it's already in the correct format
	if strings.Contains(input, "minute") || strings.Contains(input, "second") || strings.Contains(input, "hour") || strings.Contains(input, "day") {
		return input, nil
	}

	// If the input is invalid, return an error
	return "", fmt.Errorf("invalid interval format: %s", input)
}

// RunQuery godoc
//
//	@Summary	List jobs by job type and filters
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		query_id	path	string	true	"Policy ID"
//	@Produce	json
//	@Success	200	{object}	api.RunQueryResponse
//	@Router		/schedule/api/v3/query/{query_id}/run [put]
func (h HttpServer) RunQuery(ctx echo.Context) error {
	queryId := ctx.Param("query_id")

	userID := httpserver.GetUserID(ctx)
	if userID == "" {
		userID = "system"
	}

	job := &model2.QueryRunnerJob{
		QueryId:        queryId,
		Status:         queryrunner.QueryRunnerCreated,
		CreatedBy:      userID,
		FailureMessage: "",
		RetryCount:     0,
	}
	jobId, err := h.DB.CreateQueryRunnerJob(job)
	if err != nil {
		h.Scheduler.logger.Error("failed to create query runner job", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create query runner job")
	}

	response := api.RunQueryResponse{
		ID:        jobId,
		QueryId:   queryId,
		CreatedAt: job.CreatedAt,
		CreatedBy: userID,
		Status:    job.Status,
	}
	return ctx.JSON(http.StatusOK, response)
}

// GetIntegrationDiscoveryProgress godoc
//
//	@Summary	Get Integration discovery progress (number of jobs in different states)
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		request	body	api.GetIntegrationDiscoveryProgressRequest	true	"List jobs request"
//	@Produce	json
//	@Success	200	{object}	api.GetIntegrationDiscoveryProgressResponse
//	@Router		/schedule/api/v3/discovery/status [post]
func (h HttpServer) GetIntegrationDiscoveryProgress(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	var request api.GetIntegrationDiscoveryProgressRequest
	if err := ctx.Bind(&request); err != nil {
		ctx.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	var integrations []integrationapi.Integration
	for _, info := range request.IntegrationInfo {
		if info.IntegrationID != nil {
			connection, err := h.Scheduler.integrationClient.GetIntegration(clientCtx, *info.IntegrationID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			if connection != nil {
				integrations = append(integrations, *connection)
			}
			continue
		}
		var integrationTypes []string
		if info.IntegrationType != nil {
			integrationTypes = append(integrationTypes, *info.IntegrationType)
		}
		connectionsTmp, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx,
			integrationapi.ListIntegrationsRequest{
				IntegrationType: integrationTypes,
				NameRegex:       info.Name,
				ProviderIDRegex: info.ProviderID,
			})
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrations = append(integrations, connectionsTmp.Integrations...)
	}
	var err error
	if len(integrations) == 0 {
		integrationsTmp, err := h.Scheduler.integrationClient.ListIntegrations(clientCtx, nil)
		if err != nil {
			h.Scheduler.logger.Error("failed to list connections", zap.Error(err))
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		integrations = integrationsTmp.Integrations
	}

	connectionInfo := make(map[string]api.IntegrationInfo)
	var IntegrationIDs []string
	for _, c := range integrations {
		connectionInfo[c.IntegrationID] = api.IntegrationInfo{
			IntegrationID:   c.IntegrationID,
			IntegrationType: string(c.IntegrationType),
			Name:            c.Name,
			ProviderID:      c.ProviderID,
		}
		IntegrationIDs = append(IntegrationIDs, c.IntegrationID)
	}

	integrationDiscoveries, err := h.DB.ListIntegrationDiscovery(request.TriggerID, IntegrationIDs)
	if err != nil {
		h.Scheduler.logger.Error("cannot find integration discoveries", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "cannot find integration discoveries")
	}
	var integrationDiscoveriesIds []string
	for _, i := range integrationDiscoveries {
		integrationDiscoveriesIds = append(integrationDiscoveriesIds, strconv.Itoa(int(i.ID)))
	}

	jobs, err := h.DB.ListDescribeJobsByFilters(integrationDiscoveriesIds, nil, nil,
		nil, nil, nil, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	triggerIdProgressBreakdown := &api.DiscoveryProgressStatusBreakdown{}
	triggerIdProgressSummary := &api.DiscoveryProgressStatusSummary{}
	integrationsDiscoveryProgressStatus := make(map[string]api.IntegrationDiscoveryProgressStatus)
	for _, j := range jobs {
		if _, ok := integrationsDiscoveryProgressStatus[j.IntegrationID]; !ok {
			integrationsDiscoveryProgressStatus[j.IntegrationID] = api.IntegrationDiscoveryProgressStatus{
				Integration:             connectionInfo[j.IntegrationID],
				ProgressStatusBreakdown: &api.DiscoveryProgressStatusBreakdown{},
				ProgressStatusSummary:   &api.DiscoveryProgressStatusSummary{},
			}
		}
		switch j.Status {
		case api.DescribeResourceJobCreated:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.CreatedCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.CreatedCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			triggerIdProgressBreakdown.CreatedCount = triggerIdProgressBreakdown.CreatedCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
		case api.DescribeResourceJobQueued:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.QueuedCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.QueuedCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			triggerIdProgressBreakdown.QueuedCount = triggerIdProgressBreakdown.QueuedCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
		case api.DescribeResourceJobInProgress:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.InProgressCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.InProgressCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			triggerIdProgressBreakdown.InProgressCount = triggerIdProgressBreakdown.InProgressCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
		case api.DescribeResourceJobOldResourceDeletion:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.OldResourceDeletionCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.OldResourceDeletionCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			triggerIdProgressBreakdown.OldResourceDeletionCount = triggerIdProgressBreakdown.OldResourceDeletionCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
		case api.DescribeResourceJobTimeout:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.TimeoutCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.TimeoutCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount + 1
			triggerIdProgressBreakdown.TimeoutCount = triggerIdProgressBreakdown.TimeoutCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
			triggerIdProgressSummary.ProcessedCount = triggerIdProgressSummary.ProcessedCount + 1
		case api.DescribeResourceJobFailed:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.FailedCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.FailedCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount + 1
			triggerIdProgressBreakdown.FailedCount = triggerIdProgressBreakdown.FailedCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
			triggerIdProgressSummary.ProcessedCount = triggerIdProgressSummary.ProcessedCount + 1
		case api.DescribeResourceJobSucceeded:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.SucceededCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.SucceededCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount + 1
			triggerIdProgressBreakdown.SucceededCount = triggerIdProgressBreakdown.SucceededCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
			triggerIdProgressSummary.ProcessedCount = triggerIdProgressSummary.ProcessedCount + 1
		case api.DescribeResourceJobRemovingResources:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.RemovingResourcesCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.RemovingResourcesCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			triggerIdProgressBreakdown.RemovingResourcesCount = triggerIdProgressBreakdown.RemovingResourcesCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
		case api.DescribeResourceJobCanceled:
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.CanceledCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusBreakdown.CanceledCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.TotalCount + 1
			integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount = integrationsDiscoveryProgressStatus[j.IntegrationID].ProgressStatusSummary.ProcessedCount + 1
			triggerIdProgressBreakdown.CanceledCount = triggerIdProgressBreakdown.CanceledCount + 1
			triggerIdProgressSummary.TotalCount = triggerIdProgressSummary.TotalCount + 1
			triggerIdProgressSummary.ProcessedCount = triggerIdProgressSummary.ProcessedCount + 1
		}
	}

	var integrationsDiscoveryProgressStatusResult []api.IntegrationDiscoveryProgressStatus
	for _, v := range integrationsDiscoveryProgressStatus {
		integrationsDiscoveryProgressStatusResult = append(integrationsDiscoveryProgressStatusResult, v)
	}

	response := api.GetIntegrationDiscoveryProgressResponse{
		IntegrationProgress:        integrationsDiscoveryProgressStatusResult,
		TriggerIdProgressBreakdown: triggerIdProgressBreakdown,
		TriggerIdProgressSummary:   triggerIdProgressSummary,
	}

	return ctx.JSON(http.StatusOK, response)
}

// ListComplianceJobsHistory godoc
//
//	@Summary	List jobs by job type and filters
//	@Security	BearerToken
//	@Tags		scheduler
//	@Param		interval		query	string		true	"Time Interval to filter by"
//	@Param		trigger_type	query	string		true	"Trigger Type: (all(default), manual, system)"
//	@Param		created_by		query	string		true	"Created By User ID"
//	@Param		benchmark_ids	query	[]string	true	"Created By User ID"
//	@Param		cursor			query	int			true	"cursor"
//	@Param		per_page		query	int			true	"per page"
//	@Produce	json
//	@Success	200	{object}	api.ListComplianceJobsHistoryResponse
//	@Router		/schedule/api/v3/jobs/history/compliance [get]
func (h HttpServer) ListComplianceJobsHistory(ctx echo.Context) error {
	clientCtx := &httpclient.Context{UserRole: apiAuth.AdminRole}

	interval := ctx.QueryParam("interval")
	triggerType := ctx.QueryParam("trigger_type")
	createdBy := ctx.QueryParam("created_by")
	benchmarkIDs := httpserver.QueryArrayParam(ctx, "benchmark_ids")

	var cursor, perPage int64
	var err error
	cursorStr := ctx.QueryParam("cursor")
	if cursorStr != "" {
		cursor, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil {
			return err
		}
	}
	perPageStr := ctx.QueryParam("per_page")
	if perPageStr != "" {
		perPage, err = strconv.ParseInt(perPageStr, 10, 64)
		if err != nil {
			return err
		}
	}

	convertedInterval, err := convertInterval(interval)
	if err != nil {
		h.Scheduler.logger.Error("invalid interval", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid interval")
	}

	var items []api.ListComplianceJobsHistoryItem
	connectionIdsMap := make(map[string]bool)

	jobs, err := h.DB.ListComplianceJobsWithSummaryJob(aws.Bool(true), convertedInterval, triggerType, createdBy, benchmarkIDs)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for _, j := range jobs {
		var integrations []api.IntegrationInfo
		for _, c := range j.ConnectionIDs {
			connectionIdsMap[c] = true
			integrations = append(integrations, api.IntegrationInfo{
				IntegrationID: c,
			})
		}
		items = append(items, api.ListComplianceJobsHistoryItem{
			BenchmarkId:    j.BenchmarkID,
			JobId:          strconv.Itoa(int(j.ID)),
			SummarizerJobs: j.SummarizerJobs,
			TriggerType:    string(j.TriggerType),
			CreatedBy:      j.CreatedBy,
			JobStatus:      string(j.Status),
			CreatedAt:      j.CreatedAt,
			UpdatedAt:      j.UpdatedAt,
		})
	}
	var connectionIds []string
	for k, _ := range connectionIdsMap {
		connectionIds = append(connectionIds, k)
	}

	integrations, err := h.Scheduler.integrationClient.ListIntegrationsByFilters(clientCtx, integrationapi.ListIntegrationsRequest{
		IntegrationID: connectionIds,
	})
	if err != nil {
		h.Scheduler.logger.Error("failed to get connections", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get integrations info")
	}
	integrationsMap := make(map[string]api.IntegrationInfo)
	for _, c := range integrations.Integrations {
		integrationsMap[c.IntegrationID] = api.IntegrationInfo{
			IntegrationType: string(c.IntegrationType),
			ProviderID:      c.ProviderID,
			Name:            c.Name,
			IntegrationID:   c.IntegrationID,
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].JobId > items[j].JobId
	})

	totalCount := len(items)
	if perPage != 0 {
		if cursor == 0 {
			items = utils.Paginate(1, perPage, items)
		} else {
			items = utils.Paginate(cursor, perPage, items)
		}
	}

	for i, j := range items {
		for ii, integration := range j.Integrations {
			if integrationData, ok := integrationsMap[integration.IntegrationID]; ok {
				items[i].Integrations[ii] = integrationData
			}
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].JobId > items[j].JobId
	})
	return ctx.JSON(http.StatusOK, api.ListComplianceJobsHistoryResponse{
		TotalCount: totalCount,
		Items:      items,
	})
}

func parseTimeInterval(intervalStr string) (*time.Time, *time.Time, error) {
	// Define regex patterns to extract the time components
	patterns := map[string]*regexp.Regexp{
		"days":    regexp.MustCompile(`(\d+)\s*days?`),
		"hours":   regexp.MustCompile(`(\d+)\s*hours?`),
		"minutes": regexp.MustCompile(`(\d+)\s*minutes?`),
		"seconds": regexp.MustCompile(`(\d+)\s*seconds?`),
	}

	// Variables to store the extracted values
	days, hours, minutes, seconds := 0, 0, 0, 0

	// Extract and convert the values from the string
	for key, pattern := range patterns {
		match := pattern.FindStringSubmatch(intervalStr)
		if len(match) > 1 {
			value, err := strconv.Atoi(match[1])
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing %s: %v", key, err)
			}
			switch key {
			case "days":
				days = value
			case "hours":
				hours = value
			case "minutes":
				minutes = value
			case "seconds":
				seconds = value
			}
		}
	}

	// Calculate total duration based on extracted values
	duration := time.Duration(days)*24*time.Hour +
		time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second

	// Calculate endTime as now and startTime by subtracting the duration
	endTime := time.Now()
	startTime := endTime.Add(-duration)

	return &startTime, &endTime, nil
}

// CreateComplianceQuickSequence godoc
//
//	@Summary		Create Compliance Quick Sequence
//	@Description	Create Compliance Quick Sequence
//	@Security		BearerToken
//	@Tags			workspace
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/schedule/api/v3/compliance/quick/sequence [post]
func (h HttpServer) CreateComplianceQuickSequence(c echo.Context) error {
	var request api.CreateAuditJobRequest
	if err := c.Bind(&request); err != nil {
		c.Logger().Errorf("bind the request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	userID := httpserver.GetUserID(c)
	if userID == "" {
		userID = "system"
	}

	jobId, err := h.DB.CreateQuickScanSequence(&model2.QuickScanSequence{
		FrameworkID:    request.FrameworkID,
		IntegrationIDs: request.IntegrationIDs,
		IncludeResults: request.IncludeResults,
		Status:         model2.QuickScanSequenceCreated,
		CreatedBy:      userID,
	})
	if err != nil {
		h.Scheduler.logger.Error("failed to create quick scan sequence", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create quick scan sequence")
	}

	return c.String(http.StatusOK, strconv.Itoa(int(jobId)))
}

// GetComplianceQuickSequence godoc
//
//	@Summary		Get Compliance Quick Sequence by run id
//	@Description	Get Compliance Quick Sequence by run id
//	@Security		BearerToken
//	@Tags			audit
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/schedule/api/v3/compliance/quick/sequence/{run_id} [get]
func (h HttpServer) GetComplianceQuickSequence(c echo.Context) error {
	jobIdStr := c.Param("run_id")

	var jobId int64
	var err error
	if jobIdStr != "" {
		jobId, err = strconv.ParseInt(jobIdStr, 10, 64)
		if err != nil {
			return err
		}
	}

	job, err := h.DB.GetQuickScanSequenceByID(uint(jobId))
	if err != nil {
		h.Scheduler.logger.Error("failed to get compliance quick run", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get compliance quick run")
	}

	jobApi := job.ToAPI()

	if jobApi.Status == api.QuickScanSequenceFinished || jobApi.Status == api.QuickScanSequenceComplianceRunning ||
		jobApi.Status == api.QuickScanSequenceFailed {
		quickScan, err := h.DB.GetComplianceJobByCreatedByAndParentID("QuickScanSequencer", job.ID)
		if err != nil {
			h.Scheduler.logger.Error("failed to get compliance quick run", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get compliance quick run")
		}
		jobApi.ComplianceQuickRunID = aws.String(strconv.Itoa(int(quickScan.ID)))
	}

	return c.JSON(http.StatusOK, jobApi)
}

// PurgeSampleData godoc
//
//	@Summary		Delete integrations with SAMPLE_INTEGRATION state
//	@Description	Delete integrations with SAMPLE_INTEGRATION state
//	@Security		BearerToken
//	@Tags			credentials
//	@Param			integrations	query	[]string	false	"Sample Integrations"
//	@Produce		json
//	@Success		200
//	@Router			/schedule/api/v3/sample/purge [put]
func (h HttpServer) PurgeSampleData(c echo.Context) error {
	integrations := httpserver.QueryArrayParam(c, "integrations")

	err := h.DB.CleanupAllDescribeIntegrationJobsForIntegrations(integrations)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete describe connection jobs")
	}

	complianceJobs, err := h.DB.ListComplianceJobsByFilters(nil, integrations, nil, nil, nil, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get compliance jobs")
	}

	var complianceJobsIds []uint
	for _, job := range complianceJobs {
		complianceJobsIds = append(complianceJobsIds, job.ID)
	}

	summaryJobs, err := h.DB.ListAllComplianceSummarizerJobsByComplianceJobs(complianceJobsIds)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get summary jobs")
	}
	var summaryJobsIds []uint
	for _, job := range summaryJobs {
		summaryJobsIds = append(summaryJobsIds, job.ID)
	}

	err = h.DB.CleanupAllComplianceJobsForIntegrations(integrations)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete compliance jobs")
	}
	err = h.DB.CleanupAllComplianceSummarizerJobsByComplianceJobs(complianceJobsIds)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete compliance summarizer jobs")
	}
	err = h.DB.CleanupAllComplianceRunnersByComplianceJobs(complianceJobsIds)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete compliance runners")
	}

	maxID := uint(0)
	for _, i := range summaryJobsIds {
		maxID = max(i, maxID)
	}

	var ids []uint
	for i := uint(1); i <= maxID; i++ {
		ids = append(ids, i)
	}

	_, err = h.Scheduler.es.ES().Indices.Delete([]string{types.BenchmarkSummaryIndex}, h.Scheduler.es.ES().Indices.Delete.WithContext(context.Background()))
	if err != nil {
		h.Scheduler.logger.Error("failed to delete summary job", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete summary jobs")
	}

	return c.NoContent(http.StatusOK)
}
