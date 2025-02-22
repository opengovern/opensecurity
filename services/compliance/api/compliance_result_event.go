package api

import (
	"github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opensecurity/pkg/types"
	"time"
)

type GetSingleResourceFindingResponse struct {
	Resource                    es.Resource                  `json:"resource"`
	ComplianceResultDriftEvents []ComplianceResultDriftEvent `json:"complianceResultDriftEvents"`
	ControlComplianceResults    []ComplianceResult           `json:"controls"`
}

type ComplianceResultDriftEvent struct {
	ID                       string           `json:"id" example:"8e0f8e7a1b1c4e6fb7e49c6af9d2b1c8"`
	ComplianceResultID       string           `json:"complianceResultID"`
	ParentComplianceJobID    uint             `json:"parentComplianceJobID"`
	ComplianceJobID          uint             `json:"complianceJobID"`
	PreviousComplianceStatus ComplianceStatus `json:"previousComplianceStatus"`
	ComplianceStatus         ComplianceStatus `json:"complianceStatus"`
	PreviousStateActive      bool             `json:"previousStateActive"`
	StateActive              bool             `json:"stateActive"`
	EvaluatedAt              time.Time        `json:"evaluatedAt"`
	Reason                   string           `json:"reason"`

	BenchmarkID        string                         `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID          string                         `json:"controlID" example:"azure_cis_v140_7_5"`
	IntegrationID      string                         `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	IntegrationType    integration.Type               `json:"integrationType" example:"Azure"`
	Severity           types.ComplianceResultSeverity `json:"severity" example:"low"`
	PlatformResourceID string                         `json:"platformResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceID         string                         `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceType       string                         `json:"resourceType" example:"Microsoft.Compute/virtualMachines"`

	// Fake fields (won't be stored in ES)
	ResourceTypeName string `json:"resourceTypeName" example:"Virtual Machine"`
	ProviderID       string `json:"providerID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	IntegrationName  string `json:"integrationName" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	ResourceName     string `json:"resourceName" example:"vm-1"`
	ResourceLocation string `json:"resourceLocation" example:"eastus"`

	SortKey []any `json:"sortKey"`
}

func GetAPIComplianceResultDriftEventFromESComplianceResultDriftEvent(complianceResultDriftEvent types.ComplianceResultDriftEvent) ComplianceResultDriftEvent {
	f := ComplianceResultDriftEvent{
		ID:                       complianceResultDriftEvent.EsID,
		ComplianceResultID:       complianceResultDriftEvent.ComplianceResultEsID,
		ParentComplianceJobID:    complianceResultDriftEvent.ParentComplianceJobID,
		ComplianceJobID:          complianceResultDriftEvent.ComplianceJobID,
		PreviousComplianceStatus: "",
		ComplianceStatus:         "",
		PreviousStateActive:      complianceResultDriftEvent.PreviousStateActive,
		StateActive:              complianceResultDriftEvent.StateActive,
		EvaluatedAt:              time.UnixMilli(complianceResultDriftEvent.EvaluatedAt),
		Reason:                   complianceResultDriftEvent.Reason,

		BenchmarkID:        complianceResultDriftEvent.BenchmarkID,
		ControlID:          complianceResultDriftEvent.ControlID,
		IntegrationID:      complianceResultDriftEvent.IntegrationID,
		IntegrationType:    complianceResultDriftEvent.IntegrationType,
		Severity:           complianceResultDriftEvent.Severity,
		PlatformResourceID: complianceResultDriftEvent.PlatformResourceID,
		ResourceID:         complianceResultDriftEvent.ResourceID,
		ResourceType:       complianceResultDriftEvent.ResourceType,
	}
	if complianceResultDriftEvent.PreviousComplianceStatus.IsPassed() {
		f.PreviousComplianceStatus = ComplianceStatusPassed
	} else {
		f.PreviousComplianceStatus = ComplianceStatusFailed
	}
	if complianceResultDriftEvent.ComplianceStatus.IsPassed() {
		f.ComplianceStatus = ComplianceStatusPassed
	} else {
		f.ComplianceStatus = ComplianceStatusFailed
	}

	return f
}

type ComplianceResultDriftEventFilters struct {
	IntegrationType    []integration.Type               `json:"integrationType" example:"Azure"`
	ResourceType       []string                         `json:"resourceType" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines"`
	IntegrationID      []string                         `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	NotIntegrationID   []string                         `json:"notIntegrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	IntegrationGroup   []string                         `json:"integrationGroup" example:"active"`
	BenchmarkID        []string                         `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID          []string                         `json:"controlID" example:"azure_cis_v140_7_5"`
	Severity           []types.ComplianceResultSeverity `json:"severity" example:"low"`
	ComplianceStatus   []ComplianceStatus               `json:"complianceStatus" example:"alarm"`
	StateActive        []bool                           `json:"stateActive" example:"true"`
	ComplianceResultID []string                         `json:"complianceResultID" example:"8e0f8e7a1b1c4e6fb7e49c6af9d2b1c8"`
	PlatformResourceID []string                         `json:"platformResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	EvaluatedAt        struct {
		From *int64 `json:"from"`
		To   *int64 `json:"to"`
	} `json:"evaluatedAt"`
}

type ComplianceResultDriftEventFiltersWithMetadata struct {
	Connector          []FilterWithMetadata `json:"connector"`
	BenchmarkID        []FilterWithMetadata `json:"benchmarkID"`
	ControlID          []FilterWithMetadata `json:"controlID"`
	ResourceTypeID     []FilterWithMetadata `json:"resourceTypeID"`
	IntegrationID      []FilterWithMetadata `json:"integrationID"`
	ResourceCollection []FilterWithMetadata `json:"resourceCollection"`
	Severity           []FilterWithMetadata `json:"severity"`
	ComplianceStatus   []FilterWithMetadata `json:"complianceStatus"`
	StateActive        []FilterWithMetadata `json:"stateActive"`
}

type ComplianceResultDriftEventsSort struct {
	IntegrationType    *SortDirection `json:"integrationType"`
	PlatformResourceID *SortDirection `json:"platformResourceID"`
	ResourceType       *SortDirection `json:"resourceType"`
	IntegrationID      *SortDirection `json:"integrationID"`
	BenchmarkID        *SortDirection `json:"benchmarkID"`
	ControlID          *SortDirection `json:"controlID"`
	Severity           *SortDirection `json:"severity"`
	ComplianceStatus   *SortDirection `json:"complianceStatus"`
	StateActive        *SortDirection `json:"stateActive"`
}

type GetComplianceResultDriftEventsRequest struct {
	Filters      ComplianceResultDriftEventFilters `json:"filters"`
	Sort         []ComplianceResultDriftEventsSort `json:"sort"`
	Limit        int                               `json:"limit" example:"100"`
	AfterSortKey []any                             `json:"afterSortKey"`
}

type GetComplianceResultDriftEventsResponse struct {
	ComplianceResultDriftEvents []ComplianceResultDriftEvent `json:"complianceResultDriftEvents"`
	TotalCount                  int64                        `json:"totalCount" example:"100"`
}

type CountComplianceResultDriftEventsResponse struct {
	Count int64 `json:"count"`
}
