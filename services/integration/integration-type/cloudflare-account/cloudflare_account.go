package cloudflare_account

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	cloudflareDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
)

type CloudFlareAccountIntegration struct{}

func (i *CloudFlareAccountIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   cloudflareDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      cloudflareDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           cloudflareDescriberLocal.StreamName,
		NatsConsumerGroup:        cloudflareDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: cloudflareDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "cloudflare",

		UISpecFileName: "cloudflare-account.json",

		DescriberRunCommand:     cloudflareDescriberLocal.DescriberRunCommand,
		DescriberDeploymentName: cloudflareDescriberLocal.DescriberDeploymentName,
	}
}

func (i *CloudFlareAccountIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials cloudflareDescriberLocal.IntegrationCredentials
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

func (i *CloudFlareAccountIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials cloudflareDescriberLocal.IntegrationCredentials
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

func (i *CloudFlareAccountIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range cloudflareDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *CloudFlareAccountIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := cloudflareDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *CloudFlareAccountIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range cloudflareDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
