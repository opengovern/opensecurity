package models

import (
	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/integration"
)

type Manifest struct {
	PluginID        string           `json:"plugin_id" yaml:"plugin_id"`
	IntegrationType integration.Type `json:"integration_type" yaml:"integration_type"`
}

type IntegrationPluginInstallState string
type IntegrationPluginOperationalStatus string

const (
	IntegrationTypeInstallStateNotInstalled IntegrationPluginInstallState = "not_installed"
	IntegrationTypeInstallStateInstalled    IntegrationPluginInstallState = "installed"
)

const (
	IntegrationPluginOperationalStatusEnabled  IntegrationPluginOperationalStatus = "enabled"
	IntegrationPluginOperationalStatusDisabled IntegrationPluginOperationalStatus = "disabled"
	IntegrationPluginOperationalStatusFailed   IntegrationPluginOperationalStatus = "failed"
)

type IntegrationPlugin struct {
	PluginID          string `gorm:"primaryKey"`
	IntegrationType   integration.Type
	InstallState      IntegrationPluginInstallState
	OperationalStatus IntegrationPluginOperationalStatus
	URL               string

	OperationalStatusUpdates pq.StringArray `gorm:"type:text[]"`

	IntegrationPlugin []byte `gorm:"type:bytea;not null"`
	CloudQlPlugin     []byte `gorm:"type:bytea;not null"`
}
