package core

import (
	"context"
	"fmt"
	"os"

	"github.com/opengovern/og-util/pkg/httpserver"

	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
	config2 "github.com/opengovern/opensecurity/services/core/config"
	vault2 "github.com/opengovern/opensecurity/services/core/vault"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	PostgresPluginHost     = os.Getenv("POSTGRESPLUGIN_HOST")
	PostgresPluginPort     = os.Getenv("POSTGRESPLUGIN_PORT")
	PostgresPluginUsername = os.Getenv("POSTGRESPLUGIN_USERNAME")
	PostgresPluginPassword = os.Getenv("POSTGRESPLUGIN_PASSWORD")
	SchedulerBaseUrl       = os.Getenv("SCHEDULER_BASE_URL")
	IntegrationBaseUrl     = os.Getenv("INTEGRATION_BASE_URL")
	ComplianceBaseUrl      = os.Getenv("COMPLIANCE_BASE_URL")
	AuthBaseUrl            = os.Getenv("AUTH_BASE_URL")
	ComplianceEnabled      = os.Getenv("COMPLIANCE_ENABLED")
)

func Command() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			var cnf config2.Config
			config.ReadFromEnv(&cnf, nil)

			return start(cmd.Context(), cnf)
		},
	}
}

func start(ctx context.Context, cnf config2.Config) error {
	cfg := koanf.Provide("core", config2.Config{})

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}
	if cfg.Vault.Provider == vault.HashiCorpVault {
		sealHandler, err := vault2.NewSealHandler(ctx, logger, cfg)
		if err != nil {
			return fmt.Errorf("new seal handler: %w", err)
		}
		// This blocks until vault is inited and unsealed
		sealHandler.Start(ctx)
	}

	handler, err := InitializeHttpHandler(
		cfg,
		SchedulerBaseUrl, IntegrationBaseUrl, ComplianceBaseUrl, AuthBaseUrl,
		logger,
		cnf.ElasticSearch,
		ComplianceEnabled,
	)
	if err != nil {
		return fmt.Errorf("init http handler: %w", err)
	}

	return httpserver.RegisterAndStart(ctx, logger, cfg.Http.Address, handler)
}
