package interfaces

import "github.com/opengovern/opencomply/services/integration/models"

type IntegrationConfiguration struct {
	NatsScheduledJobsTopic   string
	NatsManualJobsTopic      string
	NatsStreamName           string
	NatsConsumerGroup        string
	NatsConsumerGroupManuals string

	SteampipePluginName string

	UISpecFileName string

	DescriberDeploymentName string
	DescriberRunCommand     string
}

type IntegrationType interface {
	GetConfiguration() IntegrationConfiguration
	GetResourceTypesByLabels(map[string]string) (map[string]*ResourceTypeConfiguration, error)
	GetTablesByLabels(map[string]string) ([]string, error)
	HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error)
	DiscoverIntegrations(jsonData []byte) ([]models.Integration, error)
	GetResourceTypeFromTableName(tableName string) string
}

// IntegrationCreator IntegrationType interface, credentials, error
type IntegrationCreator func() IntegrationType
