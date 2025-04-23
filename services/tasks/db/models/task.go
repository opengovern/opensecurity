package models

import (
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TaskSecretHealthStatus string

const (
	TaskSecretHealthStatusUnknown   TaskSecretHealthStatus = "unknown"
	TaskSecretHealthStatusHealthy   TaskSecretHealthStatus = "healthy"
	TaskSecretHealthStatusUnhealthy TaskSecretHealthStatus = "unhealthy"
)

type Task struct {
	gorm.Model
	ID                  string `gorm:"primarykey"`
	Name                string `gorm:"unique;not null"`
	IsEnabled           bool   `gorm:"not null"`
	Description         string
	ImageUrl            string
	SteampipePluginName string
	ArtifactsUrl        string
	Command             string
	Timeout             float64
	NatsConfig          pgtype.JSONB
	ScaleConfig         pgtype.JSONB
	EnvVars             pgtype.JSONB
	Params              pq.StringArray `gorm:"type:text[]"`
	Configs             pq.StringArray `gorm:"type:text[]"`
}

type TaskBinary struct {
	TaskID string `gorm:"primaryKey"`

	CloudQlPlugin []byte `gorm:"type:bytea"`
}

type TaskConfigSecret struct {
	TaskID       string `gorm:"primarykey"`
	Secret       string
	HealthStatus TaskSecretHealthStatus
}

type TaskRunSchedule struct {
	SchedulerID string `gorm:"primarykey"`
	TaskID      string `gorm:"primarykey"`
	Params      pgtype.JSONB
	Frequency   float64
}
