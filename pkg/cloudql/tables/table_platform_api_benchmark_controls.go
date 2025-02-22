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

func tablePlatformApiBenchmarkControls(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "platform_api_benchmark_controls",
		Description: "Wrapper for benchmark summary api",
		Cache: &plugin.TableCacheOptions{
			Enabled: false,
		},
		List: &plugin.ListConfig{
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
					Name:      "connection_id",
					Operators: []string{quals.QualOperatorEqual},
					Require:   plugin.Optional,
				},
			},
			Hydrate: og_client.ListBenchmarkControls,
		},
		Columns: []*plugin.Column{
			{
				Name:        "benchmark_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the benchmark in the platform",
				Transform:   transform.FromQual("benchmark_id"),
			},
			{
				Name:        "connection_id",
				Type:        proto.ColumnType_STRING,
				Description: "The connection IDs included in the benchmark controls",
				Transform:   transform.FromQual("connection_id"),
			},
			{
				Name:        "time_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The timestamp of the benchmark controls record",
				Transform:   transform.FromQual("time_at"),
				Default:     time.Now(),
			},
			{
				Name:        "control_id",
				Type:        proto.ColumnType_STRING,
				Description: "Control id",
				Transform:   transform.FromField("Control.ID"),
			},
			{
				Name:        "control_title",
				Type:        proto.ColumnType_STRING,
				Description: "Control title",
				Transform:   transform.FromField("Control.Title"),
			},
			{
				Name:        "control_description",
				Type:        proto.ColumnType_STRING,
				Description: "Control description",
				Transform:   transform.FromField("Control.Description"),
			},
			{
				Name:        "control",
				Type:        proto.ColumnType_JSON,
				Description: "The control object",
				Transform:   transform.FromField("Control"),
			},
			{
				Name:        "is_passed",
				Type:        proto.ColumnType_BOOL,
				Description: "The status of the control",
				Transform:   transform.FromField("Passed"),
			},
			{
				Name:        "failed_resources_count",
				Type:        proto.ColumnType_INT,
				Description: "The count of failed resources",
				Transform:   transform.FromField("FailedResourcesCount"),
			},
			{
				Name:        "total_resources_count",
				Type:        proto.ColumnType_INT,
				Description: "The total count of resources",
				Transform:   transform.FromField("TotalResourcesCount"),
			},
			{
				Name:        "failed_integration_count",
				Type:        proto.ColumnType_INT,
				Description: "The count of failed connections",
				Transform:   transform.FromField("FailedIntegrationCount"),
			},
			{
				Name:        "total_integration_count",
				Type:        proto.ColumnType_INT,
				Description: "The total count of connections",
				Transform:   transform.FromField("TotalIntegrationCount"),
			},
			{
				Name:        "cost_impact",
				Type:        proto.ColumnType_DOUBLE,
				Description: "The cost impact",
				Transform:   transform.FromField("CostImpact"),
			},
			{
				Name:        "evaluated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The timestamp of the evaluation",
				Transform:   transform.FromField("EvaluatedAt"),
			},
			{
				Name:        "api_result",
				Type:        proto.ColumnType_JSON,
				Description: "The result of the benchmark control summary",
				Transform:   transform.From(getOpenGovernanceApiBenchmarkControlResult),
			},
		},
	}
}

func getOpenGovernanceApiBenchmarkControlResult(_ context.Context, d *transform.TransformData) (any, error) {
	return d.HydrateItem.(compliance.ControlSummary), nil
}
