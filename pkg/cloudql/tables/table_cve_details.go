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
				Name:        "cve_id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique identifier for the CVE (e.g., CVE-2023-12345).",
				Transform:   transform.FromField("Description.CveID"),
			},
			{
				Name:        "source_identifier",
				Type:        proto.ColumnType_STRING,
				Description: "Identifier of the source that reported the CVE.",
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
				Name:        "descriptions",
				Type:        proto.ColumnType_JSON,
				Description: "A list of descriptions for the CVE, often in different languages.",
				Transform:   transform.FromField("Description.Descriptions"),
			},
			{
				Name:        "cvss_version",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS version used for scoring (e.g., '2.0', '3.0', '3.1').",
				Transform:   transform.FromField("Description.CvssVersion"),
			},
			{
				Name:        "cvss_score",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The base CVSS score.",
				Transform:   transform.FromField("Description.CvssScore"),
			},
			{
				Name:        "cvss_severity",
				Type:        proto.ColumnType_STRING,
				Description: "The severity rating based on the CVSS score (e.g., LOW, MEDIUM, HIGH, CRITICAL).",
				Transform:   transform.FromField("Description.CvssSeverity"),
			},
			{
				Name:        "cvss_attack_vector",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS attack vector metric.",
				Transform:   transform.FromField("Description.CvssAttackVector"),
			},
			{
				Name:        "cvss_attack_complexity",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS attack complexity metric.",
				Transform:   transform.FromField("Description.CvssAttackComplexity"),
			},
			{
				Name:        "cvss_privileges_required",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS privileges required metric.",
				Transform:   transform.FromField("Description.CvssPrivilegesRequired"),
			},
			{
				Name:        "cvss_user_interaction",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS user interaction metric.",
				Transform:   transform.FromField("Description.CvssUserInteraction"),
			},
			{
				Name:        "cvss_conf_impact",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS confidentiality impact metric.",
				Transform:   transform.FromField("Description.CvssConfImpact"),
			},
			{
				Name:        "cvss_integ_impact",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS integrity impact metric.",
				Transform:   transform.FromField("Description.CvssIntegImpact"),
			},
			{
				Name:        "cvss_avail_impact",
				Type:        proto.ColumnType_STRING,
				Description: "The CVSS availability impact metric.",
				Transform:   transform.FromField("Description.CvssAvailImpact"),
			},
			{
				Name:        "metrics",
				Type:        proto.ColumnType_JSON,
				Description: "Detailed CVSS metrics data.",
				Transform:   transform.FromField("Description.Metrics"),
			},
			{
				Name:        "weaknesses",
				Type:        proto.ColumnType_JSON,
				Description: "A list of associated weaknesses (CWEs - Common Weakness Enumeration).",
				Transform:   transform.FromField("Description.Weaknesses"),
			},
			{
				Name:        "cisa_exploit_add",
				Type:        proto.ColumnType_STRING, // Use TIMESTAMP if you parse the string
				Description: "Date the vulnerability was added to CISA's Known Exploited Vulnerabilities (KEV) catalog.",
				Transform:   transform.FromField("Description.CisaExploitAdd"),
			},
			{
				Name:        "cisa_action_due",
				Type:        proto.ColumnType_STRING, // Use TIMESTAMP if you parse the string
				Description: "The due date for required remediation action according to CISA KEV.",
				Transform:   transform.FromField("Description.CisaActionDue"),
			},
			{
				Name:        "cisa_required_action",
				Type:        proto.ColumnType_STRING,
				Description: "The action required by CISA for federal agencies.",
				Transform:   transform.FromField("Description.CisaRequiredAction"),
			},
			{
				Name:        "cisa_vulnerability_name",
				Type:        proto.ColumnType_STRING,
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
