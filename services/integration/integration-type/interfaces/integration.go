package interfaces

import (
	"github.com/hashicorp/go-plugin"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/models"
	"net/rpc"
)

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "platform-integration-plugin",
	MagicCookieValue: "integration",
}

type IntegrationConfiguration struct {
	NatsScheduledJobsTopic   string
	NatsManualJobsTopic      string
	NatsStreamName           string
	NatsConsumerGroup        string
	NatsConsumerGroupManuals string

	SteampipePluginName string

	UISpec []byte

	DescriberDeploymentName string
	DescriberRunCommand     string
}

type IntegrationType interface {
	GetIntegrationType() integration.Type
	GetConfiguration() IntegrationConfiguration
	GetResourceTypesByLabels(map[string]string) (map[string]ResourceTypeConfiguration, error)
	HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error)
	DiscoverIntegrations(jsonData []byte) ([]models.Integration, error)
	GetResourceTypeFromTableName(tableName string) string
}

// IntegrationCreator IntegrationType interface, credentials, error
type IntegrationCreator func() IntegrationType

type IntegrationTypeRPC struct {
	client *rpc.Client
}

func (i *IntegrationTypeRPC) GetIntegrationType() integration.Type {
	var integrationType integration.Type
	err := i.client.Call("Plugin.GetIntegrationType", struct{}{}, &integrationType)
	if err != nil {
		panic(err)
	}
	return integrationType
}

func (i *IntegrationTypeRPC) GetConfiguration() IntegrationConfiguration {
	var configuration IntegrationConfiguration
	err := i.client.Call("Plugin.GetConfiguration", struct{}{}, &configuration)
	if err != nil {
		panic(err)
	}
	return configuration
}

func (i *IntegrationTypeRPC) GetResourceTypesByLabels(labels map[string]string) (map[string]ResourceTypeConfiguration, error) {
	var resourceTypes map[string]ResourceTypeConfiguration
	err := i.client.Call("Plugin.GetResourceTypesByLabels", labels, &resourceTypes)
	return resourceTypes, err
}

type HealthCheckRequest struct {
	JsonData    []byte
	ProviderId  string
	Labels      map[string]string
	Annotations map[string]string
}

func (i *IntegrationTypeRPC) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var result bool
	err := i.client.Call("Plugin.HealthCheck", HealthCheckRequest{
		JsonData:    jsonData,
		ProviderId:  providerId,
		Labels:      labels,
		Annotations: annotations,
	}, &result)
	return result, err
}

func (i *IntegrationTypeRPC) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	var integrations []models.Integration
	err := i.client.Call("Plugin.DiscoverIntegrations", jsonData, &integrations)
	return integrations, err
}

func (i *IntegrationTypeRPC) GetResourceTypeFromTableName(tableName string) string {
	var resourceType string
	err := i.client.Call("Plugin.GetResourceTypeFromTableName", tableName, &resourceType)
	if err != nil {
		panic(err)
	}
	return resourceType
}

type IntegrationTypeRPCServer struct {
	Impl IntegrationType
}

func (i *IntegrationTypeRPCServer) GetIntegrationType(_ struct{}, integrationType *integration.Type) error {
	*integrationType = i.Impl.GetIntegrationType()
	return nil
}

func (i *IntegrationTypeRPCServer) GetConfiguration(_ struct{}, configuration *IntegrationConfiguration) error {
	*configuration = i.Impl.GetConfiguration()
	return nil
}

func (i *IntegrationTypeRPCServer) GetResourceTypesByLabels(labels map[string]string, resourceTypes *map[string]ResourceTypeConfiguration) error {
	var err error
	*resourceTypes, err = i.Impl.GetResourceTypesByLabels(labels)
	return err
}

func (i *IntegrationTypeRPCServer) HealthCheck(request HealthCheckRequest, result *bool) error {
	var err error
	*result, err = i.Impl.HealthCheck(request.JsonData, request.ProviderId, request.Labels, request.Annotations)
	return err
}

func (i *IntegrationTypeRPCServer) DiscoverIntegrations(jsonData []byte, integrations *[]models.Integration) error {
	var err error
	*integrations, err = i.Impl.DiscoverIntegrations(jsonData)
	return err
}

func (i *IntegrationTypeRPCServer) GetResourceTypeFromTableName(tableName string, resourceType *string) error {
	*resourceType = i.Impl.GetResourceTypeFromTableName(tableName)
	return nil
}

type IntegrationTypePlugin struct {
	Impl IntegrationType
}

func (p *IntegrationTypePlugin) Server(*plugin.MuxBroker) (any, error) {
	return &IntegrationTypeRPCServer{Impl: p.Impl}, nil
}

func (IntegrationTypePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &IntegrationTypeRPC{client: c}, nil
}
