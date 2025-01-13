package integration_type

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/hashicorp/go-getter"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/jobs/post-install-job/job/migrations/integration-type/models"
	"go.uber.org/zap"
	"os"
)

type GitParser struct {
	IntegrationBinaries []models.IntegrationTypeBinaries
}

type manifest struct {
	IntegrationType integration.Type `json:"IntegrationType" yaml:"IntegrationType"`
}

type ExtraIntegrations struct {
	URLs []string `json:"extraIntegrations" yaml:"extraIntegrations"`
}

func (g *GitParser) ExtractIntegrationBinaries(logger *zap.Logger) error {
	var extraIntegrations ExtraIntegrations
	// read file from path
	f, err := os.Open(config.IntegrationTypesYamlPath)
	if err != nil {
		logger.Error("failed to open file", zap.Error(err))
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(&extraIntegrations); err != nil {
		logger.Error("failed to decode json", zap.Error(err))
		return fmt.Errorf("decode json: %w", err)
	}

	baseDir := "/integration-types"

	// create tmp directory if not exists
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.Mkdir(baseDir, os.ModePerm); err != nil {
			logger.Error("failed to create tmp directory", zap.Error(err))
			return fmt.Errorf("create tmp directory: %w", err)
		}
	}

	// download files from urls
	for _, url := range extraIntegrations.URLs {
		// remove existing files
		if err := os.RemoveAll(baseDir + "/integarion_type"); err != nil {
			logger.Error("failed to remove existing files", zap.Error(err), zap.String("url", url), zap.String("path", baseDir+"/integarion_type"))
			return fmt.Errorf("remove existing files for url %s: %w", url, err)
		}

		downloader := getter.Client{
			Src:  url,
			Dst:  baseDir + "/integarion_type",
			Mode: getter.ClientModeDir,
		}
		err := downloader.Get()
		if err != nil {
			logger.Error("failed to get integration binaries", zap.Error(err), zap.String("url", url))
			return fmt.Errorf("get integration binaries for url %s: %w", url, err)
		}

		// read manifest file
		manifestFile, err := os.Open(baseDir + "/integarion_type/manifest.yaml")
		if err != nil {
			logger.Error("failed to open manifest file", zap.Error(err))
			return fmt.Errorf("open manifest file: %w", err)
		}
		defer manifestFile.Close()
		var m manifest
		// decode yaml
		if err := yaml.NewDecoder(manifestFile).Decode(&m); err != nil {
			logger.Error("failed to decode manifest", zap.Error(err), zap.String("url", url))
			return fmt.Errorf("decode manifest for url %s: %w", url, err)
		}

		// read integration-plugin file
		integrationPlugin, err := os.ReadFile(baseDir + "/integarion_type/integration-plugin")
		if err != nil {
			logger.Error("failed to open integration-plugin file", zap.Error(err), zap.String("url", url))
			return fmt.Errorf("open integration-plugin file for url %s: %w", url, err)
		}
		cloudqlPlugin, err := os.ReadFile(baseDir + "/integarion_type/cloudql-plugin")
		if err != nil {
			logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("url", url))
			return fmt.Errorf("open cloudql-plugin file for url %s: %w", url, err)
		}

		logger.Info("done reading files", zap.String("url", url), zap.String("integrationType", m.IntegrationType.String()), zap.Int("integrationPluginSize", len(integrationPlugin)), zap.Int("cloudqlPluginSize", len(cloudqlPlugin)))
		g.IntegrationBinaries = append(g.IntegrationBinaries, models.IntegrationTypeBinaries{
			IntegrationType:   m.IntegrationType,
			IntegrationPlugin: integrationPlugin,
			CloudQlPlugin:     cloudqlPlugin,
		})
	}

	return nil
}
