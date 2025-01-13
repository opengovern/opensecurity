package linode_account

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	linodeDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/linode-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/linode-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/linode-account/healthcheck"
	"github.com/opengovern/opencomply/services/integration/models"
)

type LinodeAccountIntegration struct{}

func (i *LinodeAccountIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   linodeDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      linodeDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           linodeDescriberLocal.StreamName,
		NatsConsumerGroup:        linodeDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: linodeDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "linode",

		UISpecFileName: "linode-account.json",

		DescriberDeploymentName: linodeDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     linodeDescriberLocal.DescriberRunCommand,
	}
}

func (i *LinodeAccountIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials linodeDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	isHealthy, err := healthcheck.LinodeIntegrationHealthcheck(credentials.Token)
	return isHealthy, err
}

func (i *LinodeAccountIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials linodeDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []models.Integration
	account, err := discovery.LinodeIntegrationDiscovery(credentials.Token)
	if err != nil {
		return nil, err
	}
	labels := map[string]string{
		"State": account.State,
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
		ProviderID: account.Euuid,
		Name:       account.Email,
		Labels:     integrationLabelsJsonb,
	})

	return integrations, nil
}

func (i *LinodeAccountIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range linodeDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *LinodeAccountIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := linodeDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}

	return ""
}

func (i *LinodeAccountIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range linodeDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
