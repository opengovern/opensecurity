package model

import (
	checkupapi "github.com/opengovern/opencomply/jobs/checkup-job/api"
	"gorm.io/gorm"
)

type CheckupJob struct {
	gorm.Model
	Status         checkupapi.CheckupJobStatus
	FailureMessage string
}
