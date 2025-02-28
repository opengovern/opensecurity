package models

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type TaskRunStatus string

const (
	TaskRunStatusCreated    TaskRunStatus = "CREATED"
	TaskRunStatusQueued     TaskRunStatus = "QUEUED"
	TaskRunStatusInProgress TaskRunStatus = "IN_PROGRESS"
	TaskRunStatusFinished   TaskRunStatus = "FINISHED"
	TaskRunStatusFailed     TaskRunStatus = "FAILED"
	TaskRunStatusTimeout    TaskRunStatus = "TIMEOUT"
)

type TaskRun struct {
	gorm.Model
	TaskID         string
	Params         pgtype.JSONB
	Status         TaskRunStatus
	Result         pgtype.JSONB
	FailureMessage string
}
