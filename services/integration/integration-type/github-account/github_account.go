package github_account

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	githubDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/github-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/github-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/github-account/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
	"strconv"
)

type GithubAccountIntegration struct{}

func (i *GithubAccountIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   githubDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      githubDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           githubDescriberLocal.StreamName,
		NatsConsumerGroup:        githubDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: githubDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "github",

		UISpecFileName: "github-account.json",

		DescriberDeploymentName: githubDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     githubDescriberLocal.DescriberRunCommand,
	}
}

func (i *GithubAccountIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials githubDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	var name string
	if v, ok := labels["OrganizationName"]; ok {
		name = v
	}
	isHealthy, err := healthcheck.GithubIntegrationHealthcheck(healthcheck.Config{
		Token:            credentials.PatToken,
		OrganizationName: name,
	})
	return isHealthy, err
}

func (i *GithubAccountIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials githubDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []models.Integration
	accounts, err := discovery.GithubIntegrationDiscovery(discovery.Config{
		Token: credentials.PatToken,
	})
	if err != nil {
		return nil, err
	}
	for _, a := range accounts {
		labels := map[string]string{
			"OrganizationName": a.Login,
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
			ProviderID: strconv.FormatInt(a.ID, 10),
			Name:       a.Login,
			Labels:     integrationLabelsJsonb,
		})
	}
	return integrations, nil
}

func (i *GithubAccountIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range githubDescriberLocal.ResourceTypesList {
		if v, ok := githubDescriberLocal.ResourceTypeConfigs[resourceType]; ok {
			resourceTypesMap[resourceType] = v
		} else {
			resourceTypesMap[resourceType] = nil
		}
	}
	return resourceTypesMap, nil
}

func (i *GithubAccountIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := githubDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *GithubAccountIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range githubDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
