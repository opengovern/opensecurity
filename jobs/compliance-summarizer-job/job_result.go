package summarizer

import (
	"time"

	"github.com/opengovern/opencomply/jobs/compliance-summarizer-job/types"
)

type ComplianceSummarizerStatus string

const (
	ComplianceSummarizerCreated    ComplianceSummarizerStatus = "CREATED"
	ComplianceSummarizerInProgress ComplianceSummarizerStatus = "IN_PROGRESS"
	ComplianceSummarizerSucceeded  ComplianceSummarizerStatus = "SUCCEEDED"
	ComplianceSummarizerFailed     ComplianceSummarizerStatus = "FAILED"
)

type JobResult struct {
	Job       types.Job
	StartedAt time.Time
	Status    ComplianceSummarizerStatus
	Error     string
}
