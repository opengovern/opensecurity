package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"time"

)

type CostTrendDatapoint struct {
	Cost                                    float64           `json:"cost" minimum:"0"`
	CostStacked                             []CostStackedItem `json:"costStacked" minimum:"0"`
	TotalDescribedConnectionCount           int64             `json:"totalIntegrationCount"`
	TotalSuccessfulDescribedConnectionCount int64             `json:"totalSuccessfulDescribedConnectionCount"`
	Date                                    time.Time         `json:"date" format:"date-time"`
}

type CostStackedItem struct {
	MetricID   string   `json:"metricID"`
	MetricName string   `json:"metricName"`
	Category   []string `json:"category"`
	Cost       float64  `json:"cost"`
}

type ResourceCountStackedItem struct {
	MetricID   string   `json:"metricID"`
	MetricName string   `json:"metricName"`
	Category   []string `json:"category"`
	Count      int      `json:"count"`
}

type ListServicesCostTrendDatapoint struct {
	ServiceName string               `json:"serviceName" example:"EC2-Service-Example"`
	CostTrend   []CostTrendDatapoint `json:"costTrend"`
}

type CostMetric struct {
	IntegrationType          []integration.Type `json:"integration_types" example:"Azure"`
	CostDimensionName        string             `json:"cost_dimension_name" example:"microsoft.compute/disks"`
	CostDimensionID          string             `json:"cost_dimension_id" example:"microsoft_compute_disks"`
	TotalCost                *float64           `json:"total_cost,omitempty" example:"621041.2436112489" minimum:"0"`
	DailyCostAtStartTime     *float64           `json:"daily_cost_at_start_time,omitempty" example:"21232.10443638001" minimum:"0"`
	DailyCostAtEndTime       *float64           `json:"daily_cost_at_end_time,omitempty" example:"14118.815231085681" minimum:"0"`
	FinderQuery              string             `json:"finderQuery"`
	FinderPerConnectionQuery string             `json:"finderPerConnectionQuery"`
}

type ListCostMetricsResponse struct {
	TotalCount int          `json:"total_count" example:"10" minimum:"0"`
	TotalCost  float64      `json:"total_cost" example:"1000" minimum:"0"`
	Metrics    []CostMetric `json:"metrics"`
}



type ListCostCompositionResponse struct {
	TotalCount     int                `json:"total_count" example:"10" minimum:"0"`
	TotalCostValue float64            `json:"total_cost_value" example:"1000" minimum:"0"`
	TopValues      map[string]float64 `json:"top_values"`
	Others         float64            `json:"others" example:"100" minimum:"0"`
}

type TableGranularityType string

const (
	TableGranularityTypeDaily   TableGranularityType = "daily"
	TableGranularityTypeMonthly TableGranularityType = "monthly"
	TableGranularityTypeYearly  TableGranularityType = "yearly"
)

type DimensionType string

const (
	DimensionTypeMetric     DimensionType = "metric"
	DimensionTypeConnection DimensionType = "connection"
)

type SpendTableRow struct {
	DimensionID     string             `json:"dimensionId" example:"compute"`
	AccountID       string             `json:"accountID" example:"1239042"`
	IntegrationType integration.Type   `json:"integrationType" example:"AWS"`
	Category        string             `json:"category" example:"Compute"`
	DimensionName   string             `json:"dimensionName" example:"Compute"`
	CostValue       map[string]float64 `json:"costValue"`
}

type AssetTableRow struct {
	DimensionID     string             `json:"dimensionId" example:"compute"`
	DimensionName   string             `json:"dimensionName" example:"Compute"`
	ResourceCount   map[string]float64 `json:"resourceCount"`
	IntegrationType integration.Type   `json:"integrationType"`
}
