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
	dbPlugin, _, err := ExtractIntegrationBinaries(logger, plugin)
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

	logger.Info("integration created")

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
	}

	logger.Info("integration binaries loaded")

	operationalStatusUpdates := pgtype.JSONB{}
	operationalStatusUpdates.Set("[]")

	var supportedPlatformVersion string
	if plugin.SupportedPlatformVersions != nil {
		supportedPlatformVersion = strings.Join(plugin.SupportedPlatformVersions, ", ")
	}

	return &models.IntegrationPlugin{
			PluginID:                 plugin.IntegrationType.String(),
			IntegrationType:          plugin.IntegrationType,
			Name:                     plugin.Name,
			Description:              plugin.Metadata.Description,
			Version:                  plugin.Version,
			Icon:                     plugin.Metadata.Icon,
			PackageType:              plugin.Type,
			InstallState:             installState,
			OperationalStatus:        operationalStatus,
			OperationalStatusUpdates: operationalStatusUpdates,
			SupportedPlatformVersion: supportedPlatformVersion,
			URL:                      url,
			DescriberURL:             describerURL,
			DemoDataURL:              demoDataUrl,
			DemoDataLoaded:           false,
			DescriberTag:             describerTags,
			DiscoveryType:            discoveryType,
			Tags:                     tagsJsonb,
		}, nil,
		nil
}
