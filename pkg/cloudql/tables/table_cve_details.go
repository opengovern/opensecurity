package opengovernance

import (
	"context"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCveDetails(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cve_details",
		Description: "CVE details",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
			Hydrate: og_client.ListCveDetails,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id", // Matches json:"id"
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the CVE entry (often the CVE ID itself).",
				Transform:   transform.FromField("Description.ID"), // Maps to the 'ID' field in the struct
			},
			{
				Name:        "source_identifier",
				Type:        proto.ColumnType_STRING,
				Description: "The source that reported the CVE.",
				Transform:   transform.FromField("Description.SourceIdentifier"),
			},
			{
				Name:        "published",
				Type:        proto.ColumnType_STRING, // Use TIMESTAMP if you parse the string to time.Time
				Description: "The date the CVE was published.",
				Transform:   transform.FromField("Description.Published"),
			},
			{
				Name:        "last_modified",
				Type:        proto.ColumnType_STRING, // Use TIMESTAMP if you parse the string to time.Time
				Description: "The date the CVE record was last modified.",
				Transform:   transform.FromField("Description.LastModified"),
			},
			{
				Name:        "vuln_status",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the vulnerability (e.g., Analyzed, Modified).",
				Transform:   transform.FromField("Description.VulnStatus"),
			},
			{
				Name:        "description",           // Matches json:"description"
				Type:        proto.ColumnType_STRING, // This is now a simple string
				Description: "A text description of the vulnerability.",
				Transform:   transform.FromField("Description.Description"), // Maps to the 'Description' field
			},
			{
				Name:        "metrics",             // Matches json:"metrics"
				Type:        proto.ColumnType_JSON, // Still complex, best represented as JSON
				Description: "Vulnerability metrics information (content represented as JSON).",
				Transform:   transform.FromField("Description.Metrics"),
			},
			{
				Name:        "weaknesses",          // Matches json:"weaknesses"
				Type:        proto.ColumnType_JSON, // Slice of structs, best as JSON
				Description: "A list of associated weaknesses (CWEs), represented as JSON.",
				Transform:   transform.FromField("Description.Weaknesses"),
			},
			{
				Name:        "cisa_exploit_add",
				Type:        proto.ColumnType_STRING, // Handles *string correctly (null if pointer is nil)
				Description: "Date the vulnerability was added to CISA's Known Exploited Vulnerabilities (KEV) catalog.",
				Transform:   transform.FromField("Description.CisaExploitAdd"),
			},
			{
				Name:        "cisa_action_due",
				Type:        proto.ColumnType_STRING, // Handles *string correctly
				Description: "The due date for required remediation action according to CISA KEV.",
				Transform:   transform.FromField("Description.CisaActionDue"),
			},
			{
				Name:        "cisa_required_action",
				Type:        proto.ColumnType_STRING, // Handles *string correctly
				Description: "The action required by CISA for federal agencies.",
				Transform:   transform.FromField("Description.CisaRequiredAction"),
			},
			{
				Name:        "cisa_vulnerability_name",
				Type:        proto.ColumnType_STRING, // Handles *string correctly
				Description: "The name assigned to the vulnerability by CISA.",
				Transform:   transform.FromField("Description.CisaVulnerabilityName"),
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
