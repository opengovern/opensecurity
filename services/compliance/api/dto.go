package api

import (
	"time"

	"github.com/opengovern/og-util/pkg/integration"
	coreClient "github.com/opengovern/opencomply/services/core/api"
	integrationapi "github.com/opengovern/opencomply/services/integration/api/models"
)

type PolicyLanguage string

const (
	PolicyLanguageSQL  PolicyLanguage = "sql"
	PolicyLanguageRego PolicyLanguage = "rego"
)

type BenchmarkAssignment struct {
	BenchmarkId          string    `json:"benchmarkId" example:"azure_cis_v140"`                         // Benchmark ID
	IntegrationId        *string   `json:"integrationId" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"` // Connection ID
	ResourceCollectionId *string   `json:"resourceCollectionId" example:"example-rc"`                    // Resource Collection ID
	AssignedAt           time.Time `json:"assignedAt"`                                                   // Unix timestamp
}

type AssignedBenchmark struct {
	Benchmark Benchmark `json:"benchmarkId"`
	Status    bool      `json:"status" example:"true"` // Status
}

type BenchmarkAssignedIntegration struct {
	IntegrationID   string           `json:"integrationID" example:"8e0f8e7a-1b1c-4e6f-b7e4-9c6af9d2b1c8"` // Connection ID
	ProviderID      string           `json:"providerID" example:"1283192749"`                              // Provider Connection ID
	IntegrationName string           `json:"integrationName"`                                              // Provider Connection Name
	IntegrationType integration.Type `json:"integrationType" example:"Azure"`                              // Clout Provider
	Status          bool             `json:"status" example:"true"`                                        // Status
}

type BenchmarkAssignedResourceCollection struct {
	ResourceCollectionID   string `json:"resourceCollectionID"`   // Resource Collection ID
	ResourceCollectionName string `json:"resourceCollectionName"` // Resource Collection Name
	Status                 bool   `json:"status" example:"true"`  // Status
}

type BenchmarkAssignedEntities struct {
	Integrations []BenchmarkAssignedIntegration `json:"integrations"`
}

type TopFieldRecord struct {
	Integration  *integrationapi.Integration
	ResourceType *coreClient.ResourceType
	Control      *Control
	Service      *string

	Field *string `json:"field"`

	ControlCount      *int `json:"controlCount,omitempty"`
	ControlTotalCount *int `json:"controlTotalCount,omitempty"`

	ResourceCount      *int `json:"resourceCount,omitempty"`
	ResourceTotalCount *int `json:"resourceTotalCount,omitempty"`

	Count      int `json:"count"`
	TotalCount int `json:"totalCount"`
}

type BenchmarkRemediation struct {
	Remediation string `json:"remediation"`
}

type AccountsComplianceResultsSummary struct {
	AccountName     string  `json:"accountName"`
	AccountId       string  `json:"accountId"`
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
		Error  int `json:"error"`
		Info   int `json:"info"`
		Skip   int `json:"skip"`
	} `json:"complianceStatusesCount"`
	LastCheckTime time.Time `json:"lastCheckTime"`
}

type GetAccountsComplianceResultsSummaryResponse struct {
	Accounts []AccountsComplianceResultsSummary `json:"accounts"`
}

type SortDirection string

const (
	SortDirectionAscending  SortDirection = "asc"
	SortDirectionDescending SortDirection = "desc"
)

type FilterWithMetadata struct {
	Key         string `json:"key" example:"key"`
	DisplayName string `json:"displayName" example:"displayName"`
	Count       *int   `json:"count" example:"10"`
}

type QueryParameter struct {
	Key      string `json:"key" example:"key"`
	Required bool   `json:"required" example:"true"`
}

type Policy struct {
	ID              string         `json:"id" example:"azure_ad_manual_control"`
	Language        PolicyLanguage `json:"language" example:"sql"`
	Definition      string         `json:"definition" example:"select\n  -- Required Columns\n  'active_directory' as resource,\n  'info' as status,\n  'Manual verification required.' as reason;\n"`
	IntegrationType []string       `json:"integrationType" example:"Azure"`
	PrimaryResource *string        `json:"primaryResource" example:"null"`
	ListOfResources []string       `json:"listOfResources" example:"null"`

	// CloudQL Fields
	Parameters []QueryParameter `json:"parameters"`

	// Rego Fields
	RegoPolicies []string `json:"regoPolicies" example:"null"`

	CreatedAt time.Time `json:"createdAt" example:"2023-06-07T14:00:15.677558Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2023-06-16T14:58:08.759554Z"`
}
