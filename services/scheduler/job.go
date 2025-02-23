package describe

import (
	"time"

	"github.com/opengovern/og-util/pkg/integration"

	"github.com/opengovern/og-util/pkg/describe/enums"
	"github.com/opengovern/opensecurity/services/scheduler/api"
)

const (
	InventorySummaryIndex = "inventory_summary"
	describeJobTimeout    = 3 * 60 * time.Minute
)

type DescribeJob struct {
	JobID           uint // DescribeResourceJob ID
	ScheduleJobID   uint
	ParentJobID     uint // DescribeSourceJob ID
	ResourceType    string
	IntegrationID   string
	ProviderID      string
	DescribedAt     int64
	IntegrationType integration.Type
	CipherText      string
	TriggerType     enums.DescribeTriggerType
	RetryCounter    uint
}

type DescribeJobResult struct {
	JobID                uint
	ParentJobID          uint
	Status               api.DescribeResourceJobStatus
	Error                string
	DescribeJob          DescribeJob
	DescribedResourceIDs []string
	ErrorCode            string
}

func (r DescribeJobResult) KeysAndIndex() ([]string, string) {
	return []string{}, ""
}
