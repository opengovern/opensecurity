package api

import (
	"time"

	queryrunner "github.com/opengovern/opencomply/jobs/query-runner-job"
)

type QuickScanSequenceStatus string

const (
	QuickScanSequenceCreated              QuickScanSequenceStatus = "CREATED"
	QuickScanSequenceStarted              QuickScanSequenceStatus = "STARTED"
	QuickScanSequenceFetchingDependencies QuickScanSequenceStatus = "FETCH_DEPENDENCIES"
	QuickScanSequenceComplianceRunning    QuickScanSequenceStatus = "RUNNING_COMPLIANCE_QUICK_SCAN"
	QuickScanSequenceFailed               QuickScanSequenceStatus = "FAILED"
	QuickScanSequenceFinished             QuickScanSequenceStatus = "FINISHED"
)

type JobType string

const (
	JobType_Discovery  JobType = "discovery"
	JobType_Compliance JobType = "compliance"
)

type JobStatus string

type JobSort string

const (
	JobSort_ByJobID         = "id"
	JobSort_ByJobType       = "job_type"
	JobSort_ByIntegrationID = "integration_id"
	JobSort_ByStatus        = "status"
	JobSort_ByCreatedAt     = "created_at"
	JobSort_ByUpdatedAt     = "updated_at"
)

type JobSortOrder string

const (
	JobSortOrder_ASC  = "ASC"
	JobSortOrder_DESC = "DESC"
)

type ListJobsRequest struct {
	Interval     *string  `json:"interval"`
	From         *int64   `json:"from"`
	To           *int64   `json:"to"`
	PageStart    int      `json:"pageStart"`
	PageEnd      int      `json:"pageEnd"`
	TypeFilters  []string `json:"typeFilters"`
	StatusFilter []string `json:"statusFilter"`

	SortBy    JobSort      `json:"sortBy"`
	SortOrder JobSortOrder `json:"sortOrder"`
}

type Job struct {
	ID                     uint      `json:"id"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	Type                   JobType   `json:"type"`
	ConnectionID           string    `json:"connectionID"`
	ConnectionProviderID   string    `json:"connectionProviderID"`
	ConnectionProviderName string    `json:"connectionProviderName"`
	Title                  string    `json:"title"`
	Status                 string    `json:"status"`
	FailureReason          string    `json:"failureReason"`
}

type JobSummary struct {
	Type   JobType `json:"type"`
	Status string  `json:"status"`
	Count  int64   `json:"count"`
}

type ListJobsResponse struct {
	Jobs      []Job        `json:"jobs"`
	Summaries []JobSummary `json:"summaries"`
}

type JobSeqCheckResponse struct {
	IsRunning bool `json:"isRunning"`
}

type GetDescribeJobsHistoryRequest struct {
	ResourceType  []string   `json:"resource_type"`
	DiscoveryType []string   `json:"discovery_type"`
	JobStatus     []string   `json:"job_status"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	SortBy        *string    `json:"sort_by"`
	Cursor        *int64     `json:"cursor"`
	PerPage       *int64     `json:"per_page"`
}

type GetDescribeJobsHistoryResponse struct {
	JobId           uint                      `json:"job_id"`
	ResourceType    string                    `json:"resource_type"`
	JobStatus       DescribeResourceJobStatus `json:"job_status"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	Title 		 string              `json:"title"`
	FailureMessage string              `json:"failure_message"`
	IntegrationInfo *IntegrationInfo          `json:"integration_info"`
}
type GetDescribeJobsHistoryFinalResponse struct {
	Items      []GetDescribeJobsHistoryResponse `json:"items"`
	TotalCount int                             `json:"total_count"`
}

type GetComplianceJobsHistoryRequest struct {
	BenchmarkId []string   `json:"benchmark_id"`
	JobStatus   []string   `json:"job_status"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	SortBy      *string    `json:"sort_by"`
	Cursor      *int64     `json:"cursor"`
	PerPage     *int64     `json:"per_page"`
}

