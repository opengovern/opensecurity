package api

import (
	"gitlab.com/keibiengine/keibi-engine/pkg/compliance/es"
	"gitlab.com/keibiengine/keibi-engine/pkg/source"
	"gitlab.com/keibiengine/keibi-engine/pkg/types"
)

type BenchmarkAssignment struct {
	BenchmarkId  string `json:"benchmarkId"`
	ConnectionId string `json:"sourceId"`
	AssignedAt   int64  `json:"assignedAt"`
}

type BenchmarkAssignedSource struct {
	Connection types.FullConnection `json:"connection"`
	AssignedAt int64                `json:"assignedAt"`
}

type FindingFilters struct {
	Connector      []source.Type            `json:"connector"`
	ResourceID     []string                 `json:"resourceID"`
	ResourceTypeID []string                 `json:"resourceTypeID"`
	ConnectionID   []string                 `json:"connectionID"`
	BenchmarkID    []string                 `json:"benchmarkID"`
	PolicyID       []string                 `json:"policyID"`
	Severity       []string                 `json:"severity"`
	Status         []types.ComplianceResult `json:"status"`
}

type FindingResponseFilters struct {
	Provider      []source.Type            `json:"provider"`
	ResourceType  []types.FullResourceType `json:"resourceTypeID"`
	Connections   []types.FullConnection   `json:"connections"`
	FindingStatus []types.ComplianceResult `json:"findingStatus"`
	Benchmarks    []types.FullBenchmark    `json:"benchmarks"`
	Policies      []types.FullPolicy       `json:"policies"`
	Severity      []string                 `json:"severity"`
}

type DirectionType string

const (
	DirectionAscending  DirectionType = "asc"
	DirectionDescending DirectionType = "desc"
)

type SortFieldType string

const (
	FieldResourceID             SortFieldType = "resourceID"
	FieldResourceName           SortFieldType = "resourceName"
	FieldResourceType           SortFieldType = "resourceType"
	FieldServiceName            SortFieldType = "serviceName"
	FieldCategory               SortFieldType = "category"
	FieldResourceLocation       SortFieldType = "resourceLocation"
	FieldStatus                 SortFieldType = "status"
	FieldDescribedAt            SortFieldType = "describedAt"
	FieldEvaluatedAt            SortFieldType = "evaluatedAt"
	FieldSourceID               SortFieldType = "sourceID"
	FieldConnectionProviderID   SortFieldType = "connectionProviderID"
	FieldConnectionProviderName SortFieldType = "connectionProviderName"
	FieldSourceType             SortFieldType = "sourceType"
	FieldBenchmarkID            SortFieldType = "benchmarkID"
	FieldPolicyID               SortFieldType = "policyID"
	FieldPolicySeverity         SortFieldType = "policySeverity"
)

type FindingSortItem struct {
	Field     SortFieldType `json:"field" enums:"resourceID,resourceName,resourceType,serviceName,category,resourceLocation,status,describedAt,evaluatedAt,sourceID,connectionProviderID,connectionProviderName,sourceType,benchmarkID,policyID,policySeverity"`
	Direction DirectionType `json:"direction" enums:"asc,desc"`
}

type Page struct {
	No   int `json:"no,omitempty"`
	Size int `json:"size,omitempty"`
}

type GetFindingsRequest struct {
	Filters FindingFilters    `json:"filters"`
	Sorts   []FindingSortItem `json:"sorts"`
	Page    Page              `json:"page" validate:"required"`
}

type TopField = string

const (
	TopField_ResourceType TopField = "resourceType"
	TopField_CloudService TopField = "serviceName"
	TopField_CloudAccount TopField = "sourceID"
	TopField_Resources    TopField = "resourceID"
)

type GetTopFieldRequest struct {
	Field   TopField       `json:"field" enums:"resourceType,serviceName,sourceID,resourceID"`
	Filters FindingFilters `json:"filters"`
	Count   int            `json:"count"`
}

type TopFieldRecord struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type GetTopFieldResponse struct {
	Records []TopFieldRecord `json:"records"`
}

type GetFindingsResponse struct {
	Findings   []es.Finding `json:"findings"`
	TotalCount int64        `json:"totalCount"`
}

type GetFindingsFiltersResponse struct {
	Filters FindingResponseFilters `json:"filters"`
}

type Datapoint struct {
	Time  int64 `json:"time"`
	Value int64 `json:"value"`
}

type StatusCount struct {
	Passed int64 `json:"passed"`
	Failed int64 `json:"failed"`
}

type BenchmarkSummaryPolicySummary struct {
	Policy       types.FullPolicy                   `json:"policy"`
	ShortSummary types.ComplianceResultShortSummary `json:"shortSummary"`
}

type BenchmarkSummaryResourceSummary struct {
	Resource     types.FullResource                 `json:"resource"`
	ShortSummary types.ComplianceResultShortSummary `json:"shortSummary"`
}

