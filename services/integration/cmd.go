package integration

import (
	"errors"
	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	api3 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opencomply/pkg/utils"
	core "github.com/opengovern/opencomply/services/core/client"
	"github.com/opengovern/opencomply/services/integration/api"
	"github.com/opengovern/opencomply/services/integration/config"
	"github.com/opengovern/opencomply/services/integration/db"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
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

			isAKSString := os.Getenv("ELASTICSEARCH_ISONAKS")
			var isAKSPtr *bool
			if isAKSString != "" {
				isAKS, err := strconv.ParseBool(isAKSString)
				if err != nil {
					logger.Error("failed to parse isAKS", zap.Error(err))
					return err
				}
				isAKSPtr = &isAKS
			}

			isOpenSearchString := os.Getenv("ELASTICSEARCH_ISOPENSEARCH")
			var isOpenSearchPtr *bool
			if isOpenSearchString != "" {
				isOpenSearch, err := strconv.ParseBool(isOpenSearchString)
				if err != nil {
					logger.Error("failed to parse isOpenSearch", zap.Error(err))
					return err
				}
				isOpenSearchPtr = &isOpenSearch
			}

			elastic, err := opengovernance.NewClient(opengovernance.ClientConfig{
				Addresses:     []string{os.Getenv("ELASTICSEARCH_ADDRESS")},
				Username:      utils.GetPointer(os.Getenv("ELASTICSEARCH_USERNAME")),
				Password:      utils.GetPointer(os.Getenv("ELASTICSEARCH_PASSWORD")),
				IsOpenSearch:  isOpenSearchPtr,
				IsOnAks:       isAKSPtr,
				AwsRegion:     utils.GetPointer(os.Getenv("ELASTICSEARCH_AWS_REGION")),
				AssumeRoleArn: utils.GetPointer(os.Getenv("ELASTICSEARCH_ASSUME_ROLE_ARN")),
			})
			if err != nil {
				logger.Error("failed to create es client due to", zap.Error(err))
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

			kubeClient, err := NewKubeClient()
			if err != nil {
				return err
			}

			inClusterConfig, err := rest.InClusterConfig()
			if err != nil {
				logger.Error("failed to get in cluster config", zap.Error(err))
			}
			// creates the clientset
			clientset, err := kubernetes.NewForConfig(inClusterConfig)
			if err != nil {
				logger.Error("failed to create clientset", zap.Error(err))
			}

			var metricsClient *metricsv.Clientset
			if isMetricsAPIAvailable(clientset) {
				metricsClient, err = metricsv.NewForConfig(inClusterConfig)
				if err != nil {
					logger.Error("failed to create metricsClient", zap.Error(err))
				}
			}

			typeManager := integration_type.NewIntegrationTypeManager(logger, db, integrationTypesDb, kubeClient, clientset, metricsClient, cnf.IntegrationPlugins)

			cmd.SilenceUsage = true

			steampipeOption := steampipe.Option{
				Host: cnf.Steampipe.Host,
				Port: cnf.Steampipe.Port,
				User: cnf.Steampipe.Username,
				Pass: cnf.Steampipe.Password,
				Db:   cnf.Steampipe.DB,
			}

			return httpserver.RegisterAndStart(
				cmd.Context(),
				logger,
				cnf.Http.Address,
				api.New(logger, db, vaultSc, &steampipeOption, kubeClient, typeManager, elastic),
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

func isMetricsAPIAvailable(clientset *kubernetes.Clientset) bool {
	discoveryClient := clientset.Discovery()
	apiGroups, err := discoveryClient.ServerGroups()
	if err != nil {
		return false
	}

	// Check if "metrics.k8s.io" is present
	for _, group := range apiGroups.Groups {
		if strings.Contains(group.Name, "metrics.k8s.io") {
			apiResources, err := discoveryClient.ServerResourcesForGroupVersion(group.PreferredVersion.GroupVersion)
			if err == nil {
				for _, resource := range apiResources.APIResources {
					if resource.Name == "pods" {
						return true
					}
				}
			}
		}
	}
	return false
}
