package models

import "time"

type OperationalStatusUpdate struct {
	Time      time.Time
	OldStatus string
	NewStatus string
	Reason    string
}

type IntegrationPlugin struct {
	ID                int                             `json:"id"`
	PluginID          string                          `json:"plugin_id"`
	IntegrationType   string                          `json:"integration_type"`
	Name              string                          `json:"name"`
	Tier              string                          `json:"tier"`
	Description       string                          `json:"description"`
	Icon              string                          `json:"icon"`
	Availability      string                          `json:"availability"`
	SourceCode        string                          `json:"source_code"`
	PackageType       string                          `json:"package_type"`
	InstallState      string                          `json:"install_state"`
	OperationalStatus string                          `json:"operational_status"`
	URL               string                          `json:"url"`
	DescriberURL      string                          `json:"describer_url"`
	DescriberTag      string                          `json:"describer_tag"`
	Count             IntegrationTypeIntegrationCount `json:"count"`

	OperationalStatusUpdates []OperationalStatusUpdate `json:"operational_status_updates"`
}

type IntegrationPluginListResponse struct {
	Items      []IntegrationPlugin `json:"items"`
	TotalCount int                 `json:"total_count"`
}