type BenchmarkSummary struct {
	ID                       string                             `json:"id"`
	Title                    string                             `json:"title"`
	Description              string                             `json:"description"`
	Result                   map[types.ComplianceResult]int     `json:"result"`
	ShortSummary             types.ComplianceResultShortSummary `json:"shortSummary"`
	Policies                 []BenchmarkSummaryPolicySummary    `json:"policies"`
	Resources                []BenchmarkSummaryResourceSummary  `json:"resources"`
	CompliancyTrend          []Datapoint                        `json:"trend"`
	AssignedConnectionsCount int64                              `json:"assignedConnectionsCount"`
	TotalConnectionResources int64                              `json:"totalConnectionResources"`
	Tags                     map[string]string                  `json:"tags"`
	Enabled                  bool                               `json:"enabled"`
}

type BenchmarkSummaryConnectionSummary struct {
	Connection   types.FullConnection               `json:"connection"`
	ShortSummary types.ComplianceResultShortSummary `json:"shortSummary"`
}

type NewBenchmarkSummary struct {
	ID         string            `json:"benchmark_id"`
	Title      string            `json:"title"`
	Connectors []source.Type     `json:"connectors"`
	Tags       map[string]string `json:"tags"`
	Enabled    bool              `json:"enabled"`

	PassedResourceIDs []string                       `json:"passed_resource_ids"`
	FailedResourceIDs []string                       `json:"failed_resource_ids"`
	Result            map[types.ComplianceResult]int `json:"result"`

	CompliancyTrend []Datapoint `json:"trend"`
}

type GetBenchmarksSummaryResponse struct {
	ShortSummary types.ComplianceResultShortSummary  `json:"shortSummary"`
	TotalAssets  int64                               `json:"totalAssets"`
	Connections  []BenchmarkSummaryConnectionSummary `json:"connections"`
	Benchmarks   []BenchmarkSummary                  `json:"benchmarks"`
}

type BenchmarkShortSummary struct {
	ID         string                        `json:"id"`
	Title      string                        `json:"title"`
	Connectors []source.Type                 `json:"connectors"`
	Tags       map[string]string             `json:"tags"`
	Enabled    bool                          `json:"enabled"`
	Result     types.ComplianceResultSummary `json:"result"`
	Coverage   float64                       `json:"coverage"`

	PassedResources int64 `json:"passedResources"`
	FailedResources int64 `json:"failedResources"`
}

type GetShortSummaryResponse struct {
	BenchmarkShortSummary []BenchmarkShortSummary `json:"benchmarkShortSummary"`

	PassedResources int64 `json:"passedResources"`
	FailedResources int64 `json:"failedResources"`
	TotalAssets     int64 `json:"totalAssets"`
}

type PolicySummary struct {
	Title       string             `json:"title"`
	Category    string             `json:"category"`
	Subcategory string             `json:"subcategory"`
	Severity    types.Severity     `json:"severity"`
	Status      types.PolicyStatus `json:"status"`
	CreatedAt   int64              `json:"createdAt"`
}

type GetPoliciesSummaryResponse struct {
	BenchmarkTitle       string                        `json:"title"`
	BenchmarkDescription string                        `json:"description"`
	ComplianceSummary    types.ComplianceResultSummary `json:"complianceSummary"`
	PolicySummary        []PolicySummary               `json:"policySummary"`
	Tags                 map[string]string             `json:"tags"`
	Enabled              bool                          `json:"enabled"`
}

type GetFindingsMetricsResponse struct {
	TotalFindings   int64 `json:"totalFindings"`
	FailedFindings  int64 `json:"failedFindings"`
	PassedFindings  int64 `json:"passedFindings"`
	UnknownFindings int64 `json:"unknownFindings"`

	LastTotalFindings   int64 `json:"lastTotalFindings"`
	LastFailedFindings  int64 `json:"lastFailedFindings"`
	LastPassedFindings  int64 `json:"lastPassedFindings"`
	LastUnknownFindings int64 `json:"lastUnknownFindings"`
}

type Alarms struct {
	Policy    types.FullPolicy       `json:"policy"`
	CreatedAt int64                  `json:"created_at"`
	Status    types.ComplianceResult `json:"status"`
}

type GetFindingDetailsResponse struct {
	Connection        types.FullConnection   `json:"connection"`
	Resource          types.FullResource     `json:"resource"`
	ResourceType      types.FullResourceType `json:"resourceType"`
	State             types.ComplianceResult `json:"state"`
	CreatedAt         int64                  `json:"createdAt"`
	PolicyTags        map[string]string      `json:"policyTags"`
	PolicyDescription string                 `json:"policyDescription"`
	Reason            string                 `json:"reason"`
	Alarms            []Alarms               `json:"alarms"`
}

type InsightRecord struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type GetBenchmarkInsightResponse struct {
	TopResourceType []InsightRecord `json:"topResourceType"`
	TopCategory     []InsightRecord `json:"topCategory"`
	TopAccount      []InsightRecord `json:"topAccount"`
	Severity        []InsightRecord `json:"severity"`
}
