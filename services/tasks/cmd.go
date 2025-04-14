package tasks

import (
	"context"
	"errors"
	"fmt"
	kedav1alpha1 "github.com/kedacore/keda/v2/apis/keda/v1alpha1"
	api3 "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/vault"
	core "github.com/opengovern/opensecurity/services/core/client"
	"github.com/opengovern/opensecurity/services/tasks/config"
	"github.com/opengovern/opensecurity/services/tasks/db"
	"github.com/opengovern/opensecurity/services/tasks/scheduler"
	"github.com/opengovern/opensecurity/services/tasks/worker"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
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

	db := db.Database{Orm: orm}
	fmt.Println("Connected to the postgres database: ", cfg.Postgres.DB)

	err = db.Initialize()
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}

	kubeClient, err := NewKubeClient()
	if err != nil {
		return err
	}

	err = setupTasks(ctx, cfg, db, kubeClient)
	if err != nil {
		return err
	}

	mainScheduler, err := scheduler.NewMainScheduler(cfg, logger, db, vaultSc)
	if err != nil {
		return err
	}

	err = mainScheduler.Start(ctx)
	if err != nil {
		return err
	}

	return httpserver.RegisterAndStart(ctx, logger, cfg.Http.Address, &httpRoutes{
		logger: logger,
		db:     db,
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

func setupTasks(ctx context.Context, cfg config.Config, db db.Database, kubeClient client.Client) error {
	tasks, err := db.GetEnabledTaskList()
	if err != nil {
		return err
	}
	for _, task := range tasks {
		currentNamespace, ok := os.LookupEnv("CURRENT_NAMESPACE")
		if !ok {
			return fmt.Errorf("current namespace lookup failed")
		}
		err = worker.CreateWorker(ctx, cfg, kubeClient, &task, currentNamespace)
		if err != nil {
			return err
		}
	}
	return nil
}
