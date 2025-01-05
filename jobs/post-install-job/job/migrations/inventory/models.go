package inventory

import (
	"github.com/opengovern/og-util/pkg/integration"
)

type Metric struct {
	IntegrationTypes         []integration.Type  `json:"integrationType" yaml:"integrationType"`
	Name                     string              `json:"name" yaml:"name"`
	Query                    string              `json:"query" yaml:"query"`
	Tables                   []string            `json:"tables" yaml:"tables"`
	FinderQuery              string              `json:"finderQuery" yaml:"finderQuery"`
	FinderPerConnectionQuery string              `json:"finderPerConnectionQuery" yaml:"finderPerConnectionQuery"`
	Tags                     map[string][]string `json:"tags" yaml:"tags"`
	Status                   string              `json:"status" yaml:"status"`
}
