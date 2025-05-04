package db

import (
	"errors"
	"fmt"

	integrationhealthcheckapi "github.com/opengovern/opensecurity/jobs/integration-health-check/api"
	"github.com/opengovern/opensecurity/services/scheduler/db/model"
	"gorm.io/gorm"
)

func (db Database) AddIntegrationHealthCheckJob(job *model.IntegrationHealthCheckJob) error {
	tx := db.ORM.Model(&model.IntegrationHealthCheckJob{}).
		Create(job)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateIntegrationHealthCheckJobStatus(job model.IntegrationHealthCheckJob) error {
	tx := db.ORM.Model(&model.IntegrationHealthCheckJob{}).
		Where("id = ?", job.ID).
		Update("status", job.Status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateIntegrationHealthCheckJob(jobID uint, status integrationhealthcheckapi.IntegrationHealthCheckJobStatus, failedMessage string) error {
	for i := 0; i < len(failedMessage); i++ {
		if failedMessage[i] == 0 {
			failedMessage = failedMessage[:i] + failedMessage[i+1:]
		}
	}

	tx := db.ORM.Model(&model.IntegrationHealthCheckJob{}).
		Where("id = ?", jobID).
		Updates(model.IntegrationHealthCheckJob{
			Status:         status,
			FailureMessage: failedMessage,
		})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) FetchLastIntegrationHealthCheckJob() (*model.IntegrationHealthCheckJob, error) {
	var job model.IntegrationHealthCheckJob
	tx := db.ORM.Model(&model.IntegrationHealthCheckJob{}).
		Order("created_at DESC").First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) ListIntegrationHealthCheckJobs() ([]model.IntegrationHealthCheckJob, error) {
	var job []model.IntegrationHealthCheckJob
	tx := db.ORM.Model(&model.IntegrationHealthCheckJob{}).Find(&job)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return job, nil
}

// UpdateIntegrationHealthCheckJobsTimedOut updates the status of IntegrationHealthCheckJobs
// that have timed out while in the status of 'IN_PROGRESS' for longer
// than integrationhealthcheckIntervalHours hours.
func (db Database) UpdateIntegrationHealthCheckJobsTimedOut(integrationhealthcheckIntervalHours int64) error {
	tx := db.ORM.
		Model(&model.IntegrationHealthCheckJob{}).
		Where(fmt.Sprintf("created_at < NOW() - INTERVAL '%d HOURS'", integrationhealthcheckIntervalHours*2)).
		Where("status IN ?", []string{string(integrationhealthcheckapi.IntegrationHealthCheckJobInProgress)}).
		Updates(model.IntegrationHealthCheckJob{Status: integrationhealthcheckapi.IntegrationHealthCheckJobFailed, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) CleanupAllIntegrationHealthCheckJobs() error {
	tx := db.ORM.Where("1 = 1").Unscoped().Delete(&model.IntegrationHealthCheckJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
