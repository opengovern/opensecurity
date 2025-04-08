package api

import "time"



type GetUserLayoutRequest struct {
	UserID string `json:"user_id"`
}
type GetUserLayoutResponse struct {
	ID uint `json:"id"`
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
	ID 	   uint `json:"id"`
	UserID      string `json:"user_id"`
	Description string `json:"description"`
	LayoutConfig []map[string]any `json:"layout_config"`
	UpdatedAt time.Time `json:"updated_at"`
	Name 	  string `json:"name"`
	IsPrivate 	  bool `json:"is_private"`
}