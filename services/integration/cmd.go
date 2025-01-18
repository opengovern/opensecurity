package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/jackc/pgtype"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"github.com/opengovern/opencomply/services/integration/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	api3 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/og-util/pkg/vault"
	core "github.com/opengovern/opencomply/services/core/client"
	"github.com/opengovern/opencomply/services/integration/api"
	"github.com/opengovern/opencomply/services/integration/config"
	"github.com/opengovern/opencomply/services/integration/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	IntegrationsJsonFilePath string = "/integrations/integrations.json"
)

func Command() *cobra.Command {
	cnf := koanf.Provide("integration", config.IntegrationConfig{})

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			logger = logger.Named("integration")
			cfg := postgres.Config{
				Host:    cnf.Postgres.Host,
				Port:    cnf.Postgres.Port,
				User:    cnf.Postgres.Username,
				Passwd:  cnf.Postgres.Password,
				DB:      cnf.Postgres.DB,
				SSLMode: cnf.Postgres.SSLMode,
			}
			gorm, err := postgres.NewClient(&cfg, logger.Named("postgres"))

			cfg.DB = "integration_types"
			integrationTypesDb, err := postgres.NewClient(&cfg, logger.Named("integration_types"))
			if err != nil {
				return err
			}

			db := db.NewDatabase(gorm, integrationTypesDb)
			if err != nil {
				return err
			}

			err = db.Initialize()
			if err != nil {
				return err
			}

			mClient := core.NewCoreServiceClient(cnf.Core.BaseURL)

			_, err = mClient.VaultConfigured(&httpclient.Context{UserRole: api3.AdminRole})
			if err != nil && errors.Is(err, core.ErrConfigNotFound) {
				return err
			}

			var vaultSc vault.VaultSourceConfig
			switch cnf.Vault.Provider {
			case vault.AwsKMS:
				vaultSc, err = vault.NewKMSVaultSourceConfig(ctx, cnf.Vault.Aws, cnf.Vault.KeyId)
				if err != nil {
					logger.Error("failed to create vault source config", zap.Error(err))
					return err
				}
			case vault.AzureKeyVault:
				vaultSc, err = vault.NewAzureVaultClient(ctx, logger, cnf.Vault.Azure, cnf.Vault.KeyId)
				if err != nil {
					logger.Error("failed to create vault source config", zap.Error(err))
					return err
				}
			case vault.HashiCorpVault:
				vaultSc, err = vault.NewHashiCorpVaultClient(ctx, logger, cnf.Vault.HashiCorp, cnf.Vault.KeyId)
				if err != nil {
					logger.Error("failed to create vault source config", zap.Error(err))
					return err
				}
			}

			typeManager := integration_type.NewIntegrationTypeManager(logger, integrationTypesDb, cnf.IntegrationPlugins.MaxAutoRebootRetries, time.Duration(cnf.IntegrationPlugins.PingIntervalSeconds)*time.Second)

			err = IntegrationTypesMigration(logger, db, typeManager, IntegrationsJsonFilePath)
			if err != nil {
				logger.Error("failed to migrate integration types", zap.Error(err))
				return err
			}

			cmd.SilenceUsage = true

			steampipeConn, err := steampipe.NewSteampipeDatabase(steampipe.Option{
				Host: cnf.Steampipe.Host,
				Port: cnf.Steampipe.Port,
				User: cnf.Steampipe.Username,
				Pass: cnf.Steampipe.Password,
				Db:   cnf.Steampipe.DB,
			})
			if err != nil {
				return fmt.Errorf("new steampipe client: %w", err)
			}
			logger.Info("Connected to the steampipe database", zap.String("database", cnf.Steampipe.DB))
			kubeClient, err := NewKubeClient()
			if err != nil {
				return err
			}

			for _, name := range typeManager.GetIntegrationTypes() {
				setup, _ := db.GetIntegrationTypeSetup(name.String())
				if setup != nil {
					continue
				}
				err = db.CreateIntegrationTypeSetup(&models.IntegrationTypeSetup{
					IntegrationType: name,
					Enabled:         false,
				})
				if err != nil {
					return err
				}
				//if name == integration_type.IntegrationTypeAWSAccount || name == integration_type.IntegrationTypeGithubAccount ||
				//	name == integration_type.IntegrationTypeOpenAIIntegration {
				//	err = integrations.EnableIntegrationType(ctx, logger, kubeClient, db, name.String())
				//	if err != nil {
				//		return err
				//	}
				//}
			}

			return httpserver.RegisterAndStart(
				cmd.Context(),
				logger,
				cnf.Http.Address,
				api.New(logger, db, vaultSc, steampipeConn, kubeClient, typeManager),
			)
		},
	}

	return cmd
}

func NewKubeClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	if err := helmv2.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := v1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := kedav1alpha1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := appsv1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	kubeClient, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

type IntegrationType struct {
	ID               int64               `json:"id"`
	Name             string              `json:"name"`
	IntegrationType  string              `json:"integration_type"`
	Tier             string              `json:"tier"`
	Annotations      map[string][]string `json:"annotations"`
	Labels           map[string][]string `json:"labels"`
	ShortDescription string              `json:"short_description"`
	Description      string              `json:"Description"`
	Icon             string              `json:"Icon"`
	Availability     string              `json:"Availability"`
	SourceCode       string              `json:"SourceCode"`
	PackageURL       string              `json:"PackageURL"`
	PackageTag       string              `json:"PackageTag"`
	Enabled          bool                `json:"enabled"`
	SchemaIDs        []string            `json:"schema_ids"`
}

func IntegrationTypesMigration(logger *zap.Logger, dbm db.Database, typeManager *integration_type.IntegrationTypeManager, onboardFilePath string) error {
	content, err := os.ReadFile(onboardFilePath)
	if err != nil {
		return err
	}

	var integrationTypes []IntegrationType
	err = json.Unmarshal(content, &integrationTypes)
	if err != nil {
		return err
	}

	err = dbm.Orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.IntegrationType{}).Where("1 = 1").Unscoped().Delete(&models.IntegrationType{}).Error
		if err != nil {
			logger.Error("failed to delete integration types", zap.Error(err))
			return err
		}

		for _, obj := range integrationTypes {
			integrationType := models.IntegrationType{
				ID:               obj.ID,
				IntegrationType:  obj.IntegrationType,
				Name:             obj.Name,
				Label:            obj.Name,
				Tier:             obj.Tier,
				ShortDescription: obj.ShortDescription,
				Description:      obj.Description,
				Logo:             obj.Icon,
				Enabled:          obj.Enabled,
				PackageURL:       obj.PackageURL,
				PackageTag:       obj.PackageTag,
			}
			if _, ok := typeManager.GetIntegrationTypeMap()[typeManager.ParseType(integrationType.IntegrationType)]; ok {
				integrationType.Enabled = true
			} else {
				integrationType.Enabled = false
			}
			annotationsJsonData, err := json.Marshal(obj.Annotations)
			if err != nil {
				return err
			}
			integrationAnnotationsJsonb := pgtype.JSONB{}
			err = integrationAnnotationsJsonb.Set(annotationsJsonData)
			integrationType.Annotations = integrationAnnotationsJsonb

			labelsJsonData, err := json.Marshal(obj.Labels)
			if err != nil {
				return err
			}
			integrationLabelsJsonb := pgtype.JSONB{}
			err = integrationLabelsJsonb.Set(labelsJsonData)
			integrationType.Labels = integrationLabelsJsonb

			// logger.Info("integrationType", zap.Any("obj", obj))
			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&integrationType).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failure in integration types transaction: %w", err)
	}

	return nil
}
