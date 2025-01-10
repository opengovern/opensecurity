package configs

import (
	_ "embed"
	"github.com/opengovern/og-util/pkg/integration"
)

//go:embed ui-spec.json
var UISpec []byte

const (
	IntegrationTypeAwsCloudAccount = integration.Type("aws_cloud_account") // example: aws_cloud, azure_subscription
)
