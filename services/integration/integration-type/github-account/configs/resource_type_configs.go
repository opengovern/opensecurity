package configs

import "github.com/opengovern/opencomply/services/integration/integration-type/interfaces"

var ResourceTypeConfigs = map[string]*interfaces.ResourceTypeConfiguration{
	"Github/Container/Package": {
		Name:            "Github/Container/Package",
		IntegrationType: IntegrationTypeGithubAccount,
		Description:     "",
		Params: []interfaces.Param{
			{
				Name:        "organization",
				Description: `Please provide the organization name`,
				Required:    false,
			},
		},
	},
	"Github/Repository": {
		Name:            "Github/Repository",
		IntegrationType: IntegrationTypeGithubAccount,
		Description:     "",
		Params: []interfaces.Param{
			{
				Name:        "repository",
				Description: `Please provide the repo name (i.e. "internal-tools")`,
				Required:    false,
			},
			{
				Name:        "organization",
				Description: `Please provide the organization name`,
				Required:    false,
			},
		},
	},
	"Github/Artifact/DockerFile": {
		Name:            "Github/Artifact/DockerFile",
		IntegrationType: IntegrationTypeGithubAccount,
		Description:     "",
		Params: []interfaces.Param{
			{
				Name:        "repository",
				Description: `Please provide the repo name (i.e. "internal-tools")`,
				Required:    false,
			},
			{
				Name:        "organization",
				Description: `Please provide the organization name`,
				Required:    false,
			},
		},
	},
	"Github/Actions/WorkflowRun": {
		Name:            "Github/Actions/WorkflowRun",
		IntegrationType: IntegrationTypeGithubAccount,
		Description:     "",
		Params: []interfaces.Param{
			{
				Name:        "repository",
				Description: `Please provide the repo name (i.e. "internal-tools")`,
				Required:    false,
			},
			{
				Name:        "organization",
				Description: `Please provide the organization name`,
				Required:    false,
			},
			{
				Name:        "run_number",
				Description: `Please provide the run number`,
				Required:    false,
			},
		},
	},
	"Github/PullRequest": {
		Name:            "Github/PullRequest",
		IntegrationType: IntegrationTypeGithubAccount,
		Description:     "",
		Params: []interfaces.Param{
			{
				Name:        "organization",
				Description: `Please provide the organization name`,
				Required:    false,
			},
		},
	},
}
