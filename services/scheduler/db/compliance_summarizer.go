package db

import (
	"errors"
	"time"

	summarizer "github.com/opengovern/opensecurity/jobs/compliance-summarizer"
	"github.com/opengovern/opensecurity/services/scheduler/db/model"
	"gorm.io/gorm"
)

const (
	summarizerRetryCount = 0
)

func (db Database) CreateSummarizerJob(summarizer *model.ComplianceSummarizer) error {
	tx := db.ORM.
		Model(&model.ComplianceSummarizer{}).
		Create(summarizer)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) FetchCreatedSummarizers(manuals bool) ([]model.ComplianceSummarizer, error) {
	var jobs []model.ComplianceSummarizer
	tx := db.ORM.Model(&model.ComplianceSummarizer{}).
		Where("status = ?", summarizer.ComplianceSummarizerCreated)

	if manuals {
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

func (db Database) RetryFailedSummarizers() error {
	tx := db.ORM.Exec("UPDATE compliance_summarizers SET retry_count = retry_count + 1, status = 'CREATED' WHERE status = 'FAILED' AND retry_count < ? AND updated_at < NOW() - interval '7 minutes'", summarizerRetryCount)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateSummarizerJob(
	id uint, status summarizer.ComplianceSummarizerStatus, startedAt time.Time, failureMsg string) error {
	tx := db.ORM.
		Model(&model.ComplianceSummarizer{}).
		Where("id = ?", id).
		Updates(model.ComplianceSummarizer{
			Status:         status,
			StartedAt:      startedAt,
			FailureMessage: failureMsg,
		})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateSummarizerJobsTimedOut() error {
	tx := db.ORM.
		Model(&model.ComplianceSummarizer{}).
		Where("created_at < NOW() - INTERVAL '6 HOURS'").
		Where("status IN ?", []string{string(summarizer.ComplianceSummarizerCreated), string(summarizer.ComplianceSummarizerInProgress)}).
		Updates(model.ComplianceSummarizer{Status: summarizer.ComplianceSummarizerFailed, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListFailedSummarizersWithParentID(id uint) ([]model.ComplianceSummarizer, error) {
	var jobs []model.ComplianceSummarizer
	tx := db.ORM.Model(&model.ComplianceSummarizer{}).
		Where("status = ?", summarizer.ComplianceSummarizerFailed).
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

func (db Database) CleanupAllComplianceSummarizerJobsByComplianceJobs(complianceJobs []uint) error {
	tx := db.ORM.Where("parent_job_id IN ?", complianceJobs).Unscoped().Delete(&model.ComplianceSummarizer{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) ListAllComplianceSummarizerJobsByComplianceJobs(complianceJobs []uint) ([]model.ComplianceSummarizer, error) {
	var jobs []model.ComplianceSummarizer
	tx := db.ORM.Where("parent_job_id IN ?", complianceJobs).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return jobs, nil
}
