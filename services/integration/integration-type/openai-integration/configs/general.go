package configs

import "github.com/opengovern/og-util/pkg/integration"
import _ "embed"

//go:embed ui-spec.json
var UISpec []byte

const (
	IntegrationTypeOpenaiIntegration = integration.Type("openai_integration") // example: aws_cloud, azure_subscription
)
