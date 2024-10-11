package config

import "github.com/kaytu-io/kaytu-util/pkg/config"

type MigratorConfig struct {
	IsManual bool `yaml:"is_manual"`

	PostgreSQL              config.Postgres
	Steampipe               config.Postgres
	ElasticSearch           config.ElasticSearch
	Metadata                config.KaytuService
	AnalyticsGitURL         string `yaml:"analytics_git_url"`
	ControlEnrichmentGitURL string `yaml:"control_enrichment_git_url"`
	GithubToken             string `yaml:"github_token"`
	PrometheusPushAddress   string `yaml:"prometheus_push_address"`
	DexGrpcAddress          string `yaml:"dex_grpc_address"`
	DefaultDexUserID        string `yaml:"default_dex_user_id"`
	DefaultDexUserName      string `yaml:"default_dex_user_name"`
	DefaultDexUserEmail     string `yaml:"default_dex_user_email"`
	DefaultDexUserPassword  string `yaml:"default_dex_user_password"`
}
