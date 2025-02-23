package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opensecurity/pkg/types"
	"time"
)

type ResourceFinding struct {
	ID                 string           `json:"id"`
	PlatformResourceID string           `json:"platformResourceID"`
	ResourceName       string           `json:"resourceName"`
	ResourceLocation   string           `json:"resourceLocation"`
	ResourceType       string           `json:"resourceType"`
	ResourceTypeLabel  string           `json:"resourceTypeLabel"`
	IntegrationType    integration.Type `json:"integrationType"`
	ComplianceJobID    string           `json:"complianceJobID"`

	FailedCount int `json:"failedCount"`
	TotalCount  int `json:"totalCount"`

	EvaluatedAt time.Time `json:"evaluatedAt"`

	ComplianceResults []ComplianceResult `json:"complianceResults"`

	SortKey []any `json:"sortKey"`

	IntegrationID   string `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`   // Connection ID
	ProviderID      string `json:"providerID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`      // Connection ID
	IntegrationName string `json:"integrationName" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"` // Connection ID
}

func GetAPIResourceFinding(resourceFinding types.ResourceFinding) ResourceFinding {
	apiRf := ResourceFinding{
		ID:                 resourceFinding.EsID,
		PlatformResourceID: resourceFinding.PlatformResourceID,
		ResourceName:       resourceFinding.ResourceName,
		ResourceLocation:   resourceFinding.ResourceLocation,
		ResourceType:       resourceFinding.ResourceType,
		IntegrationType:    resourceFinding.IntegrationType,

		FailedCount: 0,
		TotalCount:  len(resourceFinding.ComplianceResults),

		EvaluatedAt: time.UnixMilli(resourceFinding.EvaluatedAt),

		ComplianceResults: nil,
	}

	integrationIDs := make(map[string]bool)

	for _, complianceResult := range resourceFinding.ComplianceResults {
		if !complianceResult.ComplianceStatus.IsPassed() {
			apiRf.FailedCount++
		}
		integrationIDs[complianceResult.IntegrationID] = true
		apiRf.IntegrationID = complianceResult.IntegrationID
		apiRf.ComplianceResults = append(apiRf.ComplianceResults, GetAPIComplianceResultFromESComplianceResult(complianceResult))
	}

	if len(integrationIDs) > 1 {
		apiRf.IntegrationID = "Global (Multiple)"
		apiRf.ProviderID = "Global (Multiple)"
		apiRf.IntegrationName = "Global (Multiple)"
	}

	return apiRf
}

type ResourceFindingFilters struct {
	ComplianceJobId    []string                         `json:"compliance_job_id"`
	IntegrationType    []integration.Type               `json:"integrationType" example:"Azure"`
	ResourceID         []string                         `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceTypeID     []string                         `json:"resourceTypeID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines"`
	IntegrationID      []string                         `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	NotIntegrationID   []string                         `json:"notIntegrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	IntegrationGroup   []string                         `json:"integrationGroup" example:"healthy"`
	ResourceCollection []string                         `json:"resourceCollection" example:"example-rc"`
	BenchmarkID        []string                         `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID          []string                         `json:"controlID" example:"azure_cis_v140_7_5"`
	Severity           []types.ComplianceResultSeverity `json:"severity" example:"low"`
	ComplianceStatus   []ComplianceStatus               `json:"complianceStatus" example:"alarm"`
	EvaluatedAt        struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	}
	Interval *string `json:"interval"`
}

type ResourceFindingsSort struct {
	PlatformResourceID *SortDirection `json:"platformResourceID"`
	ResourceType       *SortDirection `json:"resourceType"`
	ResourceName       *SortDirection `json:"resourceName"`
	ResourceLocation   *SortDirection `json:"resourceLocation"`
	FailedCount        *SortDirection `json:"failedCount"`
	ComplianceStatus   *SortDirection `json:"complianceStatus"`
}

type ListResourceFindingsRequest struct {
	Filters      ResourceFindingFilters `json:"filters"`
	Sort         []ResourceFindingsSort `json:"sort"`
	Limit        int                    `json:"limit" example:"100"`
	AfterSortKey []any                  `json:"afterSortKey"`
}

type ListResourceFindingsResponse struct {
	TotalCount       int               `json:"totalCount"`
	ResourceFindings []ResourceFinding `json:"resourceFindings"`
}
