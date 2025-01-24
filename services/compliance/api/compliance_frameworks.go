package api

type FrameworkAssignmentAssignmentType string

const (
	FrameworkAssignmentAssignmentTypeExplicit FrameworkAssignmentAssignmentType = "explicit"
	FrameworkAssignmentAssignmentTypeImplicit FrameworkAssignmentAssignmentType = "implicit"
	FrameworkAssignmentAssignmentTypeNone     FrameworkAssignmentAssignmentType = "none"
)

type PageInfo struct {
	CurrentPage int64 `json:"current_page"`
	PageSize    int64 `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int64 `json:"total_pages"`
}

type ListFrameworkAssignmentsResponseData struct {
	IntegrationID         string                            `json:"integration_id"`
	IntegrationName       string                            `json:"integration_name"`
	IntegrationProviderID string                            `json:"integration_provider_id"`
	PluginID              string                            `json:"plugin_id"`
	AssignmentType        FrameworkAssignmentAssignmentType `json:"assignment_type"`
}

type ListFrameworkAssignmentsResponse struct {
	Data     []ListFrameworkAssignmentsResponseData `json:"data"`
	PageInfo PageInfo                               `json:"page_info"`
}

type UpdateFrameworkSettingRequest struct {
	IsBaseline *bool `json:"is_baseline"`
	Enabled    *bool `json:"enabled"`
}
