package model

import (
	"time"
)

type Job struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	JobType        string
	ConnectionID   string
	Title          string
	FailureMessage string
	Status         string
}

type JobSummary struct {
	JobType string
	Status  string
	Count   int64
}
