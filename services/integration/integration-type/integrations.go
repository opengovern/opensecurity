package integration_type

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/integration-type/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/integration-type/aws-account"
	awsConfigs "github.com/opengovern/opencomply/services/integration/integration-type/aws-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/azure-subscription"
	azureConfigs "github.com/opengovern/opencomply/services/integration/integration-type/azure-subscription/configs"
	cloudflareaccount "github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account"
	cloudflareConfigs "github.com/opengovern/opencomply/services/integration/integration-type/cloudflare-account/configs"
	cohereaiproject "github.com/opengovern/opencomply/services/integration/integration-type/cohereai-project"
	cohereaiConfigs "github.com/opengovern/opencomply/services/integration/integration-type/cohereai-project/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team"
	digitalOceanConfigs "github.com/opengovern/opencomply/services/integration/integration-type/digitalocean-team/configs"
	doppler "github.com/opengovern/opencomply/services/integration/integration-type/doppler-account"
	dopplerConfigs "github.com/opengovern/opencomply/services/integration/integration-type/doppler-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory"
	entraidConfigs "github.com/opengovern/opencomply/services/integration/integration-type/entra-id-directory/configs"
	google_workspace_account "github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account"
	googleConfig "github.com/opengovern/opencomply/services/integration/integration-type/google-workspace-account/configs"
	"github.com/opengovern/opencomply/services/integration/integration-type/interfaces"
	linodeaccount "github.com/opengovern/opencomply/services/integration/integration-type/linode-account"
	linodeConfigs "github.com/opengovern/opencomply/services/integration/integration-type/linode-account/configs"
	oci "github.com/opengovern/opencomply/services/integration/integration-type/oci-repository"
	ociConfigs "github.com/opengovern/opencomply/services/integration/integration-type/oci-repository/configs"
	openaiproject "github.com/opengovern/opencomply/services/integration/integration-type/openai-integration"
	openaiConfigs "github.com/opengovern/opencomply/services/integration/integration-type/openai-integration/configs"
	render "github.com/opengovern/opencomply/services/integration/integration-type/render-account"
	renderConfigs "github.com/opengovern/opencomply/services/integration/integration-type/render-account/configs"
	hczap "github.com/zaffka/zap-to-hclog"
)

var integrationTypes = map[integration.Type]interfaces.IntegrationType{
	awsConfigs.IntegrationTypeAwsCloudAccount:           &aws_account.Integration{},
	azureConfigs.IntegrationTypeAzureSubscription:       &azure_subscription.Integration{},
	entraidConfigs.IntegrationTypeEntraidDirectory:      &entra_id_directory.Integration{},
	digitalOceanConfigs.IntegrationTypeDigitalOceanTeam: &digitalocean_team.Integration{},
	cloudflareConfigs.IntegrationNameCloudflareAccount:  &cloudflareaccount.Integration{},
	openaiConfigs.IntegrationTypeOpenaiIntegration:      &openaiproject.Integration{},
	linodeConfigs.IntegrationTypeLinodeProject:          &linodeaccount.Integration{},
	cohereaiConfigs.IntegrationTypeCohereaiProject:      &cohereaiproject.Integration{},
	googleConfig.IntegrationTypeGoogleWorkspaceAccount:  &google_workspace_account.Integration{},
	ociConfigs.IntegrationTypeOciRepository:             &oci.Integration{},
	renderConfigs.IntegrationTypeRenderAccount:          &render.Integration{},
	dopplerConfigs.IntegrationTypeDopplerAccount:        &doppler.Integration{},
}

type IntegrationTypeManager struct {
	logger           *zap.Logger
	hcLogger         hclog.Logger
	IntegrationTypes map[integration.Type]interfaces.IntegrationType
}

func NewIntegrationTypeManager(logger *zap.Logger, integrationTypeDb *gorm.DB) *IntegrationTypeManager {
	hcLogger := hczap.Wrap(logger)

	var types []models.IntegrationTypeBinaries
	err := integrationTypeDb.Find(&types).Error
	if err != nil {
		logger.Error("failed to fetch integration types", zap.Error(err))
		return nil
	}

	// create directory for plugins if not exists
	baseDir := "/plugins"
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		err := os.Mkdir(baseDir, os.ModePerm)
		if err != nil {
			logger.Error("failed to create plugins directory", zap.Error(err))
			return nil
		}
	}

	plugins := make(map[string]string)
	for _, t := range types {
		// write the plugin to the file system
		pluginPath := filepath.Join(baseDir, t.IntegrationType.String()+".so")
		err := os.WriteFile(pluginPath, t.IntegrationPlugin, 0755)
		if err != nil {
			logger.Error("failed to write plugin to file system", zap.Error(err), zap.String("plugin", t.IntegrationType.String()))
			continue
		}
		plugins[t.IntegrationType.String()] = pluginPath
	}

	for pluginName, pluginPath := range plugins {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: interfaces.HandshakeConfig,
			Plugins:         map[string]plugin.Plugin{pluginName: &interfaces.IntegrationTypePlugin{}},
			Cmd:             exec.Command(pluginPath),
			Logger:          hcLogger,
			Managed:         true,
		})

		rpcClient, err := client.Client()
		if err != nil {
			logger.Error("failed to create plugin client", zap.Error(err), zap.String("plugin", pluginName), zap.String("path", pluginPath))
			continue
		}

		// Request the plugin
		raw, err := rpcClient.Dispense(pluginName)
		if err != nil {
			logger.Error("failed to dispense plugin", zap.Error(err), zap.String("plugin", pluginName), zap.String("path", pluginPath))
			continue
		}

		// Cast the raw interface to the appropriate interface
		itInterface, ok := raw.(interfaces.IntegrationType)
		if !ok {
			logger.Error("failed to cast plugin to integration type", zap.String("plugin", pluginName), zap.String("path", pluginPath))
			continue
		}

		integrationTypes[itInterface.GetIntegrationType()] = itInterface
	}

	return &IntegrationTypeManager{
		logger:           logger,
		hcLogger:         hcLogger,
		IntegrationTypes: integrationTypes,
	}
}

func (m *IntegrationTypeManager) GetIntegrationTypes() []integration.Type {
	types := make([]integration.Type, 0, len(m.IntegrationTypes))
	for t := range m.IntegrationTypes {
		types = append(types, t)
	}
	return types
}

func (m *IntegrationTypeManager) GetIntegrationType(t integration.Type) interfaces.IntegrationType {
	return m.IntegrationTypes[t]
}

func (m *IntegrationTypeManager) GetIntegrationTypeMap() map[integration.Type]interfaces.IntegrationType {
	return m.IntegrationTypes
}

func (m *IntegrationTypeManager) ParseType(str string) integration.Type {
	str = strings.ToLower(str)
	for t, _ := range m.IntegrationTypes {
		if str == strings.ToLower(t.String()) {
			return t
		}
	}
	return ""
}

func (m *IntegrationTypeManager) ParseTypes(str []string) []integration.Type {
	result := make([]integration.Type, 0, len(str))
	for _, s := range str {
		t := m.ParseType(s)
		if t == "" {
			continue
		}
		result = append(result, t)
	}
	return result
}

func (m *IntegrationTypeManager) UnparseTypes(types []integration.Type) []string {
	result := make([]string, 0, len(types))
	for _, t := range types {
		result = append(result, t.String())
	}
	return result
}
