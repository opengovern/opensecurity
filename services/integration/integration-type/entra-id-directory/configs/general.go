package configs

import "github.com/opengovern/og-util/pkg/integration"
import _ "embed"

//go:embed ui-spec.json
var UISpec []byte

const (
	IntegrationTypeEntraidDirectory = integration.Type("entraid_directory") // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
)
