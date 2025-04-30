package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	dexApi "github.com/dexidp/dex/api/v2"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/services/auth/authcache"
	"github.com/opengovern/opensecurity/services/auth/db"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// --- Constants ---
const (
	dbConnectRetryInterval = 5 * time.Second      // How often to retry DB connection
	dbConnectRetryDuration = 180 * time.Second    // Total duration to keep retrying
	skipDBRetryEnvVar      = "AUTH_SKIP_DB_RETRY" // Env var to disable retries (e.g., set to "true")
	devModeEnvVar          = "AUTH_DEV_MODE"      // Env var for development logging (e.g., set to "true")
	authCacheDefaultTTL    = 10 * time.Minute     // <-- Default TTL for user auth cache
)

// --- End Constants ---

var (
	// Read env vars - values will be validated in start()
	dexPublicClientRedirectUris  = os.Getenv("DEX_PUBLIC_CLIENT_REDIRECT_URIS")
	dexPrivateClientRedirectUris = os.Getenv("DEX_PRIVATE_CLIENT_REDIRECT_URIS")
	dexAuthDomain                = os.Getenv("DEX_AUTH_DOMAIN")
	dexAuthPublicClientID        = os.Getenv("DEX_AUTH_PUBLIC_CLIENT_ID")
	dexGrpcAddress               = os.Getenv("DEX_GRPC_ADDR")
	httpServerAddress            = os.Getenv("HTTP_ADDRESS")
	platformHost                 = os.Getenv("PLATFORM_HOST")
	platformKeyEnabledStr        = os.Getenv("PLATFORM_KEY_ENABLED")
	platformPublicKeyStr         = os.Getenv("PLATFORM_PUBLIC_KEY")
	platformPrivateKeyStr        = os.Getenv("PLATFORM_PRIVATE_KEY")
)

type ServerConfig struct {
	PostgreSQL config2.Postgres
}

// Command defines the root command for the auth service application.
func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "auth-service", // Set a use command name
		Short: "OpenSecurity Authentication Service",
		Long:  `Runs the OpenSecurity Authentication Service, handling user authentication and authorization.`,
		// Define the main execution logic for the command
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize logger early for startup messages
			var logger *zap.Logger
			var err error
			// Check dev mode env var to determine logger type
			if strings.ToLower(os.Getenv(devModeEnvVar)) == "true" {
				logger, err = zap.NewDevelopment()          // Human-readable logs for dev
				fmt.Println("Development logging enabled.") // Simple stdout feedback
			} else {
				logger, err = zap.NewProduction() // JSON logs for prod
			}
			if err != nil {
				// Use fmt if logger fails
				fmt.Fprintf(os.Stderr, "Error: Failed to create logger: %v\n", err)
				return fmt.Errorf("failed to create logger: %w", err)
			}
			defer logger.Sync() // Ensure logs are flushed

			// Run the main start logic from the auth package, passing the context and logger
			logger.Info("Executing auth service start logic...")
			runErr := start(cmd.Context(), logger) // Call the exported Start function
			if runErr != nil {
				// Log the final error before exiting if start fails
				// Avoid logging if logger.Fatal was already called in start
				if runErr != nil && !strings.Contains(runErr.Error(), "database connection/initialization timed out") {
					logger.Error("Service startup failed", zap.Error(runErr))
				}
			}
			return runErr // Return the error from start
		},
	}
}

