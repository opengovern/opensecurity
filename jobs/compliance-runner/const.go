package runner

import (
	complianceApi "github.com/opengovern/opensecurity/services/compliance/api"
	"github.com/opengovern/opensecurity/services/scheduler/db/model"
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

type ExecutionPlan struct {
	Callers   []model.Caller
	Query     complianceApi.Policy
	ControlID string

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

type JobResult struct {
	Job                        Job
	StartedAt                  time.Time
	Status                     model.ComplianceRunnerStatus
	PodName                    string
	Error                      string
	TotalComplianceResultCount *int
}
