package config

import (
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
)

type Config struct {
	Postgres    koanf.Postgres              `yaml:"postgres" koanf:"postgres"`
	Http        koanf.HttpServer            `yaml:"http" koanf:"http"`
	Integration koanf.OpenGovernanceService `yaml:"integration" koanf:"integration"`
	Scheduler   koanf.OpenGovernanceService `yaml:"scheduler" koanf:"scheduler"`
	Compliance  koanf.OpenGovernanceService `yaml:"compliance" koanf:"compliance"`
	Vault       vault.Config                `yaml:"vault" koanf:"vault"`

	OpengovernanceNamespace     string               `yaml:"opengovernance_namespace" koanf:"opengovernance_namespace"`
	PrimaryDomainURL            string               `yaml:"primary_domain_url" koanf:"primary_domain_url"`
	SampledataIntegrationsCheck string               `yaml:"sampledata_integrations_check" koanf:"sampledata_integrations_check"`
	ElasticSearch               config.ElasticSearch `yaml:"elasticsearch" koanf:"elasticsearch"`
}
