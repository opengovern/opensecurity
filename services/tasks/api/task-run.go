package api

import (
	"time"
)

type TaskRun struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	TaskID         string    `json:"task_id"`
	Status         string    `json:"status"`
	Result         map[string]any    `json:"result"`
	Params 	   map[string]any    `json:"params"`
	FailureMessage string    `json:"failure_message"`
}

type ListTaskRunsResponse struct {
	TotalCount int              `json:"total_count"`
	Items      []TaskRun `json:"items"`
}
