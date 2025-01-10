package configs

import "github.com/opengovern/og-util/pkg/integration"
import _ "embed"

//go:embed ui-spec.json
var UISpec []byte

const (
	IntegrationTypeCohereaiProject = integration.Type("cohereai_project") // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
)
