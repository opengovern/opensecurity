package config

import (
	"github.com/opengovern/og-util/pkg/config"
)

type RegoConfig struct {
	Http           config.HttpServer            `json:"http,omitempty" koanf:"http"`
	ElasticSearch  config.ElasticSearch         `json:"elasticsearch,omitempty" koanf:"elasticsearch"`
	Steampipe      config.Postgres              `json:"steampipe,omitempty" koanf:"steampipe"`
	PostgresPlugin config.Postgres              `json:"postgresplugin,omitempty" koanf:"postgresplugin"`
	Integration    config.OpenGovernanceService `json:"integration,omitempty" koanf:"integration"`
}
