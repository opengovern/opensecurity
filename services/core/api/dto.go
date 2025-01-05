package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"time"
)

type DirectionType string

type SortFieldType string

type Page struct {
	No   int `json:"no,omitempty"`
	Size int `json:"size,omitempty"`
}

// ResourceFilters model
//
//	@Description	if you provide two values for same filter OR operation would be used
//	@Description	if you provide value for two filters AND operation would be used
type ResourceFilters struct {
	// if you dont need to use this filter, leave them empty. (e.g. [])
	ResourceType []string `json:"resourceType"`
	// if you dont need to use this filter, leave them empty. (e.g. [])
	Category []string `json:"category"`
	// if you dont need to use this filter, leave them empty. (e.g. [])
	Service []string `json:"service"`
	// if you dont need to use this filter, leave them empty. (e.g. [])
	Location []string `json:"location"`
	// if you dont need to use this filter, leave them empty. (e.g. [])
	Provider []string `json:"provider"`
	// if you dont need to use this filter, leave them empty. (e.g. [])
	Connections []string `json:"connections"`
	// if you dont need to use this filter, leave them empty. (e.g. [])
	TagKeys []string `json:"tagKeys"`
	// if you dont need to use this filter, leave them empty. (e.g. [])
	TagValues map[string][]string `json:"tagValues"`
}

type NamedQuerySortItem struct {
	// fill this with column name
	Field     string        `json:"field"`
	Direction DirectionType `json:"direction" enums:"asc,desc"`
}

type NamedQueryItem struct {
	ID               string             `json:"id"`                // Query Id
	IntegrationTypes []integration.Type `json:"integration_types"` // Provider
	Title            string             `json:"title"`             // Title
	Category         string             `json:"category"`          // Category (Tags[category])
	Query            string             `json:"query"`             // Query
	Tags             map[string]string  `json:"tags"`              // Tags
}

type ListQueriesV2Response struct {
	Items      []NamedQueryItemV2 `json:"items"`
	TotalCount int                `json:"total_count"`
}

type NamedQueryItemV2 struct {
	ID               string              `json:"id"`    // Query Id
	Title            string              `json:"title"` // Title
	Description      string              `json:"description"`
	IntegrationTypes []integration.Type  `json:"integration_types"` // Provider
	Query            Query               `json:"query"`             // Query
	Tags             map[string][]string `json:"tags"`              // Tags
}




type ListQueryV2Request struct {
	QueryIDs          []string            `json:"query_ids"`
	TitleFilter       string              `json:"title_filter"`
	IntegrationTypes  []string            `json:"integration_types"`
	HasParameters     *bool               `json:"has_parameters"`
	IntegrationExists bool                `json:"integration_exists"`
	Categories        []string            `json:"categories"`
	IsBookmarked      bool                `json:"is_bookmarked"`
	PrimaryTable      []string            `json:"primary_table"`
	ListOfTables      []string            `json:"list_of_tables"`
	Tags              map[string][]string `json:"tags"`
	TagsRegex         *string             `json:"tags_regex"`
	Cursor            *int64              `json:"cursor"`
	PerPage           *int64              `json:"per_page"`
}

type ListQueryRequest struct {
	TitleFilter string `json:"titleFilter"` // Specifies the Title
}

type ConnectionData struct {
	ConnectionID         string     `json:"connectionID"`
	Count                *int       `json:"count"`
	OldCount             *int       `json:"oldCount"`
	LastInventory        *time.Time `json:"lastInventory"`
	TotalCost            *float64   `json:"cost"`
	DailyCostAtStartTime *float64   `json:"dailyCostAtStartTime"`
	DailyCostAtEndTime   *float64   `json:"dailyCostAtEndTime"`
}

type CountAnalyticsMetricsResponse struct {
	ConnectionCount int `json:"connectionCount"`
	MetricCount     int `json:"metricCount"`
}

type CountAnalyticsSpendResponse struct {
	ConnectionCount int `json:"connectionCount"`
	MetricCount     int `json:"metricCount"`
}

type ParametersQueries struct {
	Parameter string             `json:"parameter"`
	Queries   []NamedQueryItemV2 `json:"queries"`
}

type GetParametersQueriesResponse struct {
	ParametersQueries []ParametersQueries `json:"parameters"`
}
