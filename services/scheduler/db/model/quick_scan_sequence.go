package model

import (
	"strconv"

	"github.com/lib/pq"
	"github.com/opengovern/opencomply/services/scheduler/api"
	"gorm.io/gorm"
)

type QuickScanSequenceStatus string

const (
	QuickScanSequenceCreated              QuickScanSequenceStatus = "CREATED"
	QuickScanSequenceStarted              QuickScanSequenceStatus = "STARTED"
	QuickScanSequenceFetchingDependencies QuickScanSequenceStatus = "FETCH_DEPENDENCIES"
	QuickScanSequenceComplianceRunning    QuickScanSequenceStatus = "RUNNING_COMPLIANCE_QUICK_SCAN"
	QuickScanSequenceFailed               QuickScanSequenceStatus = "FAILED"
	QuickScanSequenceFinished             QuickScanSequenceStatus = "FINISHED"
)

type QuickScanSequence struct {
	gorm.Model
	FrameworkID    string
	IntegrationIDs pq.StringArray `gorm:"type:text[]"`
	IncludeResults pq.StringArray `gorm:"type:text[]"`
	Status         QuickScanSequenceStatus
	FailureMessage string
	CreatedBy      string
}

func (aj *QuickScanSequence) ToAPI() api.QuickScanSequence {
	return api.QuickScanSequence{
		ID:             strconv.Itoa(int(aj.ID)),
		FrameworkID:    aj.FrameworkID,
		IntegrationIDs: aj.IntegrationIDs,
		Status:         api.QuickScanSequenceStatus(aj.Status),
		FailureMessage: aj.FailureMessage,
		CreatedBy:      aj.CreatedBy,
	}
}
