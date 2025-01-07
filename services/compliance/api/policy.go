package api

type ListPolicyItem struct {
	ID            string `gorm:"id"`
	Title         string `json:"title"`
	Language      string `json:"language"`
	ControlsCount int    `json:"controls_count"`
}

type ListPoliciesResponse struct {
	Policies   []ListPolicyItem `json:"policies"`
	TotalCount int              `json:"total_count"`
}

type GetPolicyItem struct {
	ID             string   `gorm:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Language       string   `json:"language"`
	Definition     string   `json:"definition"`
	ControlsCount  int      `json:"controls_count"`
	ListOfControls []string `json:"list_of_controls"`
}
