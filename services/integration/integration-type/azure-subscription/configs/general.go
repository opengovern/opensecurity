package configs

import (
	_ "embed"
	"github.com/opengovern/og-util/pkg/integration"
)

//go:embed ui-spec.json
var UISpec []byte

const (
	IntegrationTypeAzureSubscription = integration.Type("azure_subscription") // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
)
