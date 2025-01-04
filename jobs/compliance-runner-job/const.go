package runner

import (
	"github.com/opengovern/opencomply/pkg/types"
	complianceApi "github.com/opengovern/opencomply/services/compliance/api"
	"time"
)

const (
	JobQueueTopic        = "compliance-runner-job-queue"
	JobQueueTopicManuals = "compliance-runner-job-queue-manuals"
	ResultQueueTopic     = "compliance-runner-job-result"
	ConsumerGroup        = "compliance-runner"
	ConsumerGroupManuals = "compliance-runner-manuals"

	StreamName = "compliance-runner"
)

type Caller struct {
	RootBenchmark      string
	TracksDriftEvents  bool
	ParentBenchmarkIDs []string
	ControlID          string
	ControlSeverity    types.ComplianceResultSeverity
}

type ExecutionPlan struct {
	Callers []Caller
	Query   complianceApi.Policy

	IntegrationID *string
	ProviderID    *string
}

type Job struct {
	ID          uint
	RetryCount  int
	ParentJobID uint
	CreatedAt   time.Time

	ExecutionPlan ExecutionPlan
}

type ComplianceRunnerStatus string

const (
	ComplianceRunnerCreated    ComplianceRunnerStatus = "CREATED"
	ComplianceRunnerQueued     ComplianceRunnerStatus = "QUEUED"
	ComplianceRunnerInProgress ComplianceRunnerStatus = "IN_PROGRESS"
	ComplianceRunnerSucceeded  ComplianceRunnerStatus = "SUCCEEDED"
	ComplianceRunnerFailed     ComplianceRunnerStatus = "FAILED"
	ComplianceRunnerTimeOut    ComplianceRunnerStatus = "TIMEOUT"
	ComplianceRunnerCanceled   ComplianceRunnerStatus = "CANCELED"
)

type JobResult struct {
	Job                        Job
	StartedAt                  time.Time
	Status                     ComplianceRunnerStatus
	Error                      string
	TotalComplianceResultCount *int
}
