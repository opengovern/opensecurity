package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/pkg/types"
	"strings"
	"time"
)

type ComplianceResultFilters struct {
	JobID            []string                         `json:"jobID"`
	IntegrationType  []string                         `json:"integrationType" example:"Azure"`
	ResourceID       []string                         `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceTypeID   []string                         `json:"resourceTypeID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines"`
	IntegrationID    []string                         `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	NotIntegrationID []string                         `json:"notIntegrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	IntegrationGroup []string                         `json:"integrationGroup" example:"active"`
	BenchmarkID      []string                         `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID        []string                         `json:"controlID" example:"azure_cis_v140_7_5"`
	Severity         []types.ComplianceResultSeverity `json:"severity" example:"low"`
	ComplianceStatus []ComplianceStatus               `json:"complianceStatus" example:"alarm"`
	StateActive      []bool                           `json:"stateActive" example:"true"`
	LastEvent        struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	} `json:"lastEvent"`
	EvaluatedAt struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	} `json:"evaluatedAt"`
	Interval *string `json:"interval" example:"5m"`
}

type ComplianceResultSummaryFilters struct {
	IntegrationID  []string `json:"integrationID"`
	ResourceTypeID []string `json:"resourceTypeID"`
}

type ComplianceResultFiltersWithMetadata struct {
	IntegrationType  []FilterWithMetadata `json:"integrationType"`
	BenchmarkID      []FilterWithMetadata `json:"benchmarkID"`
	ControlID        []FilterWithMetadata `json:"controlID"`
	ResourceTypeID   []FilterWithMetadata `json:"resourceTypeID"`
	IntegrationID    []FilterWithMetadata `json:"integrationID"`
	Severity         []FilterWithMetadata `json:"severity"`
	ComplianceStatus []FilterWithMetadata `json:"complianceStatus"`
	StateActive      []FilterWithMetadata `json:"stateActive"`
}

type ComplianceResultsSort struct {
	IntegrationType    *SortDirection `json:"integrationType"`
	ResourceID         *SortDirection `json:"resourceID"`
	PlatformResourceID *SortDirection `json:"platformResourceID"`
	ResourceTypeID     *SortDirection `json:"resourceTypeID"`
	IntegrationID      *SortDirection `json:"integrationID"`
	BenchmarkID        *SortDirection `json:"benchmarkID"`
	ControlID          *SortDirection `json:"controlID"`
	Severity           *SortDirection `json:"severity"`
	ComplianceStatus   *SortDirection `json:"complianceStatus"`
	StateActive        *SortDirection `json:"stateActive"`
}

type GetComplianceResultsRequest struct {
	Filters      ComplianceResultFilters `json:"filters"`
	Sort         []ComplianceResultsSort `json:"sort"`
	Limit        int                     `json:"limit" example:"100"`
	AfterSortKey []any                   `json:"afterSortKey"`
}

type GetSingleComplianceResultRequest struct {
	BenchmarkID        string `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID          string `json:"controlID" example:"azure_cis_v140_7_5"`
	IntegrationID      string `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	PlatformResourceID string `json:"platformResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceID         string `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
}

