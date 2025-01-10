package configs

import "github.com/opengovern/og-util/pkg/integration"
import _ "embed"

//go:embed ui-spec.json
var UISpec []byte

const (
	IntegrationTypeLinodeProject = integration.Type("linode_account") // example: aws_cloud, azure_subscription
)
