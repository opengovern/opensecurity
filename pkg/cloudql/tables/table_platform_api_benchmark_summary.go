package opengovernance

import (
	"context"
	"time"

	og_client "github.com/opengovern/opensecurity/pkg/cloudql/client"
	compliance "github.com/opengovern/opensecurity/services/compliance/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/quals"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tablePlatformApiBenchmarkSummary(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_api_benchmark_summary",
		Description: "Wrapper for benchmark summary api",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:      "benchmark_id",
					Operators: []string{quals.QualOperatorEqual},
					Require:   plugin.Required,
				},
				{
					Name:      "time_at",
					Operators: []string{quals.QualOperatorEqual},
					Require:   plugin.Optional,
				},
				{
					Name:      "integration_id",
					Operators: []string{quals.QualOperatorEqual},
					Require:   plugin.Optional,
				},
			},
			Hydrate: og_client.GetBenchmarkSummary,
		},
		Columns: []*plugin.Column{
			{
				Name:        "benchmark_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the benchmark in OpenGovernance",
				Transform:   transform.FromQual("benchmark_id"),
			},
			{
				Name:        "integration_id",
				Type:        proto.ColumnType_STRING,
				Description: "The integration IDs included in the benchmark summary",
				Transform:   transform.FromQual("integration_id"),
			},
			{
				Name:        "time_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The timestamp of the benchmark summary record",
				Transform:   transform.FromQual("time_at"),
				Default:     time.Now(),
			},
			{
				Name:        "compliance_status_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of checks that passed in the benchmark summary",
				Transform:   transform.FromField("ComplianceStatusSummary.PassedCount"),
			},
			{
				Name:        "compliance_status_failed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of checks that failed in the benchmark summary",
				Transform:   transform.FromField("ComplianceStatusSummary.FailedCount"),
			},
			{
				Name:        "severity_result_critical_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of critical severity findings in the benchmark summary",
				Transform:   transform.FromField("Checks.CriticalCount"),
			},
			{
				Name:        "severity_result_high_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of high severity findings in the benchmark summary",
				Transform:   transform.FromField("Checks.HighCount"),
			},
			{
				Name:        "severity_result_medium_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of medium severity findings in the benchmark summary",
				Transform:   transform.FromField("Checks.MediumCount"),
			},
			{
				Name:        "severity_result_low_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of low severity findings in the benchmark summary",
				Transform:   transform.FromField("Checks.LowCount"),
			},
			{
				Name:        "severity_result_none_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of findings with no severity in the benchmark summary",
				Transform:   transform.FromField("Checks.NoneCount"),
			},
			{
				Name:        "controls_severity_status_total_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of controls in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Total.TotalCount"),
			},
			{
				Name:        "controls_severity_status_total_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of controls that passed in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Total.PassedCount"),
			},
			{
				Name:        "controls_severity_status_critical_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of critical severity controls in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Critical.TotalCount"),
			},
			{
				Name:        "controls_severity_status_critical_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of critical severity controls that passed in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Critical.PassedCount"),
			},
			{
				Name:        "controls_severity_status_high_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of high severity controls in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.High.TotalCount"),
			},
			{
				Name:        "controls_severity_status_high_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of high severity controls that passed in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.High.PassedCount"),
			},
			{
				Name:        "controls_severity_status_medium_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of medium severity controls in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Medium.TotalCount"),
			},
			{
				Name:        "controls_severity_status_medium_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of medium severity controls that passed in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Medium.PassedCount"),
			},
			{
				Name:        "controls_severity_status_low_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of low severity controls in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Low.TotalCount"),
			},
			{
				Name:        "controls_severity_status_low_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of low severity controls that passed in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.Low.PassedCount"),
			},
			{
				Name:        "controls_severity_status_none_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of controls with no severity in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.None.TotalCount"),
			},
			{
				Name:        "controls_severity_status_none_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of controls with no severity that passed in the benchmark summary",
				Transform:   transform.FromField("ControlsSeverityStatus.None.PassedCount"),
			},
			{
				Name:        "resources_severity_status_total_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of resources in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Total.TotalCount"),
			},
			{
				Name:        "resources_severity_status_total_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of resources that passed in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Total.PassedCount"),
			},
			{
				Name:        "resources_severity_status_critical_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of critical severity resources in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Critical.TotalCount"),
			},
			{
				Name:        "resources_severity_status_critical_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of critical severity resources that passed in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Critical.PassedCount"),
			},
			{
				Name:        "resources_severity_status_high_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of high severity resources in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.High.TotalCount"),
			},
			{
				Name:        "resources_severity_status_high_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of high severity resources that passed in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.High.PassedCount"),
			},
			{
				Name:        "resources_severity_status_medium_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of medium severity resources in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Medium.TotalCount"),
			},
			{
				Name:        "resources_severity_status_medium_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of medium severity resources that passed in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Medium.PassedCount"),
			},
			{
				Name:        "resources_severity_status_low_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of low severity resources in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Low.TotalCount"),
			},
			{
				Name:        "resources_severity_status_low_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of low severity resources that passed in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.Low.PassedCount"),
			},
			{
				Name:        "resources_severity_status_none_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of resources with no severity in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.None.TotalCount"),
			},
			{
				Name:        "resources_severity_status_none_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of resources with no severity that passed in the benchmark summary",
				Transform:   transform.FromField("ResourcesSeverityStatus.None.PassedCount"),
			},
			{
				Name:        "integrations_result_compliance_status_passed_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of checks that passed in the benchmark summary for the integration",
				Transform:   transform.FromField("IntegrationsStatus.PassedCount"),
			},
			{
				Name:        "integrations_result_compliance_status_total_count",
				Type:        proto.ColumnType_INT,
				Description: "The total number of checks in the benchmark summary for the integration",
				Transform:   transform.FromField("IntegrationsStatus.TotalCount"),
			},
			{
				Name:        "cost_impact",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The cost impact score of the benchmark summary",
				Transform:   transform.FromField("CostImpact"),
			},
			{
				Name:        "evaluated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The timestamp when the benchmark summary was evaluated",
				Transform:   transform.FromField("EvaluatedAt"),
			},
			{
				Name:        "api_result",
				Type:        proto.ColumnType_JSON,
				Description: "The result of the benchmark summary",
				Transform:   transform.From(getOpenGovernanceApiBenchmarkSummaryResult),
			},
		},
	}
}

func getOpenGovernanceApiBenchmarkSummaryResult(_ context.Context, d *transform.TransformData) (any, error) {
	return d.HydrateItem.(*compliance.BenchmarkEvaluationSummary), nil
}
