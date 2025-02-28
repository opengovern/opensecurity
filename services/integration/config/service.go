package config

import (
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
)

type IntegrationPluginsConfig struct {
	PingIntervalSeconds  int `json:"ping_interval_seconds" koanf:"ping_interval_seconds"`
	MaxAutoRebootRetries int `json:"max_auto_reboot_retries" koanf:"max_auto_reboot_retries"`
}

type IntegrationConfig struct {
	Postgres  koanf.Postgres              `json:"postgres,omitempty" koanf:"postgres"`
	Steampipe koanf.Postgres              `json:"steampipe,omitempty" koanf:"steampipe"`
	Http      koanf.HttpServer            `json:"http,omitempty" koanf:"http"`
	Vault     vault.Config                `json:"vault,omitempty" koanf:"vault"`
	Core      koanf.OpenGovernanceService `json:"core,omitempty" koanf:"core"`

	IntegrationPlugins IntegrationPluginsConfig `json:"integration_plugins,omitempty" koanf:"integration_plugins"`
}
