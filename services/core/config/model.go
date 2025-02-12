package config

import (
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
)

type Config struct {
	Postgres koanf.Postgres `yaml:"postgres" koanf:"postgres"`

	Steampipe      config.Postgres `yaml:"steampipe" koanf:"steampipe"`
	PostgresPlugin config.Postgres `yaml:"postgresplugin" koanf:"postgresplugin"`

	Http        koanf.HttpServer            `yaml:"http" koanf:"http"`
	Integration koanf.OpenGovernanceService `yaml:"integration" koanf:"integration"`
	Scheduler   koanf.OpenGovernanceService `yaml:"scheduler" koanf:"scheduler"`
	Compliance  koanf.OpenGovernanceService `yaml:"compliance" koanf:"compliance"`
	Vault       vault.Config                `yaml:"vault" koanf:"vault"`

	OpengovernanceNamespace      string               `yaml:"opengovernance_namespace" koanf:"opengovernance_namespace"`
	PrimaryDomainURL             string               `yaml:"primary_domain_url" koanf:"primary_domain_url"`
	DexGrpcAddr                  string               `yaml:"dex_grpc_addr" koanf:"dex_grpc_addr"`
	SampledataIntegrationsCheck  string               `yaml:"sampledata_integrations_check" koanf:"sampledata_integrations_check"`
	DexPublicClientRedirectUris  string               `yaml:"dex_public_client_redirect_uris" koanf:"dex_public_client_redirect_uris"`
	DexPrivateClientRedirectUris string               `yaml:"dex_private_client_redirect_uris" koanf:"dex_private_client_redirect_uris"`
	ElasticSearch                config.ElasticSearch `yaml:"elastic_search" koanf:"elastic_search"`
}
