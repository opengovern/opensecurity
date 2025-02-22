package opengovernance_client

import (
	"context"
	"runtime"
	"time"

	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/services"
	"github.com/opengovern/opensecurity/pkg/utils"
	complianceApi "github.com/opengovern/opensecurity/services/compliance/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func GetBenchmarkSummary(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("GetBenchmarkSummary")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	complianceClient, err := services.NewComplianceClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}

	benchmarkId := d.EqualsQuals["benchmark_id"].GetStringValue()

	var timeAt *time.Time
	if d.Quals["time_at"] != nil {
		timeAt = utils.GetPointer(d.EqualsQuals["time_at"].GetTimestampValue().AsTime())
	}
	var integrationIds []string
	if d.EqualsQuals["integration_id"] != nil {
		q := d.EqualsQuals["integration_id"]
		if q.GetListValue() != nil {
			for _, v := range q.GetListValue().Values {
				integrationIds = append(integrationIds, v.GetStringValue())
			}
		} else {
			integrationIds = []string{d.EqualsQuals["integration_id"].GetStringValue()}
		}
	}

	res, err := complianceClient.GetBenchmarkSummary(&httpclient.Context{UserRole: api.AdminRole}, benchmarkId, integrationIds, timeAt)
	if err != nil {
		plugin.Logger(ctx).Error("GetBenchmarkSummary compliance client call failed", "error", err)
		return nil, err
	}

	return res, nil
}

func handleBenchmarkControlSummary(ctx context.Context, d *plugin.QueryData, res complianceApi.BenchmarkControlSummary) {
	for _, control := range res.Controls {
		d.StreamListItem(ctx, control)
	}
	for _, child := range res.Children {
		handleBenchmarkControlSummary(ctx, d, child)
	}
}

func ListBenchmarkControls(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListBenchmarkControls")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	complianceClient, err := services.NewComplianceClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}

	benchmarkId := d.EqualsQuals["benchmark_id"].GetStringValue()

	var timeAt *time.Time
	if d.Quals["time_at"] != nil {
		timeAt = utils.GetPointer(d.EqualsQuals["time_at"].GetTimestampValue().AsTime())
	}
	var integrationIds []string
	if d.EqualsQuals["integration_id"] != nil {
		q := d.EqualsQuals["integration_id"]
		if q.GetListValue() != nil {
			for _, v := range q.GetListValue().Values {
				integrationIds = append(integrationIds, v.GetStringValue())
			}
		} else {
			integrationIds = []string{d.EqualsQuals["integration_id"].GetStringValue()}
		}
	}

	apiRes, err := complianceClient.GetBenchmarkControls(&httpclient.Context{UserRole: api.AdminRole}, benchmarkId, integrationIds, timeAt)
	if err != nil {
		plugin.Logger(ctx).Error("GetBenchmarkSummary compliance client call failed", "error", err)
		return nil, err
	}
	if apiRes == nil {
		return nil, nil
	}

	handleBenchmarkControlSummary(ctx, d, *apiRes)

	return nil, nil
}
