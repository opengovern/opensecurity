package api

import (
	"time"
)

type TaskListResponse struct {
	Items      []TaskResponse `json:"items"`
	TotalCount int            `json:"total_count"`
}

type TaskResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImageUrl        string `json:"image_url"`
	SchedulesNumber int    `json:"schedules_number"`
}

type TaskDetailsResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	ImageUrl     string              `json:"image_url"`
	RunSchedules []RunScheduleObject `json:"run_schedules"`
	Credentials  []string            `json:"credentials"`
	EnvVars      map[string]string   `json:"env_vars"`
	ScaleConfig  ScaleConfig         `json:"scale_config"`
}

type ScaleConfig struct {
	Stream       string `json:"stream" yaml:"stream"`
	Consumer     string `json:"consumer" yaml:"consumer"`
	LagThreshold string `json:"lag_threshold" yaml:"lag_threshold"`
	MinReplica   int32  `json:"min_replica" yaml:"min_replica"`
	MaxReplica   int32  `json:"max_replica" yaml:"max_replica"`

	PollingInterval int32 `json:"polling_interval" yaml:"polling_interval"`
	CooldownPeriod  int32 `json:"cooldown_period" yaml:"cooldown_period"`
}

type RunScheduleObject struct {
	LastRun   *time.Time     `json:"last_run"`
	Params    map[string]any `json:"params"`
	Frequency float64        `json:"frequency"`
}

type RunTaskRequest struct {
	TaskID string         `json:"task_id"`
	Params map[string]any `json:"params"`
}

type TaskConfigSecret struct {
	Credentials map[string]any `json:"credentials"`
}
