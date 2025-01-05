package model

import (
	queryvalidator "github.com/opengovern/opencomply/jobs/query-validator-job"
	"gorm.io/gorm"
)

type QueryValidatorJob struct {
	gorm.Model
	QueryId        string
	QueryType      queryvalidator.QueryType
	Status         queryvalidator.QueryValidatorStatus
	HasParams      bool
	FailureMessage string
}
