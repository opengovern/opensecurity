package integration_type

import (
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-getter"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/services/integration/models"
	"go.uber.org/zap"
	"os"
)

type GitParser struct {
	Integrations IntegrationYaml
}

type Manifest struct {
	IntegrationType integration.Type `json:"IntegrationType" yaml:"IntegrationType"`
	DescriberURL    string           `json:"DescriberURL" yaml:"DescriberURL"`
	DescriberTag    string           `json:"DescriberTag" yaml:"DescriberTag"`
}

type IntegrationPlugin struct {
	ID              int                 `json:"id" yaml:"id"`
	IntegrationType integration.Type    `json:"integration_type" yaml:"integration_type"`
	Name            string              `json:"name" yaml:"name"`
	Tier            string              `json:"tier" yaml:"tier"`
	Tags            map[string][]string `json:"tags" yaml:"tags"`
	Description     string              `json:"description" yaml:"description"`
	Icon            string              `json:"icon" yaml:"icon"`
	Availability    string              `json:"availability" yaml:"availability"`
	SourceCode      string              `json:"source_code" yaml:"source_code"`
	PackageType     string              `json:"package_type" yaml:"package_type"`
	ArtifactDetails struct {
		PackageURL string `json:"package_url" yaml:"package_url"`
		PackageTag string `json:"package_tag" yaml:"package_tag"`
	} `json:"artifact_details" yaml:"artifact_details"`
}

type IntegrationYaml struct {
	Plugins []IntegrationPlugin `json:"plugins" yaml:"plugins"`
}

func (g *GitParser) ExtractIntegrations(logger *zap.Logger) error {
	// read file from path
	f, err := os.Open(config.IntegrationTypesYamlPath)
	if err != nil {
		logger.Error("failed to open file", zap.Error(err))
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&g.Integrations); err != nil {
		logger.Error("failed to decode json", zap.Error(err))
		return fmt.Errorf("decode json: %w", err)
	}

	return nil
}

func (g *GitParser) ExtractIntegrationBinaries(logger *zap.Logger, iPlugin IntegrationPlugin) (*models.IntegrationPlugin, error) {
	baseDir := "/integration-types"

	// create tmp directory if not exists
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.Mkdir(baseDir, os.ModePerm); err != nil {
			logger.Error("failed to create tmp directory", zap.Error(err))
			return nil, fmt.Errorf("create tmp directory: %w", err)
		}
	}

	// download files from urls

	if iPlugin.ArtifactDetails.PackageURL == "" || iPlugin.ArtifactDetails.PackageTag != "" {
		return nil, nil
	}
	url := iPlugin.ArtifactDetails.PackageURL
	// remove existing files
	if err := os.RemoveAll(baseDir + "/integarion_type"); err != nil {
		logger.Error("failed to remove existing files", zap.Error(err), zap.String("url", url), zap.String("path", baseDir+"/integarion_type"))
		return nil, fmt.Errorf("remove existing files for url %s: %w", iPlugin, err)
	}

	downloader := getter.Client{
		Src:  url,
		Dst:  baseDir + "/integarion_type",
		Mode: getter.ClientModeDir,
	}
	err := downloader.Get()
	if err != nil {
		logger.Error("failed to get integration binaries", zap.Error(err), zap.String("url", url))
		return nil, fmt.Errorf("get integration binaries for url %s: %w", iPlugin, err)
	}

	//// read manifest file
	manifestFile, err := os.Open(baseDir + "/integarion_type/manifest.yaml")
	if err != nil {
		logger.Error("failed to open manifest file", zap.Error(err))
		return nil, fmt.Errorf("open manifest file: %w", err)
	}
	defer manifestFile.Close()
	var m models.Manifest
	// decode yaml
	if err := yaml.NewDecoder(manifestFile).Decode(&m); err != nil {
		logger.Error("failed to decode manifest", zap.Error(err), zap.String("url", url))
		return nil, fmt.Errorf("decode manifest for url %s: %w", iPlugin, err)
	}

	// read integration-plugin file
	integrationPlugin, err := os.ReadFile(baseDir + "/integarion_type/integration-plugin")
	if err != nil {
		logger.Error("failed to open integration-plugin file", zap.Error(err), zap.String("url", url))
		return nil, fmt.Errorf("open integration-plugin file for url %s: %w", iPlugin, err)
	}
	cloudqlPlugin, err := os.ReadFile(baseDir + "/integarion_type/cloudql-plugin")
	if err != nil {
		logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("url", url))
		return nil, fmt.Errorf("open cloudql-plugin file for url %s: %w", iPlugin.IntegrationType.String(), err)
	}

	logger.Info("done reading files", zap.String("url", url), zap.String("integrationType", iPlugin.IntegrationType.String()), zap.Int("integrationPluginSize", len(integrationPlugin)), zap.Int("cloudqlPluginSize", len(cloudqlPlugin)))

	tagsJsonData, err := json.Marshal(iPlugin.Tags)
	if err != nil {
		return nil, err
	}
	tagsJsonb := pgtype.JSONB{}
	err = tagsJsonb.Set(tagsJsonData)

	return &models.IntegrationPlugin{
		ID:                iPlugin.ID,
		PluginID:          iPlugin.IntegrationType.String(),
		IntegrationType:   iPlugin.IntegrationType,
		Name:              iPlugin.Name,
		Tier:              iPlugin.Tier,
		Description:       iPlugin.Description,
		Icon:              iPlugin.Icon,
		Availability:      iPlugin.Availability,
		SourceCode:        iPlugin.SourceCode,
		PackageType:       iPlugin.PackageType,
		InstallState:      models.IntegrationTypeInstallStateInstalled,
		OperationalStatus: models.IntegrationPluginOperationalStatusEnabled,
		URL:               url,
		DescriberURL:      m.DescriberURL,
		DescriberTag:      m.DescriberTag,
		IntegrationPlugin: integrationPlugin,
		CloudQlPlugin:     cloudqlPlugin,
		Tags:              tagsJsonb,
	}, nil
}
