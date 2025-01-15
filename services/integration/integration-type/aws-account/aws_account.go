package aws_account

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/integration-type/aws-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/aws-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/aws-account/healthcheck"
	labelsPackage "github.com/opengovern/opencomply/services/integration/integration-type/aws-account/labels"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
	"golang.org/x/net/context"
	"strconv"
)

type Integration struct{}

func (i *Integration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   configs.JobQueueTopic,
		NatsManualJobsTopic:      configs.JobQueueTopicManuals,
		NatsStreamName:           configs.StreamName,
		NatsConsumerGroup:        configs.ConsumerGroup,
		NatsConsumerGroupManuals: configs.ConsumerGroupManuals,

		SteampipePluginName: "aws",

		UISpec: configs.UISpec,

		DescriberDeploymentName: configs.DescriberDeploymentName,
		DescriberRunCommand:     configs.DescriberRunCommand,
	}
}

func (i *Integration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return false, err
	}

	return healthcheck.AWSIntegrationHealthCheck(healthcheck.AWSConfigInput{
		AccessKeyID:              credentials.AwsAccessKeyID,
		SecretAccessKey:          credentials.AwsSecretAccessKey,
		RoleNameInPrimaryAccount: credentials.RoleToAssumeInMainAccount,
		CrossAccountRoleARN:      labels["CrossAccountRoleARN"],
		ExternalID:               credentials.ExternalID,
	}, providerId)
}

func (i *Integration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	ctx := context.Background()
	var credentials configs.IntegrationCredentials
	err := json.Unmarshal(jsonData, &credentials)
	if err != nil {
		return nil, err
	}

	var integrations []models.Integration
	accounts := discovery.AWSIntegrationDiscovery(discovery.Config{
		AWSAccessKeyID:                credentials.AwsAccessKeyID,
		AWSSecretAccessKey:            credentials.AwsSecretAccessKey,
		RoleNameToAssumeInMainAccount: credentials.RoleToAssumeInMainAccount,
		CrossAccountRoleName:          credentials.CrossAccountRoleName,
		ExternalID:                    credentials.ExternalID,
	})
	for _, a := range accounts {
		if a.Details.Error != "" {
			return nil, fmt.Errorf(a.Details.Error)
		}

		isOrganizationMaster, err := labelsPackage.IsOrganizationMasterAccount(ctx, labelsPackage.AWSConfigInput{
			AccessKeyID:              credentials.AwsAccessKeyID,
			SecretAccessKey:          credentials.AwsSecretAccessKey,
			RoleNameInPrimaryAccount: credentials.RoleToAssumeInMainAccount,
			CrossAccountRoleARN:      a.Labels.CrossAccountRoleARN,
			ExternalID:               credentials.ExternalID,
		})

		labels := map[string]string{
			"RoleNameInMainAccount":               a.Labels.RoleNameInMainAccount,
			"AccountType":                         a.Labels.AccountType,
			"CrossAccountRoleARN":                 a.Labels.CrossAccountRoleARN,
			"ExternalID":                          a.Labels.ExternalID,
			"integration/aws/organization-master": strconv.FormatBool(isOrganizationMaster),
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
			ProviderID: a.AccountID,
			Name:       a.AccountName,
			Labels:     integrationLabelsJsonb,
		})
	}

	return integrations, nil
}

func (i *Integration) GetResourceTypesByLabels(labels map[string]string) (map[string]interfaces.ResourceTypeConfiguration, error) {
	resourceTypes := configs.ResourceTypesList
	if labels["integration/aws/organization-master"] == "true" {
		resourceTypes = append(resourceTypes, configs.OrganizationMasterResourceTypesList...)
	}
	resourceTypesMap := make(map[string]interfaces.ResourceTypeConfiguration)
	for _, resourceType := range resourceTypes {
		resourceTypesMap[resourceType] = interfaces.ResourceTypeConfiguration{}
	}
	return resourceTypesMap, nil
}

func (i *Integration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := configs.TablesToResourceTypes[tableName]; ok {
		return v
	}
	return ""
}

func (i *Integration) GetIntegrationType() integration.Type {
	return configs.IntegrationTypeAwsCloudAccount
}
