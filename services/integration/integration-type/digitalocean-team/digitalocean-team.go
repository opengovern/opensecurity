package digitalocean_team

import (
	"context"
	"encoding/json"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team/healthcheck"
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

		SteampipePluginName: "digitalocean",

		UISpec: configs.UISpec,

		DescriberDeploymentName: configs.DescriberDeploymentName,
		DescriberRunCommand:     configs.DescriberRunCommand,
	}
}

func (i *Integration) HealthCheck(jsonData []byte, _ string, _ map[string]string, _ map[string]string) (bool, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	return healthcheck.DigitalOceanTeamHealthcheck(context.TODO(), healthcheck.Config{
		AuthToken: credentials.AuthToken,
	})
}

func (i *Integration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}

	team, err := discovery.DigitalOceanTeamDiscovery(context.TODO(), discovery.Config{
		AuthToken: credentials.AuthToken,
	})
	if err != nil {
		return nil, err
	}

	return []models.Integration{
		{
			ProviderID: team.ID,
			Name:       team.Name,
		},
	}, nil
}

func (i *Integration) GetResourceTypesByLabels(labels map[string]string) (map[string]interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]interfaces.ResourceTypeConfiguration)
	for _, resourceType := range configs.ResourceTypesList {
		resourceTypesMap[resourceType] = interfaces.ResourceTypeConfiguration{}
	}
	return resourceTypesMap, nil
}

func (i *Integration) GetResourceTypeFromTableName(tableName string) (string, error) {
	if v, ok := configs.TablesToResourceTypes[tableName]; ok {
		return v, nil
	}

	return "", nil
}

func (i *Integration) GetIntegrationType() integration.Type {
	return configs.IntegrationTypeDigitalOceanTeam
}

func (i *Integration) ListAllTables() (map[string][]string, error) {
	return make(map[string][]string), nil
}

func (i *DigitaloceanTeamIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range digitaloceanDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
