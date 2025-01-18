package models

type  IntegrationPlugin struct {
	PluginId		string		`json:"plugin_id"`
	IntegrationType string		`json:"integration_type"`
	InstallState	string		`json:"install_state"`
	OperationalStatus string		`json:"operational_state"`
	OperationalStatusUpdates []string `json:"operational_status_updates"`
	URL string `json:"url"`

}

type IntegrationPluginListResponse struct {
	Items []IntegrationPlugin `json:"items"`
	TotalCount int `json:"total_count"`
}