package opengovernance_client

import (
	"context"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/pg"
	integration "github.com/opengovern/opensecurity/services/integration/models"
	"runtime"

	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type IntegrationCredentialRow struct {
	IntegrationID   string `json:"resource_type"`
	IntegrationType string `json:"integration_type"`
	CredentialID    string `json:"credential_id"`
	CredentialType  string `json:"credential_type"`
	Secret          string `json:"secret"`
}

func getIntegrationCredentialRowFromIntegrationCredential(integrationId string, c integration.Credential) IntegrationCredentialRow {
	row := IntegrationCredentialRow{
		IntegrationID:   integrationId,
		IntegrationType: string(c.IntegrationType),
		CredentialID:    c.ID.String(),
		CredentialType:  c.CredentialType,
		Secret:          c.Secret,
	}

	return row
}

func ListIntegrationCredentials(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListIntegrationCredentials")
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
		credential, err := k.PG.GetCredential(ctx, i.CredentialID)
		if err != nil {
			plugin.Logger(ctx).Error("ListIntegrations", "integration", i, "error", err)
			continue
		}
		if credential == nil || credential.Secret == "" {
			continue
		}
		row := getIntegrationCredentialRowFromIntegrationCredential(i.IntegrationID.String(), *credential)

		d.StreamListItem(ctx, row)
	}

	return nil, nil
}

func GetIntegrationCredential(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {
	plugin.Logger(ctx).Trace("ListIntegrationResourceTypes")
	runtime.GC()
	cfg := config.GetConfig(d.Connection)
	ke, err := pg.NewClientCached(cfg, d.ConnectionCache, ctx)
	if err != nil {
		return nil, err
	}
	k := Client{PG: ke}

	opengovernanceId := d.EqualsQuals["integration_id"].GetStringValue()
	i, err := k.PG.GetIntegrationByID(ctx, opengovernanceId, "")
	if err != nil {
		return nil, err
	}

	credential, err := k.PG.GetCredential(ctx, i.CredentialID)
	if err != nil {
		plugin.Logger(ctx).Error("ListIntegrations", "integration", i, "error", err)
		return nil, err
	}
	if credential == nil || credential.Secret == "" {
		return nil, nil
	}

	row := getIntegrationCredentialRowFromIntegrationCredential(i.IntegrationID.String(), *credential)

	return row, nil
}
