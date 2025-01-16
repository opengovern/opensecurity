package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgtype"
	"strings"
	"time"

	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/opengovern/og-util/pkg/source"
	opengovernanceTrace "github.com/opengovern/og-util/pkg/trace"
	"github.com/opengovern/opencomply/services/scheduler/api"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db Database) CountDescribeJobsByDate(includeCost *bool, start time.Time, end time.Time) (int64, error) {
	var count int64
	costStmt := ""
	if includeCost != nil {
		if *includeCost {
			costStmt = "resource_type like '%Cost%' AND "
		} else {
			costStmt = "NOT(resource_type like '%Cost%') AND "
		}
	}
	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).
		Where(costStmt+"status = ? AND updated_at >= ? AND updated_at < ?", api.DescribeResourceJobSucceeded, start, end).Count(&count)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, tx.Error
	}
	return count, nil
}

func (db Database) CountQueuedDescribeIntegrationJobs(manuals bool) (int64, error) {
	var count int64
	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).
		Where("status = ? AND created_at > now() - interval '1 day'", api.DescribeResourceJobQueued)
	if manuals {
		tx = tx.Where("trigger_type = ?", enums.DescribeTriggerTypeManual)
	} else {
		tx = tx.Where("trigger_type <> ?", enums.DescribeTriggerTypeManual)
	}
	tx = tx.Count(&count)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, tx.Error
	}
	return count, nil
}

func (db Database) CountDescribeIntegrationJobsRunOverLast10Minutes(manuals bool) (int64, error) {
	var count int64
	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).
		Where("status != ? AND updated_at > now() - interval '10 minutes'", api.DescribeResourceJobCreated)
	if manuals {
		tx = tx.Where("trigger_type = ?", enums.DescribeTriggerTypeManual)
	} else {
		tx = tx.Where("trigger_type <> ?", enums.DescribeTriggerTypeManual)
	}
	tx = tx.Count(&count)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, tx.Error
	}
	return count, nil
}

type ResourceTypeCount struct {
	ResourceType string
	Count        int
}

func (db Database) CountRunningDescribeJobsPerResourceType(manuals bool) ([]ResourceTypeCount, error) {
	var count []ResourceTypeCount
	runningJobs := []api.DescribeResourceJobStatus{api.DescribeResourceJobQueued, api.DescribeResourceJobInProgress, api.DescribeResourceJobOldResourceDeletion}
	query := `select resource_type, count(*) as count from describe_integration_jobs where status in ?`
	if manuals {
		query = query + ` AND trigger_type = ?`
	} else {
		query = query + ` AND trigger_type <> ?`
	}
	query = query + ` group by 1`
	tx := db.ORM.Raw(query, runningJobs, enums.DescribeTriggerTypeManual)

	tx = tx.Find(&count)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return count, nil
}

func (db Database) GetLastDescribeIntegrationJob(integrationId, resourceType string, parameters pgtype.JSONB) (*model.DescribeIntegrationJob, error) {
	var job model.DescribeIntegrationJob

	// Ensure the parameters are correctly encoded as JSON
	var parametersBytes []byte
	if parameters.Status == pgtype.Present {
		parametersBytes = parameters.Bytes
	} else {
		return nil, fmt.Errorf("parameters are not in a valid state: %v", parameters.Status)
	}

	// Query with JSONB condition
	tx := db.ORM.Preload(clause.Associations).
		Where("integration_id = ? AND resource_type = ? AND convert_from(parameters, 'UTF8')::JSONB @> ?::JSONB", integrationId, resourceType, parametersBytes).
		Order("updated_at DESC").
		First(&job)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &job, nil
}

func (db Database) GetDescribeIntegrationJobByIntegrationID(integrationId string) ([]model.DescribeIntegrationJob, error) {
	var jobs []model.DescribeIntegrationJob
	tx := db.ORM.Preload(clause.Associations).Where("integration_id = ?", integrationId).Find(&jobs)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) GetDescribeIntegrationJobByID(id uint) (*model.DescribeIntegrationJob, error) {
	var job model.DescribeIntegrationJob
	tx := db.ORM.Preload(clause.Associations).Where("id = ?", id).First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return &job, nil
}

