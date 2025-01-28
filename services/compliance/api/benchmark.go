package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/pkg/types"
	"time"
)

type BenchmarkAssignmentStatus string

const (
	BenchmarkAssignmentStatusEnabled    BenchmarkAssignmentStatus = "enabled"
	BenchmarkAssignmentStatusDisabled   BenchmarkAssignmentStatus = "disabled"
	BenchmarkAssignmentStatusAutoEnable BenchmarkAssignmentStatus = "auto-enable"
)

type Benchmark struct {
	ID                string              `json:"id" example:"azure_cis_v140"`                                                                                                                                                       // Benchmark ID
	Title             string              `json:"title" example:"Azure CIS v1.4.0"`                                                                                                                                                  // Benchmark title
	ReferenceCode     string              `json:"referenceCode" example:"CIS 1.4.0"`                                                                                                                                                 // Benchmark display code
	Description       string              `json:"description" example:"The CIS Microsoft Azure Foundations Security Benchmark provides prescriptive guidance for establishing a secure baseline configuration for Microsoft Azure."` // Benchmark description
	LogoURI           string              `json:"logoURI"`                                                                                                                                                                           // Benchmark logo URI
	Category          string              `json:"category"`                                                                                                                                                                          // Benchmark category
	DocumentURI       string              `json:"documentURI" example:"benchmarks/azure_cis_v140.md"`                                                                                                                                // Benchmark document URI
	IsBaseline        bool                `json:"isBaseline" example:"true"`                                                                                                                                                         // Whether the benchmark is auto assigned or not
	Enabled           bool                `json:"enabled" example:"true"`
	TracksDriftEvents bool                `json:"tracksDriftEvents" example:"true"` // Whether the benchmark tracks drift events or not
	Tags              map[string][]string `json:"tags" `                            // Benchmark tags
	IntegrationTypes  []string            `json:"integrationTypes"`                 // Benchmark connectors
	Children          []string            `json:"children"`                         // Benchmark children
	Controls          []string            `json:"controls"`                         // Benchmark controls
	CreatedAt         time.Time           `json:"createdAt"`                        // Benchmark creation date
	UpdatedAt         time.Time           `json:"updatedAt"`                        // Benchmark last update date
}

type NestedBenchmark struct {
	ID               string              `json:"id" example:"azure_cis_v140"`                                                                                                                                                       // Benchmark ID
	Title            string              `json:"title" example:"Azure CIS v1.4.0"`                                                                                                                                                  // Benchmark title
	ReferenceCode    string              `json:"referenceCode" example:"CIS 1.4.0"`                                                                                                                                                 // Benchmark display code
	Description      string              `json:"description" example:"The CIS Microsoft Azure Foundations Security Benchmark provides prescriptive guidance for establishing a secure baseline configuration for Microsoft Azure."` // Benchmark description
	LogoURI          string              `json:"logoURI"`                                                                                                                                                                           // Benchmark logo URI
	Category         string              `json:"category"`                                                                                                                                                                          // Benchmark category
	DocumentURI      string              `json:"documentURI" example:"benchmarks/azure_cis_v140.md"`                                                                                                                                // Benchmark document URI
	AutoAssign       bool                `json:"autoAssign" example:"true"`                                                                                                                                                         // Whether the benchmark is auto assigned or not
	Tags             map[string][]string `json:"tags" `                                                                                                                                                                             // Benchmark tags
	IntegrationTypes []integration.Type  `json:"integrationTypes" example:"[azure]"`                                                                                                                                                // Benchmark connectors
	Children         []NestedBenchmark   `json:"children" example:"[azure_cis_v140_1, azure_cis_v140_2]"`                                                                                                                           // Benchmark children
	Controls         []string            `json:"controls" example:"[azure_cis_v140_1_1, azure_cis_v140_1_2]"`                                                                                                                       // Benchmark controls
	CreatedAt        time.Time           `json:"createdAt" example:"2020-01-01T00:00:00Z"`                                                                                                                                          // Benchmark creation date
	UpdatedAt        time.Time           `json:"updatedAt" example:"2020-01-01T00:00:00Z"`                                                                                                                                          // Benchmark last update date
}

type BenchmarkTrendDatapoint struct {
	Timestamp               time.Time                       `json:"timestamp" example:"1686346668"`
	ComplianceStatusSummary ComplianceStatusSummary         `json:"complianceStatusSummary"`
	Checks                  types.SeverityResult            `json:"checks"`
	ControlsSeverityStatus  BenchmarkControlsSeverityStatus `json:"controlsSeverityStatus"`
}

