package opengovernance

import (
	"context"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformTasks(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "tasks",
		Description: "Platform Tasks",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id"}),
			Hydrate:    og_client.GetTask,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListTasks,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING},
			{Name: "title", Type: proto.ColumnType_STRING},
			{Name: "description", Type: proto.ColumnType_STRING},
			{Name: "image_url", Type: proto.ColumnType_STRING},
			{Name: "LastRun", Type: proto.ColumnType_TIMESTAMP},
			{Name: "params", Type: proto.ColumnType_JSON},
		},
	}
}
