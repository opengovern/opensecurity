package opengovernance_client

import (
	"context"
	"github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/services"
	"github.com/opengovern/opensecurity/services/integration/api/models"
	"runtime"

	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type IntegrationResourceTypeRow struct {
	ResourceType    string      `json:"resource_type"`
	IntegrationType string      `json:"integration_type"`
	Description     string      `json:"description"`
	Params          []Parameter `json:"params"`
	Table           string      `json:"table"`
}

type Parameter struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Required    bool    `json:"required"`
	Default     *string `json:"default"`
}

func getIntegrationResourceTypeRowFromIntegrationResourceType(rt models.ResourceTypeConfiguration) IntegrationResourceTypeRow {
	var params []Parameter

	for _, param := range rt.Params {
		params = append(params, Parameter{
			Name:        param.Name,
			Description: param.Description,
			Required:    param.Required,
			Default:     param.Default,
		})
	}

	row := IntegrationResourceTypeRow{
		ResourceType:    rt.Name,
		IntegrationType: string(rt.IntegrationType),
		Description:     rt.Description,
		Params:          params,
		Table:           rt.Table,
	}

	return row
}

func ListIntegrationResourceTypes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListIntegrationResourceTypes")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	integrationClient, err := services.NewIntegrationClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}

	integrationType := d.EqualsQuals["integration_type"].GetStringValue()

	rts, err := integrationClient.ListIntegrationTypeResourceTypes(&httpclient.Context{UserRole: api.AdminRole}, integrationType)
	if err != nil {
		plugin.Logger(ctx).Error("GetBenchmarkSummary compliance client call failed", "error", err)
		return nil, err
	}
	if rts == nil || len(rts.ResourceTypes) == 0 {
		return nil, nil
	}

	for _, rt := range rts.ResourceTypes {
		row := getIntegrationResourceTypeRowFromIntegrationResourceType(rt)
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func GetIntegrationResourceType(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListIntegrationResourceTypes")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	integrationClient, err := services.NewIntegrationClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}

	integrationType := d.EqualsQuals["integration_type"].GetStringValue()
	resourceType := d.EqualsQuals["resource_type"].GetStringValue()

	rt, err := integrationClient.GetIntegrationTypeResourceType(&httpclient.Context{UserRole: api.AdminRole}, integrationType, resourceType)
	if err != nil {
		plugin.Logger(ctx).Error("GetBenchmarkSummary compliance client call failed", "error", err)
		return nil, err
	}
	if rt == nil {
		return nil, nil
	}

	return getIntegrationResourceTypeRowFromIntegrationResourceType(*rt), nil
}
