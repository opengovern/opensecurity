package api



type GetUserLayoutRequest struct {
	UserID string `json:"user_id"`
}
type ChangePrivacyRequest struct {
	UserID string `json:"user_id"`
	IsPrivate bool `json:"is_private"`
}

type SetUserLayout struct {
	UserID      string `json:"user_id"`
	LayoutConfig []map[string]any `json:"layout_config"`
}