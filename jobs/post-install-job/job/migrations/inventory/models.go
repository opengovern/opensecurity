package inventory

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/shared"
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

type NamedPolicy struct {
	ID               string              `json:"id" yaml:"id"`
	Title            string              `json:"title" yaml:"title"`
	Description      string              `json:"description" yaml:"description"`
	IntegrationTypes []integration.Type  `json:"integration_type" yaml:"integration_type"`
	Policy           shared.Policy       `json:"policy" yaml:"policy"`
	Tags             map[string][]string `json:"tags" yaml:"tags"`
}
