package google_workspace_account

import (
	"encoding/json"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	renderDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/render-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/render-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/render-account/healthcheck"
	"github.com/opengovern/opencomply/services/integration/models"
)

type RenderAccountIntegration struct{}

func (i *RenderAccountIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   renderDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      renderDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           renderDescriberLocal.StreamName,
		NatsConsumerGroup:        renderDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: renderDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "render",

		UISpecFileName: "render-account.json",

		DescriberDeploymentName: renderDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     renderDescriberLocal.DescriberRunCommand,
	}
}

func (i *RenderAccountIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials renderDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	isHealthy, err := healthcheck.RenderIntegrationHealthcheck(healthcheck.Config{
		APIKey: credentials.APIKey,
	})
	return isHealthy, err
}

func (i *RenderAccountIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials renderDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []models.Integration
	user, err := discovery.RenderIntegrationDiscovery(discovery.Config{
		APIKey: credentials.APIKey,
	})
	integrations = append(integrations, models.Integration{
		ProviderID: user.Email,
		Name:       user.Name,
	})
	return integrations, nil
}

func (i *RenderAccountIntegration) GetResourceTypesByLabels(map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range renderDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *RenderAccountIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := renderDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *RenderAccountIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range renderDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
