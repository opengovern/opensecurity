package model

import (
	"time"

	"github.com/jackc/pgtype"

	"github.com/lib/pq"
	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opensecurity/services/scheduler/api"
	"gorm.io/gorm"
)

type DiscoveryType string
type IntegrationDiscoveryStatus string

const (
	DiscoveryType_Fast DiscoveryType = "FAST"
	DiscoveryType_Full DiscoveryType = "FULL"
	DiscoveryType_Cost DiscoveryType = "COST"

	IntegrationDiscoveryStatusInProgress IntegrationDiscoveryStatus = "IN_PROGRESS"
	IntegrationDiscoveryStatusCompleted  IntegrationDiscoveryStatus = "COMPLETED"
)

type IntegrationDiscovery struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TriggerID     uint
	ConnectionID  string                    `json:"connectionID"`
	AccountID     string                    `json:"accountID"`
	TriggerType   enums.DescribeTriggerType `json:"triggerType"`
	TriggeredBy   string                    `json:"triggeredBy"`
	DiscoveryType DiscoveryType
	ResourceTypes pq.StringArray `gorm:"type:text[]" json:"resourceTypes"`
}

type ManualDiscoverySchedule struct {
	gorm.Model
	ResourceType    string
	IntegrationID   string
	IntegrationType integration.Type
	Parameters      pgtype.JSONB
	CreatedBy       string
}

type ManualDiscoveryScheduleUnique struct {
	ResourceType string
	Parameters   pgtype.JSONB
}

type DescribeIntegrationJob struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time      `gorm:"index:,sort:desc"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	QueuedAt       time.Time
	InProgressedAt time.Time
	CreatedBy      string
	ParentID       *uint `gorm:"index:,sort:desc"`

	IntegrationID   string `gorm:"index:idx_source_id_full_discovery;index"`
	IntegrationType integration.Type
	ProviderID      string
	TriggerType     enums.DescribeTriggerType

	Parameters pgtype.JSONB // map[string]string

	ResourceType           string                        `gorm:"index:idx_resource_type_status;index"`
	Status                 api.DescribeResourceJobStatus `gorm:"index:idx_resource_type_status;index"`
	RetryCount             int
	FailureMessage         string // Should be NULLSTRING
	ErrorCode              string // Should be NULLSTRING
	DescribedResourceCount int64
	DeletingCount          int64

	NatsSequenceNumber uint64
}

type ResourceTypeDescribedCount struct {
	ResourceType           string
	TableName              string
	IntegrationID          string
	DescribedResourceCount int64
	UpdatedAt              time.Time
}
