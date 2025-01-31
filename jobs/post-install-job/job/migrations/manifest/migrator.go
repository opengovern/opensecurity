package manifest

import (
	"context"
	"fmt"
	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/jobs/post-install-job/db"
	models2 "github.com/opengovern/opencomply/services/core/db/models"
	"github.com/opengovern/opencomply/services/integration/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Migration struct {
}

func (m Migration) IsGitBased() bool {
	return true
}
func (m Migration) AttachmentFolderPath() string {
	return config.ConfigzGitPath
}

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "integration_types",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{ORM: orm}

	err = dbm.ORM.AutoMigrate(&models.IntegrationPlugin{}, &models.IntegrationPluginBinary{})
	if err != nil {
		logger.Error("failed to auto migrate integration binaries", zap.Error(err))
		return err
	}

	parser := GitParser{}
	err = parser.ExtractIntegrations(logger)
	if err != nil {
		return err
	}

	if parser.Manifest.Type != "platform_manifest" {
		return fmt.Errorf("manifest type %s is not supported, should be platform_manifest", parser.Manifest.Type)
	}

	valid, platformVersion, err := m.CheckManifestVersion(parser.Manifest.SupportedPlatformVersions)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("manifest version %v is not supported by this platform version %s",
			parser.Manifest.SupportedPlatformVersions, platformVersion)
	}

	for _, iPlugin := range parser.Manifest.Integrations.Plugins {
		plugin, pluginBinary, err := parser.ExtractIntegrationBinaries(logger, iPlugin)
		if err != nil {
			return err
		}
		if plugin == nil {
			continue
		}

		err = dbm.ORM.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "plugin_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"id", "integration_type", "name", "tier", "description", "icon",
				"availability", "source_code", "package_type", "url", "tags"}),
		}).Create(plugin).Error
		if err != nil {
			logger.Error("failed to create integration binary", zap.Error(err))
			return err
		}

		err = dbm.ORM.Clauses(clause.OnConflict{
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
	}

	err = m.GetCoreConfigs(conf, logger, parser)
	if err != nil {
		logger.Error("failed to get core configs", zap.Error(err))
		return err
	}

	return nil
}

func (m Migration) GetCoreConfigs(conf config.MigratorConfig, logger *zap.Logger, parser GitParser) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "core",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{ORM: orm}

	if parser.Manifest.Schedules.IntegrationDiscoveryFrequency != nil {
		configMetadata := models2.ConfigMetadata{
			Key:   "full_discovery_job_interval",
			Value: *parser.Manifest.Schedules.IntegrationDiscoveryFrequency,
		}
		err = dbm.ORM.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&configMetadata).Error
		if err != nil {
			return err
		}
	}
	if parser.Manifest.Schedules.ComplianceEvaluationFrequency != nil {
		configMetadata := models2.ConfigMetadata{
			Key:   "compliance_job_interval",
			Value: *parser.Manifest.Schedules.ComplianceEvaluationFrequency,
		}
		err = dbm.ORM.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).Create(&configMetadata).Error
		if err != nil {
			return err
		}
	}

	for _, param := range parser.Manifest.Parameters {
		err := dbm.ORM.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}, {Name: "control_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"value": gorm.Expr("CASE WHEN policy_parameter_values.value = '' THEN ? ELSE policy_parameter_values.value END", param.Value),
			}),
		}).Create(&models2.PolicyParameterValues{
			Key:   param.Key,
			Value: param.Value,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Migration) CheckManifestVersion(versions []string) (bool, string, error) {
	kubeClient, err := NewKubeClient()
	if err != nil {
		return false, "", err
	}

	version := ""
	var opengovernanceVersionConfig corev1.ConfigMap
	err = kubeClient.Get(context.Background(), k8sclient.ObjectKey{
		Namespace: os.Getenv("PLATFORM_NAMESPACE"),
		Name:      "platform-version",
	}, &opengovernanceVersionConfig)
	if err == nil {
		version = opengovernanceVersionConfig.Data["version"]
	} else {
		fmt.Printf("failed to load version due to %v\n", err)
	}

	for _, v := range versions {
		if compareVersions(version, v) {
			return true, version, nil
		}
	}
	return false, version, nil
}

func NewKubeClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := helmv2.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := v1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	kubeClient, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

func compareVersions(platformVersion, validVersion string) bool {
	validRegex := "^" + regexp.MustCompile(`\.`).ReplaceAllString(validVersion, `\.`) + "$"
	validRegex = regexp.MustCompile(`X`).ReplaceAllString(validRegex, `\d+`)

	match, _ := regexp.MatchString(validRegex, platformVersion)
	return match
}
