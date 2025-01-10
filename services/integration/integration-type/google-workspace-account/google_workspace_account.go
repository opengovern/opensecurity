package google_workspace_account

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	googleWorkspaceDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
)

type GoogleWorkspaceAccountIntegration struct{}

func (i *GoogleWorkspaceAccountIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   googleWorkspaceDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      googleWorkspaceDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           googleWorkspaceDescriberLocal.StreamName,
		NatsConsumerGroup:        googleWorkspaceDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: googleWorkspaceDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "googleworkspace",

		UISpecFileName: "google-workspace-account.json",

		DescriberDeploymentName: googleWorkspaceDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     googleWorkspaceDescriberLocal.DescriberRunCommand,
	}
}

func (i *GoogleWorkspaceAccountIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials googleWorkspaceDescriberLocal.IntegrationCredentials
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

func (i *GoogleWorkspaceAccountIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials googleWorkspaceDescriberLocal.IntegrationCredentials
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

func (i *GoogleWorkspaceAccountIntegration) GetResourceTypesByLabels(map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range googleWorkspaceDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *GoogleWorkspaceAccountIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := googleWorkspaceDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *GoogleWorkspaceAccountIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range googleWorkspaceDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
