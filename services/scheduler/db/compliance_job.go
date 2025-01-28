package db

import (
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"math/rand"
	"time"

	"github.com/lib/pq"

	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"gorm.io/gorm"
)

func (db Database) CountComplianceJobsByDate(withIncidents bool, start time.Time, end time.Time) (int64, error) {
	var count int64
	tx := db.ORM.Model(&model.ComplianceJob{}).
		Where("with_incidents = ?", withIncidents).
		Where("status = ? AND updated_at >= ? AND updated_at < ?", model.ComplianceJobSucceeded, start, end).Count(&count)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, tx.Error
	}
	return count, nil
}

func (db Database) CreateComplianceJob(tx *gorm.DB, job *model.ComplianceJob) error {
	if tx == nil {
		tx = db.ORM
	}
	tx = tx.
		Model(&model.ComplianceJob{}).
		Create(job)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateComplianceJob(
	id uint, status model.ComplianceJobStatus, failureMsg string, stepFailed *model.ComplianceJobStatus) error {
	update := model.ComplianceJob{
		Status:         status,
		FailureMessage: failureMsg,
	}
	if stepFailed != nil {
		update.StepFailed = *stepFailed
	}
	tx := db.ORM.
		Model(&model.ComplianceJob{}).
		Where("id = ?", id).
		Updates(update)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateComplianceJobAreAllRunnersQueued(id uint, areAllRunnersQueued bool) error {
	tx := db.ORM.
		Model(&model.ComplianceJob{}).
		Where("id = ?", id).
		Updates(model.ComplianceJob{
			AreAllRunnersQueued: areAllRunnersQueued,
		})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateComplianceJobsTimedOut(withIncidents bool, complianceIntervalMinutes int64) error {
	tx := db.ORM.
		Model(&model.ComplianceJob{}).
		Where("with_incidents = ?", withIncidents).
		Where(fmt.Sprintf("created_at < NOW() - INTERVAL '%d MINUTES'", complianceIntervalMinutes)).
		Where("status IN ?", []string{string(model.ComplianceJobCreated),
			string(model.ComplianceJobRunnersInProgress),
			string(model.ComplianceJobQueued),
			string(model.ComplianceJobInProgress),
			string(model.ComplianceJobSummarizerInProgress),
		}).
		Updates(model.ComplianceJob{Status: model.ComplianceJobTimeOut, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) GetComplianceJobByID(ID uint) (*model.ComplianceJob, error) {
	var job model.ComplianceJob
	tx := db.ORM.Where("id = ?", ID).Find(&job)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &job, nil
}

func (db Database) GetComplianceJobByCreatedByAndParentID(createdBy string, ParentID uint) (*model.ComplianceJob, error) {
	var job model.ComplianceJob
	tx := db.ORM.Where("parent_id = ?", ParentID).Where("created_by = ?", createdBy).Find(&job)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &job, nil
}

func (db Database) CleanupComplianceJobsOlderThan(t time.Time) error {
	tx := db.ORM.Where("updated_at < ?", t).Unscoped().Delete(&model.ComplianceJob{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) GetLastComplianceJob(withIncidents bool, frameworkID string) (*model.ComplianceJob, error) {
	var job model.ComplianceJob
	tx := db.ORM.Model(&model.ComplianceJob{}).
		Where("with_incidents = ?", withIncidents).
		Where("framework_ids @> ?", pq.StringArray{frameworkID}).Order("created_at DESC").First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) ListComplianceJobs(withIncidents bool) ([]model.ComplianceJob, error) {
	var job []model.ComplianceJob
	tx := db.ORM.Model(&model.ComplianceJob{}).
		Where("with_incidents = ?", withIncidents).
		First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) ListCreatedComplianceJobs(withIncidents bool) ([]model.ComplianceJob, error) {
	var job []model.ComplianceJob
	tx := db.ORM.Model(&model.ComplianceJob{}).
		Where("with_incidents = ?", withIncidents).
		Where("status = ?", model.ComplianceJobCreated).
		First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) ListComplianceJobsForInterval(withIncidents *bool, interval, triggerType, createdBy string) ([]model.ComplianceJob, error) {
	var job []model.ComplianceJob

	tx := db.ORM.Model(&model.ComplianceJob{})

	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
	}

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

func (db Database) UpdateComplianceJobRunnersStatus(id uint, RunnersStatus pgtype.JSONB) error {
	tx := db.ORM.Model(&model.ComplianceJob{}).Where("id = ?", id).Update("runners_status", RunnersStatus)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) ListComplianceJobsWithSummaryJob(withIncidents *bool, interval, triggerType, createdBy string, frameworkIDs []string) ([]model.ComplianceJobWithSummarizerJob, error) {
	var result []model.ComplianceJobWithSummarizerJob

	// Base query
	tx := db.ORM.Table("compliance_jobs").
		Select(`
			compliance_jobs.id, 
			compliance_jobs.created_at, 
			compliance_jobs.updated_at, 
			compliance_jobs.framework_ids, 
			compliance_jobs.status, 
			compliance_jobs.integration_ids, 
			compliance_jobs.trigger_type, 
			compliance_jobs.created_by,
			COALESCE(array_agg(COALESCE(compliance_summarizers.id::text, '')), '{}') as summarizer_jobs
		`).
		Joins("LEFT JOIN compliance_summarizers ON compliance_jobs.id = compliance_summarizers.parent_job_id").
		Group("compliance_jobs.id")

	// Apply filters
	if withIncidents != nil {
		tx = tx.Where("compliance_jobs.with_incidents = ?", *withIncidents)
	}

	if interval != "" {
		tx = tx.Where(fmt.Sprintf("NOW() - compliance_jobs.updated_at < INTERVAL '%s'", interval))
	}
	if triggerType != "" {
		tx = tx.Where("compliance_jobs.trigger_type = ?", triggerType)
	}
	if createdBy != "" {
		tx = tx.Where("compliance_jobs.created_by = ?", createdBy)
	}
	if len(frameworkIDs) > 0 {
		tx = tx.Where("compliance_jobs.framework_ids && ?", pq.Array(frameworkIDs))
	}

	// Execute the query
	if err := tx.Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (db Database) ListComplianceJobsByIntegrationID(withIncidents *bool, integrationIds []string) ([]model.ComplianceJob, error) {
	var job []model.ComplianceJob
	tx := db.ORM.Model(&model.ComplianceJob{}).Where("integration_ids && ?", pq.Array(integrationIds))
	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
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

func (db Database) ListPendingComplianceJobsByIntegrationID(withIncidents *bool, integrationIds []string) ([]model.ComplianceJob, error) {
	var job []model.ComplianceJob
	tx := db.ORM.Model(&model.ComplianceJob{}).
		Where("integration_ids && ?", integrationIds).
		Where("status IN ?", []model.ComplianceJobStatus{
			model.ComplianceJobCreated,
			model.ComplianceJobQueued,
			model.ComplianceJobInProgress,
			model.ComplianceJobRunnersInProgress})
	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
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

func (db Database) ListComplianceJobsByFrameworkID(withIncidents *bool, frameworkIDs []string) ([]model.ComplianceJob, error) {
	var job []model.ComplianceJob
	tx := db.ORM.Model(&model.ComplianceJob{}).Where("framework_ids && ?", pq.Array(frameworkIDs))
	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
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

func (db Database) ListComplianceJobsByStatus(withIncidents *bool, status model.ComplianceJobStatus) ([]model.ComplianceJob, error) {
	var jobs []model.ComplianceJob
	tx := db.ORM.Where("status = ?", status)
	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
	}
	tx = tx.Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) ListComplianceJobsByIds(withIncidents *bool, ids []string) ([]model.ComplianceJob, error) {
	var jobs []model.ComplianceJob
	tx := db.ORM.Where("id IN ?", ids)
	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
	}
	tx = tx.Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) ListComplianceRunnersWithStatus(withIncidents *bool, status model.ComplianceJobStatus) ([]model.ComplianceJob, error) {
	var jobs []model.ComplianceJob
	tx := db.ORM.Where("status = ?", status)
	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
	}
	tx = tx.Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) ListComplianceJobsWithUnqueuedRunners(withIncidents bool) ([]model.ComplianceJob, error) {
	var jobs []model.ComplianceJob
	tx := db.ORM.Where("are_all_runners_queued = ?", false).Where("with_incidents = ?", withIncidents).
		Where("status IN ?", []string{string(model.ComplianceJobCreated), string(model.ComplianceJobRunnersInProgress)}).
		Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	// shuffle jobs
	rand.Shuffle(len(jobs), func(i, j int) {
		jobs[i], jobs[j] = jobs[j], jobs[i]
	})
	return jobs, nil
}

func (db Database) SetJobToRunnersInProgress() error {
	tx := db.ORM.Exec(`
UPDATE compliance_jobs j SET status = 'RUNNERS_IN_PROGRESS' WHERE status = 'CREATED' AND
	(select count(*) from compliance_runners where parent_job_id = j.id) > 0
`)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListJobsWithRunnersCompleted(manuals bool) ([]model.ComplianceJob, error) {
	var jobs []model.ComplianceJob

	query := `
SELECT * FROM compliance_jobs j WHERE status IN ('RUNNERS_IN_PROGRESS', 'SINK_IN_PROGRESS') AND with_incidents = true AND are_all_runners_queued = TRUE AND
	(select count(*) from compliance_runners where parent_job_id = j.id AND 
	                                               NOT (status = 'SUCCEEDED' OR status = 'TIMEOUT' OR (status = 'FAILED' and retry_count >= ?))
	                                         ) = 0
`
	if manuals {
		query = query + ` AND trigger_type = ?`
	} else {
		query = query + ` AND trigger_type <> ?`
	}
	tx := db.ORM.Raw(query, runnerRetryCount, model.ComplianceTriggerTypeManual).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) GetLastUpdatedRunnerForParent(jobId uint) (model.ComplianceRunner, error) {
	var runner model.ComplianceRunner
	tx := db.ORM.Where("parent_job_id = ?", jobId).Order("updated_at DESC").First(&runner)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return model.ComplianceRunner{}, nil
		}
		return model.ComplianceRunner{}, tx.Error
	}

	return runner, nil
}

func (db Database) GetRunnersByParentJobID(jobID uint) ([]model.ComplianceRunner, error) {
	var runners []model.ComplianceRunner
	tx := db.ORM.Where("parent_job_id = ?", jobID).Find(&runners)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return runners, nil
}

func (db Database) FetchTotalFindingCountForComplianceJob(jobID uint) (int, error) {
	var count int
	tx := db.ORM.Raw(`SELECT coalesce(sum(coalesce(total_finding_count,0)), 0) FROM compliance_runners WHERE parent_job_id = ?`, jobID).Scan(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return count, nil
}

func (db Database) ListJobsToFinish() ([]model.ComplianceJob, error) {
	var jobs []model.ComplianceJob
	tx := db.ORM.Raw(`
SELECT * FROM compliance_jobs j WHERE status = 'SUMMARIZER_IN_PROGRESS' AND with_incidents = true AND
	(select count(*) from compliance_summarizers where parent_job_id = j.id AND (status = 'SUCCEEDED' OR (status = 'FAILED' and retry_count >= ?))) > 0
`, summarizerRetryCount).Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) ListComplianceJobsByFilters(withIncidents *bool, integrationId []string, frameworkId []string, status []string,
	startTime, endTime *time.Time) ([]model.ComplianceJob, error) {
	var jobs []model.ComplianceJob
	tx := db.ORM.Model(&model.ComplianceJob{})

	if withIncidents != nil {
		tx = tx.Where("with_incidents = ?", *withIncidents)
	}

	if len(integrationId) > 0 {
		tx = tx.Where("integration_ids && ?", pq.Array(integrationId))
	}

	if len(frameworkId) > 0 {
		tx = tx.Where("framework_ids && ?", pq.Array(frameworkId))
	}
	if len(status) > 0 {
		tx = tx.Where("status IN ?", status)
	}
	if startTime != nil {
		tx = tx.Where("updated_at >= ?", *startTime)
	}
	if endTime != nil {
		tx = tx.Where("updated_at <= ?", *endTime)
	}

	tx = tx.Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) GetComplianceJobsIntegrations() ([]string, error) {
	var uniqueIntegrationIDs []string
	query := `
		SELECT DISTINCT unnest(integration_ids) AS integration
		FROM compliance_jobs
	`
	if err := db.ORM.Raw(query).Pluck("integration", &uniqueIntegrationIDs).Error; err != nil {
		return nil, err
	}
	return uniqueIntegrationIDs, nil
}

func (db Database) CleanupAllComplianceJobsForIntegrations(integrationId []string) error {
	tx := db.ORM.Where("integration_ids && ?", pq.Array(integrationId)).Unscoped().Delete(&model.ComplianceJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
