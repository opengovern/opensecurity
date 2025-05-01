package summarizer

import (
	"time"

	"github.com/opengovern/opensecurity/jobs/compliance-summarizer/types"
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
