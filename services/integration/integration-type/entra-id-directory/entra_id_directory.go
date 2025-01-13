package entra_id_directory

import (
	"encoding/json"
	entraidDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
)

type EntraIdDirectoryIntegration struct{}

func (i *EntraIdDirectoryIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   entraidDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      entraidDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           entraidDescriberLocal.StreamName,
		NatsConsumerGroup:        entraidDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: entraidDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "entraid",

		UISpecFileName: "entraid-directory.json",

		DescriberDeploymentName: entraidDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     entraidDescriberLocal.DescriberRunCommand,
	}
}

func (i *EntraIdDirectoryIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var configs entraidDescriberLocal.IntegrationCredentials
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

func (i *EntraIdDirectoryIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var configs entraidDescriberLocal.IntegrationCredentials
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

func (i *EntraIdDirectoryIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range entraidDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *EntraIdDirectoryIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := entraidDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *EntraIdDirectoryIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range entraidDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
