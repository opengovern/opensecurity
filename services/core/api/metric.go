package api

import (
	"github.com/opengovern/og-util/pkg/integration"
	"time"

)

type Metric struct {
	ID                       string              `json:"id" example:"vms"`
	FinderQuery              string              `json:"finderQuery" example:"select * from platform_resources where resource_type = 'aws::ec2::instance'"`
	FinderPerConnectionQuery string              `json:"finderPerConnectionQuery" example:"select * from platform_resources where resource_type = 'aws::ec2::instance' AND connection_id IN <CONNECTION_ID_LIST>"`
	IntegrationTypes         []integration.Type  `json:"integrationTypes" example:"[Azure]"`            // Cloud Provider
	Name                     string              `json:"name" example:"VMs"`                            // Resource Type
	Tags                     map[string][]string `json:"tags,omitempty"`                                // Tags
	LastEvaluated            *time.Time          `json:"last_evaluated" example:"2020-01-01T00:00:00Z"` // Last time the metric was evaluated

	Count    *int `json:"count" example:"100"`    // Number of Resources of this Resource Type - Metric
	OldCount *int `json:"old_count" example:"90"` // Number of Resources of this Resource Type in the past - Metric
}


type ListMetricsResponse struct {
	TotalCount    *int     `json:"total_count"`
	TotalOldCount *int     `json:"total_old_count"`
	TotalMetrics  int      `json:"total_metrics"`
	Metrics       []Metric `json:"metrics"`
}
