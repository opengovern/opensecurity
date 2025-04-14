package models

import (
	"github.com/jackc/pgtype"
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
	ID          string `gorm:"primarykey"`
	Name        string `gorm:"unique;not null"`
	IsEnabled   bool   `gorm:"not null"`
	Description string
	ImageUrl    string
	Command     string
	Timeout     float64
	NatsConfig  pgtype.JSONB
	ScaleConfig pgtype.JSONB
}

type TaskConfigSecret struct {
	TaskID       string `gorm:"primarykey"`
	Secret       string
	HealthStatus TaskSecretHealthStatus
}

type TaskRunSchedule struct {
	TaskID    string
	Params    pgtype.JSONB
	Frequency float64
}
