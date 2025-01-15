package google_workspace_account

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account/healthcheck"
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

		SteampipePluginName: "googleworkspace",

		UISpec: configs.UISpec,

		DescriberDeploymentName: configs.DescriberDeploymentName,
		DescriberRunCommand:     configs.DescriberRunCommand,
	}
}

func (i *Integration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	isHealthy, err := healthcheck.GoogleWorkspaceIntegrationHealthcheck(healthcheck.Config{
		AdminEmail: credentials.AdminEmail,
		CustomerID: credentials.CustomerID,
		KeyFile:    credentials.KeyFile,
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
	customer, err := discovery.GoogleWorkspaceIntegrationDiscovery(discovery.Config{
		AdminEmail: credentials.AdminEmail,
		CustomerID: credentials.CustomerID,
		KeyFile:    credentials.KeyFile,
	})
	labels := map[string]string{
		"Domain": customer.Domain,
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
		ProviderID: customer.ID,
		Name:       customer.Domain,
		Labels:     integrationLabelsJsonb,
	})
	return integrations, nil
}

func (i *Integration) GetResourceTypesByLabels(map[string]string) (map[string]interfaces.ResourceTypeConfiguration, error) {
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
	return configs.IntegrationTypeGoogleWorkspaceAccount
}