// ListDescribeJobsWithParamsByIntegrationID List only describe jobs that contain parameters
func (db Database) ListDescribeJobsWithParamsByIntegrationID(integrationId string) ([]model.DescribeIntegrationJobWithParams, error) {
	var jobs []model.DescribeIntegrationJobWithParams
	query := `
SELECT resource_type, parameters FROM describe_integration_jobs 
                                 WHERE integration_id = ? AND parameters <> '\x7b7d'
                                 GROUP BY resource_type, parameters 
`

	tx := db.ORM.Raw(query, integrationId).Find(&jobs)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) RetryDescribeIntegrationJob(id uint) error {
	tx := db.ORM.Exec("update describe_integration_jobs set status = ? where id = ?", api.DescribeResourceJobCreated, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
func (db Database) QueueDescribeIntegrationJob(id uint) error {
	tx := db.ORM.Exec("update describe_integration_jobs set status = ?, queued_at = NOW(), retry_count = retry_count + 1 where id = ?", api.DescribeResourceJobQueued, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateDescribeIntegrationJobNatsSeqNum(id uint, seqNum uint64) error {
	tx := db.ORM.Exec("update describe_integration_jobs set nats_sequence_number = ? where id = ?", seqNum, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListRandomCreatedDescribeIntegrationJobs(ctx context.Context, limit int, manuals bool) ([]model.DescribeIntegrationJob, error) {
	ctx, span := otel.Tracer(opengovernanceTrace.JaegerTracerName).Start(ctx, opengovernanceTrace.GetCurrentFuncName())
	defer span.End()

	var job []model.DescribeIntegrationJob

	query := `
SELECT
	*, random() as r
FROM
	describe_integration_jobs dr
WHERE
	status = ?`

	if manuals {
		query = query + ` AND trigger_type = ?`
	} else {
		query = query + ` AND trigger_type <> ?`
	}

	query = query + ` ORDER BY r DESC
LIMIT ?
`
	tx := db.ORM.Raw(query, api.DescribeResourceJobCreated, enums.DescribeTriggerTypeManual, limit).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) ListAllJobs(pageStart, pageEnd int, interval *string, from *time.Time, to *time.Time, typeFilter []string,
	statusFilter []string, sortBy, sortOrder string) ([]model.Job, error) {
	var job []model.Job

	whereQuery := ""
	var values []interface{}
	if from != nil && to != nil && interval == nil {
		values = append(values, *from, *to, *from, *to, *from, *to)
	}

	if len(typeFilter) > 0 || len(statusFilter) > 0 {
		var queries []string
		if len(typeFilter) > 0 {
			queries = append(queries, "job_type IN ?")
			values = append(values, typeFilter)
		}

		if len(statusFilter) > 0 {
			queries = append(queries, "status IN ?")
			values = append(values, statusFilter)
		}

		whereQuery = "WHERE " + strings.Join(queries, " AND ")
	}

	var rawQuery string
	if interval != nil {
		rawQuery = fmt.Sprintf(`
SELECT * FROM (
(
(SELECT id, created_at, updated_at, 'discovery' AS job_type, integration_id, resource_type AS title, status, failure_message FROM describe_integration_jobs WHERE created_at > now() - interval '%[1]s')
UNION ALL 
(SELECT id, created_at, updated_at, 'compliance' AS job_type, 'all' AS integration_id, framework_id::text AS title, status, failure_message FROM compliance_jobs WHERE created_at > now() - interval '%[1]s')
)
) AS t %s ORDER BY %s %s LIMIT ? OFFSET ?;
`, *interval, whereQuery, sortBy, sortOrder)
	} else if from != nil && to != nil {
		rawQuery = fmt.Sprintf(`
SELECT * FROM (
(
(SELECT id, created_at, updated_at, 'discovery' AS job_type, integration_id, resource_type AS title, status, failure_message FROM describe_integration_jobs WHERE created_at >= ? AND created_at <= ?)
UNION ALL 
(SELECT id, created_at, updated_at, 'compliance' AS job_type, 'all' AS integration_id, framework_id::text AS title, status, failure_message FROM compliance_jobs WHERE created_at >= ? AND created_at <= ?)
)
) AS t %s ORDER BY %s %s LIMIT ? OFFSET ?;
`, whereQuery, sortBy, sortOrder)
	} else {
		rawQuery = fmt.Sprintf(`
SELECT * FROM (
(SELECT id, created_at, updated_at, 'discovery' AS job_type, integration_id, resource_type AS title, status, failure_message FROM describe_integration_jobs)
UNION ALL
(SELECT id, created_at, updated_at, 'compliance' AS job_type, 'all' AS integration_id, framework_id::text AS title, status, failure_message FROM compliance_jobs)
) AS t %s ORDER BY %s %s LIMIT ? OFFSET ?;
`, whereQuery, sortBy, sortOrder)
	}

	values = append(values, pageEnd-pageStart, pageStart)
	tx := db.ORM.Raw(rawQuery, values...).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) ListSummaryJobs(complianceJobIDs []string) ([]model.ComplianceSummarizer, error) {
	var jobs []model.ComplianceSummarizer

	tx := db.ORM.Model(model.ComplianceSummarizer{}).Where("parent_job_id IN ?", complianceJobIDs).Find(&jobs)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) GetAllJobSummary(interval *string, from *time.Time, to *time.Time, typeFilter []string, statusFilter []string) ([]model.JobSummary, error) {
	var job []model.JobSummary

	whereQuery := ""
	var values []interface{}
	if from != nil && to != nil && interval == nil {
		values = append(values, *from, *to, *from, *to, *from, *to)
	}
	if len(typeFilter) > 0 || len(statusFilter) > 0 {
		var queries []string
		if len(typeFilter) > 0 {
			queries = append(queries, "job_type IN ?")
			values = append(values, typeFilter)
		}

		if len(statusFilter) > 0 {
			queries = append(queries, "status IN ?")
			values = append(values, statusFilter)
		}

		whereQuery = "WHERE " + strings.Join(queries, " AND ")
	}

	var rawQuery string
	if interval != nil {
		rawQuery = fmt.Sprintf(`
SELECT * FROM (
(
(SELECT 'discovery' AS job_type, status, count(*) AS count FROM describe_integration_jobs WHERE created_at > now() - interval '%[1]s' GROUP BY status )
UNION ALL 
(SELECT 'compliance' AS job_type, status, count(*) AS count FROM compliance_jobs WHERE created_at > now() - interval '%[1]s' GROUP BY status )
)
) AS t %s;
`, *interval, whereQuery)
	} else if from != nil && to != nil {
		rawQuery = fmt.Sprintf(`
SELECT * FROM (
(
(SELECT 'discovery' AS job_type, status, count(*) AS count FROM describe_integration_jobs WHERE created_at >= ? AND created_at <= ? GROUP BY status )
UNION ALL 
(SELECT 'compliance' AS job_type, status, count(*) AS count FROM compliance_jobs WHERE created_at >= ? AND created_at <= ? GROUP BY status )
)
) AS t %s;`, whereQuery)
	} else {
		rawQuery = fmt.Sprintf(`
SELECT * FROM (
	(
		(SELECT 'discovery' AS job_type, status, count(*) AS count FROM describe_integration_jobs GROUP BY status)
		UNION ALL
		(SELECT 'compliance' AS job_type, status, count(*) AS count FROM compliance_jobs GROUP BY status)
	)
) AS t %s;`, whereQuery)

	}
	tx := db.ORM.Raw(rawQuery, values...).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) ListDescribeJobs(integrationId string) (*model.DescribeIntegrationJob, error) {
	var job model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).
		Where("integration_id = ?", integrationId).
		Order("updated_at DESC").
		Limit(1).
		First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) ListDescribeJobsByStatus(status api.DescribeResourceJobStatus) ([]model.DescribeIntegrationJob, error) {
	var job []model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).Where("status = ?", status).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) LatestDescribeJobForIntegration(status api.DescribeResourceJobStatus) ([]model.DescribeIntegrationJob, error) {
	var job []model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).Where("status = ?", status).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) ListDescribeJobsByIds(ids []string) ([]model.DescribeIntegrationJob, error) {
	var job []model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).Where("id IN ?", ids).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) ListDescribeJobsForInterval(interval, triggerType, createdBy string) ([]model.DescribeIntegrationJob, error) {
	var job []model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{})

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

func (db Database) ListDescribeJobsByFilters(parentIds []string, integrationIds []string, resourceType []string,
	discoveryType []string, jobStatus []string, startTime *time.Time, endTime *time.Time) ([]model.DescribeIntegrationJob, error) {
	var job []model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{})

	if len(parentIds) > 0 {
		tx = tx.Where("parent_id IN ?", parentIds)
	}

	if len(integrationIds) > 0 {
		tx = tx.Where("integration_id IN ?", integrationIds)
	}

	if len(resourceType) > 0 {
		tx = tx.Where("resource_type IN ?", resourceType)
	}
	if len(discoveryType) > 0 {
		tx = tx.Where("discovery_type IN ?", discoveryType)
	}
	if len(jobStatus) > 0 {
		tx = tx.Where("status IN ?", jobStatus)
	}
	if startTime != nil {
		tx = tx.Where("updated_at >= ?", *startTime)
	}
	if endTime != nil {
		tx = tx.Where("updated_at <= ?", *endTime)
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

func (db Database) ListPendingDescribeJobsByFilters(integrationIds []string, resourceType []string,
	jobStatus []string, startTime *time.Time, endTime *time.Time) ([]model.DescribeIntegrationJob, error) {
	var job []model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{})

	if len(integrationIds) > 0 {
		tx = tx.Where("integration_id IN ?", integrationIds)
	}

	if len(resourceType) > 0 {
		tx = tx.Where("resource_type IN ?", resourceType)
	}
	if len(jobStatus) > 0 {
		tx = tx.Where("status IN ?", jobStatus)
	}
	if startTime != nil {
		tx = tx.Where("updated_at >= ?", startTime)
	}
	if endTime != nil {
		tx = tx.Where("updated_at <= ?", *endTime)
	}

	tx = tx.Where("status IN ?", []api.DescribeResourceJobStatus{api.DescribeResourceJobCreated, api.DescribeResourceJobQueued})

	tx = tx.Find(&job)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) GetDescribeJobById(jobId string) (*model.DescribeIntegrationJob, error) {
	var job model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{})

	tx = tx.Where("id = ?", jobId)

	tx = tx.Find(&job)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) GetFailedDescribeIntegrationJobs(ctx context.Context) ([]model.DescribeIntegrationJob, error) {
	ctx, span := otel.Tracer(opengovernanceTrace.JaegerTracerName).Start(ctx, opengovernanceTrace.GetCurrentFuncName())
	defer span.End()

	var job []model.DescribeIntegrationJob

	tx := db.ORM.Raw(`
SELECT
	*
FROM
	describe_integration_jobs dr
WHERE
    trigger_type <> ? AND
	(status = ? OR status = ?) AND
	created_at > now() - interval '2 day' AND
    updated_at < now() - interval '5 minutes' AND
	NOT(error_code IN ('InvalidApiVersionParameter', 'AuthorizationFailed', 'AccessDeniedException', 'InvalidAuthenticationToken', 'AccessDenied', 'InsufficientPrivilegesException', '403', '404', '401', '400')) AND
	(retry_count < 1 OR retry_count IS NULL)
	ORDER BY id DESC
`, enums.DescribeTriggerTypeManual, api.DescribeResourceJobFailed, api.DescribeResourceJobTimeout).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) CreateDescribeIntegrationJob(job *model.DescribeIntegrationJob) error {
	tx := db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Create(job)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) CleanupManualDescribeIntegrationJobsOlderThan(t time.Time) error {
	tx := db.ORM.Where("created_at < ?", t).Where("trigger_type = ?", enums.DescribeTriggerTypeManual).Unscoped().Delete(&model.DescribeIntegrationJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) CleanupScheduledDescribeIntegrationJobsOlderThan(t time.Time) error {
	tx := db.ORM.Where("created_at < ?", t).Where("trigger_type <> ?", enums.DescribeTriggerTypeManual).Unscoped().Delete(&model.DescribeIntegrationJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// UpdateDescribeIntegrationJobsTimedOut updates the status of DescribeResourceJobs
// that have timed out while in the status of 'CREATED' or 'QUEUED' for longer
// than 4 hours.
func (db Database) UpdateDescribeIntegrationJobsTimedOut(describeIntervalHours int64) error {
	tx := db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("updated_at < NOW() - INTERVAL '20 minutes'").
		Where("status IN ?", []string{string(api.DescribeResourceJobInProgress)}).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobTimeout, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	tx = db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("updated_at < NOW() - INTERVAL '30 minutes'").
		Where("status IN ?", []string{string(api.DescribeResourceJobOldResourceDeletion)}).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobTimeout, FailureMessage: "Job timed out"})
	if tx.Error != nil {
		return tx.Error
	}

	tx = db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where(fmt.Sprintf("updated_at < NOW() - INTERVAL '2 hour'")).
		Where("status IN ?", []string{string(api.DescribeResourceJobQueued)}).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobTimeout, FailureMessage: "Queued job didn't run"})
	if tx.Error != nil {
		return tx.Error
	}

	tx = db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where(fmt.Sprintf("updated_at < NOW() - INTERVAL '%d hours'", describeIntervalHours)).
		Where("status IN ?", []string{string(api.DescribeResourceJobCreated)}).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobTimeout, FailureMessage: "Job is aborted"})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// UpdateResourceTypeDescribeIntegrationJobsTimedOut updates the status of DescribeResourceJobs
