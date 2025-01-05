package api

import (
	"time"
)

type QueryEngine string

const (
	QueryEngineCloudQL     = "cloudql"
	QueryEngineCloudQLRego = "cloudql-rego"
)

type RunQueryRequest struct {
	Page      Page                 `json:"page" validate:"required"`
	Query     *string              `json:"query"`
	AccountId *string              `json:"account_id"`
	SourceId  *string              `json:"source_id"`
	Engine    *QueryEngine         `json:"engine"`
	Sorts     []NamedQuerySortItem `json:"sorts"`
}

type RunQueryResponse struct {
	Title   string   `json:"title"`   // Query Title
	Query   string   `json:"query"`   // Query
	Headers []string `json:"headers"` // Column names
	Result  [][]any  `json:"result"`  // Result of query. in order to access a specific cell please use Result[Row][Column]
}

type NamedQueryHistory struct {
	Query      string    `json:"query"`
	ExecutedAt time.Time `json:"executed_at"`
}

type NamedQueryTagsResult struct {
	Key          string
	UniqueValues []string
}

type RunQueryByIDRequest struct {
	Page        Page                 `json:"page" validate:"required"`
	Type        string               `json:"type"`
	ID          string               `json:"id"`
	Sorts       []NamedQuerySortItem `json:"sorts"`
	QueryParams map[string]string    `json:"query_params"`
}

type ListQueriesFiltersResponse struct {
	Providers []string               `json:"providers"`
	Tags      []NamedQueryTagsResult `json:"tags"`
}

type GetAsyncQueryRunResultResponse struct {
	RunId       string           `json:"runID"`
	QueryID     string           `json:"queryID"`
	Parameters  []QueryParameter `json:"parameters"`
	ColumnNames []string         `json:"columnNames"`
	CreatedBy   string           `json:"createdBy"`
	TriggeredAt int64            `json:"triggeredAt"`
	EvaluatedAt int64            `json:"evaluatedAt"`
	Result      [][]string       `json:"result"`
}