// Start runs the main application logic after initial setup.
// It now takes the logger as an argument. Exported for main.go.
func start(ctx context.Context, logger *zap.Logger) error {
	logger = logger.Named("auth_startup")
	logger.Info("Starting Auth Service...")

	// --- Configuration Loading and Validation ---
	var conf ServerConfig
	// Load PostgreSQL specific config (e.g., host, port, user, dbname, sslmode)
	config2.ReadFromEnv(&conf, nil)
	logger.Info("PostgreSQL configuration loaded from environment")

	// Define environment variables that are absolutely critical for startup.
	// The service will fatally exit if any of these are missing.
	fatalEnvVars := map[string]string{
		"DEX_PUBLIC_CLIENT_REDIRECT_URIS":  dexPublicClientRedirectUris,
		"DEX_PRIVATE_CLIENT_REDIRECT_URIS": dexPrivateClientRedirectUris,
		"DEX_AUTH_DOMAIN":                  dexAuthDomain,
		"DEX_AUTH_PUBLIC_CLIENT_ID":        dexAuthPublicClientID,
		"DEX_GRPC_ADDR":                    dexGrpcAddress,
		"PLATFORM_HOST":                    platformHost,
	}
	validationErrorMessages := []string{} // Collects messages for the final fatal log summary
	// Check each mandatory variable
	for name, value := range fatalEnvVars {
		if value == "" {
			// *** ADDED DEBUG LOGGING HERE ***
			// Log the specific error including the (empty) value for immediate debugging context.
			logger.Error("Required environment variable is missing or empty",
				zap.String("variable", name),
				zap.String("value", value), // Explicitly log the empty value found
			)
			// Add a user-friendly message to the list for the final fatal error summary.
			validationErrorMessages = append(validationErrorMessages, fmt.Sprintf("%s is required", name))
		}
	}

	// If any critical variables were missing, log details and exit fatally.
	if len(validationErrorMessages) > 0 {
		errMsg := "Missing critical configuration. See preceding errors for details."
		// Log Fatal with a summary message. The individual errors with values were already logged above.
		logger.Fatal(errMsg,
			zap.Strings("missing_variables_summary", validationErrorMessages),
		)
		// Fatal should exit, but return an error in case it doesn't (e.g., in tests)
		return errors.New(errMsg + ": " + strings.Join(validationErrorMessages, ", "))
	}
	logger.Info("Successfully validated critical environment variables.")

	// --- Optional/Format Validation for other variables ---
	// Check non-fatal but potentially problematic variables or formats.
	formatValidationErrors := []string{}
	// Validate DEX_AUTH_DOMAIN format (must start with http:// or https://)
	if dexAuthDomain != "" && !strings.HasPrefix(dexAuthDomain, "http://") && !strings.HasPrefix(dexAuthDomain, "https://") {
		// Note: dexAuthDomain emptiness was already checked by fatalEnvVars, this focuses on format.
		errDetail := fmt.Sprintf("DEX_AUTH_DOMAIN ('%s') must start with http:// or https://", dexAuthDomain)
		formatValidationErrors = append(formatValidationErrors, errDetail)
		logger.Warn("Configuration format issue", zap.String("detail", errDetail)) // Log as warning/info
	}
	// Warn if HTTP_ADDRESS is not set, as the service might be unreachable.
	if httpServerAddress == "" {
		logger.Warn("HTTP_ADDRESS environment variable is not set, service might not be reachable externally", zap.String("variable", "HTTP_ADDRESS"))
		// Consider setting a default like ":8080" if appropriate, but for now just warn.
	}
	// Validate PG details (even if not fatal if missing, they are needed for DB connection)
	if conf.PostgreSQL.Host == "" {
		formatValidationErrors = append(formatValidationErrors, "POSTGRES_HOST is required for database connection")
	}
	if conf.PostgreSQL.Port == "" {
		formatValidationErrors = append(formatValidationErrors, "POSTGRES_PORT is required for database connection")
	} else if _, convErr := strconv.Atoi(conf.PostgreSQL.Port); convErr != nil {
		formatValidationErrors = append(formatValidationErrors, fmt.Sprintf("POSTGRES_PORT ('%s') is not a valid integer", conf.PostgreSQL.Port))
	}
	if conf.PostgreSQL.Username == "" {
		formatValidationErrors = append(formatValidationErrors, "POSTGRES_USER is required for database connection")
	}
	if conf.PostgreSQL.DB == "" {
		formatValidationErrors = append(formatValidationErrors, "POSTGRES_DB is required for database connection")
	}

	// If format errors exist for non-fatal vars, log Error and exit gracefully (not Fatal)
	if len(formatValidationErrors) > 0 {
		errMsg := "Invalid configuration format"
		logger.Error(errMsg, zap.Strings("details", formatValidationErrors))
		return errors.New(errMsg + ": " + strings.Join(formatValidationErrors, ", "))
	}

	// Log loaded configuration
	logger.Info("Loaded configuration",
		zap.String("dexAuthDomain", dexAuthDomain),
		zap.String("dexAuthPublicClientID", dexAuthPublicClientID),
		zap.String("dexGrpcAddress", dexGrpcAddress),
		zap.String("httpServerAddress", httpServerAddress), // Log even if optional
		zap.String("platformHost", platformHost),
		zap.String("platformKeyEnabled", platformKeyEnabledStr), // Log the string value
		zap.String("postgresHost", conf.PostgreSQL.Host),
		zap.String("postgresPort", conf.PostgreSQL.Port),
		zap.String("postgresUser", conf.PostgreSQL.Username),
		zap.String("postgresDB", conf.PostgreSQL.DB),
		zap.String("postgresSSLMode", conf.PostgreSQL.SSLMode),
	)
	// --- End Configuration Loading and Validation ---

	// --- Initialize Dex Verifier ---
	logger.Info("Initializing Dex OIDC Verifier...")
	dexVerifier, err := newDexOidcVerifier(ctx, dexAuthDomain, dexAuthPublicClientID)
	if err != nil {
		logger.Error("Failed to initialize Dex OIDC Verifier",
			zap.String("domain", dexAuthDomain),
			zap.String("clientId", dexAuthPublicClientID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to initialize Dex OIDC Verifier: %w", err)
	}
	logger.Info("Dex OIDC Verifier initialized successfully.")
	// --- End Initialize Dex Verifier ---

	// --- Database Connection and Initialization with Retry ---
	var orm *gorm.DB
	var adb db.Database
	var lastDbErr error

	skipRetryVal := strings.ToLower(os.Getenv(skipDBRetryEnvVar))
	skipRetry := skipRetryVal == "true" || skipRetryVal == "1"
	dbReady := false
	startTime := time.Now()

	cfg := postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      conf.PostgreSQL.DB,
		SSLMode: conf.PostgreSQL.SSLMode,
	}

	if skipRetry {
		logger.Info("Attempting database connection (retries skipped)", zap.String("envVar", skipDBRetryEnvVar))
		orm, err = postgres.NewClient(&cfg, logger)
		if err != nil {
			logger.Error("Database connection failed (retries skipped)", zap.Error(err))
			return fmt.Errorf("database connection failed: %w", err)
		}
		adb = db.Database{Orm: orm, Logger: logger.Named("database")}
		logger.Info("Attempting database initialization (retries skipped)...")
		err = adb.Initialize()
		if err != nil {
			logger.Error("Database initialization failed (retries skipped)", zap.Error(err))
			return fmt.Errorf("database initialization failed: %w", err)
		}
		dbReady = true
		logger.Info("Database connected and initialized successfully (retries skipped).")
	} else {
		logger.Info("Attempting to connect and initialize database with retry",
			zap.Duration("timeout", dbConnectRetryDuration),
			zap.Duration("interval", dbConnectRetryInterval))

		for time.Since(startTime) < dbConnectRetryDuration {
			logger.Info("Attempting database connection...")
			orm, err = postgres.NewClient(&cfg, logger)
			if err != nil {
				lastDbErr = err
				logger.Warn("Database connection attempt failed, retrying...", zap.Error(err), zap.Duration("retryIn", dbConnectRetryInterval))
				time.Sleep(dbConnectRetryInterval)
				continue
			}

			logger.Info("Database connection successful. Attempting initialization...")
			adb = db.Database{Orm: orm, Logger: logger.Named("database")}
			err = adb.Initialize()
			if err != nil {
				lastDbErr = err
				logger.Warn("Database initialization attempt failed, retrying...", zap.Error(err), zap.Duration("retryIn", dbConnectRetryInterval))
				sqlDB, dbErr := orm.DB()
				if dbErr == nil && sqlDB != nil {
					sqlDB.Close()
				}
				orm = nil
				time.Sleep(dbConnectRetryInterval)
				continue
			}

			dbReady = true
			logger.Info("Database connected and initialized successfully.")
			break
		}

		if !dbReady {
			logger.Fatal("Failed to connect to and initialize database after timeout",
				zap.Duration("duration", dbConnectRetryDuration),
				zap.NamedError("lastError", lastDbErr),
			)
			return fmt.Errorf("database connection/initialization timed out after %v: %w", dbConnectRetryDuration, lastDbErr)
		}
	}
	// --- End Database Connection and Initialization ---

	// --- Initialize Auth Cache ---
	logger.Info("Initializing Auth Cache service...")
	authCacheSvc, err := authcache.NewAuthCacheService(logger) // Uses DefaultTTL from authcache pkg
	if err != nil {
		logger.Error("Failed to initialize Auth Cache service", zap.Error(err))
		return fmt.Errorf("failed to initialize auth cache: %w", err)
	}
	defer func() {
		logger.Info("Attempting to close Auth Cache service...")
		if closeErr := authCacheSvc.Close(); closeErr != nil {
			logger.Error("Error closing auth cache", zap.Error(closeErr))
		} else {
			logger.Info("Auth Cache service closed successfully.")
		}
	}()
	logger.Info("Auth Cache service initialized successfully.")
	// --- End Initialize Auth Cache ---

	// --- Platform Key Handling ---
	logger.Info("Setting up platform keys...")
	var platformKeyEnabled bool
	var keyErr error

	// Default platformKeyEnabled to false if string is empty or invalid
	platformKeyEnabled, keyErr = strconv.ParseBool(platformKeyEnabledStr)
	if keyErr != nil {
		// Log as Info/Warn since default is false, not a fatal error if missing/invalid
		logger.Info("PLATFORM_KEY_ENABLED not set or invalid, defaulting to false",
			zap.String("value", platformKeyEnabledStr),
			zap.Error(keyErr))
		platformKeyEnabled = false
	}

	var platformPublicKey *rsa.PublicKey
	var platformPrivateKey *rsa.PrivateKey

	if platformKeyEnabled {
		logger.Info("Platform key loading from environment variables enabled")
		// Public Key loading (ensure PLATFORM_PUBLIC_KEY is set if enabled)
		if platformPublicKeyStr == "" {
			logger.Fatal("PLATFORM_KEY_ENABLED is true, but PLATFORM_PUBLIC_KEY is not set")
			return errors.New("PLATFORM_PUBLIC_KEY is required when PLATFORM_KEY_ENABLED is true")
		}
		bPub, err := base64.StdEncoding.DecodeString(platformPublicKeyStr)
		if err != nil {
			logger.Error("Failed to decode base64 public key from env", zap.Error(err))
			return fmt.Errorf("platform public key decode error: %w", err)
		}
		blockPub, _ := pem.Decode(bPub)
		if blockPub == nil || !strings.Contains(blockPub.Type, "PUBLIC KEY") {
			logger.Error("Failed to decode PEM block from public key or invalid type", zap.String("type", blockPub.Type))
			return fmt.Errorf("invalid platform public key PEM block or type")
		}
		pubInterface, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
		if err != nil {
			logger.Error("Failed to parse PKIX public key", zap.Error(err))
			return fmt.Errorf("failed to parse platform public key: %w", err)
		}
		var ok bool
		platformPublicKey, ok = pubInterface.(*rsa.PublicKey)
		if !ok {
			logger.Error("Platform public key is not an RSA public key")
			return fmt.Errorf("platform public key is not RSA")
		}

		// Private Key loading (ensure PLATFORM_PRIVATE_KEY is set if enabled)
		if platformPrivateKeyStr == "" {
			logger.Fatal("PLATFORM_KEY_ENABLED is true, but PLATFORM_PRIVATE_KEY is not set")
			return errors.New("PLATFORM_PRIVATE_KEY is required when PLATFORM_KEY_ENABLED is true")
		}
		bPriv, err := base64.StdEncoding.DecodeString(platformPrivateKeyStr)
		if err != nil {
			logger.Error("Failed to decode base64 private key from env", zap.Error(err))
			return fmt.Errorf("platform private key decode error: %w", err)
		}
		blockPriv, _ := pem.Decode(bPriv)
		if blockPriv == nil || !strings.Contains(blockPriv.Type, "PRIVATE KEY") {
			logger.Error("Failed to decode PEM block from private key or invalid type", zap.String("type", blockPriv.Type))
			return fmt.Errorf("invalid platform private key PEM block or type")
		}
		privInterface, err := x509.ParsePKCS8PrivateKey(blockPriv.Bytes)
		if err != nil {
			logger.Warn("Failed to parse PKCS8 private key, attempting PKCS1", zap.Error(err))
			privInterface, err = x509.ParsePKCS1PrivateKey(blockPriv.Bytes)
			if err != nil {
				logger.Error("Failed to parse private key as PKCS8 or PKCS1", zap.Error(err))
				return fmt.Errorf("failed to parse platform private key: %w", err)
			}
		}
		platformPrivateKey, ok = privInterface.(*rsa.PrivateKey)
		if !ok {
			logger.Error("Platform private key is not an RSA private key")
			return fmt.Errorf("platform private key is not RSA")
		}
		logger.Info("Successfully loaded platform keys from environment variables.")

	} else {
		logger.Info("Platform key loading from environment disabled, checking database...")
		keyPair, err := adb.GetKeyPair()
		if err != nil {
			logger.Error("Failed to get key pair from database", zap.Error(err))
			return fmt.Errorf("failed to get platform key pair from database: %w", err)
		}

		if len(keyPair) < 2 {
			if len(keyPair) == 1 {
				logger.Warn("Found only one platform key in database, generating new pair.")
			}
			logger.Info("Generating new platform RSA key pair...")
			platformPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				logger.Error("Error generating platform RSA key", zap.Error(err))
				return fmt.Errorf("error generating platform RSA key: %w", err)
			}
			platformPublicKey = &platformPrivateKey.PublicKey

			pubBytes, err := x509.MarshalPKIXPublicKey(platformPublicKey)
			if err != nil {
				logger.Error("Failed to marshal generated public key", zap.Error(err))
				return fmt.Errorf("failed to marshal generated public key: %w", err)
			}
			pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
			pubStr := base64.StdEncoding.EncodeToString(pubPEM)
			if err = adb.AddConfiguration(&db.Configuration{Key: "public_key", Value: pubStr}); err != nil {
				logger.Error("Failed to save generated public key to database", zap.Error(err))
				return fmt.Errorf("failed to save generated public key: %w", err)
			}

			privBytes, err := x509.MarshalPKCS8PrivateKey(platformPrivateKey)
			if err != nil {
				logger.Error("Failed to marshal generated private key", zap.Error(err))
				return fmt.Errorf("failed to marshal generated private key: %w", err)
			}
			privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
			privStr := base64.StdEncoding.EncodeToString(privPEM)
			if err = adb.AddConfiguration(&db.Configuration{Key: "private_key", Value: privStr}); err != nil {
				logger.Error("Failed to save generated private key to database", zap.Error(err))
				return fmt.Errorf("failed to save generated private key: %w", err)
			}
			logger.Info("Successfully generated and saved new platform key pair to database.")

		} else {
			logger.Info("Loading platform keys from database...")
			foundPub, foundPriv := false, false
			for _, k := range keyPair {
				keyBytes, err := base64.StdEncoding.DecodeString(k.Value)
				if err != nil {
					logger.Error("Failed to decode key from database", zap.String("key_name", k.Key), zap.Error(err))
					return fmt.Errorf("failed to decode platform key '%s': %w", k.Key, err)
				}
				block, _ := pem.Decode(keyBytes)
				if block == nil {
					logger.Error("Failed to decode PEM block from database key", zap.String("key_name", k.Key))
					return fmt.Errorf("failed to decode PEM block for platform key '%s'", k.Key)
				}

				if k.Key == "public_key" {
					if !strings.Contains(block.Type, "PUBLIC KEY") {
						logger.Error("Invalid PEM type for public key in DB", zap.String("type", block.Type))
						return fmt.Errorf("invalid public key type in DB: %s", block.Type)
					}
					pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
					if err != nil {
						logger.Error("Failed to parse public key from DB", zap.Error(err))
						return fmt.Errorf("failed to parse public key from DB: %w", err)
					}
					var ok bool
					platformPublicKey, ok = pubInterface.(*rsa.PublicKey)
					if !ok {
						logger.Error("Public key from DB is not an RSA key")
						return fmt.Errorf("public key from DB is not RSA")
					}
					foundPub = true
				} else if k.Key == "private_key" {
					if !strings.Contains(block.Type, "PRIVATE KEY") {
						logger.Error("Invalid PEM type for private key in DB", zap.String("type", block.Type))
						return fmt.Errorf("invalid private key type in DB: %s", block.Type)
					}
					privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
					if err != nil {
						logger.Warn("Failed to parse PKCS8 private key from DB, attempting PKCS1", zap.Error(err))
						privInterface, err = x509.ParsePKCS1PrivateKey(block.Bytes)
						if err != nil {
							logger.Error("Failed to parse private key from DB as PKCS8 or PKCS1", zap.Error(err))
							return fmt.Errorf("failed to parse private key from DB: %w", err)
						}
					}
					var ok bool
					platformPrivateKey, ok = privInterface.(*rsa.PrivateKey)
					if !ok {
						logger.Error("Private key from DB is not an RSA key")
						return fmt.Errorf("private key from DB is not RSA")
					}
					foundPriv = true
				}
			}
			if !foundPub || !foundPriv {
				logger.Error("Could not load both platform keys from database records")
				return fmt.Errorf("incomplete platform key pair found in database")
			}
			logger.Info("Successfully loaded platform keys from database.")
		}
	}
	// --- End Platform Key Handling ---

	// --- Dex Client Setup ---
	logger.Info("Setting up Dex gRPC client...")
	dexClient, err := newDexClient(dexGrpcAddress)
	if err != nil {
		logger.Error("Failed to create Dex gRPC client", zap.String("address", dexGrpcAddress), zap.Error(err))
		return fmt.Errorf("failed to create Dex gRPC client: %w", err)
	}
	logger.Info("Dex gRPC client connected.", zap.String("address", dexGrpcAddress))

	logger.Info("Ensuring Dex OAuth clients...")
	err = ensureDexClients(ctx, logger, dexClient)
	if err != nil {
		logger.Error("Failed to ensure Dex OAuth clients", zap.Error(err))
		return fmt.Errorf("failed to ensure Dex OAuth clients: %w", err)
	}
	logger.Info("Dex OAuth clients ensured.")
	// --- End Dex Client Setup ---

	// --- Server Initialization ---
	logger.Info("Initializing main application server...")
	if !dbReady {
		logger.Error("Database is not ready, cannot initialize server")
		return errors.New("cannot initialize server, database not ready")
	}
	authServer := &Server{
		host:              platformHost,
		platformPublicKey: platformPublicKey,
		dexVerifier:       dexVerifier,
		dexClient:         dexClient,
		logger:            logger.Named("authServer"),
		db:                adb,
		authCache:         authCacheSvc,            // Inject AuthCacheService
		updateLogin:       make(chan User, 100000), // TODO: Remove if loop is removed
	}

	// TODO: Remove this goroutine call if UpdateLastLoginLoop is removed
	go authServer.UpdateLastLoginLoop()
	logger.Info("Application server initialized.")
	// --- End Server Initialization ---

	// --- Start HTTP Server ---
	logger.Info("Starting HTTP server...", zap.String("address", httpServerAddress))
	httpServerErrors := make(chan error, 1)
	go func() {
		routes := httpRoutes{
			logger:             logger.Named("httpRoutes"),
			platformPrivateKey: platformPrivateKey,
			db:                 adb,
			authCache:          authCacheSvc, // Inject AuthCacheService
			authServer:         authServer,
		}
		httpErr := httpserver.RegisterAndStart(ctx, logger.Named("httpServer"), httpServerAddress, &routes)
		if httpErr != nil && !errors.Is(httpErr, http.ErrServerClosed) {
			httpServerErrors <- fmt.Errorf("http server error: %w", httpErr)
		} else {
			httpServerErrors <- httpErr
		}
	}()
	// --- End Start HTTP Server ---

	// Wait for server exit
	logger.Info("Auth Service started successfully.")
	err = <-httpServerErrors

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("HTTP server stopped unexpectedly", zap.Error(err))
		return err
	}

	logger.Info("HTTP server stopped gracefully.")
	return nil
}

