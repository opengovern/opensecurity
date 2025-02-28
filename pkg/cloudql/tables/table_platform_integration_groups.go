package opengovernance

import (
	"context"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformIntegrationGroups(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_integration_groups",
		Description: "OpenGovernance Integrations",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"name"}),
			Hydrate:    og_client.GetIntegrationGroup,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListIntegrationsGroup,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The ID of the integration in OpenGovernance"},
			{Name: "query", Type: proto.ColumnType_STRING, Description: "The ID of the integration in the provider"},
		},
	}
}
