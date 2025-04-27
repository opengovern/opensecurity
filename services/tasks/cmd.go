package tasks

import (
	"context"
	"errors"
	"fmt"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	api3 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opensecurity/pkg/utils"
	core "github.com/opengovern/opensecurity/services/core/client"
	"github.com/opengovern/opensecurity/services/tasks/config"
	"github.com/opengovern/opensecurity/services/tasks/db"
	"github.com/opengovern/opensecurity/services/tasks/scheduler"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Command() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context())
		},
	}
}

// start runs both HTTP and GRPC server.
// GRPC server has Check method to ensure user is
// authenticated and authorized to perform an action.
// HTTP server has multiple endpoints to view and update
// the user roles.
func start(ctx context.Context) error {
	cfg := koanf.Provide("tasks", config.Config{})

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	logger = logger.Named("tasks")

	coreClient := core.NewCoreServiceClient(cfg.Core.BaseURL)

	_, err = coreClient.VaultConfigured(&httpclient.Context{UserRole: api3.AdminRole})
	if err != nil && errors.Is(err, core.ErrConfigNotFound) {
		return err
	}

	var vaultSc vault.VaultSourceConfig
	switch cfg.Vault.Provider {
	case vault.HashiCorpVault:
		vaultSc, err = vault.NewHashiCorpVaultClient(ctx, logger, cfg.Vault.HashiCorp, cfg.Vault.KeyId)
		if err != nil {
			logger.Error("failed to create vault source config", zap.Error(err))
			return err
		}
	}

	// setup postgres connection
	postgresCfg := postgres.Config{
		Host:    cfg.Postgres.Host,
		Port:    cfg.Postgres.Port,
		User:    cfg.Postgres.Username,
		Passwd:  cfg.Postgres.Password,
		DB:      cfg.Postgres.DB,
		SSLMode: cfg.Postgres.SSLMode,
	}
	orm, err := postgres.NewClient(&postgresCfg, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}

	dbm := db.Database{Orm: orm}
	fmt.Println("Connected to the postgres database: ", cfg.Postgres.DB)

	// setup postgres connection
	itPostgresConfig := postgres.Config{
		Host:    cfg.Postgres.Host,
		Port:    cfg.Postgres.Port,
		User:    cfg.Postgres.Username,
		Passwd:  cfg.Postgres.Password,
		DB:      "integration_types",
		SSLMode: cfg.Postgres.SSLMode,
	}
	itOrm, err := postgres.NewClient(&itPostgresConfig, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}

	itDbm := db.Database{Orm: itOrm}
	fmt.Println("Connected to the postgres database: ", cfg.Postgres.DB)

	err = dbm.Initialize()
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}

	kubeClient, err := NewKubeClient()
	if err != nil {
		return err
	}

	jq, err := jq.New(cfg.NATS.URL, logger)
	if err != nil {
		logger.Error("Failed to create job queue", zap.Error(err))
		return err
	}

	mainScheduler, err := scheduler.NewMainScheduler(cfg, logger, dbm, kubeClient, vaultSc, jq)
	if err != nil {
		return err
	}

	err = mainScheduler.Start(ctx)
	if err != nil {
		return err
	}

	utils.EnsureRunGoroutine(func() {
		mainScheduler.CreateTaskScheduler(ctx)
	})

	return httpserver.RegisterAndStart(ctx, logger, cfg.Http.Address, &httpRoutes{
		logger: logger,
		db:     dbm,
		itDb:   itDbm,
		jq:     jq,
		vault:  vaultSc,
	})
}

func NewKubeClient() (client.Client, error) {
	scheme := runtime.NewScheme()
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
