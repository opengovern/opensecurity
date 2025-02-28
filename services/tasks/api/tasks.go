package api


type TaskListResponse struct {
	Items      []TaskResponse `json:"items"`
	TotalCount int           `json:"total_count"`
}

type TaskResponse struct {
	ID 		string `json:"id"`
	Name 	string `json:"name"`
	ResultType  string `json:"result_type"`
	Description string `json:"description"`
	ImageUrl    string `json:"image_url"`
	Interval    uint64 `json:"interval"`
	Timeout     uint64 `json:"timeout"`

}
	


type RunTaskRequest struct {
	TaskID string         `json:"task_id"`
	Params map[string]any `json:"params"`
}
