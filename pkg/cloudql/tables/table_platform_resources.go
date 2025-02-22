package opengovernance

import (
	"context"
	"time"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tablePlatformResources(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_resources",
		Description: "OpenGovernance Resources",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListResources,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "resource_type",
					Require: "required",
				},
			},
		},
		Columns: []*plugin.Column{
			{Name: "platform_id", Type: proto.ColumnType_STRING},
			{Name: "resource_id", Type: proto.ColumnType_STRING},
			{Name: "integration_id", Type: proto.ColumnType_STRING},
			{Name: "integration_type", Type: proto.ColumnType_STRING},
			{Name: "resource_type", Type: proto.ColumnType_STRING},
			{Name: "resource_name", Type: proto.ColumnType_STRING},
			{Name: "described_by", Type: proto.ColumnType_STRING},
			{Name: "described_at", Transform: transform.From(fixTime), Type: proto.ColumnType_TIMESTAMP},
			{Name: "description", Type: proto.ColumnType_JSON},
		},
	}
}

func fixTime(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	resource := d.HydrateItem.(og_client.Resource)
	t := time.UnixMilli(resource.DescribedAt)
	return t.Format("2006-01-02T15:04:05"), nil
}

func getResourceRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	resource := d.HydrateItem.(og_client.Resource)
	if len(resource.Metadata.Region) > 0 {
		return resource.Metadata.Region, nil
	}
	return resource.Metadata.Location, nil
}
