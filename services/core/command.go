package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/opengovern/og-util/pkg/httpserver"

	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
	config2 "github.com/opengovern/opensecurity/services/core/config" // Aliased standard config import
	vault2 "github.com/opengovern/opensecurity/services/core/vault"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Read environment variables into package-level vars
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

// Command remains the same
func Command() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			var cnf config2.Config
			// Read config - note this might potentially load values over env vars depending on precedence
			config.ReadFromEnv(&cnf, nil)
			return start(cmd.Context(), cnf) // Pass potentially populated cnf
		},
	}
}

// start contains the main application startup logic.
func start(ctx context.Context, cnf config2.Config) error {
	// Load config using Koanf (can potentially override env vars if set in files/flags)
	cfg := koanf.Provide("core", config2.Config{})

	// Setup Logger
	logger, err := zap.NewProduction()
	if err != nil {
		// Cannot log error here as logger failed, return directly
		return fmt.Errorf("failed to initialize zap logger: %w", err)
	}
	defer logger.Sync() // Flushes buffer, important for reliable logging
	logger.Info("Core service starting...")

	// --- Environment Variable Check ---
	logger.Debug("Checking required environment variables...")
	requiredEnvVars := map[string]string{ // Use map for easy value access
		"POSTGRESPLUGIN_HOST":     PostgresPluginHost,
		"POSTGRESPLUGIN_PORT":     PostgresPluginPort,
		"POSTGRESPLUGIN_USERNAME": PostgresPluginUsername,
		"POSTGRESPLUGIN_PASSWORD": PostgresPluginPassword,
		"SCHEDULER_BASE_URL":      SchedulerBaseUrl,
		"INTEGRATION_BASE_URL":    IntegrationBaseUrl,
		"COMPLIANCE_BASE_URL":     ComplianceBaseUrl,
		"AUTH_BASE_URL":           AuthBaseUrl,
		"COMPLIANCE_ENABLED":      ComplianceEnabled,
	}
	var missingVars []string
	for name, value := range requiredEnvVars {
		// Check if the value read at startup is empty
		if strings.TrimSpace(value) == "" {
			missingVars = append(missingVars, name)
		}
	}

	if len(missingVars) > 0 {
		errMsg := fmt.Sprintf("critical configuration error: required environment variable(s) not set: %s", strings.Join(missingVars, ", "))
		// Log before panic, although panic might prevent flushing
		logger.Error(errMsg)
		// Panic as requested
		panic(errMsg)
	}
	logger.Debug("Required environment variables check passed.")
	// --- END Environment Variable Check ---

	logger.Debug("Koanf configuration loaded", zap.Any("config", cfg)) // Be careful logging sensitive config
	var sealHandler *vault.HashiCorpVaultSealHandler                   // <<< UPDATED TYPE

	// --- Vault Initialization and Unsealing ---
	if cfg.Vault.Provider == vault.HashiCorpVault {
		logger.Info("HashiCorp Vault provider configured, initializing SealHandler...")
		// Pass the koanf loaded cfg here, as it likely contains vault specifics
		sealHandler, err := vault2.NewSealHandler(ctx, logger, cfg)
		if err != nil {
			logger.Error("Failed to create Vault seal handler", zap.Error(err))
			return fmt.Errorf("new seal handler: %w", err)
		}

		logger.Info("Starting Vault seal handler (blocks until ready or error)...")
		err = sealHandler.Start(ctx) // Start now returns an error
		if err != nil {
			// Log and return error - main() will handle exit
			logger.Error("Vault seal handler failed to start or timed out", zap.Error(err))
			return fmt.Errorf("vault seal handler start failed: %w", err)
		}
		logger.Info("Vault seal handler reported Vault is ready.")

	} else {
		logger.Info("HashiCorp Vault provider not configured, skipping Vault seal handling.")
	}

	// --- Initialize HTTP Handler ---
	logger.Info("Initializing HTTP handler...")
	// *** Pass the concrete pointer type ***
	handler, err := InitializeHttpHandler(
		cfg,
		SchedulerBaseUrl, IntegrationBaseUrl, ComplianceBaseUrl, AuthBaseUrl,
		sealHandler,
		logger,
		cnf.ElasticSearch,
		ComplianceEnabled,
	)
	if err != nil {
		logger.Error("Failed to initialize HTTP handler", zap.Error(err))
		return fmt.Errorf("init http handler: %w", err)
	}
	logger.Info("HTTP handler initialized.")

	// --- Start HTTP Server ---
	logger.Info("Registering HTTP handler and starting server...", zap.String("address", cfg.Http.Address))
	err = httpserver.RegisterAndStart(ctx, logger, cfg.Http.Address, handler)
	if err != nil {
		logger.Error("HTTP server failed", zap.Error(err))
		return fmt.Errorf("http server failed: %w", err)
	}

	logger.Info("Core service stopped.")
	return nil // Normal exit
}
