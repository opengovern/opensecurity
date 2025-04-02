package opengovernance

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformPackageVulnerabilities(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_package_vulnerabilities",
		Description: "Platform Package Vulnerabilities",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListPackageWithVulnIDs,
		},
		Columns: []*plugin.Column{
			{
				Name:      "package_identifier",
				Transform: transform.FromField("Description.PackageIdentifier"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "package_name",
				Transform: transform.FromField("Description.PackageName"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "ecosystem",
				Transform: transform.FromField("Description.Ecosystem"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "version",
				Transform: transform.FromField("Description.Version"),
				Type:      proto.ColumnType_STRING,
			},
			{
				Name:      "vulnerabilities",
				Transform: transform.FromField("Description.Vulnerabilities"),
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
