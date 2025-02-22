package opengovernance_client

import (
	"context"
	"runtime"

	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/pg"
	integration "github.com/opengovern/opensecurity/services/integration/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type IntegrationGroupRow struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

func getIntegrationGroupRowFromIntegrationGroup(ctx context.Context, integrationGroup integration.IntegrationGroup) (*IntegrationGroupRow, error) {
	row := IntegrationGroupRow{
		Name:  integrationGroup.Name,
		Query: integrationGroup.Query,
	}

	return &row, nil
}

func ListIntegrationsGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListIntegrationsGroup")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	ke, err := pg.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{PG: ke}

	integrationGroups, err := k.PG.ListIntegrationGroups(ctx)
	if err != nil {
		return nil, err
	}

	for _, i := range integrationGroups {
		row, err := getIntegrationGroupRowFromIntegrationGroup(ctx, i)
		if err != nil {
			plugin.Logger(ctx).Error("ListIntegrations", "integration", i, "error", err)
			continue
		}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func GetIntegrationGroup(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("GetIntegrationGroup")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	ke, err := pg.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{PG: ke}

	name := d.EqualsQuals["name"].GetStringValue()
	i, err := k.PG.GetIntegrationGroupByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, nil
	}

	row, err := getIntegrationGroupRowFromIntegrationGroup(ctx, *i)
	if err != nil {
		plugin.Logger(ctx).Error("GetIntegration", "integration", i, "error", err)
		return nil, err
	}
	return row, nil
}
