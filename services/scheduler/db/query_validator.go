package db

import (
	"errors"
	"fmt"

	queryrunner "github.com/opengovern/opensecurity/jobs/query-runner-job"

	queryvalidator "github.com/opengovern/opensecurity/jobs/query-validator-job"
	"github.com/opengovern/opensecurity/services/scheduler/db/model"
	"gorm.io/gorm"
)

func (db Database) CreateQueryValidatorJob(job *model.QueryValidatorJob) (uint, error) {
	tx := db.ORM.Create(job)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return job.ID, nil
}

func (db Database) GetQueryValidatorJob(id uint) (*model.QueryValidatorJob, error) {
	var job model.QueryValidatorJob
	tx := db.ORM.Model(&model.QueryValidatorJob{}).Where("id = ?", id).First(&job)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) ListQueryValidatorJobsById(ids []string) ([]model.QueryValidatorJob, error) {
	var jobs []model.QueryValidatorJob
	tx := db.ORM.Model(&model.QueryValidatorJob{}).Where("id IN ?", ids).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) ListQueryValidatorJobs() ([]model.QueryValidatorJob, error) {
	var jobs []model.QueryValidatorJob
	tx := db.ORM.Model(&model.QueryValidatorJob{}).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) FetchCreatedQueryValidatorJobs(limit int64) ([]model.QueryValidatorJob, error) {
	var jobs []model.QueryValidatorJob
	tx := db.ORM.Model(&model.QueryValidatorJob{}).Where("status = ?", queryvalidator.QueryValidatorCreated).Limit(int(limit)).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) GetInProgressJobsCount() (int64, error) {
	var count int64
	tx := db.ORM.Model(&model.QueryValidatorJob{}).Where("status IN ?", []string{string(queryvalidator.QueryValidatorInProgress),
		string(queryvalidator.QueryValidatorQueued)}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}

func (db Database) DeleteQueryValidatorJob(id uint) error {
	tx := db.ORM.Model(&model.QueryValidatorJob{}).Delete(&model.QueryValidatorJob{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateQueryValidatorJobStatus(jobId uint, status queryvalidator.QueryValidatorStatus, failureReason string) error {
	tx := db.ORM.Model(&model.QueryValidatorJob{}).Where("id = ?", jobId).
		Updates(model.QueryValidatorJob{Status: status, FailureMessage: failureReason})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateTimedOutInProgressQueryValidators() error {
	tx := db.ORM.
		Model(&model.QueryValidatorJob{}).
		Where("status = ?", queryrunner.QueryRunnerInProgress).
		Where("updated_at < NOW() - INTERVAL '5 MINUTES'").
		Updates(model.QueryValidatorJob{Status: queryvalidator.QueryValidatorTimeOut, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateTimedOutQueuedQueryValidators() error {
	tx := db.ORM.
		Model(&model.QueryValidatorJob{}).
		Where("status = ?", queryrunner.QueryRunnerQueued).
		Where("updated_at < NOW() - INTERVAL '12 HOURS'").
		Updates(model.QueryValidatorJob{Status: queryvalidator.QueryValidatorTimeOut, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListQueryValidatorJobForInterval(interval, triggerType, createdBy string) ([]model.QueryValidatorJob, error) {
	var job []model.QueryValidatorJob

	tx := db.ORM.Model(&model.QueryValidatorJob{})

	if interval != "" {
		tx = tx.Where(fmt.Sprintf("NOW() - updated_at < INTERVAL '%s'", interval))
	}
	if triggerType != "" {
		tx = tx.Where("trigger_type = ?", triggerType)
	}
	if createdBy != "" {
		tx = tx.Where("created_by = ?", createdBy)
	}

	tx = tx.Find(&job)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) CleanupAllQueryValidatorJobs() error {
	tx := db.ORM.Where("1 = 1").Unscoped().Delete(&model.QueryValidatorJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
