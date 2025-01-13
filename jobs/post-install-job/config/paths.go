package config

const (
	ControlEnrichmentGitPath = "/tmp/loader-control-enrichment-git"

	ConfigzGitPath            = "/tmp/loader-analytics-git"
	AssetsGitPath             = ConfigzGitPath + "/analytics/cloud-infra"
	SpendGitPath              = ConfigzGitPath + "/analytics/cloud-spend"
	QueriesGitPath            = ConfigzGitPath + "/queries"
	ComplianceGitPath         = ConfigzGitPath + "/compliance"
	QueryViewsGitPath         = ConfigzGitPath + "/views"
	ResourceCollectionGitPath = ConfigzGitPath + "/resource-collections"
	IntegrationsGitPath       = ConfigzGitPath + "/integrations"
	IntegrationTypesYamlPath  = ConfigzGitPath + "/integrations-types/extraIntegrations.yaml"
)
