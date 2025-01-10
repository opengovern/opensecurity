package configs

import "github.com/opengovern/og-util/pkg/integration"
import _ "embed"

//go:embed ui-spec.json
var UISpec []byte

const (
	IntegrationTypeLower             = "cloudflare"                                    // example: aws, azure
	IntegrationNameCloudflareAccount = integration.Type("cloudflare_account")          // example: AWS_ACCOUNT, AZURE_SUBSCRIPTION
	OGPluginRepoURL                  = "github.com/opengovern/og-describer-cloudflare" // example: github.com/opengovern/og-describer-aws
)
