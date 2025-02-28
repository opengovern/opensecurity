package config

import (
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
)

type Config struct {
	Postgres      koanf.Postgres   `yaml:"postgres" koanf:"postgres"`
	Http          koanf.HttpServer `yaml:"http" koanf:"http"`
	NATS          config.NATS      `yaml:"nats" koanf:"nats"`
	Vault         vault.Config     `yaml:"vault" koanf:"vault"`
	ElasticSearch config.ElasticSearch

	ESSinkEndpoint string `yaml:"essink_endpoint" koanf:"essink_endpoint"`
}
