package model

import (
	"github.com/jackc/pgtype"
	"time"
)

type JobsStatus string

const (
	JobStatusCompleted  JobsStatus = "SUCCEEDED"
	JobStatusPending    JobsStatus = "PENDING"
	JobStatusInProgress JobsStatus = "IN_PROGRESS"
	JobStatusFailed     JobsStatus = "FAILED"
)

type JobInfo struct {
	MigrationJobName string     `json:"migrationJobName"`
	Status           JobsStatus `json:"status"`
	FailureReason    string     `json:"failureReason"`
}

type Migration struct {
	ID         string `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     string
	JobsStatus pgtype.JSONB

	AdditionalInfo string
}
