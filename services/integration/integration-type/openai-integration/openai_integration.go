package openai_integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	openaiDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/openai-integration/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/openai-integration/healthcheck"
	"github.com/opengovern/opencomply/services/integration/models"
)

type OpenAIIntegration struct{}

func (i *OpenAIIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   openaiDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      openaiDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           openaiDescriberLocal.StreamName,
		NatsConsumerGroup:        openaiDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: openaiDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "openai",

		UISpecFileName: "openai-integration.json",

		DescriberDeploymentName: openaiDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     openaiDescriberLocal.DescriberRunCommand,
	}
}

func (i *OpenAIIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials openaiDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	isHealthy, err := healthcheck.OpenAIIntegrationHealthcheck(credentials.APIKey)
	return isHealthy, err
}

func (i *OpenAIIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials openaiDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []models.Integration
	_, err = healthcheck.OpenAIIntegrationHealthcheck(credentials.APIKey)
	if err != nil {
		return nil, err
	}
	labels := map[string]string{
		"OrganizationID": credentials.OrganizationID,
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
	providerID := hashSHA256(credentials.APIKey)
	integrations = append(integrations, models.Integration{
		ProviderID: providerID,
		Name:       credentials.ProjectName,
		Labels:     integrationLabelsJsonb,
	})

	return integrations, nil
}

func (i *OpenAIIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range openaiDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *OpenAIIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := openaiDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func hashSHA256(input string) string {
	hash := sha256.New()

	hash.Write([]byte(input))

	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func (i *OpenAIIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range openaiDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
