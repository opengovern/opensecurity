package opengovernance_client

import (
	"context"
	"encoding/json"
	"runtime"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/pg"
	integration "github.com/opengovern/opensecurity/services/integration/models"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type IntegrationRow struct {
	IntegrationID   string            `json:"integration_id"`
	ProviderID      string            `json:"provider_id"`
	Name            string            `json:"name"`
	IntegrationType string            `json:"integration_type"`
	State           string            `json:"state"`
	Annotations     map[string]string `json:"annotations"`
	Labels          map[string]string `json:"labels"`
}

func getIntegrationRowFromIntegration(ctx context.Context, integration integration.Integration) (*IntegrationRow, error) {
	row := IntegrationRow{
		IntegrationID:   integration.IntegrationID.String(),
		ProviderID:      integration.ProviderID,
		Name:            integration.Name,
		IntegrationType: strings.ToLower(integration.IntegrationType.String()),
		State:           string(integration.State),
		Annotations:     make(map[string]string),
		Labels:          make(map[string]string),
	}
	var annotations map[string]string
	if integration.Annotations.Status == pgtype.Present {
		if err := json.Unmarshal(integration.Annotations.Bytes, &annotations); err != nil {
			return nil, err
		}
	}
	row.Annotations = annotations

	var labels map[string]string
	if integration.Labels.Status == pgtype.Present {
		if err := json.Unmarshal(integration.Labels.Bytes, &labels); err != nil {
			return nil, err
		}
	}
	row.Labels = labels

	return &row, nil
}

func ListIntegrations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListIntegrations")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	ke, err := pg.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{PG: ke}

	integrations, err := k.PG.ListIntegrations(ctx)
	if err != nil {
		return nil, err
	}

	for _, i := range integrations {
		row, err := getIntegrationRowFromIntegration(ctx, i)
		if err != nil {
			plugin.Logger(ctx).Error("ListIntegrations", "integration", i, "error", err)
			continue
		}
		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func GetIntegration(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("GetIntegration")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	ke, err := pg.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{PG: ke}

	opengovernanceId := d.EqualsQuals["integration_id"].GetStringValue()
	id := d.EqualsQuals["provider_id"].GetStringValue()
	i, err := k.PG.GetIntegrationByID(ctx, opengovernanceId, id)
	if err != nil {
		return nil, err
	}
	if i == nil {
		return nil, nil
	}

	row, err := getIntegrationRowFromIntegration(ctx, *i)
	if err != nil {
		plugin.Logger(ctx).Error("GetIntegration", "integration", i, "error", err)
		return nil, err
	}
	return row, nil
}
