package doppler_account

import (
	"encoding/json"
	"github.com/jackc/pgtype"
	dopplerDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/doppler-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/doppler-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/doppler-account/healthcheck"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
)

type DopplerAccountIntegration struct{}

func (i *DopplerAccountIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   dopplerDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      dopplerDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           dopplerDescriberLocal.StreamName,
		NatsConsumerGroup:        dopplerDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: dopplerDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "doppler",

		UISpecFileName: "doppler-account.json",

		DescriberDeploymentName: dopplerDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     dopplerDescriberLocal.DescriberRunCommand,
	}
}

func (i *DopplerAccountIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials dopplerDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	isHealthy, err := healthcheck.DopplerIntegrationHealthcheck(healthcheck.Config{
		Token: credentials.Token,
	})
	return isHealthy, err
}

func (i *DopplerAccountIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var credentials dopplerDescriberLocal.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}
	var integrations []models.Integration
	workplace, err := discovery.DopplerIntegrationDiscovery(discovery.Config{
		Token: credentials.Token,
	})
	labels := map[string]string{
		"BillingEmail":  workplace.BillingEmail,
		"SecurityEmail": workplace.SecurityEmail,
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
		ProviderID: workplace.ID,
		Name:       workplace.Name,
		Labels:     integrationLabelsJsonb,
	})
	return integrations, nil
}

func (i *DopplerAccountIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range dopplerDescriberLocal.ResourceTypesList {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *DopplerAccountIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := dopplerDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}
	return ""
}

func (i *DopplerAccountIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range dopplerDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