// that have timed out while in the status of 'CREATED' or 'QUEUED' for longer
// than time interval for the specific resource type.
func (db Database) UpdateResourceTypeDescribeIntegrationJobsTimedOut(resourceType string, describeIntervalHours time.Duration) (int, error) {
	totalCount := 0
	tx := db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("updated_at < NOW() - INTERVAL '20 minutes'").
		Where("status IN ?", []string{string(api.DescribeResourceJobInProgress)}).
		Where("resource_type = ?", resourceType).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobTimeout, FailureMessage: "Job timed out", ErrorCode: "JobTimeOut"})
	if tx.Error != nil {
		return totalCount, tx.Error
	}
	tx = db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("updated_at < NOW() - INTERVAL '30 minutes'").
		Where("status IN ?", []string{string(api.DescribeResourceJobOldResourceDeletion)}).
		Where("resource_type = ?", resourceType).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobTimeout, FailureMessage: "Job timed out", ErrorCode: "JobTimeOut"})
	if tx.Error != nil {
		return totalCount, tx.Error
	}
	tx = db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where(fmt.Sprintf("updated_at < NOW() - INTERVAL '1 hours'")).
		Where("status IN ?", []string{string(api.DescribeResourceJobQueued)}).
		Where("resource_type = ?", resourceType).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobFailed, FailureMessage: "Queued job didn't run", ErrorCode: "JobTimeOut"})
	if tx.Error != nil {
		return totalCount, tx.Error
	}
	totalCount += int(tx.RowsAffected)
	tx = db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where(fmt.Sprintf("updated_at < NOW() - INTERVAL '%d hours'", int(describeIntervalHours.Hours()))).
		Where("status IN ?", []string{string(api.DescribeResourceJobCreated)}).
		Where("resource_type = ?", resourceType).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobFailed, FailureMessage: "Job is aborted", ErrorCode: "JobTimeOut"})
	if tx.Error != nil {
		return totalCount, tx.Error
	}
	totalCount += int(tx.RowsAffected)
	return totalCount, nil
}

