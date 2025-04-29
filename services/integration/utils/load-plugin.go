package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/platformspec"
	"github.com/opengovern/opensecurity/services/integration/models"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"os"
	"strings"
)

func ValidateAndLoadPlugin(itOrm *gorm.DB, logger *zap.Logger, data []byte) error {
	validator := platformspec.NewDefaultValidator()

	// --- Process the Specification (Full Validation including Artifacts) ---
	validatedSpecInterface, err := validator.ProcessSpecification(
		data,
		"",
		"",
		platformspec.ArtifactTypeAll,
		false,
	)
	if err != nil {
		return err
	}
	switch spec := validatedSpecInterface.(type) {
	case *platformspec.PluginSpecification:
		if spec == nil {
			return errors.New("nil plugin specification")
		}
		err = LoadPlugin(itOrm, logger, *spec)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid type for ValidateAndLoadPlugin")
	}

	return nil
}

func LoadPlugin(itOrm *gorm.DB, logger *zap.Logger, plugin platformspec.PluginSpecification) error {
	dbPlugin, pluginBinary, err := ExtractIntegrationBinaries(logger, plugin)
	if err != nil {
		return err
	}

	err = itOrm.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "plugin_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"id", "integration_type", "name", "tier", "description", "icon",
			"availability", "source_code", "package_type", "url", "demo_data_url", "discovery_type", "tags"}),
	}).Create(dbPlugin).Error
	if err != nil {
		logger.Error("failed to create integration binary", zap.Error(err))
		return err
	}

	err = itOrm.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "plugin_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"integration_plugin": gorm.Expr(
				"CASE WHEN ? <> '' THEN CAST(? AS bytea) ELSE integration_plugin_binaries.integration_plugin END",
				pluginBinary.IntegrationPlugin,
				pluginBinary.IntegrationPlugin,
			),
			"cloud_ql_plugin": gorm.Expr(
				"CASE WHEN ? <> '' THEN CAST(? AS bytea) ELSE integration_plugin_binaries.cloud_ql_plugin END",
				pluginBinary.CloudQlPlugin,
				pluginBinary.CloudQlPlugin,
			),
		}),
	}).Create(pluginBinary).Error
	if err != nil {
		logger.Error("failed to create integration binary", zap.Error(err))
		return err
	}

	return nil
}