// newServerCredentials loads TLS transport credentials for the GRPC server.
func newServerCredentials(certPath string, keyPath string, caPath string) (credentials.TransportCredentials, error) {
	srv, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	p := x509.NewCertPool()

	if caPath != "" {
		ca, err := os.ReadFile(caPath)
		if err != nil {
			return nil, err
		}

		p.AppendCertsFromPEM(ca)
	}

	return credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{srv},
		RootCAs:      p,
	}), nil
}

func ensureDexClients(ctx context.Context, logger *zap.Logger, dexClient dexApi.DexClient) error {
	publicUris := strings.Split(dexPublicClientRedirectUris, ",")

	publicClientResp, _ := dexClient.GetClient(ctx, &dexApi.GetClientReq{
		Id: "public-client",
	})

	logger.Info("public URIS", zap.Any("uris", publicUris))

	if publicClientResp != nil && publicClientResp.Client != nil {
		publicClientReq := dexApi.UpdateClientReq{
			Id:           "public-client",
			Name:         "Public Client",
			RedirectUris: publicUris,
		}

		_, err := dexClient.UpdateClient(ctx, &publicClientReq)
		if err != nil {
			logger.Error("Auth Migrator: failed to create dex public client", zap.Error(err))
			return err
		}
	} else {
		publicClientReq := dexApi.CreateClientReq{
			Client: &dexApi.Client{
				Id:           "public-client",
				Name:         "Public Client",
				RedirectUris: publicUris,
				Public:       true,
			},
		}

		_, err := dexClient.CreateClient(ctx, &publicClientReq)
		if err != nil {
			logger.Error("Auth Migrator: failed to create dex public client", zap.Error(err))
			return err
		}
	}

	privateUris := strings.Split(dexPrivateClientRedirectUris, ",")

	logger.Info("private URIS", zap.Any("uris", privateUris))

	privateClientResp, _ := dexClient.GetClient(ctx, &dexApi.GetClientReq{
		Id: "private-client",
	})
	if privateClientResp != nil && privateClientResp.Client != nil {
		privateClientReq := dexApi.UpdateClientReq{
			Id:           "private-client",
			Name:         "Private Client",
			RedirectUris: privateUris,
		}

		_, err := dexClient.UpdateClient(ctx, &privateClientReq)
		if err != nil {
			logger.Error("Auth Migrator: failed to create dex private client", zap.Error(err))
			return err
		}
	} else {
		privateClientReq := dexApi.CreateClientReq{
			Client: &dexApi.Client{
				Id:           "private-client",
				Name:         "Private Client",
				RedirectUris: privateUris,
				Secret:       "secret",
			},
		}

		_, err := dexClient.CreateClient(ctx, &privateClientReq)
		if err != nil {
			logger.Error("Auth Migrator: failed to create dex private client", zap.Error(err))
			return err
		}
	}
	return nil
}

// --- Helpers  ---

func newDexOidcVerifier(ctx context.Context, domain, clientId string) (*oidc.IDTokenVerifier, error) {
	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 10,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	provider, err := oidc.NewProvider(
		oidc.InsecureIssuerURLContext(
			oidc.ClientContext(ctx, httpClient),
			domain,
		), domain,
	)
	if err != nil {
		return nil, err
	}

	return provider.Verifier(&oidc.Config{
		ClientID:          clientId,
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}), nil
}

func newDexClient(hostAndPort string) (dexApi.DexClient, error) {
	// TODO: Add TLS credentials if Dex gRPC requires them

	conn, err := grpc.NewClient(hostAndPort, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %v", err)
	}
	return dexApi.NewDexClient(conn), nil
}
