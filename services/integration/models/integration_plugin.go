package models

import (
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/integration"
	"time"
)

type Manifest struct {
	IntegrationType integration.Type `json:"IntegrationType" yaml:"IntegrationType"`
	DescriberURL    string           `json:"DescriberURL" yaml:"DescriberURL"`
	DescriberTag    string           `json:"DescriberTag" yaml:"DescriberTag"`
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

type OperationalStatusUpdate struct {
	Time      time.Time
	OldStatus IntegrationPluginOperationalStatus
	NewStatus IntegrationPluginOperationalStatus
	Reason    string
}

type IntegrationPlugin struct {
	ID                int
	PluginID          string `gorm:"primaryKey"`
	IntegrationType   integration.Type
	Name              string
	Tier              string
	Description       string
	Icon              string
	Availability      string
	SourceCode        string
	PackageType       string
	InstallState      IntegrationPluginInstallState
	OperationalStatus IntegrationPluginOperationalStatus
	URL               string
	DescriberURL      string
	DescriberTag      string

	OperationalStatusUpdates pq.StringArray `gorm:"type:text[]"`

	IntegrationPlugin []byte `gorm:"type:bytea;not null"`
	CloudQlPlugin     []byte `gorm:"type:bytea;not null"`

	Tags pgtype.JSONB
}