func ExtractIntegrationBinaries(logger *zap.Logger, plugin platformspec.PluginSpecification) (*models.IntegrationPlugin, *models.IntegrationPluginBinary, error) {
	baseDir := "/integration-types"

	// create tmp directory if not exists
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		if err := os.Mkdir(baseDir, os.ModePerm); err != nil {
			logger.Error("failed to create tmp directory", zap.Error(err))
			return nil, nil, fmt.Errorf("create tmp directory: %w", err)
		}
	}

	tagsJsonData, err := json.Marshal(plugin.Tags)
	if err != nil {
		return nil, nil, err
	}
	tagsJsonb := pgtype.JSONB{}
	err = tagsJsonb.Set(tagsJsonData)

	// download files from urls
	var integrationPlugin []byte
	var cloudqlPlugin []byte
	var describerURL, describerTags, demoDataUrl, url string
	discoveryType := models.IntegrationPluginDiscoveryTypeClassic
	installState := models.IntegrationTypeInstallStateNotInstalled
	operationalStatus := models.IntegrationPluginOperationalStatusDisabled

	if err := os.RemoveAll(baseDir + "/integarion_type"); err != nil {
		logger.Error("failed to remove existing files", zap.Error(err), zap.String("path", baseDir+"/integarion_type"))
		return nil, nil, fmt.Errorf("remove existing files: %w", err)
	}

	if plugin.Components.PlatformBinary.URI != "" && plugin.Components.PlatformBinary.PathInArchive != "" {
		downloader := getter.Client{
			Src:  plugin.Components.PlatformBinary.URI,
			Dst:  baseDir + "/integarion_type",
			Mode: getter.ClientModeDir,
		}
		err = downloader.Get()
		if err != nil {
			logger.Error("failed to get integration binaries", zap.Error(err), zap.String("url", url))
			return nil, nil, fmt.Errorf("get integration binaries for url %s: %w", plugin, err)
		}

		//// read manifest file
		manifestFile, err := os.ReadFile(baseDir + "/integarion_type/manifest.yaml")
		if err != nil {
			logger.Error("failed to open manifest file", zap.Error(err))
			return nil, nil, fmt.Errorf("open manifest file: %w", err)
		}

		var m models.Manifest
		// decode yaml
		if err := yaml.Unmarshal(manifestFile, &m); err != nil {
			logger.Error("failed to decode manifest", zap.Error(err), zap.String("url", plugin.Components.PlatformBinary.URI))
			return nil, nil, fmt.Errorf("decode manifest for url %s: %w", plugin, err)
		}

		logger.Info("manifestFile", zap.String("file", string(manifestFile)), zap.Any("manifest", m))
		describerURL = m.DescriberURL
		describerTags = m.DescriberTag
		demoDataUrl = m.DemoDataURL
		if strings.ToLower(m.DiscoveryType) == "task" {
			discoveryType = models.IntegrationPluginDiscoveryTypeTask
		}
		logger.Info("platform binary path", zap.String("path", baseDir+"/integarion_type/"+plugin.Components.PlatformBinary.PathInArchive))
		integrationPlugin, err = os.ReadFile(baseDir + "/integarion_type/" + plugin.Components.PlatformBinary.PathInArchive)
		if err != nil {
			logger.Error("failed to open integration-plugin file", zap.Error(err), zap.String("url", plugin.Components.PlatformBinary.URI))
			return nil, nil, fmt.Errorf("open integration-plugin file for url %s: %w", plugin, err)
		}

		installState = models.IntegrationTypeInstallStateInstalled
		operationalStatus = models.IntegrationPluginOperationalStatusEnabled
	}
	if plugin.Components.CloudQLBinary.URI != "" && plugin.Components.CloudQLBinary.PathInArchive != "" {
		if plugin.Components.CloudQLBinary.URI != plugin.Components.PlatformBinary.URI {
			downloader := getter.Client{
				Src:  plugin.Components.CloudQLBinary.URI,
				Dst:  baseDir + "/integarion_type",
				Mode: getter.ClientModeDir,
			}
			err = downloader.Get()
			if err != nil {
				logger.Error("failed to get integration binaries", zap.Error(err), zap.String("url", url))
				return nil, nil, fmt.Errorf("get integration binaries for url %s: %w", plugin, err)
			}
		} else {
			url = plugin.Components.CloudQLBinary.URI
		}

		logger.Info("cloudql binary path", zap.String("path", baseDir+"/integarion_type/"+plugin.Components.CloudQLBinary.PathInArchive))
		cloudqlPlugin, err = os.ReadFile(baseDir + "/integarion_type/" + plugin.Components.CloudQLBinary.PathInArchive)
		if err != nil {
			logger.Error("failed to open cloudql-plugin file", zap.Error(err), zap.String("url", plugin.Components.CloudQLBinary.URI))
			return nil, nil, fmt.Errorf("open cloudql-plugin file for url %s: %w", plugin.IntegrationType.String(), err)
		}
	}

	operationalStatusUpdates := pgtype.JSONB{}
	operationalStatusUpdates.Set("[]")

	return &models.IntegrationPlugin{
			PluginID:                 plugin.Name,
			IntegrationType:          plugin.IntegrationType,
			Name:                     plugin.Name,
			Description:              plugin.Metadata.Description,
			Icon:                     plugin.Icon,
			PackageType:              plugin.Type,
			InstallState:             installState,
			OperationalStatus:        operationalStatus,
			OperationalStatusUpdates: operationalStatusUpdates,
			URL:                      url,
			DescriberURL:             describerURL,
			DemoDataURL:              demoDataUrl,
			DemoDataLoaded:           false,
			DescriberTag:             describerTags,
			DiscoveryType:            discoveryType,
			Tags:                     tagsJsonb,
		}, &models.IntegrationPluginBinary{
			PluginID:          plugin.IntegrationType.String(),
			IntegrationPlugin: integrationPlugin,
			CloudQlPlugin:     cloudqlPlugin},
		nil
}
