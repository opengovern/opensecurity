package model

import (
	"encoding/json"
	"fmt"
	"github.com/opengovern/opencomply/pkg/types"
	"github.com/opengovern/opencomply/services/scheduler/api"
	"gorm.io/gorm"
	"strconv"
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

	QueuedAt          time.Time
	ExecutedAt        time.Time
	CompletedAt       time.Time
	TotalFindingCount *int
	Status            ComplianceRunnerStatus
	FailureMessage    string
	RetryCount        int
	TriggerType       ComplianceTriggerType

	NatsSequenceNumber uint64
	WorkerPodName      string
}

func (cr ComplianceRunner) ToAPI() api.ComplianceJobRunner {
	integrationId := ""
	if cr.IntegrationID != nil {
		integrationId = *cr.IntegrationID
	}
	return api.ComplianceJobRunner{
		RunnerId:        strconv.Itoa(int(cr.ID)),
		ComplianceJobId: strconv.Itoa(int(cr.ParentJobID)),
		ControlId:       cr.ControlID,
		IntegrationId:   integrationId,
		WorkerPodName:   cr.WorkerPodName,
		QueuedAt:        cr.QueuedAt,
		ExecutedAt:      cr.ExecutedAt,
		CompletedAt:     cr.CompletedAt,
		Status:          cr.Status.ToAPI(),
		FailureMessage:  cr.FailureMessage,
		TriggerType:     string(cr.TriggerType),
	}
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
