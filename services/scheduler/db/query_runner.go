package db

import (
	"errors"
	"fmt"

	queryrunner "github.com/opengovern/opensecurity/jobs/query-runner-job"
	"github.com/opengovern/opensecurity/services/scheduler/db/model"
	"gorm.io/gorm"
)

func (db Database) CreateQueryRunnerJob(job *model.QueryRunnerJob) (uint, error) {
	tx := db.ORM.Create(job)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return job.ID, nil
}

func (db Database) GetQueryRunnerJob(id uint) (*model.QueryRunnerJob, error) {
	var job model.QueryRunnerJob
	tx := db.ORM.Model(&model.QueryRunnerJob{}).Where("id = ?", id).First(&job)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) ListQueryRunnerJobsById(ids []string) ([]model.QueryRunnerJob, error) {
	var jobs []model.QueryRunnerJob
	tx := db.ORM.Model(&model.QueryRunnerJob{}).Where("id IN ?", ids).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) ListQueryRunnerJobs() ([]model.QueryRunnerJob, error) {
	var jobs []model.QueryRunnerJob
	tx := db.ORM.Model(&model.QueryRunnerJob{}).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) FetchCreatedQueryRunnerJobs() ([]model.QueryRunnerJob, error) {
	var jobs []model.QueryRunnerJob
	tx := db.ORM.Model(&model.QueryRunnerJob{}).Where("status = ?", queryrunner.QueryRunnerCreated).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) DeleteQueryRunnerJob(id uint) error {
	tx := db.ORM.Model(&model.QueryRunnerJob{}).Delete(&model.QueryRunnerJob{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateQueryRunnerJobStatus(jobId uint, status queryrunner.QueryRunnerStatus, failureReason string) error {
	tx := db.ORM.Model(&model.QueryRunnerJob{}).Where("id = ?", jobId).
		Updates(model.QueryRunnerJob{Status: status, FailureMessage: failureReason})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateQueryRunnerJobNatsSeqNum(
	id uint, seqNum uint64) error {
	tx := db.ORM.
		Model(&model.QueryRunnerJob{}).
		Where("id = ?", id).
		Updates(model.QueryRunnerJob{
			NatsSequenceNumber: seqNum,
		})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateTimedOutInProgressQueryRunners() error {
	tx := db.ORM.
		Model(&model.QueryRunnerJob{}).
		Where("status = ?", queryrunner.QueryRunnerInProgress).
		Where("updated_at < NOW() - INTERVAL '5 MINUTES'").
		Updates(model.QueryRunnerJob{Status: queryrunner.QueryRunnerTimeOut, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateTimedOutQueuedQueryRunners() error {
	tx := db.ORM.
		Model(&model.QueryRunnerJob{}).
		Where("status = ?", queryrunner.QueryRunnerQueued).
		Where("updated_at < NOW() - INTERVAL '12 HOURS'").
		Updates(model.QueryRunnerJob{Status: queryrunner.QueryRunnerTimeOut, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListQueryRunnerJobForInterval(interval, triggerType, createdBy string) ([]model.QueryRunnerJob, error) {
	var job []model.QueryRunnerJob

	tx := db.ORM.Model(&model.QueryRunnerJob{})

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

func (db Database) CleanupAllQueryRunnerJobs() error {
	tx := db.ORM.Where("1 = 1").Unscoped().Delete(&model.QueryRunnerJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
