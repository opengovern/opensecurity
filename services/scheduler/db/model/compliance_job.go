package model

import (
	"github.com/jackc/pgtype"
	"time"

	"github.com/lib/pq"
	summarizer "github.com/opengovern/opencomply/jobs/compliance-summarizer-job"
	"github.com/opengovern/opencomply/services/scheduler/api"
	"gorm.io/gorm"
)

type ComplianceJobStatus string
type ComplianceTriggerType string

const (
	ComplianceJobCreated              ComplianceJobStatus = "CREATED"
	ComplianceJobRunnersInProgress    ComplianceJobStatus = "RUNNERS_IN_PROGRESS"
	ComplianceJobSinkInProgress       ComplianceJobStatus = "SINK_IN_PROGRESS"
	ComplianceJobSummarizerInProgress ComplianceJobStatus = "SUMMARIZER_IN_PROGRESS"
	ComplianceJobFailed               ComplianceJobStatus = "FAILED"
	ComplianceJobSucceeded            ComplianceJobStatus = "SUCCEEDED"
	ComplianceJobTimeOut              ComplianceJobStatus = "TIMEOUT"
	ComplianceJobCanceled             ComplianceJobStatus = "CANCELED"

	ComplianceJobQueued     ComplianceJobStatus = "QUEUED"      // for quick audit
	ComplianceJobInProgress ComplianceJobStatus = "IN_PROGRESS" // for quick audit

	ComplianceTriggerTypeScheduled ComplianceTriggerType = "scheduled" // default
	ComplianceTriggerTypeManual    ComplianceTriggerType = "manual"
	ComplianceTriggerTypeEmpty     ComplianceTriggerType = ""
)

func (c ComplianceJobStatus) ToApi() api.ComplianceJobStatus {
	return api.ComplianceJobStatus(c)
}

type ComplianceRunnersStatus struct {
	RunnersCreated   int64 `json:"runners_created"`
	RunnersQueued    int64 `json:"runners_queued"`
	RunnersRunning   int64 `json:"runners_running"`
	RunnersFailed    int64 `json:"runners_failed"`
	RunnersSucceeded int64 `json:"runners_succeeded"`
	RunnersTimedOut  int64 `json:"runners_timed_out"`
	TotalCount       int64 `json:"total_count"`
}

type ComplianceJob struct {
	gorm.Model
	FrameworkIds        pq.StringArray `gorm:"type:text[]"`
	WithIncidents       bool
	Status              ComplianceJobStatus
	RunnersStatus       pgtype.JSONB
	IncludeResults      pq.StringArray `gorm:"type:text[]"`
	AreAllRunnersQueued bool
	IntegrationIDs      pq.StringArray `gorm:"type:text[]"`
	StepFailed          ComplianceJobStatus
	FailureMessage      string
	TriggerType         ComplianceTriggerType
	ParentID            *uint
	CreatedBy           string
}

func (c ComplianceJob) ToApi() api.ComplianceJob {
	return api.ComplianceJob{
		ID:             c.ID,
		FrameworkIds:   c.FrameworkIds,
		Status:         c.Status.ToApi(),
		FailureMessage: c.FailureMessage,
	}
}

type ComplianceSummarizer struct {
	gorm.Model

	BenchmarkID    string
	ParentJobID    uint
	IntegrationIDs pq.StringArray `gorm:"type:text[]"`

	StartedAt      time.Time
	RetryCount     int
	Status         summarizer.ComplianceSummarizerStatus
	FailureMessage string

	TriggerType ComplianceTriggerType
}

type ComplianceJobWithSummarizerJob struct {
	ID             uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
	BenchmarkID    string
	Status         ComplianceJobStatus
	ConnectionIDs  pq.StringArray `gorm:"type:text[]"`
	SummarizerJobs pq.StringArray `gorm:"type:text[]"`
	TriggerType    ComplianceTriggerType
	CreatedBy      string
}
