package opengovernance

import (
	"context"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformIntegrationResourceTypes(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_integration_resource_types",
		Description: "OpenGovernance Integration Resource Types",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{{
				Name:    "integration_type",
				Require: "required",
			}, {
				Name:    "resource_type",
				Require: "required",
			}},
			Hydrate: og_client.GetIntegrationResourceType,
		},
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{{
				Name:    "integration_type",
				Require: "required",
			}},
			Hydrate: og_client.ListIntegrationResourceTypes,
		},
		Columns: []*plugin.Column{
			{Name: "resource_type", Type: proto.ColumnType_STRING, Description: "Name of the resource type."},
			{Name: "integration_type", Type: proto.ColumnType_STRING},
			{Name: "description", Type: proto.ColumnType_STRING},
			{Name: "params", Type: proto.ColumnType_JSON},
			{Name: "table", Type: proto.ColumnType_STRING},
		},
	}
}
