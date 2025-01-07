package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"time"
)

type QueryParameter struct {
	Key           string `json:"key"`
	ControlID     string `json:"control_id"`
	Value         string `json:"value"`
	Required      bool   `json:"required" example:"true"`
	ControlsCount int    `json:"controls_count"`
	QueriesCount  int    `json:"queries_count"`
}

type SetQueryParameterRequest struct {
	QueryParameters []QueryParameter `json:"query_parameters"`
}

type ListQueryParametersResponse struct {
	Items      []QueryParameter `json:"items"`
	TotalCount int              `json:"total_count"`
}

type ListQueryParametersRequest struct {
	Cursor   int64    `json:"cursor"`
	PerPage  int64    `json:"per_page"`
	KeyRegex *string  `json:"key_regex"`
	Controls []string `json:"controls"`
	Queries  []string `json:"queries"`
}

type ControlQueryParameter struct {
	Key      string `json:"key" example:"key"`
	Required bool   `json:"required" example:"true"`
}
type PolicyLanguage string

type Policy struct {
	ID              string             `json:"id" example:"azure_ad_manual_control"`
	Language        PolicyLanguage     `json:"language" example:"sql"`
	Definition      string             `json:"definition" example:"select\n  -- Required Columns\n  'active_directory' as resource,\n  'info' as status,\n  'Manual verification required.' as reason;\n"`
	IntegrationType []integration.Type `json:"integrationType" example:"Azure"`
	PrimaryResource *string            `json:"primaryResource" example:"null"`
	ListOfResources []string           `json:"listOfResources" example:"null"`

	// CloudQL Fields
	Parameters []ControlQueryParameter `json:"parameters"`

	// Rego Fields
	RegoPolicies []string `json:"regoPolicies" example:"null"`

	CreatedAt time.Time `json:"createdAt" example:"2023-06-07T14:00:15.677558Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2023-06-16T14:58:08.759554Z"`
}

type ComplianceResultSeverity string

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

	IntegrationType    []string                 `json:"integration_type" example:"Azure"`
	Enabled            bool                     `json:"enabled" example:"true"`
	DocumentURI        string                   `json:"documentURI" example:"benchmarks/azure_cis_v140_1_1.md"`
	Policy             *Policy                  `json:"policy"`
	Severity           ComplianceResultSeverity `json:"severity" example:"low"`
	ManualVerification bool                     `json:"manualVerification" example:"true"`
	Managed            bool                     `json:"managed" example:"true"`
	CreatedAt          time.Time                `json:"createdAt" example:"2020-01-01T00:00:00Z"`
	UpdatedAt          time.Time                `json:"updatedAt" example:"2020-01-01T00:00:00Z"`
}

type GetQueryParamDetailsResponse struct {
	Key      string             `json:"key"`
	Value    string             `json:"value"`
	Controls []Control          `json:"controls"`
	Queries  []NamedQueryItemV2 `json:"queries"`
}
