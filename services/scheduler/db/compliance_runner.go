package db

import (
	"errors"
	"time"

	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"gorm.io/gorm"
)

const (
	runnerRetryCount = 0
)

func (db Database) CreateRunnerJobs(tx *gorm.DB, runners []*model.ComplianceRunner) error {
	if tx == nil {
		tx = db.ORM
	}
	tx = tx.
		Model(&model.ComplianceRunner{}).
		CreateInBatches(runners, 500)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) DeleteOldRunnerJob(parentJobId *uint) error {
	tx := db.ORM.Model(&model.ComplianceRunner{})
	if parentJobId != nil {
		tx = tx.Where("parent_job_id = ?", *parentJobId)
	} else {
		tx = tx.Where("created_at < ?", time.Now().Add(-time.Hour*24*2))
	}
	tx = tx.Unscoped().Delete(&model.ComplianceRunner{})
	if tx.Error != nil {
		return tx.Error
	}

	tx = db.ORM.Model(&model.ComplianceRunner{}).
		Where("created_at < ?", time.Now().Add(-time.Hour*24*7)).
		Unscoped().Delete(&model.ComplianceRunner{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) FetchCreatedRunners(manual bool) ([]model.ComplianceRunner, error) {
	var jobs []model.ComplianceRunner
	tx := db.ORM.Model(&model.ComplianceRunner{}).
		Where("status = ?", model.ComplianceRunnerCreated)

	if manual {
		tx = tx.Where("trigger_type = ?", model.ComplianceTriggerTypeManual)
	} else {
		tx = tx.Where("trigger_type <> ?", model.ComplianceTriggerTypeManual)
	}

	tx = tx.Order("created_at ASC").Limit(1000).Find(&jobs)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) UpdateTimedOutInProgressRunners() error {
	tx := db.ORM.
		Model(&model.ComplianceRunner{}).
		Where("status = ?", model.ComplianceRunnerInProgress).
		Where("updated_at < NOW() - INTERVAL '10 MINUTES'").
		Updates(model.ComplianceRunner{Status: model.ComplianceRunnerTimeOut, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) RetryFailedRunners() error {
	tx := db.ORM.Exec("UPDATE compliance_runners SET retry_count = retry_count + 1, status = 'CREATED', updated_at = NOW() WHERE status = 'FAILED' AND retry_count < ? AND updated_at < NOW() - interval '5 minutes'", runnerRetryCount)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateRunnerJob(
	id uint, status model.ComplianceRunnerStatus, startedAt time.Time, totalFindingCount *int, failureMsg string, podName *string) error {
	updated := model.ComplianceRunner{
		Status:            status,
		StartedAt:         startedAt,
		FailureMessage:    failureMsg,
		TotalFindingCount: totalFindingCount,
	}
	if podName != nil {
		updated.WorkerPodName = *podName
	}

	tx := db.ORM.
		Model(&model.ComplianceRunner{}).
		Where("id = ?", id).
		Updates(updated)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateRunnerJobNatsSeqNum(
	id uint, seqNum uint64) error {
	tx := db.ORM.
		Model(&model.ComplianceRunner{}).
		Where("id = ?", id).
		Updates(model.ComplianceRunner{
			NatsSequenceNumber: seqNum,
		})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateTimeoutQueuedRunnerJobs() error {
	tx := db.ORM.
		Model(&model.ComplianceRunner{}).
		Where("created_at < NOW() - INTERVAL '60 MINUTES'").
		Where("status IN ?", []string{string(model.ComplianceRunnerCreated), string(model.ComplianceRunnerQueued)}).
		Updates(model.ComplianceRunner{Status: model.ComplianceRunnerTimeOut, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListRunnersWithID(ids []int64) ([]model.ComplianceRunner, error) {
	var jobs []model.ComplianceRunner
	tx := db.ORM.Where("id IN ?", ids).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) ListFailedRunnersWithParentID(id uint) ([]model.ComplianceRunner, error) {
	var jobs []model.ComplianceRunner
	tx := db.ORM.Model(&model.ComplianceRunner{}).
		Where("status = ?", model.ComplianceRunnerFailed).
		Where("parent_job_id = ?", id).
		Find(&jobs)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) ListComplianceJobRunnersWithParentID(parentId uint) ([]model.ComplianceRunner, error) {
	var jobs []model.ComplianceRunner
	tx := db.ORM.Where("parent_job_id = ?", parentId).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) CleanupAllComplianceRunnersByComplianceJobs(complianceJobs []uint) error {
	tx := db.ORM.Where("parent_job_id IN ?", complianceJobs).Unscoped().Delete(&model.ComplianceRunner{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
