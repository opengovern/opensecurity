// /Users/anil/workspace/opensecurity/jobs/app-init/command.go
package app_init

import (
	"errors" // Required for errors.Is potentially
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opengovern/opensecurity/jobs/app-init/configurators"
	initTypes "github.com/opengovern/opensecurity/jobs/app-init/types"
)

// Environment variable names used by the application.
const (
	// --- Auth Service (HTTP Endpoint) ---
	envAuthServiceName = "AUTH_SERVICE_NAME" // K8s service name (e.g., "auth-service")
	envAuthNamespace   = "AUTH_NAMESPACE"    // K8s namespace (e.g., "opensecurity")
	envAuthPort        = "AUTH_SERVICE_PORT" // Service port (e.g., "8251")
	envAuthHealthPath  = "AUTH_HEALTH_PATH"  // Optional: Health check path (defaults to "/health")

	// --- PostgreSQL Database ---
	envPgHost     = "PGHOST"     // DB hostname/IP (e.g., "postgres-primary")
	envPgPort     = "PGPORT"     // DB port (e.g., "5432")
	envPgDatabase = "PGDATABASE" // Target DB name (e.g., "auth")
	envPgUser     = "PGUSER"     // Optional: DB user (defaults to "postgres")
	envPgPassword = "PGPASSWORD" // Required: DB password
	envPgSslMode  = "PGSSLMODE"  // Required: DB SSL mode (e.g., "require", "disable")

	// --- Dex gRPC & Client Config ---
	envDexGrpcAddress               = "DEX_GRPC_ADDR"                    // Required: Dex gRPC endpoint (e.g., "dex-server:5557")
	envDexPublicClientRedirectUris  = "DEX_PUBLIC_CLIENT_REDIRECT_URIS"  // Required: Comma-separated URIs
	envDexPrivateClientRedirectUris = "DEX_PRIVATE_CLIENT_REDIRECT_URIS" // Required: Comma-separated URIs
	envDexPublicClientID            = "DEX_PUBLIC_CLIENT_ID"             // Required: ID for the public client (e.g., "public-client")
	envDexPrivateClientID           = "DEX_PRIVATE_CLIENT_ID"            // Required: ID for the private client (e.g., "private-client")
	envDexPrivateClientSecret       = "DEX_PRIVATE_CLIENT_SECRET"        // Required: Secret for the private client
	// --- New: Dex HTTP Health Check ---
	envDexHTTPHealthURL = "DEX_HTTP_HEALTH_URL" // Required: Full URL for Dex HTTP health check (e.g., "http://dex-service:5558/healthz/ready")

	// --- Default Admin/First User ---
	envDefaultUserEmail    = "DEFAULT_DEX_USER_EMAIL"    // Required: Email for the default user
	envDefaultUserName     = "DEFAULT_DEX_USER_NAME"     // Required: Username for the default user
	envDefaultUserPassword = "DEFAULT_DEX_USER_PASSWORD" // Required: Password for the default user
)

// config holds all the configuration values read from the environment.
type config struct {
	authServiceName        string
	authNamespace          string
	authPort               string
	authHealthPath         string
	pgHost                 string
	pgPort                 string
	pgDatabase             string
	pgUser                 string
	pgPassword             string
	pgSslMode              string
	dexGrpcAddr            string
	dexPublicUris          string
	dexPrivateUris         string
	dexPublicClientID      string
	dexPrivateClientID     string
	dexPrivateClientSecret string
	dexHTTPHealthURL       string // New field
	defaultUserEmail       string
	defaultUserName        string
	defaultUserPassword    string
}

// readConfig reads all environment variables, applies defaults, validates, and returns a config struct or an error.
func readConfig() (*config, error) {
	cfg := &config{
		// Read values
		authServiceName:        os.Getenv(envAuthServiceName),
		authNamespace:          os.Getenv(envAuthNamespace),
		authPort:               os.Getenv(envAuthPort),
		authHealthPath:         os.Getenv(envAuthHealthPath),
		pgHost:                 os.Getenv(envPgHost),
		pgPort:                 os.Getenv(envPgPort),
		pgDatabase:             os.Getenv(envPgDatabase),
		pgUser:                 os.Getenv(envPgUser),
		pgPassword:             os.Getenv(envPgPassword),
		pgSslMode:              os.Getenv(envPgSslMode),
		dexGrpcAddr:            os.Getenv(envDexGrpcAddress),
		dexPublicUris:          os.Getenv(envDexPublicClientRedirectUris),
		dexPrivateUris:         os.Getenv(envDexPrivateClientRedirectUris),
		dexPublicClientID:      os.Getenv(envDexPublicClientID),
		dexPrivateClientID:     os.Getenv(envDexPrivateClientID),
		dexPrivateClientSecret: os.Getenv(envDexPrivateClientSecret),
		dexHTTPHealthURL:       os.Getenv(envDexHTTPHealthURL), // Read new var
		defaultUserEmail:       os.Getenv(envDefaultUserEmail),
		defaultUserName:        os.Getenv(envDefaultUserName),
		defaultUserPassword:    os.Getenv(envDefaultUserPassword),
	}

	// Apply defaults
	if cfg.pgUser == "" {
		cfg.pgUser = "postgres"
		log.Printf("INFO: Environment variable %s not set, defaulting to '%s'", envPgUser, cfg.pgUser)
	}
	if cfg.authHealthPath == "" {
		cfg.authHealthPath = "/health"
		log.Printf("INFO: Environment variable %s not set, using default path: %s", envAuthHealthPath, cfg.authHealthPath)
	}

	// Validate required fields
	var missingVars []string
	validate := func(val, name string) {
		if val == "" {
			missingVars = append(missingVars, name)
		}
	}

	validate(cfg.authServiceName, envAuthServiceName)
	validate(cfg.authNamespace, envAuthNamespace)
	validate(cfg.authPort, envAuthPort)
	validate(cfg.pgHost, envPgHost)
	validate(cfg.pgPort, envPgPort)
	validate(cfg.pgDatabase, envPgDatabase)
	validate(cfg.pgPassword, envPgPassword)
	validate(cfg.pgSslMode, envPgSslMode)
	validate(cfg.dexGrpcAddr, envDexGrpcAddress)
	validate(cfg.dexPublicUris, envDexPublicClientRedirectUris)
	validate(cfg.dexPrivateUris, envDexPrivateClientRedirectUris)
	validate(cfg.dexPublicClientID, envDexPublicClientID)
	validate(cfg.dexPrivateClientID, envDexPrivateClientID)
	validate(cfg.dexPrivateClientSecret, envDexPrivateClientSecret)
	validate(cfg.dexHTTPHealthURL, envDexHTTPHealthURL) // Validate new var
	validate(cfg.defaultUserEmail, envDefaultUserEmail)
	validate(cfg.defaultUserName, envDefaultUserName)
	validate(cfg.defaultUserPassword, envDefaultUserPassword)

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variable(s): %s", strings.Join(missingVars, ", "))
	}

	// Optional: Add URL format validation for dexHTTPHealthURL
	if _, err := url.ParseRequestURI(cfg.dexHTTPHealthURL); err != nil {
		return nil, fmt.Errorf("invalid URL format for %s ('%s'): %w", envDexHTTPHealthURL, cfg.dexHTTPHealthURL, err)
	}

	return cfg, nil
}