type GetSingleResourceFindingRequest struct {
	PlatformResourceID string  `json:"platformResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceType       *string `json:"resourceType" example:"Microsoft.Compute/virtualMachines"`
}

type ComplianceStatus string

const (
	ComplianceStatusFailed ComplianceStatus = "failed"
	ComplianceStatusPassed ComplianceStatus = "passed"
)

func ListComplianceStatuses() []ComplianceStatus {
	return []ComplianceStatus{ComplianceStatusFailed, ComplianceStatusPassed}
}

func (cs ComplianceStatus) GetEsComplianceStatuses() []types.ComplianceStatus {
	switch cs {
	case ComplianceStatusFailed:
		return types.GetFailedComplianceStatuses()
	case ComplianceStatusPassed:
		return types.GetPassedComplianceStatuses()
	}
	return nil
}

func ParseComplianceStatuses(complianceStatuses []string) []ComplianceStatus {
	var result []ComplianceStatus
	for _, cs := range complianceStatuses {
		switch strings.ToLower(cs) {
		case strings.ToLower(string(ComplianceStatusFailed)):
			result = append(result, ComplianceStatusFailed)
		case strings.ToLower(string(ComplianceStatusPassed)):
			result = append(result, ComplianceStatusPassed)
		}
	}
	return result
}

type ComplianceResult struct {
	ID                 string                         `json:"id" example:"1"`
	BenchmarkID        string                         `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID          string                         `json:"controlID" example:"azure_cis_v140_7_5"`
	IntegrationID      string                         `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	EvaluatedAt        int64                          `json:"evaluatedAt" example:"1589395200"`
	StateActive        bool                           `json:"stateActive" example:"true"`
	ComplianceStatus   ComplianceStatus               `json:"complianceStatus" example:"alarm"`
	Severity           types.ComplianceResultSeverity `json:"severity" example:"low"`
	IntegrationType    integration.Type               `json:"integrationType" example:"Azure"`
	PlatformResourceID string                         `json:"platformResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceID         string                         `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceName       string                         `json:"resourceName" example:"vm-1"`
	ResourceType       string                         `json:"resourceType" example:"Microsoft.Compute/virtualMachines"`
	Reason             string                         `json:"reason" example:"The VM is not using managed disks"`
	CostImpact         *float64                       `json:"costImpact" example:"0.5"`
	RunnerID           uint                           `json:"runnerID" example:"1"`
	ComplianceJobID    uint                           `json:"complianceJobID" example:"1"`
	ControlPath        string                         `json:"controlPath" example:"aws_cis2/aws_cis2_1/unsecure_http"`
	LastEvent          time.Time                      `json:"lastEvent" example:"1589395200"`

	ResourceTypeName     string   `json:"resourceTypeName" example:"Virtual Machine"`
	ParentBenchmarkNames []string `json:"parentBenchmarkNames" example:"Azure CIS v1.4.0"`
	ControlTitle         string   `json:"controlTitle"`
	ProviderID           string   `json:"providerID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`      // Connection ID
	IntegrationName      string   `json:"integrationName" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"` // Connection ID

	SortKey []any `json:"sortKey"`
}

func GetAPIComplianceResultFromESComplianceResult(complianceResult types.ComplianceResult) ComplianceResult {
	f := ComplianceResult{
		ID:                 complianceResult.EsID,
		BenchmarkID:        complianceResult.BenchmarkID,
		ControlID:          complianceResult.ControlID,
		IntegrationID:      complianceResult.IntegrationID,
		EvaluatedAt:        complianceResult.EvaluatedAt,
		StateActive:        complianceResult.StateActive,
		ComplianceStatus:   "",
		Severity:           complianceResult.Severity,
		IntegrationType:    complianceResult.IntegrationType,
		PlatformResourceID: complianceResult.PlatformResourceID,
		ResourceID:         complianceResult.ResourceID,
		ResourceName:       complianceResult.ResourceName,
		ResourceType:       complianceResult.ResourceType,
		Reason:             complianceResult.Reason,
		CostImpact:         complianceResult.CostImpact,
		RunnerID:           complianceResult.RunnerID,
		ComplianceJobID:    complianceResult.ComplianceJobID,
		ControlPath:        complianceResult.ControlPath,
		LastEvent:          time.UnixMilli(complianceResult.LastUpdatedAt),
	}
	if complianceResult.ComplianceStatus.IsPassed() {
		f.ComplianceStatus = ComplianceStatusPassed
	} else {
		f.ComplianceStatus = ComplianceStatusFailed
	}
	if f.ResourceType == "" {
		f.ResourceType = "Unknown"
		f.ResourceTypeName = "Unknown"
	}

	return f
}

