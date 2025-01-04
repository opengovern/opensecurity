package api

import (
	"time"

	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/pkg/types"
	inventoryApi "github.com/opengovern/opencomply/services/inventory/api"
)

type Control struct {
	ID          string              `json:"id" example:"azure_cis_v140_1_1"`
	Title       string              `json:"title" example:"1.1 Ensure that multi-factor authentication status is enabled for all privileged users"`
	Description string              `json:"description" example:"Enable multi-factor authentication for all user credentials who have write access to Azure resources. These include roles like 'Service Co-Administrators', 'Subscription Owners', 'Contributors'."`
	Tags        map[string][]string `json:"tags"`

	Explanation       string `json:"explanation" example:"Multi-factor authentication adds an additional layer of security by requiring users to enter a code from a mobile device or phone in addition to their username and password when signing into Azure."`
	NonComplianceCost string `json:"nonComplianceCost" example:"Non-compliance to this control could result in several costs including..."`
	UsefulExample     string `json:"usefulExample" example:"Access to resources must be closely controlled to prevent malicious activity like data theft..."`

	CliRemediation          string `json:"cliRemediation" example:"To enable multi-factor authentication for a user, run the following command..."`
	ManualRemediation       string `json:"manualRemediation" example:"To enable multi-factor authentication for a user, run the following command..."`
	GuardrailRemediation    string `json:"guardrailRemediation" example:"To enable multi-factor authentication for a user, run the following command..."`
	ProgrammaticRemediation string `json:"programmaticRemediation" example:"To enable multi-factor authentication for a user, run the following command..."`

	IntegrationType    []string                       `json:"integration_type" example:"Azure"`
	Enabled            bool                           `json:"enabled" example:"true"`
	DocumentURI        string                         `json:"documentURI" example:"benchmarks/azure_cis_v140_1_1.md"`
	Policy             *Policy                        `json:"policy"`
	Severity           types.ComplianceResultSeverity `json:"severity" example:"low"`
	ManualVerification bool                           `json:"manualVerification" example:"true"`
	Managed            bool                           `json:"managed" example:"true"`
	CreatedAt          time.Time                      `json:"createdAt" example:"2020-01-01T00:00:00Z"`
	UpdatedAt          time.Time                      `json:"updatedAt" example:"2020-01-01T00:00:00Z"`
}

type ControlSummary struct {
	Control      Control                    `json:"control"`
	ResourceType *inventoryApi.ResourceType `json:"resourceType"`

	Benchmarks []Benchmark `json:"benchmarks"`

	Passed                 bool      `json:"passed"`
	FailedResourcesCount   int       `json:"failedResourcesCount"`
	TotalResourcesCount    int       `json:"totalResourcesCount"`
	FailedIntegrationCount int       `json:"failedIntegrationCount"`
	TotalIntegrationCount  int       `json:"totalIntegrationCount"`
	CostImpact             *float64  `json:"costImpact"`
	EvaluatedAt            time.Time `json:"evaluatedAt"`
}

type ControlTrendDatapoint struct {
	Timestamp              int `json:"timestamp"` // Time
	FailedResourcesCount   int `json:"failedResourcesCount"`
	TotalResourcesCount    int `json:"totalResourcesCount"`
	FailedIntegrationCount int `json:"failedIntegrationCount"`
	TotalIntegrationCount  int `json:"totalIntegrationCount"`
}

type ControlsFilterSummaryRequest struct {
	IntegrationTypes        []string                 `json:"integration_types"`
	Severity                []string                 `json:"severity"`
	RootBenchmark           []string                 `json:"root_benchmark"`
	ParentBenchmark         []string                 `json:"parent_benchmark"`
	HasParameters           *bool                    `json:"has_parameters"`
	PrimaryResource         []string                 `json:"primary_resource"`
	ListOfResources         []string                 `json:"list_of_resources"`
	Tags                    map[string][]string      `json:"tags"`
	TagsRegex               *string                  `json:"tags_regex"`
	ComplianceResultFilters *ComplianceResultFilters `json:"compliance_result_filters"`
}