type ListBenchmarksSummaryResponse struct {
	BenchmarkSummary []BenchmarkEvaluationSummary `json:"benchmarkSummary"`

	TotalComplianceStatusSummary ComplianceStatusSummary `json:"totalComplianceStatusSummary"`
	TotalChecks                  types.SeverityResult    `json:"totalChecks"`
}

type BenchmarkStatusResult struct {
	PassedCount int `json:"passed"`
	TotalCount  int `json:"total"`
}

type BenchmarkStatusResultV2 struct {
	TotalCount  int `json:"total"`
	PassedCount int `json:"passed"`
	FailedCount int `json:"failed"`
}

type BenchmarkControlsSeverityStatus struct {
	Total BenchmarkStatusResult `json:"total"`

	Critical BenchmarkStatusResult `json:"critical"`
	High     BenchmarkStatusResult `json:"high"`
	Medium   BenchmarkStatusResult `json:"medium"`
	Low      BenchmarkStatusResult `json:"low"`
	None     BenchmarkStatusResult `json:"none"`
}

type BenchmarkResourcesSeverityStatus struct {
	Total BenchmarkStatusResult `json:"total"`

	Critical BenchmarkStatusResult `json:"critical"`
	High     BenchmarkStatusResult `json:"high"`
	Medium   BenchmarkStatusResult `json:"medium"`
	Low      BenchmarkStatusResult `json:"low"`
	None     BenchmarkStatusResult `json:"none"`
}

type BenchmarkControlsSeverityStatusV2 struct {
	Total BenchmarkStatusResultV2 `json:"total"`

	Critical BenchmarkStatusResultV2 `json:"critical"`
	High     BenchmarkStatusResultV2 `json:"high"`
	Medium   BenchmarkStatusResultV2 `json:"medium"`
	Low      BenchmarkStatusResultV2 `json:"low"`
	None     BenchmarkStatusResultV2 `json:"none"`
}

type BenchmarkResourcesSeverityStatusV2 struct {
	Total BenchmarkStatusResultV2 `json:"total"`

	Critical BenchmarkStatusResultV2 `json:"critical"`
	High     BenchmarkStatusResultV2 `json:"high"`
	Medium   BenchmarkStatusResultV2 `json:"medium"`
	Low      BenchmarkStatusResultV2 `json:"low"`
	None     BenchmarkStatusResultV2 `json:"none"`
}

type ComplianceStatusSummary struct {
	PassedCount int `json:"passed"`
	FailedCount int `json:"failed"`
}

type ComplianceStatusSummaryV2 struct {
	TotalCount  int `json:"total_count"`
	PassedCount int `json:"passed"`
	FailedCount int `json:"failed"`
}

func (c *ComplianceStatusSummary) AddESComplianceStatusMap(summary map[types.ComplianceStatus]int) {
	c.PassedCount += summary[types.ComplianceStatusOK]
	c.FailedCount += summary[types.ComplianceStatusALARM]
	c.PassedCount += summary[types.ComplianceStatusINFO]
	c.PassedCount += summary[types.ComplianceStatusSKIP]
	c.FailedCount += summary[types.ComplianceStatusERROR]
}

func (c *ComplianceStatusSummaryV2) AddESComplianceStatusMap(summary map[types.ComplianceStatus]int) {
	c.PassedCount += summary[types.ComplianceStatusOK]
	c.FailedCount += summary[types.ComplianceStatusALARM]
	c.PassedCount += summary[types.ComplianceStatusINFO]
	c.PassedCount += summary[types.ComplianceStatusSKIP]
	c.FailedCount += summary[types.ComplianceStatusERROR]

	c.TotalCount = c.FailedCount + c.PassedCount
}

type BenchmarkEvaluationSummary struct {
	Benchmark
	ComplianceStatusSummary ComplianceStatusSummary          `json:"complianceStatusSummary"`
	Checks                  types.SeverityResult             `json:"checks"`
	ControlsSeverityStatus  BenchmarkControlsSeverityStatus  `json:"controlsSeverityStatus"`
	ResourcesSeverityStatus BenchmarkResourcesSeverityStatus `json:"resourcesSeverityStatus"`
	IntegrationsStatus      BenchmarkStatusResult            `json:"IntegrationsStatus"`
	CostImpact              *float64                         `json:"costImpact"`
	EvaluatedAt             *time.Time                       `json:"evaluatedAt" example:"2020-01-01T00:00:00Z"`
	LastJobStatus           string                           `json:"lastJobStatus" example:"success"`
	TopIntegrations         []TopFieldRecord                 `json:"topIntegrations"`
}

