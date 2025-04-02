package opengovernance

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tablePlatformOsvVulnerabilityDetails(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_osv_vulnerability_details",
		Description: "Provides detailed information about OSV vulnerabilities.",
		List: &plugin.ListConfig{
			Hydrate: og_client.ListOsvVulnerabilityDetail,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the vulnerability entry.",
				Transform:   transform.FromField("Description.ID"),
			},
			{
				Name:        "schema_version",
				Type:        proto.ColumnType_STRING,
				Description: "The version of the OSV schema.",
				Transform:   transform.FromField("Description.SchemaVersion"),
			},
			{
				Name:        "modified",
				Type:        proto.ColumnType_TIMESTAMP, // Assuming this can be parsed as a timestamp
				Description: "The timestamp of when the vulnerability entry was last modified.",
				Transform:   transform.FromField("Description.Modified").Transform(transform.NullIfZeroValue), // Or use STRING if it's not always a standard timestamp format
			},
			{
				Name:        "published",
				Type:        proto.ColumnType_TIMESTAMP, // Assuming this can be parsed as a timestamp
				Description: "The timestamp of when the vulnerability entry was published.",
				Transform:   transform.FromField("Description.Published").Transform(transform.NullIfZeroValue), // Or use STRING
			},
			{
				Name:        "withdrawn",
				Type:        proto.ColumnType_TIMESTAMP, // Assuming this can be parsed as a timestamp
				Description: "The timestamp of when the vulnerability entry was withdrawn.",
				Transform:   transform.FromField("Description.Withdrawn").Transform(transform.NullIfZeroValue), // Or use STRING
			},
			{
				Name:        "aliases",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IDs for the same vulnerability from other databases (e.g., CVEs).",
				Transform:   transform.FromField("Description.Aliases"),
			},
			{
				Name:        "related",
				Type:        proto.ColumnType_JSON,
				Description: "A list of IDs of related vulnerabilities.",
				Transform:   transform.FromField("Description.Related"),
			},
			{
				Name:        "summary",
				Type:        proto.ColumnType_STRING,
				Description: "A one-line summary of the vulnerability.",
				Transform:   transform.FromField("Description.Summary"),
			},
			{
				Name:        "details",
				Type:        proto.ColumnType_STRING,
				Description: "Detailed information about the vulnerability.",
				Transform:   transform.FromField("Description.Details"),
			},
			{
				Name:        "severity",
				Type:        proto.ColumnType_JSON,
				Description: "Severity information for the vulnerability (schema-dependent).",
				Transform:   transform.FromField("Description.Severity"),
			},
			{
				Name:        "affected",
				Type:        proto.ColumnType_JSON,
				Description: "A list of affected packages, versions, and ecosystems.",
				Transform:   transform.FromField("Description.Affected"),
			},
			{
				Name:        "references",
				Type:        proto.ColumnType_JSON,
				Description: "A list of references for the vulnerability (e.g., advisories, reports).",
				Transform:   transform.FromField("Description.References"),
			},
			{
				Name:        "credits",
				Type:        proto.ColumnType_JSON,
				Description: "Information about who contributed to the vulnerability information.",
				Transform:   transform.FromField("Description.Credits"),
			},
			{
				Name:        "database_specific",
				Type:        proto.ColumnType_JSON,
				Description: "Database-specific information.",
				Transform:   transform.FromField("Description.DatabaseSpecific"),
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
