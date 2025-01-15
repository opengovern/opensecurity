package cloudflare_account

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
)

type Integration struct{}

func (i *Integration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   configs.JobQueueTopic,
		NatsManualJobsTopic:      configs.JobQueueTopicManuals,
		NatsStreamName:           configs.StreamName,
		NatsConsumerGroup:        configs.ConsumerGroup,
		NatsConsumerGroupManuals: configs.ConsumerGroupManuals,

		SteampipePluginName: "cloudflare",

		UISpec: configs.UISpec,

		DescriberRunCommand:     configs.DescriberRunCommand,
		DescriberDeploymentName: configs.DescriberDeploymentName,
	}
}

func (i *Integration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	isHealthy, err := healthcheck.CloudflareIntegrationHealthcheck(healthcheck.Config{
		Token:     credentials.Token,
		AccountID: credentials.AccountID,
	})
	return isHealthy, err
}

func (i *Integration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []models.Integration
	account, err := discovery.CloudflareIntegrationDiscovery(discovery.Config{
		Token:     credentials.Token,
		AccountID: credentials.AccountID,
	})
	labels := map[string]string{
		"Type": account.Type,
	}
	labelsJsonData, err := json.Marshal(labels)
	if err != nil {
		return nil, err
	}
	integrationLabelsJsonb := pgtype.JSONB{}
	err = integrationLabelsJsonb.Set(labelsJsonData)
	if err != nil {
		return nil, err
	}
	integrations = append(integrations, models.Integration{
		ProviderID: account.ID,
		Name:       account.Name,
		Labels:     integrationLabelsJsonb,
	})
	return integrations, nil
}

func (i *Integration) GetResourceTypesByLabels(labels map[string]string) (map[string]interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]interfaces.ResourceTypeConfiguration)
	for _, resourceType := range configs.ResourceTypesList {
		resourceTypesMap[resourceType] = interfaces.ResourceTypeConfiguration{}
	}
	return resourceTypesMap, nil
}

func (i *Integration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := configs.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *Integration) GetIntegrationType() integration.Type {
	return configs.IntegrationNameCloudflareAccount
}

func (i *Integration) ListAllTables() map[string][]string {
	return make(map[string][]string)
}
