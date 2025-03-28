package opengovernance

import (
	"context"
	"encoding/json"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformArtifactSboms(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_artifact_sbom",
		Description: "Platform Artifact SBOMs",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListArtifactSboms,
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
				Name:      "sbom_format",
				Transform: transform.FromField("Description.SbomFormat"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "sbom",
				Transform: transform.FromField("Description.Sbom"),
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
func marshalJSON(_ context.Context, d *transform.TransformData) (interface{}, error) {
	b, err := json.Marshal(d.Value)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