type GetComplianceJobsHistoryResponse struct {
	JobId           uint                `json:"job_id"`
	WithIncidents   bool                `json:"with_incidents"`
	BenchmarkId     string              `json:"benchmark_id"`
	JobStatus       ComplianceJobStatus `json:"job_status"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	Title 		 string              `json:"title"`
	FailureMessage string              `json:"failure_message"`
	IntegrationInfo []IntegrationInfo   `json:"integration_info"`
}
type GetComplianceJobsHistoryFinalResponse struct {
	Items      []GetComplianceJobsHistoryResponse `json:"items"`
	TotalCount int                               `json:"total_count"`
}

type BenchmarkAuditHistoryItem struct {
	JobId                uint                `json:"job_id"`
	WithIncidents        bool                `json:"with_incidents"`
	BenchmarkId          string              `json:"benchmark_id"`
	JobStatus            ComplianceJobStatus `json:"job_status"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	IntegrationInfo      []IntegrationInfo   `json:"integration_info"`
	NumberOfIntegrations int                 `json:"number_of_integrations"`
}

type BenchmarkAuditHistoryResponse struct {
	Items      []BenchmarkAuditHistoryItem `json:"items"`
	TotalCount int                         `json:"total_count"`
}

type RunBenchmarkByIdRequest struct {
	WithIncidents   bool `json:"with_incidents"`
	IntegrationInfo []struct {
		IntegrationType *string `json:"integration_type"`
		ProviderID      *string `json:"provider_id"`
		Name            *string `json:"name"`
		IntegrationID   *string `json:"integration_id"`
	} `json:"integration_info"`
}

type RunBenchmarkRequest struct {
	BenchmarkIds    []string                `json:"benchmark_ids"`
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
}

type IntegrationInfo struct {
	IntegrationType string `json:"integration_type"`
	ProviderID      string `json:"provider_id"`
	Name            string `json:"name"`
	IntegrationID   string `json:"integration_id"`
}

type IntegrationInfoFilter struct {
	IntegrationType *string `json:"integration_type"`
	ProviderID      *string `json:"provider_id"`
	Name            *string `json:"name"`
	IntegrationID   *string `json:"integration_id"`
}

type RunBenchmarkItem struct {
	JobId           uint              `json:"job_id"`
	WithIncident    bool              `json:"with_incident"`
	BenchmarkId     string            `json:"benchmark_id"`
	IntegrationInfo []IntegrationInfo `json:"integration_info"`
}

type RunBenchmarkResponse struct {
	Jobs []RunBenchmarkItem `json:"jobs"`
}

type ResourceTypeRunDiscoveryRequest struct {
	ResourceType string              `json:"resource_type"`
	Parameters   map[string]string `json:"parameters"`
}

type RunDiscoveryRequest struct {
	ResourceTypes   []ResourceTypeRunDiscoveryRequest `json:"resource_types"`
	IntegrationInfo []IntegrationInfoFilter           `json:"integration_info"`
}

type RunDiscoveryJob struct {
	JobId           uint            `json:"job_id"`
	ResourceType    string          `json:"resource_type"`
	Status          string          `json:"status"`
	FailureReason   string          `json:"failure_reason"`
	IntegrationInfo IntegrationInfo `json:"integration_info"`
}

type RunDiscoveryResponse struct {
	Jobs      []RunDiscoveryJob `json:"jobs"`
	TriggerID uint              `json:"trigger_id"`
}

