package api

import (
	"github.com/opengovern/og-util/pkg/es"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/source"
	"github.com/opengovern/opengovernance/pkg/types"
	"time"
)

type GetSingleResourceFindingResponse struct {
	Resource                    es.Resource                  `json:"resource"`
	ComplianceResultDriftEvents []ComplianceResultDriftEvent `json:"complianceResultDriftEvents"`
	ControlComplianceResults    []ComplianceResult           `json:"controls"`
}

type ComplianceResultDriftEvent struct {
	ID                        string            `json:"id" example:"8e0f8e7a1b1c4e6fb7e49c6af9d2b1c8"`
	ComplianceResultID        string            `json:"complianceResultID"`
	ParentComplianceJobID     uint              `json:"parentComplianceJobID"`
	ComplianceJobID           uint              `json:"complianceJobID"`
	PreviousConformanceStatus ConformanceStatus `json:"previousConformanceStatus"`
	ConformanceStatus         ConformanceStatus `json:"conformanceStatus"`
	PreviousStateActive       bool              `json:"previousStateActive"`
	StateActive               bool              `json:"stateActive"`
	EvaluatedAt               time.Time         `json:"evaluatedAt"`
	Reason                    string            `json:"reason"`

	BenchmarkID               string                         `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID                 string                         `json:"controlID" example:"azure_cis_v140_7_5"`
	IntegrationID             string                         `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	IntegrationType           integration.Type               `json:"integrationType" example:"Azure"`
	Severity                  types.ComplianceResultSeverity `json:"severity" example:"low"`
	OpenGovernanceResourceID  string                         `json:"opengovernanceResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceID                string                         `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceType              string                         `json:"resourceType" example:"Microsoft.Compute/virtualMachines"`
	ParentBenchmarkReferences []string                       `json:"parentBenchmarkReferences"`

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
		ID:                        complianceResultDriftEvent.EsID,
		ComplianceResultID:        complianceResultDriftEvent.ComplianceResultEsID,
		ParentComplianceJobID:     complianceResultDriftEvent.ParentComplianceJobID,
		ComplianceJobID:           complianceResultDriftEvent.ComplianceJobID,
		PreviousConformanceStatus: "",
		ConformanceStatus:         "",
		PreviousStateActive:       complianceResultDriftEvent.PreviousStateActive,
		StateActive:               complianceResultDriftEvent.StateActive,
		EvaluatedAt:               time.UnixMilli(complianceResultDriftEvent.EvaluatedAt),
		Reason:                    complianceResultDriftEvent.Reason,

		BenchmarkID:               complianceResultDriftEvent.BenchmarkID,
		ControlID:                 complianceResultDriftEvent.ControlID,
		IntegrationID:             complianceResultDriftEvent.IntegrationID,
		IntegrationType:           complianceResultDriftEvent.IntegrationType,
		Severity:                  complianceResultDriftEvent.Severity,
		OpenGovernanceResourceID:  complianceResultDriftEvent.OpenGovernanceResourceID,
		ResourceID:                complianceResultDriftEvent.ResourceID,
		ResourceType:              complianceResultDriftEvent.ResourceType,
		ParentBenchmarkReferences: complianceResultDriftEvent.ParentBenchmarkReferences,
	}
	if complianceResultDriftEvent.PreviousConformanceStatus.IsPassed() {
		f.PreviousConformanceStatus = ConformanceStatusPassed
	} else {
		f.PreviousConformanceStatus = ConformanceStatusFailed
	}
	if complianceResultDriftEvent.ConformanceStatus.IsPassed() {
		f.ConformanceStatus = ConformanceStatusPassed
	} else {
		f.ConformanceStatus = ConformanceStatusFailed
	}

	return f
}

type ComplianceResultDriftEventFilters struct {
	Connector                []source.Type                    `json:"connector" example:"Azure"`
	ResourceType             []string                         `json:"resourceType" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines"`
	IntegrationID            []string                         `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	NotIntegrationID         []string                         `json:"notIntegrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	ConnectionGroup          []string                         `json:"connectionGroup" example:"healthy"`
	BenchmarkID              []string                         `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID                []string                         `json:"controlID" example:"azure_cis_v140_7_5"`
	Severity                 []types.ComplianceResultSeverity `json:"severity" example:"low"`
	ConformanceStatus        []ConformanceStatus              `json:"conformanceStatus" example:"alarm"`
	StateActive              []bool                           `json:"stateActive" example:"true"`
	ComplianceResultID       []string                         `json:"complianceResultID" example:"8e0f8e7a1b1c4e6fb7e49c6af9d2b1c8"`
	OpenGovernanceResourceID []string                         `json:"opengovernanceResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	EvaluatedAt              struct {
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
	ConformanceStatus  []FilterWithMetadata `json:"conformanceStatus"`
	StateActive        []FilterWithMetadata `json:"stateActive"`
}

type ComplianceResultDriftEventsSort struct {
	Connector                *SortDirection `json:"connector"`
	OpenGovernanceResourceID *SortDirection `json:"opengovernanceResourceID"`
	ResourceType             *SortDirection `json:"resourceType"`
	IntegrationID            *SortDirection `json:"integrationID"`
	BenchmarkID              *SortDirection `json:"benchmarkID"`
	ControlID                *SortDirection `json:"controlID"`
	Severity                 *SortDirection `json:"severity"`
	ConformanceStatus        *SortDirection `json:"conformanceStatus"`
	StateActive              *SortDirection `json:"stateActive"`
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