type GetComplianceResultsResponse struct {
	ComplianceResults []ComplianceResult `json:"complianceResults"`
	TotalCount        int64              `json:"totalCount" example:"100"`
}

type ComplianceResultKPIResponse struct {
	FailedComplianceResultsCount int64 `json:"failedComplianceResultsCount"`
	FailedResourceCount          int64 `json:"failedResourceCount"`
	FailedControlCount           int64 `json:"failedControlCount"`
	FailedIntegrationCount       int64 `json:"failedIntegrationCount"`
}

type ServiceComplianceResultsSummary struct {
	ServiceName     string  `json:"serviceName"`
	ServiceLabel    string  `json:"serviceLabel"`
	SecurityScore   float64 `json:"securityScore"`
	SeveritiesCount struct {
		Critical int `json:"critical"`
		High     int `json:"high"`
		Medium   int `json:"medium"`
		Low      int `json:"low"`
		None     int `json:"none"`
	} `json:"severitiesCount"`
	ComplianceStatusesCount struct {
		Passed int `json:"passed"`
		Failed int `json:"failed"`
	} `json:"complianceStatusesCount"`
}

type GetServicesComplianceResultsSummaryResponse struct {
	Services []ServiceComplianceResultsSummary `json:"services"`
}

type GetFieldCountResponse struct {
	Controls []struct {
		ControlName string           `json:"controlName"`
		FieldCounts []TopFieldRecord `json:"fieldCounts"`
	} `json:"controls"`
}

type GetTopFieldResponse struct {
	TotalCount int              `json:"totalCount" example:"100"`
	Records    []TopFieldRecord `json:"records"`
}

type CountComplianceResultsResponse struct {
	Count int64 `json:"count"`
}

type GetComplianceResultDriftEventsByComplianceResultIDResponse struct {
	ComplianceResultDriftEvents []ComplianceResultDriftEvent `json:"complianceResultDriftEvents"`
}

type ComplianceResultFiltersV2 struct {
	Integration     []IntegrationInfoFilter          `json:"integration"`
	BenchmarkID     []string                         `json:"benchmark_id"`
	NotBenchmarkID  []string                         `json:"not_benchmark_id"`
	ControlID       []string                         `json:"control_id"`
	NotControlID    []string                         `json:"not_control_id"`
	Severity        []types.ComplianceResultSeverity `json:"severity"`
	NotSeverity     []types.ComplianceResultSeverity `json:"not_severity"`
	ResourceID      []string                         `json:"resource_id"`
	NotResourceID   []string                         `json:"not_resource_id"`
	IsCompliant     []bool                           `json:"is_compliant"`
	IsActive        []bool                           `json:"is_active"`
	ResourceType    []string                         `json:"resource_type"`
	NotResourceType []string                         `json:"not_resource_type"`
	LastUpdated     struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	} `json:"last_event"`
	NotLastUpdated struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	} `json:"not_last_event"`
}

type IntegrationInfoFilter struct {
	IntegrationType *string `json:"integration_type"`
	ProviderID      *string `json:"provider_id"`
	Name            *string `json:"name"`
	IntegrationID   *string `json:"integration_id"`
}

type ComplianceResultsSortV2 struct {
	BenchmarkID      *SortDirection `json:"benchmark_id"`
	ControlID        *SortDirection `json:"control_id"`
	Severity         *SortDirection `json:"severity"`
	ComplianceStatus *SortDirection `json:"compliance_status"`
	ResourceType     *SortDirection `json:"resource_type"`
	LastUpdated      *SortDirection `json:"last_updated"`
}

type GetComplianceResultsRequestV2 struct {
	Filters      ComplianceResultFiltersV2 `json:"filters"`
	Sort         []ComplianceResultsSortV2 `json:"sort"`
	Limit        int                       `json:"limit" example:"100"`
	AfterSortKey []any                     `json:"afterSortKey"`
}