type BenchmarkControlSummary struct {
	Benchmark Benchmark                 `json:"benchmark"`
	Controls  []ControlSummary          `json:"control"`
	Children  []BenchmarkControlSummary `json:"children"`
}

type GetBenchmarkDetailsRequest struct {
	TagsRegex               *string                         `json:"tags_regex"`
	ComplianceResultFilters *ComplianceResultSummaryFilters `json:"compliance_result_filters"`
	BenchmarkChildren       bool                            `json:"benchmark_children"`
}

type GetBenchmarkDetailsMetadata struct {
	ID                string              `json:"id"`
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	Enabled           bool                `json:"enabled"`
	IntegrationTypes  []integration.Type  `json:"integration_types"`
	NumberOfControls  int                 `json:"number_of_controls"`
	SupportedControls []string            `json:"supported_controls"`
	PrimaryResources  []string            `json:"primary_resources"`
	ListOfResources   []string            `json:"list_of_resources"`
	Tags              map[string][]string `json:"tags"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

type GetBenchmarkDetailsComplianceResults struct {
	Results         map[types.ComplianceStatus]int `json:"results"`
	LastEvaluatedAt time.Time                      `json:"lastEvaluated_at"`
	IntegrationIDs  []string                       `json:"integration_ids"`
}

type GetBenchmarkDetailsChildren struct {
	ID                string                               `json:"id"`
	Title             string                               `json:"title"`
	Tags              map[string][]string                  `json:"tags"`
	ControlIDs        []string                             `json:"control_ids"`
	ComplianceResults GetBenchmarkDetailsComplianceResults `json:"compliance_results"`
	Children          []GetBenchmarkDetailsChildren        `json:"children"`
}

type GetBenchmarkDetailsResponse struct {
	Metadata          GetBenchmarkDetailsMetadata          `json:"metadata"`
	ComplianceResults GetBenchmarkDetailsComplianceResults `json:"compliance_results"`
	Children          []GetBenchmarkDetailsChildren        `json:"children"`
}

type GetFrameworkListRequest struct {
	TitleRegex              *string                         `json:"title_regex"`
	ParentBenchmarkID       []string                        `json:"parent_benchmark_id"`
	Tags                    map[string][]string             `json:"tags"`
	TagsRegex               *string                         `json:"tags_regex"`
	PrimaryResource         []string                        `json:"primary_resource"`
	ListOfResources         []string                        `json:"list_of_resources"`
	Controls                []string                        `json:"controls"`
	Integration             []IntegrationFilter             `json:"integration"`
	IntegrationTypes        []string                        `json:"integration_types"`
	Root                    *bool                           `json:"root"`
	Assigned                *bool                           `json:"assigned"`
	IsBaseline              *bool                           `json:"is_baseline"`
	ComplianceResultFilters *ComplianceResultSummaryFilters `json:"compliance_result_filters"`
	SortBy                  string                          `json:"sort_by"`
	Cursor                  *int64                          `json:"cursor"`
	PerPage                 *int64                          `json:"per_page"`
}

type GetFrameworkSummaryListRequest struct {
	TitleRegex       *string  `json:"title_regex"`
	IntegrationTypes []string `json:"integration_types"`
	IsBaseline       *bool    `json:"is_baseline"`
	SortBy           string   `json:"sort_by"`
	Cursor           *int64   `json:"cursor"`
	PerPage          *int64   `json:"per_page"`
	Root             *bool    `json:"root"`
	Assigned         *bool    `json:"assigned"`
}

type GetBenchmarkListMetadata struct {
	ID                  string              `json:"id"`
	Title               string              `json:"title"`
	Description         string              `json:"description"`
	IntegrationType     []string            `json:"connectors"`
	NumberOfControls    int                 `json:"number_of_controls"`
	Enabled             bool                `json:"enabled"`
	TrackDriftEvents    bool                `json:"track_drift_events"`
	AutoAssigned        bool                `json:"auto_assigned"`
	NumberOfAssignments int                 `json:"number_of_assignments"`
	PrimaryTables       []string            `json:"primary_tables"`
	Tags                map[string][]string `json:"tags"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

type GetBenchmarkListItem struct {
	Benchmark     GetBenchmarkListMetadata `json:"benchmark"`
	IncidentCount int                      `json:"incident_count"`
}

type GetBenchmarkListResponse struct {
	Items      []GetBenchmarkListItem `json:"items"`
	TotalCount int                    `json:"total_count"`
}

type GetBenchmarkListSummaryMetadata struct {
	ID                         string                             `json:"id"`
	Title                      string                             `json:"title"`
	Description                string                             `json:"description"`
	IntegrationType            []string                           `json:"connectors"`
	Enabled                    bool                               `json:"enabled"`
	CreatedAt                  time.Time                          `json:"created_at"`
	UpdatedAt                  time.Time                          `json:"updated_at"`
	ComplianceScore            float64                            `json:"compliance_score"`
	SeveritySummaryByControl   BenchmarkControlsSeverityStatusV2  `json:"severity_summary_by_control"`
	SeveritySummaryByResource  BenchmarkResourcesSeverityStatusV2 `json:"severity_summary_by_resource"`
	SeveritySummaryByIncidents types.SeverityResultV2             `json:"severity_summary_by_incidents"`
	CostImpact                 *float64                           `json:"cost_impact"`
	ComplianceResultsSummary   ComplianceStatusSummaryV2          `json:"compliance_results_summary"`
	IssuesCount                int                                `json:"issues_count"`
	LastEvaluatedAt            *time.Time                         `json:"last_evaluated_at"`
	LastJobStatus              string                             `json:"last_job_status"`
}

type GetBenchmarkSummaryListResponse struct {
	Items      []GetBenchmarkListSummaryMetadata `json:"items"`
	TotalCount int                               `json:"total_count"`
}

type GetBenchmarkAssignmentsResponse struct {
	Items  []GetBenchmarkAssignmentsItem `json:"items"`
	Status BenchmarkAssignmentStatus     `json:"status"`
}

type GetBenchmarkAssignmentsItem struct {
	Integration              IntegrationInfo `json:"integration"`
	Assigned                 bool            `json:"assigned"`
	AssignmentType           *string         `json:"assignment_type"`
	AssignmentChangePossible bool            `json:"assignment_change_possible"`
}

type IntegrationInfo struct {
	IntegrationType string  `json:"integration_type"`
	ProviderID      *string `json:"provider_id"`
	Name            *string `json:"name"`
	IntegrationID   *string `json:"integration_id"`
}

type IntegrationFilter struct {
	IntegrationType *string `json:"integration_type"`
	ProviderID      *string `json:"provider_id"`
	Name            *string `json:"name"`
	IntegrationID   *string `json:"integration_id"`
}

type IntegrationFilterRequest struct {
	Integration []IntegrationFilter `json:"integration"`
	AutoEnable  bool                `json:"auto_enable"`
	Disable     bool                `json:"disable"`
}

type SeveritySummary struct {
	Issues    int `json:"issues"`
	Resources int `json:"resources"`
	Controls  int `json:"controls"`
}
type TopIntegration struct {
	IntegrationInfo IntegrationInfo `json:"integration_info"`
	Issues          int             `json:"issues"`
}

type ComplianceSummaryOfIntegrationRequest struct {
	BenchmarkId string            `json:"benchmark_id"'`
	Integration IntegrationFilter `json:"integration"`
	ShowTop     int               `json:"show_top"`
}

type TopFiledRecordV2 struct {
	Field  string `json:"field"`
	Key    string `json:"key"`
	Issues int    `json:"issues"`
}

type ComplianceSummaryOfIntegrationResponse struct {
	BenchmarkID                string                             `json:"benchmark_id"`
	Integration                IntegrationInfo                    `json:"integration"`
	ComplianceScore            float64                            `json:"compliance_score"`
	SeveritySummaryByControl   BenchmarkControlsSeverityStatusV2  `json:"severity_summary_by_control"`
	SeveritySummaryByResource  BenchmarkResourcesSeverityStatusV2 `json:"severity_summary_by_resource"`
	ComplianceResultsSummary   ComplianceStatusSummaryV2          `json:"compliance_results_summary"`
	IssuesCount                int                                `json:"issues_count"`
	TopResourcesWithIssues     []TopFiledRecordV2                 `json:"top_resources_with_issues"`
	TopResourceTypesWithIssues []TopFiledRecordV2                 `json:"top_resource_types_with_issues"`
	TopControlsWithIssues      []TopFiledRecordV2                 `json:"top_controls_with_issues"`
	LastEvaluatedAt            *time.Time                         `json:"last_evaluated_at"`
	LastJobStatus              string                             `json:"last_job_status"`
	LastJobId                  string                             `json:"last_job_id"`
}

type ComplianceSummaryOfBenchmarkRequest struct {
	Benchmarks []string `json:"benchmarks"`
	IsRoot     *bool    `json:"is_root"`
	ShowTop    int      `json:"show_top"`
}

type ComplianceSummaryOfBenchmarkResponse struct {
	BenchmarkID                string                             `json:"benchmark_id"`
	BenchmarkTitle             string                             `json:"benchmark_title"`
	ComplianceScore            float64                            `json:"compliance_score"`
	IntegrationTypes           []string                           `json:"connectors"` // Benchmark connectors
	SeveritySummaryByControl   BenchmarkControlsSeverityStatusV2  `json:"severity_summary_by_control"`
	SeveritySummaryByResource  BenchmarkResourcesSeverityStatusV2 `json:"severity_summary_by_resource"`
	SeveritySummaryByIncidents types.SeverityResultV2             `json:"severity_summary_by_incidents"`
	CostImpact                 *float64                           `json:"cost_impact"`
	ComplianceResultsSummary   ComplianceStatusSummaryV2          `json:"compliance_results_summary"`
	IssuesCount                int                                `json:"issues_count"`
	TopIntegrations            []TopIntegration                   `json:"top_integrations"`
	TopResourcesWithIssues     []TopFiledRecordV2                 `json:"top_resources_with_issues"`
	TopResourceTypesWithIssues []TopFiledRecordV2                 `json:"top_resource_types_with_issues"`
	TopControlsWithIssues      []TopFiledRecordV2                 `json:"top_controls_with_issues"`
	LastEvaluatedAt            *time.Time                         `json:"last_evaluated_at"`
	LastJobStatus              string                             `json:"last_job_status"`
	LastJobId                  string                             `json:"last_job_id"`
}

type ListComplianceJobsHistoryItem struct {
	BenchmarkId              string                    `json:"benchmark_id"`
	Integrations             []IntegrationInfo         `json:"integrations"`
	JobId                    string                    `json:"job_id"`
	ComplianceResultsSummary ComplianceStatusSummaryV2 `json:"compliance_results_summary"`
	ComplianceScore          float64                   `json:"compliance_score"`
	TriggerType              string                    `json:"trigger_type"`
	CreatedBy                string                    `json:"created_by"`
	JobStatus                string                    `json:"job_status"`
	CreatedAt                time.Time                 `json:"created_at"`
	UpdatedAt                time.Time                 `json:"updated_at"`
}

type ListComplianceJobsHistoryResponse struct {
	Items      []ListComplianceJobsHistoryItem `json:"items"`
	TotalCount int                             `json:"total_count"`
}

type ListBenchmarksFiltersResponse struct {
	ParentBenchmarkID []string              `json:"parent_benchmark_id"`
	PrimaryResource   []string              `json:"primary_resource"`
	ListOfResources   []string              `json:"list_of_resources"`
	Tags              []BenchmarkTagsResult `json:"tags"`
}

type GetBenchmarkTrendV3Request struct {
	Integration []IntegrationFilter `json:"integration"`
	StartTime   *int64              `json:"start_time"`
	EndTime     *int64              `json:"end_time"`
	Granularity *int64              `json:"granularity"`
}

type BenchmarkTrendDatapointV3 struct {
	Timestamp                time.Time `json:"timestamp"`
	ComplianceResultsSummary *struct {
		Incidents    int `json:"incidents"`
		NonIncidents int `json:"non_incidents"`
	} `json:"compliance_results_summary"`
	IncidentsSeverityBreakdown *types.SeverityResult `json:"incidents_severity_breakdown"`
}

type GetBenchmarkTrendV3Response struct {
	Datapoints    []BenchmarkTrendDatapointV3 `json:"datapoints"`
	MaximumValues BenchmarkTrendDatapointV3   `json:"maximum_values"`
	MinimumValues BenchmarkTrendDatapointV3   `json:"minimum_values"`
}
