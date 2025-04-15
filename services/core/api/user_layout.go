package api

import "time"



type GetUserLayoutRequest struct {
	UserID string `json:"user_id"`
}
type GetUserLayoutResponse struct {
	ID string `json:"id"`
	IsDefault bool `json:"is_default"`
	UserID string `json:"user_id"`
	LayoutConfig []map[string]any `json:"layout_config"`
	Name string `json:"name"`
	Description string `json:"description"`
	UpdatedAt time.Time `json:"updated_at"`
	IsPrivate bool `json:"is_private"`


}
type ChangePrivacyRequest struct {
	UserID string `json:"user_id"`
	IsPrivate bool `json:"is_private"`
}

type SetUserLayoutRequest struct {
	ID 	   string `json:"id"`
	UserID      string `json:"user_id"`
	IsDefault bool `json:"is_default"`

	Description string `json:"description"`
	WidgetIDs []string `json:"widget_ids"`
	UpdatedAt time.Time `json:"updated_at"`
	Name 	  string `json:"name"`
	IsPrivate 	  bool `json:"is_private"`
}

type UpdateWidgetDashboardsRequest struct {
	WidgetID string `json:"widget_id"`
	Dashboards []string `json:"dashboards"`
}
type UpdateDashboardWidgetsRequest struct {
	DashboardID string `json:"dashboard_id"`
	Widgets []string `json:"widgets"`
}
type SetUserWidgetRequest struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	WidgetType string `json:"widget_type"`
	WidgetProps map[string]any `json:"widget_props"`
	RowSpan int `json:"row_span"`
	ColumnSpan int `json:"column_span"`
	ColumnOffset int `json:"column_offset"`
	IsPublic bool `json:"is_public"`
	UserID string `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

}

type GetUserWidgetRequest struct {
	UserID string `json:"user_id"`
}