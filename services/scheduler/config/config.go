package config

import (
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/vault"
)

type ServerlessProviderType string

const (
	ServerlessProviderTypeLocal ServerlessProviderType = "local"
)

func (s ServerlessProviderType) String() string {
	return string(s)
}

type SchedulerConfig struct {
	ComplianceIntervalHours int    `yaml:"compliance_interval_hours"`
	ServerlessProvider      string `yaml:"serverless_provider"`
	ElasticSearch           config.ElasticSearch
	Integration             config.OpenGovernanceService
	NATS                    config.NATS
	Vault                   vault.Config `yaml:"vault" koanf:"vault"`
	QueryRunnerEnabled      string       `yaml:"query_runner_enabled" koanf:"query_runner_enabled"`
}
