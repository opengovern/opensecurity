package api

type TaskListResponse struct {
	Items      []TaskResponse `json:"items"`
	TotalCount int            `json:"total_count"`
}

type TaskResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    string  `json:"image_url"`
	Timeout     float64 `json:"timeout"`
}

type RunTaskRequest struct {
	TaskID string         `json:"task_id"`
	Params map[string]any `json:"params"`
}

type TaskConfigSecret struct {
	Credentials map[string]any `json:"credentials"`
}
