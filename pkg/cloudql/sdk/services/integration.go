package services

import (
	"context"
	"errors"
	"github.com/opengovern/opensecurity/pkg/cloudql/sdk/config"
	integrationClient "github.com/opengovern/opensecurity/services/integration/client"
	"github.com/turbot/steampipe-plugin-sdk/v5/connection"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func NewIntegrationClientCached(c config.ClientConfig, cache *connection.ConnectionCache, ctx context.Context) (integrationClient.IntegrationServiceClient, error) {
	value, ok := cache.Get(ctx, "opengovernance-integration-service-client")
	if ok {
		return value.(integrationClient.IntegrationServiceClient), nil
	}

	plugin.Logger(ctx).Warn("integration service client is not cached, creating a new one")

	if c.IntegrationServiceBaseURL == nil {
		plugin.Logger(ctx).Error("integration service base url is not set")
		return nil, errors.New("integration service base url is not set")
	}
	client := integrationClient.NewIntegrationServiceClient(*c.IntegrationServiceBaseURL)

	cache.Set(ctx, "opengovernance-integration-service-client", client)

	return client, nil
}
