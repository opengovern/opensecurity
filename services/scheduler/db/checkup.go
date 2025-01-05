package db

import (
	"errors"
	"fmt"

	checkupapi "github.com/opengovern/opencomply/jobs/checkup-job/api"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"gorm.io/gorm"
)

func (db Database) AddCheckupJob(job *model.CheckupJob) error {
	tx := db.ORM.Model(&model.CheckupJob{}).
		Create(job)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateCheckupJobStatus(job model.CheckupJob) error {
	tx := db.ORM.Model(&model.CheckupJob{}).
		Where("id = ?", job.ID).
		Update("status", job.Status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateCheckupJob(jobID uint, status checkupapi.CheckupJobStatus, failedMessage string) error {
	for i := 0; i < len(failedMessage); i++ {
		if failedMessage[i] == 0 {
			failedMessage = failedMessage[:i] + failedMessage[i+1:]
		}
	}

	tx := db.ORM.Model(&model.CheckupJob{}).
		Where("id = ?", jobID).
		Updates(model.CheckupJob{
			Status:         status,
			FailureMessage: failedMessage,
		})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) FetchLastCheckupJob() (*model.CheckupJob, error) {
	var job model.CheckupJob
	tx := db.ORM.Model(&model.CheckupJob{}).
		Order("created_at DESC").First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) ListCheckupJobs() ([]model.CheckupJob, error) {
	var job []model.CheckupJob
	tx := db.ORM.Model(&model.CheckupJob{}).Find(&job)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return job, nil
}

// UpdateCheckupJobsTimedOut updates the status of CheckupJobs
// that have timed out while in the status of 'IN_PROGRESS' for longer
// than checkupIntervalHours hours.
func (db Database) UpdateCheckupJobsTimedOut(checkupIntervalHours int64) error {
	tx := db.ORM.
		Model(&model.CheckupJob{}).
		Where(fmt.Sprintf("created_at < NOW() - INTERVAL '%d HOURS'", checkupIntervalHours*2)).
		Where("status IN ?", []string{string(checkupapi.CheckupJobInProgress)}).
		Updates(model.CheckupJob{Status: checkupapi.CheckupJobFailed, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) CleanupAllCheckupJobs() error {
	tx := db.ORM.Where("1 = 1").Unscoped().Delete(&model.CheckupJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
