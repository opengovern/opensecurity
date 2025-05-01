package config

import "github.com/opengovern/og-util/pkg/config"

type MigratorConfig struct {
	IsManual bool `yaml:"is_manual" koanf:"is_manual"`

	PostgreSQL              config.Postgres
	Steampipe               config.Postgres
	ElasticSearch           config.ElasticSearch
	Core                    config.OpenGovernanceService
	Integration             config.OpenGovernanceService
	AnalyticsGitURL         string `yaml:"analytics_git_url"`
	ControlEnrichmentGitURL string `yaml:"control_enrichment_git_url"`
	GithubToken             string `yaml:"github_token"`
	PrometheusPushAddress   string `yaml:"prometheus_push_address"`
	DexGrpcAddress          string `yaml:"dex_grpc_address"`
	DefaultDexUserName      string `yaml:"default_dex_user_name"`
	DefaultDexUserEmail     string `yaml:"default_dex_user_email"`
	DefaultDexUserPassword  string `yaml:"default_dex_user_password"`
}
