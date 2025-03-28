package inventory

import (
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opensecurity/jobs/post-install-job/job/migrations/shared"
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

type QueryView struct {
	ID          string              `json:"id" yaml:"id"`
	Title       string              `json:"title" yaml:"title"`
	Description string              `json:"description" yaml:"description"`
	Query       string              `json:"query" yaml:"query"`
	Tags        map[string][]string `json:"tags" yaml:"tags"`
}

type NamedQuery struct {
	ID               string                    `json:"id" yaml:"id"`
	Title            string                    `json:"title" yaml:"title"`
	Description      string                    `json:"description" yaml:"description"`
	Parameters       []shared.ControlParameter `json:"parameters" yaml:"parameters"`
	IntegrationTypes []string                  `json:"integration_type" yaml:"integration_type"`
	Query            string                    `json:"query" yaml:"query"`
	Tags             map[string][]string       `json:"tags" yaml:"tags"`
}
