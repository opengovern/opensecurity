package api
import (
	complianceapi "github.com/opengovern/opencomply/services/compliance/api"

)

type QueryParameter struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	Required bool   `json:"required" example:"true"`
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
	Cursor   int64 `json:"cursor"`
	PerPage  int64 `json:"per_page"`
	Controls []string `json:"controls"`
	Queries  []string `json:"queries"`

}

type GetQueryParamDetailsResponse struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	Controls 	[]complianceapi.Control `json:"controls"`
	Queries 	[]NamedQueryItemV2 `json:"queries"`

}