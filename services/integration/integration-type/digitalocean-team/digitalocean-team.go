package digitalocean_team

import (
	"context"
	"encoding/json"
	digitaloceanDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
)

type DigitaloceanTeamIntegration struct{}

func (i *DigitaloceanTeamIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   digitaloceanDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      digitaloceanDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           digitaloceanDescriberLocal.StreamName,
		NatsConsumerGroup:        digitaloceanDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: digitaloceanDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "digitalocean",

		UISpecFileName: "digitalocean-team.json",

		DescriberDeploymentName: digitaloceanDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     digitaloceanDescriberLocal.DescriberRunCommand,
	}
}

func (i *DigitaloceanTeamIntegration) HealthCheck(jsonData []byte, _ string, _ map[string]string, _ map[string]string) (bool, error) {
	var credentials digitaloceanDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	return healthcheck.DigitalOceanTeamHealthcheck(context.TODO(), healthcheck.Config{
		AuthToken: credentials.AuthToken,
	})
}

func (i *DigitaloceanTeamIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials digitaloceanDescriberLocal.IntegrationCredentials
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

func (i *DigitaloceanTeamIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range digitaloceanDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *DigitaloceanTeamIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := digitaloceanDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *DigitaloceanTeamIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range digitaloceanDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
