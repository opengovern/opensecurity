package types

import (
	"fmt"
	"github.com/opengovern/og-util/pkg/integration"
)

type ComplianceResultDriftEvent struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	ComplianceResultEsID     string           `json:"complianceResultEsID"`
	ParentComplianceJobID    uint             `json:"parentComplianceJobID"`
	ComplianceJobID          uint             `json:"complianceJobID"`
	PreviousComplianceStatus ComplianceStatus `json:"previousComplianceStatus"`
	ComplianceStatus         ComplianceStatus `json:"complianceStatus"`
	PreviousStateActive      bool             `json:"previousStateActive"`
	StateActive              bool             `json:"stateActive"`
	EvaluatedAt              int64            `json:"evaluatedAt"`
	Reason                   string           `json:"reason"`

	BenchmarkID        string                   `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID          string                   `json:"controlID" example:"azure_cis_v140_7_5"`
	IntegrationID      string                   `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	IntegrationType    integration.Type         `json:"integrationType" example:"Azure"`
	Severity           ComplianceResultSeverity `json:"severity" example:"low"`
	PlatformResourceID string                   `json:"platformResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceID         string                   `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceType       string                   `json:"resourceType" example:"Microsoft.Compute/virtualMachines"`
}

func (r ComplianceResultDriftEvent) KeysAndIndex() ([]string, string) {
	return []string{
		r.ComplianceResultEsID,
		fmt.Sprintf("%d", r.ComplianceJobID),
		fmt.Sprintf("%d", r.EvaluatedAt),
	}, ComplianceResultEventsIndex
}

type ComplianceResult struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	BenchmarkID        string                   `json:"benchmarkID" example:"azure_cis_v140"`
	ControlID          string                   `json:"controlID" example:"azure_cis_v140_7_5"`
	IntegrationID      string                   `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"`
	EvaluatedAt        int64                    `json:"evaluatedAt" example:"1589395200"`
	StateActive        bool                     `json:"stateActive" example:"true"`
	ComplianceStatus   ComplianceStatus         `json:"complianceStatus" example:"alarm"`
	Severity           ComplianceResultSeverity `json:"severity" example:"low"`
	IntegrationType    integration.Type         `json:"integrationType" example:"Azure"`
	PlatformResourceID string                   `json:"platformResourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceID         string                   `json:"resourceID" example:"/subscriptions/123/resourceGroups/rg-1/providers/Microsoft.Compute/virtualMachines/vm-1"`
	ResourceName       string                   `json:"resourceName" example:"vm-1"`
	ResourceType       string                   `json:"resourceType" example:"Microsoft.Compute/virtualMachines"`
	Reason             string                   `json:"reason" example:"The VM is not using managed disks"`
	CostImpact         *float64                 `json:"costImpact"`
	ControlPath        string                   `json:"controlPath" example:"aws_cis2/aws_cis2_1/unsecure_http"`
	RunnerID           uint                     `json:"runnerID" example:"1"`
	ComplianceJobID    uint                     `json:"complianceJobID" example:"1"`
	LastUpdatedAt      int64                    `json:"lastUpdatedAt" example:"1589395200"`

	ParentBenchmarks []string `json:"-"`
}

func (r ComplianceResult) KeysAndIndex() ([]string, string) {
	index := ComplianceResultsIndex
	keys := []string{
		r.PlatformResourceID,
		r.ResourceID,
		r.IntegrationID,
		r.ControlID,
		r.BenchmarkID,
	}
	keys = append(keys, r.ParentBenchmarks...)
	return keys, index
}

type ComplianceStatusPerSeverity struct {
	OkCount    SeverityResultWithTotal `json:"okCount"`
	AlarmCount SeverityResultWithTotal `json:"alarmCount"`
	InfoCount  SeverityResultWithTotal `json:"infoCount"`
	SkipCount  SeverityResultWithTotal `json:"skipCount"`
	ErrorCount SeverityResultWithTotal `json:"errorCount"`
}

type ResourceFinding struct {
	EsID    string `json:"es_id"`
	EsIndex string `json:"es_index"`

	PlatformResourceID string           `json:"platformResourceID"`
	ResourceType       string           `json:"resourceType"`
	ResourceName       string           `json:"resourceName"`
	ResourceLocation   string           `json:"resourceLocation"`
	IntegrationType    integration.Type `json:"integrationType"`

	//ComplianceStatusPerSeverity ComplianceStatusPerSeverity `json:"complianceStatusPerSeverity"`
	ComplianceResults []ComplianceResult `json:"complianceResults"`

	ResourceCollection    []string        `json:"resourceCollection" example:"azure_cis_v140_7_5"`
	ResourceCollectionMap map[string]bool `json:"-"` // for creation of the slice only

	JobId       uint  `json:"jobId" example:"1"`
	EvaluatedAt int64 `json:"evaluatedAt" example:"1589395200"`
}

func (r ResourceFinding) KeysAndIndex() ([]string, string) {
	return []string{
		r.PlatformResourceID,
		r.ResourceType,
	}, ResourceFindingsIndex
}