// buildComponentList reads config, constructs URLs/config, and creates the single AuthComponent.
func buildComponentList() ([]initTypes.InitializableComponent, error) {
	log.Println("INFO: Reading and validating configuration...")
	cfg, err := readConfig()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %w", err) // Wrap error
	}

	// --- Construct Auth URL ---
	healthPath := cfg.authHealthPath
	if !strings.HasPrefix(healthPath, "/") {
		healthPath = "/" + healthPath
	}
	authURL := fmt.Sprintf("http://%s.%s.svc.cluster.local:%s%s",
		cfg.authServiceName, cfg.authNamespace, cfg.authPort, healthPath)

	// --- Log Loaded Configuration (Masking secrets) ---
	log.Println("INFO: Configuration loaded successfully:")
	log.Printf("  Auth Service URL: %s", authURL)
	log.Printf("  PostgreSQL: Host=%s, Port=%s, DB=%s, User=%s, SSLMode=%s, Password=[HIDDEN]",
		cfg.pgHost, cfg.pgPort, cfg.pgDatabase, cfg.pgUser, cfg.pgSslMode)
	log.Printf("  Dex: gRPC=%s, HTTPHealth=%s, PublicClientID=%s, PrivateClientID=%s, Secret=[HIDDEN]", // Added HTTPHealth
		cfg.dexGrpcAddr, cfg.dexHTTPHealthURL, cfg.dexPublicClientID, cfg.dexPrivateClientID)
	log.Printf("  Dex Public Redirects: %s", cfg.dexPublicUris)
	log.Printf("  Dex Private Redirects: %s", cfg.dexPrivateUris)
	log.Printf("  Default User: Email=%s, Username=%s, Password=[HIDDEN]",
		cfg.defaultUserEmail, cfg.defaultUserName)

	// --- Create the Single Component ---
	log.Println("INFO: Creating initialization component...")
	authComp, err := configurators.NewAuthComponent(
		// Pass values from the config struct
		authURL,
		cfg.pgHost, cfg.pgPort, cfg.pgUser, cfg.pgPassword, cfg.pgDatabase, cfg.pgSslMode,
		cfg.dexGrpcAddr,
		cfg.dexPublicClientID,
		cfg.dexPublicUris,
		cfg.dexPrivateClientID,
		cfg.dexPrivateUris,
		cfg.dexPrivateClientSecret,
		cfg.dexHTTPHealthURL, // Pass new value
		cfg.defaultUserEmail, cfg.defaultUserName, cfg.defaultUserPassword,
	)
	if err != nil {
		// Error during component construction (e.g., Dex connection failure)
		return nil, fmt.Errorf("failed to create auth component: %w", err)
	}
	log.Println("INFO: Initialization component created.")

	// --- Define Initialization Order ---
	components := []initTypes.InitializableComponent{
		authComp,
	}

	log.Printf("INFO: Prepared %d component(s) for initialization.", len(components))
	return components, nil
}

// Command creates the cobra command for the app-init job.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app-init",
		Short: "Checks/Initializes prerequisites for the Auth Service.",
		Long: `This job performs essential setup tasks before the Auth Service can run reliably.
It checks external dependencies (Postgres, Auth Service HTTP, Dex HTTP/gRPC),
ensures the database exists, configures Dex OAuth clients,
and sets up the initial administrative user if necessary.`, // Updated help text
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.New(os.Stdout, "APP-INIT: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
			logger.Println("INFO: Starting App Init Job...")
			ctx := cmd.Context()

			componentsToRun, err := buildComponentList()
			if err != nil {
				logger.Printf("ERROR: Failed to build component list: %v", err)
				return err
			}
			if len(componentsToRun) == 0 {
				logger.Println("WARN: No components were configured to run. Exiting.")
				return errors.New("no components configured for initialization")
			}

			runner := NewRunner(logger)
			err = runner.Run(ctx, componentsToRun)
			if err != nil {
				logger.Printf("ERROR: Component initialization sequence failed: %v", err)
				return err
			}

			logger.Println("INFO: App Init Job finished successfully.")
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}
