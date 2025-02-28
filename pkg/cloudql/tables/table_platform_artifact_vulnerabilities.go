package opengovernance

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformArtifactVulnerabilities(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_artifact_vulnerabilities",
		Description: "Platform Artifact Vulnerabilities",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListArtifactVulnerabilities,
		},
		Columns: []*plugin.Column{
			{
				Name:      "image_url",
				Transform: transform.FromField("Description.imageUrl"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "artifact_digest",
				Transform: transform.FromField("Description.artifactDigest"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "vulnerabilities",
				Transform: transform.FromField("Description.Vulnerabilities"),
				Type:      proto.ColumnType_JSON,
			},
		},
	}
}
