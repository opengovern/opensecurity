package oci_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/og-util/pkg/integration/interfaces"
	"github.com/opengovern/opencomply/services/integration/integration-type/oci-repository/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/oci-repository/healthcheck"
)

type Integration struct{}

func (i *Integration) GetConfiguration() (interfaces.IntegrationConfiguration, error) {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   configs.JobQueueTopic,
		NatsManualJobsTopic:      configs.JobQueueTopicManuals,
		NatsStreamName:           configs.StreamName,
		NatsConsumerGroup:        configs.ConsumerGroup,
		NatsConsumerGroupManuals: configs.ConsumerGroupManuals,

		SteampipePluginName: "oci",

		UISpec: configs.UISpec,

		DescriberDeploymentName: configs.DescriberDeploymentName,
		DescriberRunCommand:     configs.DescriberRunCommand,
	}, nil
}

func (i *Integration) HealthCheck(jsonData []byte, _ string, _ map[string]string, _ map[string]string) (bool, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	return healthcheck.IntegrationHealthcheck(context.TODO(), healthcheck.Config{})
}

func (i *Integration) DiscoverIntegrations(jsonData []byte) ([]integration.Integration, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}

	switch credentials.GetRegistryType() {
	case configs.RegistryTypeDockerhub:
		if credentials.DockerhubCredentials == nil {
			return nil, fmt.Errorf("dockerhub credentials are required with registry type: %s", credentials.GetRegistryType())
		}
		return []integration.Integration{
			{
				ProviderID: fmt.Sprintf("dockerhub/%s", credentials.DockerhubCredentials.Owner),
				Name:       fmt.Sprintf("Dockerhub - %s", credentials.DockerhubCredentials.Owner),
			},
		}, nil
	case configs.RegistryTypeGHCR:
		if credentials.GhcrCredentials == nil {
			return nil, fmt.Errorf("ghcr credentials are required with registry type: %s", credentials.GetRegistryType())
		}
		return []integration.Integration{
			{
				ProviderID: fmt.Sprintf("ghcr/%s", credentials.GhcrCredentials.Owner),
				Name:       fmt.Sprintf("GitHub Container Registry - %s", credentials.GhcrCredentials.Owner),
			},
		}, nil
	case configs.RegistryTypeGCR:
		if credentials.GcrCredentials == nil {
			return nil, fmt.Errorf("gcr credentials are required with registry type: %s", credentials.GetRegistryType())
		}
		return []integration.Integration{
			{
				ProviderID: fmt.Sprintf("gcr/%s/%s", credentials.GcrCredentials.ProjectID, credentials.GcrCredentials.Location),
				Name:       fmt.Sprintf("Google Container Registry - %s/%s", credentials.GcrCredentials.ProjectID, credentials.GcrCredentials.Location),
			},
		}, nil
	case configs.RegistryTypeECR:
		if credentials.EcrCredentials == nil {
			return nil, fmt.Errorf("ecr credentials are required with registry type: %s", credentials.GetRegistryType())
		}
		return []integration.Integration{
			{
				ProviderID: fmt.Sprintf("ecr/%s/%s", credentials.EcrCredentials.AccountID, credentials.EcrCredentials.Region),
				Name:       fmt.Sprintf("AWS ECR (%s) - %s", credentials.EcrCredentials.Region, credentials.EcrCredentials.AccountID),
			},
		}, nil
	case configs.RegistryTypeACR:
		if credentials.AcrCredentials == nil {
			return nil, fmt.Errorf("acr credentials are required with registry type: %s", credentials.GetRegistryType())
		}
		return []integration.Integration{
			{
				ProviderID: fmt.Sprintf("acr/%s/%s", credentials.AcrCredentials.TenantID, credentials.AcrCredentials.LoginServer),
				Name:       fmt.Sprintf("Azure Container Registry - %s", credentials.AcrCredentials.LoginServer),
			},
		}, nil
	}

	return nil, fmt.Errorf("unknown registry type: %s", credentials.GetRegistryType())
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

func (i *Integration) GetIntegrationType() (integration.Type, error) {
	return configs.IntegrationTypeOciRepository, nil
}

func (i *Integration) ListAllTables() (map[string][]interfaces.CloudQLColumn, error) {
	return make(map[string][]interfaces.CloudQLColumn), nil
}

func (i *Integration) Ping() error {
	return nil
}
