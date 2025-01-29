package model

import (
	"encoding/json"
	"fmt"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/scheduler/api"
	"gorm.io/gorm"
	"time"
)

type Caller struct {
	RootBenchmark      string
	TracksDriftEvents  bool
	ParentBenchmarkIDs []string
	ControlID          string
	ControlSeverity    types.ComplianceResultSeverity
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
	Status            ComplianceRunnerStatus
	FailureMessage    string
	RetryCount        int
	TriggerType       ComplianceTriggerType

	NatsSequenceNumber uint64
	WorkerPodName      string
}

func (s ComplianceRunnerStatus) ToAPI() api.ComplianceRunnerStatus {
	return api.ComplianceRunnerStatus(s)
}

func (cr *ComplianceRunner) GetKeyIdentifier() string {
	cid := "all"
	if cr.IntegrationID != nil {
		cid = *cr.IntegrationID
	}
	return fmt.Sprintf("%s-%s-%s-%d", cr.FrameworkID, cr.PolicyID, cid, cr.ParentJobID)
}

func (cr *ComplianceRunner) GetCallers() ([]Caller, error) {
	var res []Caller
	err := json.Unmarshal([]byte(cr.Callers), &res)
	return res, err
}

func (cr *ComplianceRunner) SetCallers(callers []Caller) error {
	b, err := json.Marshal(callers)
	if err != nil {
		return err
	}
	cr.Callers = string(b)
	return nil
}
