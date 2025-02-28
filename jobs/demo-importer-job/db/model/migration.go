package model

import (
	"github.com/jackc/pgtype"
	"time"
)

type JobsStatus string

const (
	JobStatusCompleted JobsStatus = "SUCCEEDED"
	JobStatusFailed    JobsStatus = "FAILED"

	MigrationJobName = "import-sample-data"
)

type Migration struct {
	ID         string `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Status     string
	JobsStatus pgtype.JSONB

	AdditionalInfo string
}

type ESImportProgress struct {
	Progress float64 `json:"progress"`
}
