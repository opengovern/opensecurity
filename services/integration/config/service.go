package config

import (
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
)

type IntegrationConfig struct {
	Postgres  koanf.Postgres              `json:"postgres,omitempty" koanf:"postgres"`
	Steampipe koanf.Postgres              `json:"steampipe,omitempty" koanf:"steampipe"`
	Http      koanf.HttpServer            `json:"http,omitempty" koanf:"http"`
	Vault     vault.Config                `json:"vault,omitempty" koanf:"vault"`
	Core  koanf.OpenGovernanceService `json:"core,omitempty" koanf:"core"`
}
