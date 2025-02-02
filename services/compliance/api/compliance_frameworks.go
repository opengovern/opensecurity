package api

import (
	"github.com/opengovern/opencomply/pkg/types"
	"time"
)

type FrameworkAssignmentAssignmentType string

const (
	FrameworkAssignmentAssignmentTypeExplicit FrameworkAssignmentAssignmentType = "explicit"
	FrameworkAssignmentAssignmentTypeImplicit FrameworkAssignmentAssignmentType = "implicit"
	FrameworkAssignmentAssignmentTypeNone     FrameworkAssignmentAssignmentType = "none"
)

type PageInfo struct {
	CurrentPage int64 `json:"current_page"`
	PageSize    int64 `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int64 `json:"total_pages"`
}

type ListFrameworkAssignmentsResponseData struct {
	IntegrationID         string                            `json:"integration_id"`
	IntegrationName       string                            `json:"integration_name"`
	IntegrationProviderID string                            `json:"integration_provider_id"`
	PluginID              string                            `json:"plugin_id"`
	AssignmentType        FrameworkAssignmentAssignmentType `json:"assignment_type"`
}

type ListFrameworkAssignmentsResponse struct {
	Data     []ListFrameworkAssignmentsResponseData `json:"data"`
	PageInfo PageInfo                               `json:"page_info"`
}

type UpdateFrameworkSettingRequest struct {
	IsBaseline *bool `json:"is_baseline"`
	Enabled    *bool `json:"enabled"`
}

type FrameworkCoverage struct {
	FrameworkID      string   `json:"framework_id"`
	PrimaryResources []string `json:"primary_resources"`
	ListOfResources  []string `json:"list_of_resources"`
	Controls         []string `json:"controls"`
}

type FrameworkItem struct {
	FrameworkID                string                             `json:"framework_id"`
	FrameworkTitle             string                             `json:"framework_title"`
	ComplianceScore            float64                            `json:"compliance_score"`
	Plugins                    []string                           `json:"plugins"`
	SeveritySummaryByControl   BenchmarkControlsSeverityStatusV2  `json:"severity_summary_by_control"`
	SeveritySummaryByResource  BenchmarkResourcesSeverityStatusV2 `json:"severity_summary_by_resource"`
	SeveritySummaryByIncidents types.SeverityResultV2             `json:"severity_summary_by_incidents"`
	ComplianceResultsSummary   ComplianceStatusSummaryV2          `json:"compliance_results_summary"`
	NoOfTotalAssignments       int64                              `json:"no_of_total_assignments"`
	IssuesCount                int                                `json:"issues_count"`
	IsBaseline                 bool                               `json:"is_baseline"`
	Enabled                    bool                               `json:"enabled"`
	LastEvaluatedAt            *time.Time                         `json:"last_evaluated_at"`
}
