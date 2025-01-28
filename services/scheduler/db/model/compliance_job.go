package model

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"time"

	"github.com/lib/pq"
	runner "github.com/opengovern/opencomply/jobs/compliance-runner-job"
	summarizer "github.com/opengovern/opencomply/jobs/compliance-summarizer-job"
	"github.com/opengovern/opencomply/services/scheduler/api"
	"gorm.io/gorm"
)

type ComplianceJobStatus string
type ComplianceTriggerType string

const (
	ComplianceJobCreated              ComplianceJobStatus = "CREATED"
	ComplianceJobQueued               ComplianceJobStatus = "QUEUED"      // for quick audit
	ComplianceJobInProgress           ComplianceJobStatus = "IN_PROGRESS" // for quick audit
	ComplianceJobRunnersInProgress    ComplianceJobStatus = "RUNNERS_IN_PROGRESS"
	ComplianceJobSinkInProgress       ComplianceJobStatus = "SINK_IN_PROGRESS"
	ComplianceJobSummarizerInProgress ComplianceJobStatus = "SUMMARIZER_IN_PROGRESS"
	ComplianceJobFailed               ComplianceJobStatus = "FAILED"
	ComplianceJobSucceeded            ComplianceJobStatus = "SUCCEEDED"
	ComplianceJobTimeOut              ComplianceJobStatus = "TIMEOUT"
	ComplianceJobCanceled             ComplianceJobStatus = "CANCELED"

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
	FrameworkID         string
	WithIncidents       bool
	Status              ComplianceJobStatus
	RunnersStatus       pgtype.JSONB
	IncludeResults      pq.StringArray `gorm:"type:text[]"`
	AreAllRunnersQueued bool
	IntegrationIDs      pq.StringArray `gorm:"type:text[]"`
	FailureMessage      string
	TriggerType         ComplianceTriggerType
	ParentID            *uint
	CreatedBy           string
}

func (c ComplianceJob) ToApi() api.ComplianceJob {
	return api.ComplianceJob{
		ID:             c.ID,
		BenchmarkID:    c.FrameworkID,
		Status:         c.Status.ToApi(),
		FailureMessage: c.FailureMessage,
	}
}

type ComplianceRunner struct {
	gorm.Model

	Callers              string
	FrameworkID          string
	ControlID            string
	PolicyID             string
	IntegrationID        *string
	ResourceCollectionID *string
	ParentJobID          uint `gorm:"index"`

	StartedAt         time.Time
	TotalFindingCount *int
	Status            runner.ComplianceRunnerStatus
	FailureMessage    string
	RetryCount        int

	TriggerType        ComplianceTriggerType
	NatsSequenceNumber uint64
}

func (cr *ComplianceRunner) GetKeyIdentifier() string {
	cid := "all"
	if cr.IntegrationID != nil {
		cid = *cr.IntegrationID
	}
	return fmt.Sprintf("%s-%s-%s-%d", cr.FrameworkID, cr.PolicyID, cid, cr.ParentJobID)
}

func (cr *ComplianceRunner) GetCallers() ([]runner.Caller, error) {
	var res []runner.Caller
	err := json.Unmarshal([]byte(cr.Callers), &res)
	return res, err
}

func (cr *ComplianceRunner) SetCallers(callers []runner.Caller) error {
	b, err := json.Marshal(callers)
	if err != nil {
		return err
	}
	cr.Callers = string(b)
	return nil
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