type GetDescribeJobStatusResponse struct {
	JobId           uint            `json:"job_id"`
	IntegrationInfo IntegrationInfo `json:"integration_info"`
	JobStatus       string          `json:"job_status"`
	DiscoveryType   string          `json:"discovery_type"`
	ResourceType    string          `json:"resource_type"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type GetComplianceJobStatusResponse struct {
	JobId           uint              `json:"job_id"`
	WithIncidents   bool              `json:"with_incidents"`
	SummaryJobId    *uint             `json:"summary_job_id"`
	IntegrationInfo []IntegrationInfo `json:"integration_info"`
	JobStatus       string            `json:"job_status"`
	FrameworkId     string            `json:"framework_id"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type GetAsyncQueryRunJobStatusResponse struct {
	JobId          uint                          `json:"job_id"`
	QueryId        string                        `json:"query_id"`
	CreatedAt      time.Time                     `json:"created_at"`
	UpdatedAt      time.Time                     `json:"updated_at"`
	CreatedBy      string                        `json:"created_by"`
	JobStatus      queryrunner.QueryRunnerStatus `json:"job_status"`
	FailureMessage string                        `json:"failure_message"`
}

type ListDescribeJobsRequest struct {
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
	ResourceType    []string                `json:"resource_type"`
	DiscoveryType   []string                `json:"discovery_type"`
	JobStatus       []string                `json:"job_status"`
	Interval 	  *string                 `json:"interval"`
	StartTime       time.Time               `json:"start_time"`
	EndTime         *time.Time              `json:"end_time"`
	SortBy          *string                 `json:"sort_by"`
	Cursor          *int64                  `json:"cursor"`
	PerPage         *int64                  `json:"per_page"`
}

type ListComplianceJobsRequest struct {
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
	BenchmarkId     []string                `json:"benchmark_id"`
	JobStatus       []string                `json:"job_status"`
	StartTime       time.Time               `json:"start_time"`
	EndTime         *time.Time              `json:"end_time"`
	Interval 	  *string                 `json:"interval"`
	SortBy          *string                 `json:"sort_by"`
	Cursor          *int64                  `json:"cursor"`
	PerPage         *int64                  `json:"per_page"`
}

type BenchmarkAuditHistoryRequest struct {
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
	JobStatus       []string                `json:"job_status"`
	WithIncidents   *bool                   `json:"with_incidents"`
	Interval        *string                 `json:"interval"`
	StartTime       time.Time               `json:"start_time"`
	EndTime         *time.Time              `json:"end_time"`
	SortBy          *string                 `json:"sort_by"`
	Cursor          *int64                  `json:"cursor"`
	PerPage         *int64                  `json:"per_page"`
}

type GetIntegrationLastDiscoveryJobRequest struct {
	IntegrationInfo IntegrationInfoFilter `json:"integration_info"`
}

type GetDescribeJobsHistoryByIntegrationRequest struct {
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
	ResourceType    []string                `json:"resource_type"`
	DiscoveryType   []string                `json:"discovery_type"`
	JobStatus       []string                `json:"job_status"`
	StartTime       time.Time               `json:"start_time"`
	EndTime         *time.Time              `json:"end_time"`
	SortBy          *string                 `json:"sort_by"`
	Cursor          *int64                  `json:"cursor"`
	PerPage         *int64                  `json:"per_page"`
}

type GetComplianceJobsHistoryByIntegrationRequest struct {
	IntegrationInfo IntegrationInfoFilter `json:"integration_info"`
	WithIncidents   *bool                 `json:"with_incidents"`
	BenchmarkId     []string              `json:"benchmark_id"`
	JobStatus       []string              `json:"job_status"`
	StartTime       time.Time             `json:"start_time"`
	EndTime         *time.Time            `json:"end_time"`
	SortBy          *string               `json:"sort_by"`
	Cursor          *int64                `json:"cursor"`
	PerPage         *int64                `json:"per_page"`
}

type CancelJobRequest struct {
	JobType         string                  `json:"job_type"`
	Selector        string                  `json:"selector"`
	JobId           []string                `json:"job_id"`
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
	Status          []string                `json:"status"`
}

type CancelJobResponse struct {
	JobId    string `json:"job_id"`
	JobType  string `json:"job_type"`
	Canceled bool   `json:"canceled"`
	Reason   string `json:"reason"`
}

type ListJobsByTypeRequest struct {
	JobType         string                  `json:"job_type"`
	Selector        string                  `json:"selector"`
	JobId           []string                `json:"job_id"`
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
	Status          []string                `json:"status"`
	Benchmark       []string                `json:"benchmark"`
	SortBy          JobSort                 `json:"sort_by"`
	SortOrder       JobSortOrder            `json:"sort_order"`
	UpdatedAt       struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	} `json:"updated_at"`
	CreatedAt struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	} `json:"created_at"`
	Cursor  *int64 `json:"cursor"`
	PerPage *int64 `json:"per_page"`
}

