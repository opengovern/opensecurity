package opengovernance

import (
	"context"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformConnections(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_integrations",
		Description: "OpenGovernance Integrations",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"integration_id", "provider_id"}),
			Hydrate:    og_client.GetIntegration,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListIntegrations,
		},
		Columns: []*plugin.Column{
			{Name: "integration_id", Type: proto.ColumnType_STRING, Description: "The ID of the integration in OpenGovernance"},
			{Name: "provider_id", Type: proto.ColumnType_STRING, Description: "The ID of the integration in the provider"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the integration"},
			{Name: "integration_type", Type: proto.ColumnType_STRING},
			{Name: "state", Type: proto.ColumnType_STRING},
			{Name: "annotations", Type: proto.ColumnType_JSON},
			{Name: "labels", Type: proto.ColumnType_JSON},
		},
	}
}