// UpdateDescribeIntegrationJobStatus updates the status of the DescribeResourceJob to the provided status.
// If the status if 'FAILED', msg could be used to indicate the failure reason
func (db Database) UpdateDescribeIntegrationJobStatus(id uint, status api.DescribeResourceJobStatus, msg, errCode string, resourceCount, deletingCount int64) error {
	tx := db.ORM.Exec("UPDATE describe_integration_jobs SET status = ?, failure_message = ?, error_code = ?,  described_resource_count = ?, deleting_count = ? WHERE id = ?",
		status, msg, errCode, resourceCount, deletingCount, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateDescribeIntegrationJobToInProgress(id uint) error {
	tx := db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("id = ?", id).
		Where("status IN ?", []string{string(api.DescribeResourceJobCreated), string(api.DescribeResourceJobQueued)}).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobInProgress, InProgressedAt: time.Now()})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateDescribeIntegrationJobToDeletionOfOldResources(id uint) error {
	tx := db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("id = ?", id).
		Updates(model.DescribeIntegrationJob{Status: api.DescribeResourceJobOldResourceDeletion})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) GetDescribeStatus(resourceType string) ([]api.DescribeStatus, error) {
	var job []api.DescribeStatus

	tx := db.ORM.Raw(`with conns as (
    select 
        integration_id, max(updated_at) as updated_at 
    from describe_integration_jobs 
    where lower(resource_type) = ? and status in ('SUCCEEDED', 'FAILED', 'TIMEOUT') group by 1
)
select j.integration_id, j.connector, j.status from describe_integration_jobs j inner join conns c on j.integration_id = c.integration_id where j.updated_at = c.updated_at;`, strings.ToLower(resourceType)).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) GetIntegrationDescribeStatus(integrationId string) ([]api.IntegrationDescribeStatus, error) {
	var job []api.IntegrationDescribeStatus

	tx := db.ORM.Raw(`with resourceTypes as (
    select 
        resource_type, max(updated_at) as updated_at 
    from
		describe_integration_jobs 
    where 
		integration_id = ?
	group by 1
)
select 
	j.resource_type, j.status 
from 
	describe_integration_jobs j inner join resourceTypes c on j.resource_type = c.resource_type 
where 
	integration_id = ? AND j.updated_at = c.updated_at;`,
		integrationId, integrationId).Find(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return job, nil
}

func (db Database) CountJobsWithStatus(interval int, connector source.Type, status api.DescribeResourceJobStatus) (*int64, error) {
	var count int64
	query := fmt.Sprintf("SELECT count(*) FROM describe_integration_jobs WHERE (connector = '%s' and created_at > now() - interval '%d hour' and status = '%s') AND deleted_at IS NULL", connector, interval, status)
	tx := db.ORM.Raw(query).Find(&count)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &count, nil
}

func (db Database) ListAllPendingIntegration() ([]string, error) {
	var integrationIds []string

	tx := db.ORM.Raw(`select distinct(integration_id) from describe_integration_jobs where status in ('CREATED', 'QUEDED', 'IN_PROGRESS')`).Find(&integrationIds)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return integrationIds, nil
}

func (db Database) ListAllFirstTryPendingIntegration() ([]string, error) {
	var discoveryTypes []string

	tx := db.ORM.Raw(`select distinct(discovery_type) from describe_integration_jobs where (status = 'CREATED' AND retry_count = 0) OR (status in ('QUEDED', 'IN_PROGRESS') and retry_count = 1)`).Find(&discoveryTypes)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return discoveryTypes, nil
}

func (db Database) ListAllSuccessfulDescribeJobs() ([]model.DescribeIntegrationJob, error) {
	var jobs []model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).Where("status = ?", api.DescribeResourceJobSucceeded).Find(&jobs)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return jobs, nil
}

func (db Database) GetLastSuccessfulDescribeJob() (*model.DescribeIntegrationJob, error) {
	var job model.DescribeIntegrationJob

	tx := db.ORM.Model(&model.DescribeIntegrationJob{}).
		Where("status = 'SUCCEEDED'").
		Order("updated_at DESC").First(&job)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &job, nil
}

func (db Database) CountJobsAndResources() (*int64, *int64, error) {
	var count, sum *int64
	err := db.ORM.Raw("select count(*), sum(described_resource_count) from describe_integration_jobs").Row().Scan(&count, &sum)
	if err != nil {
		return nil, nil, err
	}
	return count, sum, nil
}

func (db Database) CleanupAllDescribeIntegrationJobsForIntegrations(integrations []string) error {
	tx := db.ORM.Where("integration_id IN ?", integrations).Unscoped().Delete(&model.DescribeIntegrationJob{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) CreateIntegrationDiscovery(discovery *model.IntegrationDiscovery) error {
	tx := db.ORM.
		Model(&model.IntegrationDiscovery{}).
		Create(discovery)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListIntegrationDiscovery(triggerId string, integrationIds []string) ([]model.IntegrationDiscovery, error) {
	var jobs []model.IntegrationDiscovery
	tx := db.ORM.
		Model(&model.IntegrationDiscovery{}).
		Where("trigger_id = ?", triggerId)

	if len(integrationIds) > 0 {
		tx = tx.Where("integration_id in ?", integrationIds)
	}
	tx = tx.Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) GetDiscoveryJobsByParentID(parentId uint) ([]model.DescribeIntegrationJob, error) {
	var jobs []model.DescribeIntegrationJob
	tx := db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("parent_id = ?", parentId).
		Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}

func (db Database) CheckJobsDoneByParentIDAndCreatedBy(createdBy string, parentId uint) ([]model.DescribeIntegrationJob, error) {
	var jobs []model.DescribeIntegrationJob
	tx := db.ORM.
		Model(&model.DescribeIntegrationJob{}).
		Where("parent_id = ?", parentId).
		Where("created_by = ?", createdBy).
		Where("status IN ?", []api.DescribeResourceJobStatus{
			api.DescribeResourceJobCreated,
			api.DescribeResourceJobQueued,
			api.DescribeResourceJobInProgress,
			api.DescribeResourceJobOldResourceDeletion}).
		Find(&jobs)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return jobs, nil
}
