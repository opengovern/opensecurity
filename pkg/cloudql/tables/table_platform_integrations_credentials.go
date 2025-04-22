package opengovernance

import (
	"context"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformIntegrationsCredentials(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_integrations_credentials",
		Description: "OpenGovernance Integrations Credentials",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"integration_id"}),
			Hydrate:    og_client.GetIntegrationCredential,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListIntegrationCredentials,
		},
		Columns: []*plugin.Column{
			{Name: "integration_id", Type: proto.ColumnType_STRING, Description: "The ID of the integration in OpenGovernance"},
			{Name: "integration_type", Type: proto.ColumnType_STRING, Description: "The type of the integration"},
			{Name: "credential_id", Type: proto.ColumnType_STRING, Description: "The ID of the credential"},
			{Name: "credential_type", Type: proto.ColumnType_STRING},
			{Name: "secret", Type: proto.ColumnType_STRING},
		},
	}
}