type ListControlsFilterRequest struct {
	IntegrationTypes        []string                 `json:"integration_types"`
	Severity                []string                 `json:"severity"`
	RootBenchmark           []string                 `json:"root_benchmark"`
	ParentBenchmark         []string                 `json:"parent_benchmark"`
	HasParameters           *bool                    `json:"has_parameters"`
	PrimaryResource         []string                 `json:"primary_resource"`
	ListOfResources         []string                 `json:"list_of_resources"`
	Tags                    map[string][]string      `json:"tags"`
	TagsRegex               *string                  `json:"tags_regex"`
	ComplianceResultFilters *ComplianceResultFilters `json:"compliance_result_filters"`
	ComplianceResultSummary bool                     `json:"compliance_result_summary"`
	SortBy                  string                   `json:"sort_by"`
	SortOrder               string                   `json:"sort_order"`
	Cursor                  *int64                   `json:"cursor"`
	PerPage                 *int64                   `json:"per_page"`
}

type ListControlsFilterResponse struct {
	Items      []ListControlsFilterResultControl `json:"items"`
	TotalCount int                               `json:"total_count"`
}

type ListControlsFilterResultControl struct {
	ID              string                         `json:"id"`
	Title           string                         `json:"title"`
	Description     string                         `json:"description"`
	IntegrationType []integration.Type             `json:"integration_type"`
	Severity        types.ComplianceResultSeverity `json:"severity"`
	Tags            map[string][]string            `json:"tags"`
	Policy          struct {
		PrimaryResource string           `json:"primary_resource"`
		ListOfResources []string         `json:"list_of_resources"`
		Parameters      []QueryParameter `json:"parameters"`
	} `json:"policy"`
	ComplianceResultsSummary struct {
		IncidentCount         int64    `json:"incident_count"`
		NonIncidentCount      int64    `json:"non_incident_count"`
		NonCompliantResources int      `json:"noncompliant_resources"`
		CompliantResources    int      `json:"compliant_resources"`
		ImpactedResources     int      `json:"impacted_resources"`
		CostImpact            *float64 `json:"cost_impact"`
	} `json:"compliance_results_summary"`
}

type ControlsFilterSummaryResult struct {
	ControlsCount    int64               `json:"controls_count"`
	IntegrationTypes []string            `json:"integration_types"`
	Severity         []string            `json:"severity"`
	Tags             map[string][]string `json:"tags"`
	PrimaryResource  []string            `json:"primary_resource"`
	ListOfResources  []string            `json:"list_of_resources"`
}

type ControlTagsResult struct {
	Key          string
	UniqueValues []string
}

type BenchmarkTagsResult struct {
	Key          string
	UniqueValues []string
}

type GetControlDetailsResponse struct {
	ID              string             `json:"id"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	IntegrationType []integration.Type `json:"integrationType"`
	Severity        string             `json:"severity"`
	Policy          struct {
		Language        string           `json:"language"`
		Definition      string           `json:"definition"`
		PrimaryResource string           `json:"primaryResource"`
		ListOfResources []string         `json:"listOfResources"`
		Parameters      []QueryParameter `json:"parameters"`
	} `json:"policy"`
	Tags       map[string][]string `json:"tags"`
	Benchmarks *struct {
		Roots    []string `json:"roots"`
		FullPath []string `json:"fullPath"`
	} `json:"benchmarks"`
}

type ListControlsFiltersResponse struct {
	Provider        []string            `json:"provider"`
	Severity        []string            `json:"severity"`
	RootBenchmark   []string            `json:"root_benchmark"`
	ParentBenchmark []string            `json:"parent_benchmark"`
	PrimaryResource []string            `json:"primary_resource"`
	ListOfResources []string            `json:"list_of_resources"`
	Tags            []ControlTagsResult `json:"tags"`
}

type ServiceControls struct {
	Service  string    `json:"service"`
	Controls []Control `json:"controls"`
}

type CategoryControls struct {
	Category string            `json:"category"`
	Services []ServiceControls `json:"services"`
}

type GetCategoriesControlsResponse struct {
	Categories []CategoryControls `json:"categories"`
}

type ParametersControls struct {
	Parameter string    `json:"parameter"`
	Controls  []Control `json:"controls"`
}

type GetParametersControlsResponse struct {
	ParametersControls []ParametersControls `json:"parameters"`
}
