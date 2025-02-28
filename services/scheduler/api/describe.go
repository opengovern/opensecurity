package api

import "github.com/opengovern/og-util/pkg/source"

type DescribeSingleResourceRequest struct {
	Provider         source.Type `json:"provider"`
	ResourceType     string
	AccountID        string
	AccessKey        string
	SecretKey        string
	AdditionalFields map[string]string
}

type DescribeStatus struct {
	ConnectionID string
	Connector    string
	Status       DescribeResourceJobStatus
}

type IntegrationDescribeStatus struct {
	ResourceType string
	Status       DescribeResourceJobStatus
}

type ComplianceJobStatus string

const (
	ComplianceJobCreated              ComplianceJobStatus = "CREATED"
	ComplianceJobRunnersInProgress    ComplianceJobStatus = "RUNNERS_IN_PROGRESS"
	ComplianceJobSinkInProgress       ComplianceJobStatus = "SINK_IN_PROGRESS"
	ComplianceJobSummarizerInProgress ComplianceJobStatus = "SUMMARIZER_IN_PROGRESS"
	ComplianceJobFailed               ComplianceJobStatus = "FAILED"
	ComplianceJobSucceeded            ComplianceJobStatus = "SUCCEEDED"
	ComplianceJobTimeout              ComplianceJobStatus = "TIMEOUT"
	ComplianceJobCanceled             ComplianceJobStatus = "CANCELED"

	ComplianceJobQueued     ComplianceJobStatus = "QUEUED"      // for quick audit
	ComplianceJobInProgress ComplianceJobStatus = "IN_PROGRESS" // for quick audit
)

type ComplianceJob struct {
	ID             uint
	FrameworkIds   []string
	Status         ComplianceJobStatus
	FailureMessage string
}
