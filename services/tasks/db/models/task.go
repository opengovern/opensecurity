package models

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
	"time"
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
	ID        uint `gorm:"primarykey"`
	TaskID    string
	LastRun   *time.Time
	Params    pgtype.JSONB
	Frequency float64
}
