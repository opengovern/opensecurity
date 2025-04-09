package opengovernance

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableArtifactPackageList(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "software_packages",
		Description: "Platform Artifact SBOMs",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListArtifactPackageList,
		},
		Columns: []*plugin.Column{
			{
				Name:      "image_url",
				Transform: transform.FromField("Description.ImageURL"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "artifact_id",
				Transform: transform.FromField("Description.ArtifactID"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "packages",
				Transform: transform.FromField("Description.Packages"),
				Type:      proto.ColumnType_JSON,
			},
			{
				Name:        "platform_description",
				Type:        proto.ColumnType_JSON,
				Description: "The full model description of the resource",
				Transform:   transform.FromField("Description").Transform(marshalJSON),
			},
		},
	}
}