type ListJobsByTypeItem struct {
	JobId     string    `json:"job_id"`
	JobType   string    `json:"job_type"`
	JobStatus string    `json:"job_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListJobsIntervalResponse struct {
	Items      []ListJobsByTypeItem `json:"items"`
	TotalCount int                  `json:"total_count"`
}

type ListJobsByTypeResponse struct {
	Items      []ListJobsByTypeItem `json:"items"`
	TotalCount int                  `json:"total_count"`
}

type ListComplianceJobsHistoryItem struct {
	BenchmarkId    string            `json:"benchmark_id"`
	Integrations   []IntegrationInfo `json:"integrations"`
	JobId          string            `json:"job_id"`
	SummarizerJobs []string          `json:"summarizer_jobs"`
	TriggerType    string            `json:"trigger_type"`
	CreatedBy      string            `json:"created_by"`
	JobStatus      string            `json:"job_status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type ListComplianceJobsHistoryResponse struct {
	Items      []ListComplianceJobsHistoryItem `json:"items"`
	TotalCount int                             `json:"total_count"`
}

type RunQueryResponse struct {
	ID        uint                          `json:"id"`
	CreatedAt time.Time                     `json:"created_at"`
	QueryId   string                        `json:"query_id"`
	CreatedBy string                        `json:"created_by"`
	Status    queryrunner.QueryRunnerStatus `json:"status"`
}

type GetIntegrationDiscoveryProgressRequest struct {
	IntegrationInfo []IntegrationInfoFilter `json:"integration_info"`
	TriggerID       string                  `json:"trigger_id"`
}

type DiscoveryProgressStatusBreakdown struct {
	CreatedCount             int64 `json:"created_count"`
	QueuedCount              int64 `json:"queued_count"`
	InProgressCount          int64 `json:"in_progress_count"`
	OldResourceDeletionCount int64 `json:"old_resource_deletion"`
	TimeoutCount             int64 `json:"timeout_count"`
	FailedCount              int64 `json:"failed_count"`
	SucceededCount           int64 `json:"succeeded_count"`
	RemovingResourcesCount   int64 `json:"removing_resources_count"`
	CanceledCount            int64 `json:"canceled_count"`
}

type DiscoveryProgressStatusSummary struct {
	TotalCount     int64 `json:"total_count"`
	ProcessedCount int64 `json:"processed_count"`
}

type IntegrationDiscoveryProgressStatus struct {
	Integration             IntegrationInfo                   `json:"integration"`
	ProgressStatusBreakdown *DiscoveryProgressStatusBreakdown `json:"breakdown"`
	ProgressStatusSummary   *DiscoveryProgressStatusSummary   `json:"summary"`
}

type GetIntegrationDiscoveryProgressResponse struct {
	IntegrationProgress        []IntegrationDiscoveryProgressStatus `json:"integration_progress"`
	TriggerIdProgressSummary   *DiscoveryProgressStatusSummary      `json:"trigger_id_progress_summary"`
	TriggerIdProgressBreakdown *DiscoveryProgressStatusBreakdown    `json:"trigger_id_progress_breakdown"`
}

type QuickScanSequence struct {
	ID                   string                  `json:"id"`
	FrameworkID          string                  `json:"framework_id"`
	IntegrationIDs       []string                `json:"integration_ids"`
	Status               QuickScanSequenceStatus `json:"status"`
	FailureMessage       string                  `json:"failure_message"`
	CreatedBy            string                  `json:"created_by"`
	ComplianceQuickRunID *string                 `json:"compliance_quick_run_id"`
}

type CreateAuditJobRequest struct {
	FrameworkID    string   `json:"framework_id"`
	IntegrationIDs []string `json:"integration_ids"`
	IncludeResults []string `json:"include_results"`
}
