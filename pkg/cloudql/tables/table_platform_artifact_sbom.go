package opengovernance

import (
	"context"
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
				Transform: transform.FromField("Description.imageUrl"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "artifact_id",
				Transform: transform.FromField("Description.artifactId"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "sbom_format",
				Transform: transform.FromField("Description.sbomFormat"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "sbom",
				Transform: transform.FromField("Description.sbom"),
				Type:      proto.ColumnType_JSON,
			},
		},
	}
}
