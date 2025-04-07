package api

import "time"



type GetUserLayoutRequest struct {
	UserID string `json:"user_id"`
}
type GetUserLayoutResponse struct {
	UserID string `json:"user_id"`
	LayoutConfig []map[string]any `json:"layout_config"`
	Name string `json:"name"`
	Description string `json:"description"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ChangePrivacyRequest struct {
	UserID string `json:"user_id"`
	IsPrivate bool `json:"is_private"`
}

type SetUserLayoutRequest struct {
	UserID      string `json:"user_id"`
	LayoutConfig []map[string]any `json:"layout_config"`
	Name 	  string `json:"name"`
	IsPrivate 	  bool `json:"is_private"`
}