package cloudql_init_job

import "github.com/opengovern/og-util/pkg/config"

type Config struct {
	Postgres      config.Postgres
	ElasticSearch config.ElasticSearch
	Integration   config.OpenGovernanceService
	Steampipe     config.Postgres
}
