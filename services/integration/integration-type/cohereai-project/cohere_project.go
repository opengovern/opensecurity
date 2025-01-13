package cohereai_project

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	cohereaiDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/cohereai-project/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/cohereai-project/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/cohereai-project/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"

	"github.com/opengovern/opencomply/services/integration/models"
)

type CohereAIProjectIntegration struct{}

func (i *CohereAIProjectIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   cohereaiDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      cohereaiDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           cohereaiDescriberLocal.StreamName,
		NatsConsumerGroup:        cohereaiDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: cohereaiDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "cohereai",

		UISpecFileName: "cohereai-project.json",

		DescriberDeploymentName: cohereaiDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     cohereaiDescriberLocal.DescriberRunCommand,
	}
}

func (i *CohereAIProjectIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials cohereaiDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	isHealthy, err := healthcheck.CohereAIIntegrationHealthcheck(credentials.APIKey)
	return isHealthy, err
}

func (i *CohereAIProjectIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials cohereaiDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []models.Integration
	connectors, err1 := discovery.CohereAIIntegrationDiscovery(credentials.APIKey)
	if err1 != nil {
		return nil, err1
	}
	labels := map[string]string{
		"ClientName": credentials.ClientName,
	}
	if len(connectors) > 0 {
		labels["OrganizationID"] = connectors[0].OrganizationID
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
	// for in esponse
	for _, connector := range connectors {
		integrations = append(integrations, models.Integration{
			ProviderID: connector.ID,
			Name:       connector.Name,
			Labels:     integrationLabelsJsonb,
		})
	}

	return integrations, nil
}

func (i *CohereAIProjectIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range cohereaiDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *CohereAIProjectIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := cohereaiDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *CohereAIProjectIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range cohereaiDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
