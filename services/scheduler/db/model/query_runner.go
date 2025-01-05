package model

import (
	queryrunner "github.com/opengovern/opencomply/jobs/query-runner-job"
	"gorm.io/gorm"
)

type QueryRunnerJob struct {
	gorm.Model
	QueryId            string
	CreatedBy          string
	RetryCount         int
	Status             queryrunner.QueryRunnerStatus
	FailureMessage     string
	NatsSequenceNumber uint64
}
