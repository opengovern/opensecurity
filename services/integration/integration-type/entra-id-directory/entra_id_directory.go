package entra_id_directory

import (
	"encoding/json"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory/healthcheck"
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

		SteampipePluginName: "entraid",

		UISpec: configs.UISpec,

		DescriberDeploymentName: configs.DescriberDeploymentName,
		DescriberRunCommand:     configs.DescriberRunCommand,
	}
}

func (i *Integration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var configs configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &configs)
	if err != nil {
		return false, err
	}

	return healthcheck.EntraidIntegrationHealthcheck(healthcheck.Config{
		TenantID:     providerId,
		ClientID:     configs.ClientID,
		ClientSecret: configs.ClientPassword,
		CertPath:     "",
		CertContent:  configs.Certificate,
		CertPassword: configs.CertificatePassword,
	})
}

func (i *Integration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var configs configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &configs)
	if err != nil {
		return nil, err
	}

	var integrations []models.Integration
	directories, err := discovery.EntraidIntegrationDiscovery(discovery.Config{
		TenantID:     configs.TenantID,
		ClientID:     configs.ClientID,
		ClientSecret: configs.ClientPassword,
		CertPath:     "",
		CertContent:  configs.Certificate,
		CertPassword: configs.CertificatePassword,
	})
	if err != nil {
		return nil, err
	}
	for _, s := range directories {
		integrations = append(integrations, models.Integration{
			ProviderID: s.TenantID,
			Name:       s.Name,
		})
	}

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
	return configs.IntegrationTypeEntraidDirectory
}
