package model

import (
	checkupapi "github.com/opengovern/opensecurity/jobs/integration-health-check/api"
	"gorm.io/gorm"
)

type CheckupJob struct {
	gorm.Model
	Status         checkupapi.IntegrationHealthCheckJobStatus
	FailureMessage string
}
