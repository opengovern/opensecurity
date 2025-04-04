package api



type GetUserLayout struct {
	UserID string `json:"user_id"`
}
type SetUserLayout struct {
	UserID      string `json:"user_id"`
	LayoutConfig []map[string]any `json:"layout_config"`
}