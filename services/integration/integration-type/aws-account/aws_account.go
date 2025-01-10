package aws_account

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgtype"
	awsDescriberLocal "github.com/opengovern/opencomply/services/integration/integration-type/aws-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/aws-account/discovery"
	"github.com/opengovern/opencomply/services/integration/integration-type/aws-account/healthcheck"
	labelsPackage "github.com/opengovern/opencomply/services/integration/integration-type/aws-account/labels"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	"github.com/opengovern/opencomply/services/integration/models"
	"golang.org/x/net/context"
	"strconv"
)

type AwsCloudAccountIntegration struct{}

func (i *AwsCloudAccountIntegration) GetConfiguration() interfaces.IntegrationConfiguration {
	return interfaces.IntegrationConfiguration{
		NatsScheduledJobsTopic:   awsDescriberLocal.JobQueueTopic,
		NatsManualJobsTopic:      awsDescriberLocal.JobQueueTopicManuals,
		NatsStreamName:           awsDescriberLocal.StreamName,
		NatsConsumerGroup:        awsDescriberLocal.ConsumerGroup,
		NatsConsumerGroupManuals: awsDescriberLocal.ConsumerGroupManuals,

		SteampipePluginName: "aws",

		UISpecFileName: "aws-cloud-account.json",

		DescriberDeploymentName: awsDescriberLocal.DescriberDeploymentName,
		DescriberRunCommand:     awsDescriberLocal.DescriberRunCommand,
	}
}

func (i *AwsCloudAccountIntegration) HealthCheck(jsonData []byte, providerId string, labels map[string]string, annotations map[string]string) (bool, error) {
	var credentials awsDescriberLocal.IntegrationCredentials
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

func (i *AwsCloudAccountIntegration) DiscoverIntegrations(jsonData []byte) ([]models.Integration, error) {
	ctx := context.Background()
	var credentials awsDescriberLocal.IntegrationCredentials
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

func (i *AwsCloudAccountIntegration) GetResourceTypesByLabels(labels map[string]string) (map[string]*interfaces.ResourceTypeConfiguration, error) {
	resourceTypes := awsDescriberLocal.ResourceTypesList
	if labels["integration/aws/organization-master"] == "true" {
		resourceTypes = append(resourceTypes, awsDescriberLocal.OrganizationMasterResourceTypesList...)
	}
	resourceTypesMap := make(map[string]*interfaces.ResourceTypeConfiguration)
	for _, resourceType := range resourceTypes {
		resourceTypesMap[resourceType] = nil
	}
	return resourceTypesMap, nil
}

func (i *AwsCloudAccountIntegration) GetResourceTypeFromTableName(tableName string) string {
	if v, ok := awsDescriberLocal.TablesToResourceTypes[tableName]; ok {
		return v
	}
	return ""
}

func (i *AwsCloudAccountIntegration) GetTablesByLabels(map[string]string) ([]string, error) {
	var tables []string
	for t, _ := range awsDescriberLocal.TablesToResourceTypes {
		tables = append(tables, t)
	}
	return tables, nil
}
