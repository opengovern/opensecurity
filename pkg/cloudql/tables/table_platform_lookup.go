package opengovernance

import (
	"context"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tablePlatformLookup(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_lookup",
		Description: "OpenGovernance Resource Lookup",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListLookupResources,
		},
		Columns: []*plugin.Column{
			{Name: "platform_resource_id", Transform: transform.FromField("PlatformId"), Type: proto.ColumnType_STRING},
			{Name: "resource_id", Type: proto.ColumnType_STRING},
			{Name: "name", Type: proto.ColumnType_STRING},
			{Name: "integration_type", Transform: transform.FromField("IntegrationType"), Type: proto.ColumnType_STRING},
			{Name: "resource_type", Type: proto.ColumnType_STRING},
			{Name: "integration_id", Transform: transform.FromField("IntegrationID"), Type: proto.ColumnType_STRING},
			{Name: "described_at", Type: proto.ColumnType_INT},
			{Name: "described_by", Type: proto.ColumnType_STRING},
			{Name: "tags", Type: proto.ColumnType_JSON},
		},
	}
}
