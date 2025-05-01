package core

import (
	"context"
	"fmt" // <<< Added missing import for placeholder
	"os"
	"strings"

	"github.com/opengovern/og-util/pkg/httpserver"

	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
	config2 "github.com/opengovern/opensecurity/services/core/config" // Aliased standard config import
	vault2 "github.com/opengovern/opensecurity/services/core/vault"   // Use alias for vault package
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

// Gets secret name from VAULT_SECRET_NAME env var or returns default
func getVaultSecretName() string {
	const defaultSecretName = "vault-unseal-keys"
	envVarName := "VAULT_SECRET_NAME" // Standard env var name
	secretName := os.Getenv(envVarName)
	trimmedName := strings.TrimSpace(secretName)

	if trimmedName == "" {
		// Using fmt here as logger might not be ready yet
		fmt.Printf("INFO: Environment variable %s not set or empty, using default Vault secret name: %s\n", envVarName, defaultSecretName)
		return defaultSecretName
	}
	fmt.Printf("INFO: Using Vault secret name from environment variable %s: %s\n", envVarName, trimmedName)
	return trimmedName
}

// Command remains the same
func Command() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			var cnf config2.Config
			config.ReadFromEnv(&cnf, nil)
			return start(cmd.Context(), cnf)
		},
	}
}

// start contains the main application startup logic.
func start(ctx context.Context, cnf config2.Config) error {
	cfg := koanf.Provide("core", config2.Config{})
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}
	defer logger.Sync()
	logger.Info("Core service starting...")

	// --- Environment Variable Check ---
	// ... (check logic remains the same) ...
	requiredEnvVars := map[string]string{ /* ... vars ... */ }
	var missingVars []string
	for name, value := range requiredEnvVars {
		if strings.TrimSpace(value) == "" {
			missingVars = append(missingVars, name)
		}
	}
	if len(missingVars) > 0 {
		errMsg := fmt.Sprintf("critical config error: required env var(s) not set: %s", strings.Join(missingVars, ", "))
		logger.Error(errMsg)
		panic(errMsg)
	}
	logger.Debug("Required environment variables check passed.")
	// --- END Environment Variable Check ---

	logger.Debug("Koanf configuration loaded", zap.Any("config", cfg))
	var sealHandler *vault2.SealHandler // <<< Use the type from core/vault package

	// --- Vault Initialization and Unsealing ---
	if cfg.Vault.Provider == vault.HashiCorpVault {
		logger.Info("HashiCorp Vault provider configured, initializing SealHandler...")

		// *** ADD variable declaration by CALLING getVaultSecretName() ***
		vaultSecretName := getVaultSecretName()
		logger.Debug("Using Vault secret name", zap.String("secretName", vaultSecretName)) // Log the name being used

		// Pass the koanf loaded cfg and the determined secret name
		// Use alias vault2 for the vault package from services/core/vault
		sealHandler, err = vault2.NewSealHandler(ctx, logger, cfg, vaultSecretName) // <<< PASS variable HERE
		if err != nil {
			logger.Error("Failed to create Vault seal handler", zap.Error(err))
			return fmt.Errorf("new seal handler: %w", err)
		}

		logger.Info("Starting Vault seal handler (blocks until ready or error)...")
		err = sealHandler.Start(ctx) // Start returns an error
		if err != nil {
			logger.Error("Vault seal handler failed to start or timed out", zap.Error(err))
			return fmt.Errorf("vault seal handler start failed: %w", err)
		}
		logger.Info("Vault seal handler reported Vault is ready.")

	} else {
		logger.Info("HashiCorp Vault provider not configured, skipping Vault seal handling.")
	}

	// --- Initialize HTTP Handler ---
	logger.Info("Initializing HTTP handler...")
	// Pass sealHandler (which might be nil if Vault disabled)
	handler, err := InitializeHttpHandler(
		cfg,
		SchedulerBaseUrl, IntegrationBaseUrl, ComplianceBaseUrl, AuthBaseUrl,
		sealHandler, // Pass the initialized (or nil) sealHandler
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
