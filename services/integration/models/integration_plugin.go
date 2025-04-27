package models

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/integration"
	"time"
)

type Manifest struct {
	IntegrationType          integration.Type `json:"IntegrationType" yaml:"IntegrationType"`
	DescriberURL             string           `json:"DescriberURL" yaml:"DescriberURL"`
	DemoDataURL              string           `json:"DemoDataURL" yaml:"DemoDataURL"`
	DescriberTag             string           `json:"DescriberTag" yaml:"DescriberTag"`
	Publisher                string           `json:"Publisher" yaml:"Publisher"`
	Author                   string           `json:"Author" yaml:"Author"`
	SupportedPlatformVersion string           `json:"SupportedPlatformVersion" yaml:"SupportedPlatformVersion"`
	UpdateDate               string           `json:"UpdateDate" yaml:"UpdateDate"`
	DiscoveryType            string           `json:"DiscoveryType" yaml:"DiscoveryType"`
}

type IntegrationPluginInstallState string
type IntegrationPluginOperationalStatus string

type IntegrationPluginDiscoveryType string

const (
	IntegrationTypeInstallStateNotInstalled IntegrationPluginInstallState = "not_installed"
	IntegrationTypeInstallStateInstalling   IntegrationPluginInstallState = "installing"
	IntegrationTypeInstallStateInstalled    IntegrationPluginInstallState = "installed"
)

const (
	IntegrationPluginOperationalStatusEnabled  IntegrationPluginOperationalStatus = "enabled"
	IntegrationPluginOperationalStatusDisabled IntegrationPluginOperationalStatus = "disabled"
	IntegrationPluginOperationalStatusFailed   IntegrationPluginOperationalStatus = "failed"
)

const (
	IntegrationPluginDiscoveryTypeClassic IntegrationPluginDiscoveryType = "classic"
	IntegrationPluginDiscoveryTypeTask    IntegrationPluginDiscoveryType = "task"
)

type OperationalStatusUpdate struct {
	Time      time.Time
	OldStatus IntegrationPluginOperationalStatus
	NewStatus IntegrationPluginOperationalStatus
	Reason    string
}

type IntegrationPlugin struct {
	ID                       int
	UpdatedAt                time.Time
	PluginID                 string `gorm:"primaryKey"`
	IntegrationType          integration.Type
	Name                     string
	Tier                     string
	Description              string
	Icon                     string
	Availability             string
	SourceCode               string
	PackageType              string
	InstallState             IntegrationPluginInstallState
	OperationalStatus        IntegrationPluginOperationalStatus
	URL                      string
	DescriberURL             string
	DemoDataURL              string
	DemoDataLoaded           bool
	DescriberTag             string
	DiscoveryType            IntegrationPluginDiscoveryType `gorm:"default:classic"`
	OperationalStatusUpdates pgtype.JSONB                   `gorm:"default:'[]'"`
	Tags                     pgtype.JSONB
}

func (ip IntegrationPlugin) GetStringOperationalStatusUpdates() ([]string, error) {
	stringUpdates := make([]string, 0)
	var jsonVal []byte
	if err := ip.OperationalStatusUpdates.AssignTo(&jsonVal); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonVal, &stringUpdates); err != nil {
		return nil, err
	}

	return stringUpdates, nil
}

type IntegrationPluginBinary struct {
	PluginID string `gorm:"primaryKey"`

	IntegrationPlugin []byte `gorm:"type:bytea"`
	CloudQlPlugin     []byte `gorm:"type:bytea"`
}
