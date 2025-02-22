package entities

import (
	"github.com/opengovern/og-util/pkg/steampipe"
	api "github.com/opengovern/opensecurity/services/integration/api/models"
	"github.com/opengovern/opensecurity/services/integration/models"
	"golang.org/x/net/context"
)

func NewIntegrationGroup(ctx context.Context, steampipe *steampipe.Database, cg models.IntegrationGroup) (*api.IntegrationGroup, error) {
	apiCg := api.IntegrationGroup{
		Name:  cg.Name,
		Query: cg.Query,
	}

	if steampipe == nil || cg.Query == "" {
		return &apiCg, nil
	}

	integrationsQueryResult, err := steampipe.QueryAll(ctx, cg.Query)
	if err != nil {
		return nil, err
	}

	var integrationIds []string
	for i, header := range integrationsQueryResult.Headers {
		if header != "integration_id" {
			continue
		}
		for _, row := range integrationsQueryResult.Data {
			if len(row) <= i || row[i] == nil {
				continue
			}
			if strRow, ok := row[i].(string); ok {
				integrationIds = append(integrationIds, strRow)
			}
		}
	}

	apiCg.IntegrationIds = integrationIds

	return &apiCg, nil
}
