package opengovernance

import (
	"context"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformFindings(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_findings",
		Description: "OpenGovernance Compliance ComplianceResults",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListFindings,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING},
			{Name: "benchmark_id", Type: proto.ColumnType_STRING},
			{Name: "policy_id", Type: proto.ColumnType_STRING},
			{Name: "integration_id", Type: proto.ColumnType_STRING},
			{Name: "described_at", Type: proto.ColumnType_INT},
			{Name: "evaluated_at", Type: proto.ColumnType_INT},
			{Name: "state_active", Type: proto.ColumnType_BOOL},
			{Name: "result", Type: proto.ColumnType_STRING},
			{Name: "severity", Type: proto.ColumnType_STRING},
			{Name: "evaluator", Type: proto.ColumnType_STRING},
			{Name: "integration_type", Type: proto.ColumnType_STRING},
			{Name: "platform_resource_id", Type: proto.ColumnType_STRING},
			{Name: "resource_name", Type: proto.ColumnType_STRING},
			{Name: "resource_type", Type: proto.ColumnType_STRING},
			{Name: "reason", Type: proto.ColumnType_STRING},
		},
	}
}
