// Package auth provides the core authentication and authorization logic,
// including OIDC integration with Dex, platform token issuance/verification,
// API key management, user management, and an Envoy external authorization check service.
// This file, cmd.go, defines the command structure (using Cobra) and the main
// startup sequence for the auth service executable.
package auth

import (
	"context"
	"crypto" // For crypto.SHA256
	"crypto/rand"
	"crypto/rsa" // For crypto.SHA256 hash
	"crypto/tls" // Needed for TLS config
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors" // <-- Added import for errors.Is
	"fmt"
	"log"      // Use standard log for fatal startup errors before zap is ready
	"net/http" // <-- Added import for http.ErrServerClosed and OIDC client
	"os"
	"os/signal" // <-- Added import for signal handling
	"strconv"   // <-- Added import for strconv.ParseBool, Atoi
	"strings"
	"syscall" // <-- Added import for signal handling
	"time"    // <-- Added import for OIDC client timeout

	dexApi "github.com/dexidp/dex/api/v2"
	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/services/auth/db" // Local DB package
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"          // For TLS credentials
	"google.golang.org/grpc/credentials/insecure" // For insecure fallback

	// Use v4 as confirmed working by 'go get'
	jose "github.com/go-jose/go-jose/v4"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Environment variables read at startup to configure the service.
var (
	// Dex Configuration
	dexAuthDomain         = os.Getenv("DEX_AUTH_DOMAIN")           // Base URL for Dex OIDC issuer (e.g., http://dex.example.com:5556/dex)
	dexAuthPublicClientID = os.Getenv("DEX_AUTH_PUBLIC_CLIENT_ID") // Dex OAuth client ID used by the frontend (e.g., "public-client")
	dexGrpcAddress        = os.Getenv("DEX_GRPC_ADDR")             // Address (host:port) for Dex's gRPC API (e.g., "dex.example.com:5557")
	// Dex gRPC TLS Configuration (File Paths)
	dexGrpcTlsCertPath = os.Getenv("DEX_GRPC_TLS_CERT_PATH") // Path to client certificate PEM file for mTLS with Dex gRPC (optional).
	dexGrpcTlsKeyPath  = os.Getenv("DEX_GRPC_TLS_KEY_PATH")  // Path to client private key PEM file for mTLS with Dex gRPC (optional).
	dexGrpcTlsCaPath   = os.Getenv("DEX_GRPC_TLS_CA_PATH")   // Path to CA certificate PEM file to verify the Dex gRPC server's certificate (recommended for production TLS).
	// Dex OAuth Client Redirect URIs (Comma-separated list)
	dexPublicClientRedirectUris  = os.Getenv("DEX_PUBLIC_CLIENT_REDIRECT_URIS")  // Allowed callback URLs for the 'public-client'.
	dexPrivateClientRedirectUris = os.Getenv("DEX_PRIVATE_CLIENT_REDIRECT_URIS") // Allowed callback URLs for the 'private-client'.
	// HTTP Server Configuration
	httpServerAddress = os.Getenv("HTTP_ADDRESS") // Listen address for this auth service's HTTP API (e.g., ":5555").
	// Platform Key Configuration
	platformHost          = os.Getenv("PLATFORM_HOST")        // Hostname associated with the platform (usage needs clarification, potentially for issuer/audience).
	platformKeyEnabledStr = os.Getenv("PLATFORM_KEY_ENABLED") // If "true", load platform keys from env vars below; otherwise, use DB or generate.
	platformPublicKeyStr  = os.Getenv("PLATFORM_PUBLIC_KEY")  // Base64 encoded PEM public key (used if PLATFORM_KEY_ENABLED=true).
	platformPrivateKeyStr = os.Getenv("PLATFORM_PRIVATE_KEY") // Base64 encoded PEM private key (PKCS8 recommended) (used if PLATFORM_KEY_ENABLED=true).
)

// calculateKeyID computes the JWK thumbprint (SHA256, Base64URL encoded) for the given public key.
// This thumbprint is used as the 'kid' (Key ID) in platform-issued JWT headers, allowing
// verifiers to identify the correct public key used for signing.
func calculateKeyID(pub *rsa.PublicKey) (string, error) {
	// Create a JWK representation of the public key.
	jwk := jose.JSONWebKey{Key: pub}
	// Calculate the thumbprint using SHA256 hash algorithm.
	thumbprintBytes, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("failed to calculate JWK thumbprint: %w", err)
	}
	// Encode the resulting hash bytes using Base64 URL encoding (no padding).
	kid := base64.RawURLEncoding.EncodeToString(thumbprintBytes)
	return kid, nil
}

// Command creates the Cobra command structure for the auth service executable.
// This allows the main package (`cmd/auth-service/main.go`) to simply execute this command,
// benefiting from Cobra's context handling and integration with OS signals for graceful shutdown.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth-service",
		Short: "Starts the OpenSecurity authentication and authorization service",
		Long: `Initializes dependencies (database, Dex connection, keys)
and starts the HTTP server handling authentication logic and APIs.
Listens for OS signals (Interrupt, SIGTERM) for graceful shutdown.`,
		SilenceUsage: true, // Prevents usage printing on RunE error return
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a root context that listens for OS interrupt signals (SIGINT, SIGTERM).
			// This context will be cancelled when such a signal is received.
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			// Ensure the stop function is called when RunE exits, releasing signal resources.
			defer stop()

			// Execute the main service startup logic, passing the signal-aware context.
			err := start(ctx)

			// Handle different shutdown scenarios for clean exit codes.
			if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, http.ErrServerClosed) {
				// Log fatal error only if it's not a context cancellation or clean HTTP server shutdown.
				log.Printf("ERROR: Service failed: %v\n", err) // Use standard log as Zap might be closed/unflushed.
				return err                                     // Return error for Cobra to handle exit code (non-zero).
			}
			if errors.Is(err, context.Canceled) {
				log.Println("Service shutdown requested via signal.")
			}
			log.Println("Service shutdown complete.")
			// Return nil on clean shutdown or context cancellation for Cobra (exit code 0).
			return nil
		},
	}
	return cmd
}

// ServerConfig holds configuration read from the environment via og-util/config.
// It primarily contains PostgreSQL connection details.
type ServerConfig struct {
	PostgreSQL config2.Postgres // Assumes config2.Postgres has fields like Host, Port (string), Username, Password, DB, SSLMode
}

// start initializes all components (logger, config, db, keys, clients)
// and runs the main service logic (HTTP server, background tasks).
// It accepts a context that can be cancelled (e.g., by OS signals) for graceful shutdown.
func start(ctx context.Context) error {
	// --- Logger Setup ---
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("CRITICAL: Failed to initialize Zap logger: %v", err)
		return err
	}
	defer func() { _ = logger.Sync() }() // Ensure logs are flushed before exiting.
	logger = logger.Named("auth-service")
	logger.Info("Auth service starting...")

	// --- Configuration Reading & Validation ---
	var conf ServerConfig
	// Read config from environment variables using og-util helper.
	// NOTE: This specific ReadFromEnv function panics on internal parsing errors.
	config2.ReadFromEnv(&conf, nil)
	logger.Info("Configuration loaded from environment")

	// Validate essential configurations read directly from environment or via config struct.
	if dexAuthDomain == "" || dexAuthPublicClientID == "" || dexGrpcAddress == "" {
		return fmt.Errorf("required Dex configuration missing (DEX_AUTH_DOMAIN, DEX_AUTH_PUBLIC_CLIENT_ID, DEX_GRPC_ADDR)")
	}
	if httpServerAddress == "" {
		return fmt.Errorf("required HTTP server address missing (HTTP_ADDRESS)")
	}
	if conf.PostgreSQL.Host == "" || conf.PostgreSQL.Port == "" || conf.PostgreSQL.Username == "" || conf.PostgreSQL.DB == "" {
		return fmt.Errorf("incomplete PostgreSQL configuration provided (host, port, user, db required)")
	}
	if _, err := strconv.Atoi(conf.PostgreSQL.Port); err != nil {
		return fmt.Errorf("invalid PostgreSQL port number '%s': must be a number", conf.PostgreSQL.Port)
	}
	logger.Info("Configuration validated")

	// --- OIDC Verifier Setup ---
	// Initialize the OIDC token verifier for validating tokens issued by Dex.
	dexVerifier, err := newDexOidcVerifier(ctx, dexAuthDomain, dexAuthPublicClientID)
	if err != nil {
		return fmt.Errorf("failed to create OIDC dex verifier: %w", err)
	}
	logger.Info("Instantiated Open ID Connect verifier", zap.String("issuer", dexAuthDomain))

	// --- Database Setup ---
	// Prepare PostgreSQL connection config.
	pgCfg := postgres.Config{Host: conf.PostgreSQL.Host, Port: conf.PostgreSQL.Port, User: conf.PostgreSQL.Username, Passwd: conf.PostgreSQL.Password, DB: conf.PostgreSQL.DB, SSLMode: conf.PostgreSQL.SSLMode}
	// Create GORM DB connection pool.
	orm, err := postgres.NewClient(&pgCfg, logger.Named("postgres"))
	if err != nil {
		return fmt.Errorf("failed to create postgres client: %w", err)
	}
	// Verify database connectivity.
	sqlDB, err := orm.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB from GORM: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	// Create the concrete database handler struct.
	adb := db.Database{Orm: orm}
	logger.Info("Connected to the postgres database", zap.String("db_name", conf.PostgreSQL.DB))
	// Run database migrations. Consider making this a separate step/command for production.
	if err := adb.Initialize(); err != nil {
		return fmt.Errorf("database migration/initialization error: %w", err)
	}
	logger.Info("Database initialized successfully")

	// --- Platform Key Loading/Generation ---
	// Determine whether to load keys from env vars or DB/generate.
	if platformKeyEnabledStr == "" {
		platformKeyEnabledStr = "false"
	}
	platformKeyEnabled, err := strconv.ParseBool(platformKeyEnabledStr)
	if err != nil {
		return fmt.Errorf("invalid PLATFORM_KEY_ENABLED value [%s]: %w", platformKeyEnabledStr, err)
	}
	var platformPublicKey *rsa.PublicKey
	var platformPrivateKey *rsa.PrivateKey
	var platformKeyID string // Holds the calculated 'kid'

	if platformKeyEnabled {
		// Load keys directly from environment variables.
		logger.Info("Loading platform keys from environment variables.")
		if platformPublicKeyStr == "" || platformPrivateKeyStr == "" {
			return fmt.Errorf("PLATFORM_KEY_ENABLED=true but PLATFORM_PUBLIC_KEY or PLATFORM_PRIVATE_KEY is missing")
		}
		// Decode and parse public key.
		pubBytes, err := base64.StdEncoding.DecodeString(platformPublicKeyStr)
		if err != nil {
			return fmt.Errorf("failed to base64 decode PLATFORM_PUBLIC_KEY: %w", err)
		}
		pubBlock, _ := pem.Decode(pubBytes)
		if pubBlock == nil {
			return fmt.Errorf("failed to pem decode PLATFORM_PUBLIC_KEY")
		}
		pubParsed, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse public key from env: %w", err)
		}
		var ok bool
		platformPublicKey, ok = pubParsed.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("key parsed from PLATFORM_PUBLIC_KEY is not an RSA public key")
		}
		// Decode and parse private key (expecting PKCS8 format).
		privBytes, err := base64.StdEncoding.DecodeString(platformPrivateKeyStr)
		if err != nil {
			return fmt.Errorf("failed to base64 decode PLATFORM_PRIVATE_KEY: %w", err)
		}
		privBlock, _ := pem.Decode(privBytes)
		if privBlock == nil {
			return fmt.Errorf("failed to pem decode PLATFORM_PRIVATE_KEY")
		}
		privParsed, err := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key (PKCS8) from env: %w", err)
		}
		platformPrivateKey, ok = privParsed.(*rsa.PrivateKey)
		if !ok {
			return fmt.Errorf("key parsed from PLATFORM_PRIVATE_KEY is not an RSA private key")
		}
	} else {
		// Load keys from database, or generate if not found.
		logger.Info("Attempting to load/generate platform keys from/to database.")
		keyPair, err := adb.GetKeyPair(ctx)
		if err != nil {
			return fmt.Errorf("failed to query key pair from db: %w", err)
		}
		if len(keyPair) == 0 {
			// Generate new keys if none exist in DB.
			logger.Info("No keys found in database, generating new platform RSA key pair.")
			platformPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return fmt.Errorf("error generating RSA key: %w", err)
			}
			platformPublicKey = &platformPrivateKey.PublicKey
			// Store public key in DB.
			bPub, errPub := x509.MarshalPKIXPublicKey(platformPublicKey)
			if errPub != nil {
				return fmt.Errorf("failed to marshal generated public key: %w", errPub)
			}
			bpPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: bPub})
			strPub := base64.StdEncoding.EncodeToString(bpPub)
			errDbPub := adb.AddConfiguration(ctx, &db.Configuration{Key: "public_key", Value: strPub})
			if errDbPub != nil {
				return fmt.Errorf("failed to save generated public key to db: %w", errDbPub)
			}
			// Store private key (PKCS8 format) in DB.
			bPri, errPri := x509.MarshalPKCS8PrivateKey(platformPrivateKey)
			if errPri != nil {
				return fmt.Errorf("failed to marshal generated private key (PKCS8): %w", errPri)
			}
			bpPri := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: bPri})
			strPri := base64.StdEncoding.EncodeToString(bpPri)
			errDbPri := adb.AddConfiguration(ctx, &db.Configuration{Key: "private_key", Value: strPri})
			if errDbPri != nil {
				return fmt.Errorf("failed to save generated private key to db: %w", errDbPri)
			}
			logger.Info("Saved generated key pair to database.")
		} else {
			// Load existing keys from DB records.
			logger.Info("Loading platform key pair from database.")
			var pubFound, privFound bool
			for _, k := range keyPair {
				keyBytes, err := base64.StdEncoding.DecodeString(k.Value)
				if err != nil {
					return fmt.Errorf("failed to base64 decode key '%s' from db: %w", k.Key, err)
				}
				block, _ := pem.Decode(keyBytes)
				if block == nil {
					return fmt.Errorf("failed to pem decode key '%s' from db", k.Key)
				}
				if k.Key == "public_key" {
					pubParsed, err := x509.ParsePKIXPublicKey(block.Bytes)
					if err != nil {
						return fmt.Errorf("failed to parse public key from db: %w", err)
					}
					var ok bool
					platformPublicKey, ok = pubParsed.(*rsa.PublicKey)
					if !ok {
						return fmt.Errorf("public key from db is not RSA")
					}
					pubFound = true
				} else if k.Key == "private_key" {
					privParsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
					if err != nil {
						return fmt.Errorf("failed to parse private key (PKCS8) from db: %w", err)
					}
					var ok bool
					platformPrivateKey, ok = privParsed.(*rsa.PrivateKey)
					if !ok {
						return fmt.Errorf("private key from db is not RSA")
					}
					privFound = true
				}
			}
			if !pubFound || !privFound {
				return fmt.Errorf("could not find both public and private keys in db configuration")
			}
		}
	}

	// Calculate Key ID ('kid') from the obtained public key using JWK thumbprint.
	if platformPublicKey == nil {
		return fmt.Errorf("platform public key was not loaded or generated")
	}
	platformKeyID, err = calculateKeyID(platformPublicKey)
	if err != nil {
		logger.Error("Failed to calculate Key ID from platform public key", zap.Error(err))
		return fmt.Errorf("failed to derive platform key ID: %w", err)
	}
	logger.Info("Derived platform Key ID (kid) for JWTs", zap.String("kid", platformKeyID))
	if platformPrivateKey == nil || platformKeyID == "" {
		return fmt.Errorf("platform private key or key ID could not be initialized")
	}

	// --- Dex Client Setup ---
	// Establish gRPC connection to Dex, handling TLS conditionally.
	dexClient, conn, err := newDexClient(logger, dexGrpcAddress)
	if err != nil {
		logger.Error("Failed to create dex client", zap.Error(err))
		return err
	}
	// Ensure the gRPC connection is closed when the start function returns.
	defer func() {
		logger.Info("Closing Dex gRPC client connection...")
		if closeErr := conn.Close(); closeErr != nil {
			logger.Warn("Error closing Dex gRPC client connection", zap.Error(closeErr))
		}
	}()
	// Verify/create required Dex OAuth clients (public-client, private-client).
	if err = ensureDexClients(ctx, logger, dexClient); err != nil {
		logger.Error("Failed to ensure dex clients", zap.Error(err))
		return err
	}
	logger.Info("Dex gRPC client connected and clients ensured")

	// --- Instantiate Servers ---
	// Create the core Server instance, injecting dependencies.
	// Pass address of adb (&adb) which satisfies db.DatabaseInterface.
	// !!! This assignment requires methods in db/db.go to be updated for context/interface !!!
	authServer := &Server{
		host: platformHost, platformPublicKey: platformPublicKey, platformKeyID: platformKeyID,
		dexVerifier: dexVerifier, dexClient: dexClient, logger: logger.Named("authServer"),
		db:          &adb,                    // Assign pointer to concrete struct to interface field
		updateLogin: make(chan User, 100000), // Consider configurable buffer size
	}
	// Start the background goroutine for updating last login times.
	go authServer.UpdateLastLoginLoop()

	// Channel to receive errors from the HTTP server goroutine.
	errorsChan := make(chan error, 1)

	// Start the HTTP server in a separate goroutine.
	go func() {
		// Create the struct holding HTTP route handlers and their dependencies.
		// !!! This assignment requires methods in db/db.go to be updated for context/interface !!!
		httpRoutes := httpRoutes{
			logger: logger.Named("httpRoutes"), platformPrivateKey: platformPrivateKey,
			platformKeyID: platformKeyID, db: &adb, // Assign pointer to concrete struct to interface field
			authServer: authServer,
		}
		logger.Info("Starting HTTP server", zap.String("address", httpServerAddress))
		// httpserver.RegisterAndStart blocks until the server shuts down or encounters an error.
		serverErr := httpserver.RegisterAndStart(ctx, logger, httpServerAddress, &httpRoutes)
		// Send error to channel ONLY if it's not ErrServerClosed (graceful shutdown).
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			logger.Error("HTTP server failed", zap.Error(serverErr))
			errorsChan <- fmt.Errorf("http server error: %w", serverErr)
		} else {
			logger.Info("HTTP server shut down.")
			// Close channel to signal graceful exit without error.
			close(errorsChan)
		}
	}()

	// --- Wait for Shutdown Signal or Error ---
	logger.Info("Auth service started successfully. Waiting for shutdown signal...")
	select {
	case err, ok := <-errorsChan: // Wait for an error from the HTTP server goroutine or channel close.
		if ok && err != nil { // Check if channel is open and error is not nil.
			logger.Error("Service failed", zap.Error(err))
			// http.ErrServerClosed is expected on graceful shutdown, don't return it as failure.
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("HTTP server closed normally.")
				return nil
			}
			return err // Return the actual error.
		}
		// Channel closed without error means graceful HTTP server shutdown.
		logger.Info("Errors channel closed, service stopped gracefully.")
		return nil
	case <-ctx.Done(): // Wait for the context (listening for OS signals) to be cancelled.
		logger.Info("Service shutting down due to context cancellation signal...")
		// The context passed to RegisterAndStart should trigger its shutdown.
		// Allow time for graceful shutdown? The defer stop() in RunE handles the signal context.
		return ctx.Err() // Return context error (e.g., context.Canceled).
	}
}

// --- Helper Functions ---

// newServerCredentials loads TLS credentials from specified file paths.
// Used for establishing secure gRPC connections (e.g., to Dex).
// certPath/keyPath are for the client's certificate (for mTLS). caPath is for verifying the server.
func newServerCredentials(certPath string, keyPath string, caPath string) (credentials.TransportCredentials, error) {
	var clientCert tls.Certificate
	var err error
	// Load client certificate and key if paths are provided (for mTLS).
	if certPath != "" && keyPath != "" {
		clientCert, err = tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load client TLS key pair (cert: %s, key: %s): %w", certPath, keyPath, err)
		}
	}

	// Load CA certificate pool for server verification.
	// Handle potential error from NewSystemCertPool.
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Printf("Warning: Failed to load system certificate pool: %v. Using empty pool.", err)
		caCertPool = x509.NewCertPool() // Fallback to empty pool.
	}

	if caPath != "" {
		caBytes, err := os.ReadFile(caPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate %s: %w", caPath, err)
		}
		if !caCertPool.AppendCertsFromPEM(caBytes) {
			return nil, fmt.Errorf("failed to append CA certs from %s", caPath)
		}
	}

	// Create TLS configuration.
	tlsConfig := &tls.Config{
		RootCAs:    caCertPool, // CAs used to verify the server certificate.
		MinVersion: tls.VersionTLS12,
	}
	// Add client certificate if loaded for mTLS.
	if clientCert.Certificate != nil {
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	return credentials.NewTLS(tlsConfig), nil
}

// newDexOidcVerifier creates a verifier for OIDC ID Tokens issued by Dex.
// It uses the provided context for HTTP requests during OIDC discovery.
func newDexOidcVerifier(ctx context.Context, domain, clientId string) (*oidc.IDTokenVerifier, error) {
	// Create an HTTP client with a timeout for OIDC discovery requests.
	httpClient := &http.Client{Timeout: 10 * time.Second}
	// Create a context aware of the HTTP client for the OIDC library.
	providerCtx := oidc.ClientContext(ctx, httpClient)

	// Discover the OIDC provider's configuration (endpoints, JWKS URI).
	// IMPORTANT: For production Dex using HTTPS, remove the InsecureIssuerURLContext wrapper.
	// Use: provider, err := oidc.NewProvider(providerCtx, domain)
	provider, err := oidc.NewProvider(oidc.InsecureIssuerURLContext(providerCtx, domain), domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider for %s: %w", domain, err)
	}

	// Return a verifier configured for the provider and client ID.
	// Standard checks (issuer, audience, expiry) are enabled by default.
	return provider.Verifier(&oidc.Config{ClientID: clientId}), nil
}

// newDexClient creates a gRPC client connection to the Dex API server.
// It conditionally uses TLS based on environment variables DEX_GRPC_TLS_*.
// Returns the Dex API client, the underlying gRPC connection (for closing), and an error.
func newDexClient(logger *zap.Logger, hostAndPort string) (dexApi.DexClient, *grpc.ClientConn, error) {
	var opts []grpc.DialOption // gRPC dialing options.

	// Check if TLS environment variables are set, indicating a secure connection is desired.
	if dexGrpcTlsCertPath != "" && dexGrpcTlsKeyPath != "" || dexGrpcTlsCaPath != "" {
		logger.Info("Attempting to establish Dex gRPC connection using TLS",
			zap.String("certPath", dexGrpcTlsCertPath),
			zap.String("keyPath", dexGrpcTlsKeyPath),
			zap.String("caPath", dexGrpcTlsCaPath))

		// Load TLS credentials from the specified file paths.
		creds, err := newServerCredentials(dexGrpcTlsCertPath, dexGrpcTlsKeyPath, dexGrpcTlsCaPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load TLS credentials for Dex gRPC client: %w", err)
		}
		// Add the loaded TLS credentials to the gRPC dial options.
		opts = append(opts, grpc.WithTransportCredentials(creds))
		logger.Info("Using TLS for Dex gRPC connection.")
	} else {
		// Fallback to insecure connection if no TLS paths are provided.
		// THIS IS NOT RECOMMENDED FOR PRODUCTION.
		logger.Warn("Using insecure credentials for Dex gRPC connection. Set DEX_GRPC_TLS_*_PATH for TLS.", zap.String("address", hostAndPort))
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Establish the gRPC connection with the determined options.
	// Consider adding grpc.WithBlock() if synchronous connection is needed at startup.
	conn, err := grpc.NewClient(hostAndPort, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial dex grpc server at %s: %w", hostAndPort, err)
	}

	// Return the client stub, the connection object (for deferred closing), and nil error.
	return dexApi.NewDexClient(conn), conn, nil
}

// ensureDexClients verifies that the required Dex OAuth2 clients ('public-client', 'private-client')
// exist and have the correct redirect URIs configured via the Dex gRPC API.
// It creates or updates the clients as needed based on environment variables.
func ensureDexClients(ctx context.Context, logger *zap.Logger, dexClient dexApi.DexClient) error {
	// --- Ensure Public Client ---
	publicUrisList := strings.Split(strings.TrimSpace(dexPublicClientRedirectUris), ",")
	// Filter out empty strings that might result from trailing commas etc.
	var validPublicUris []string
	for _, uri := range publicUrisList {
		if trimmed := strings.TrimSpace(uri); trimmed != "" {
			validPublicUris = append(validPublicUris, trimmed)
		}
	}

	if len(validPublicUris) == 0 {
		logger.Warn("DEX_PUBLIC_CLIENT_REDIRECT_URIS is not set or empty, skipping public client setup.")
	} else {
		clientID := "public-client"
		clientName := "Public Client"
		// Attempt to get the existing client configuration.
		clientResp, err := dexClient.GetClient(ctx, &dexApi.GetClientReq{Id: clientID})
		// Handle potential errors from GetClient, ignoring "not found".
		if err != nil && !strings.Contains(err.Error(), "not found") { // Crude check for not found error
			logger.Error("Failed to get dex public client", zap.String("clientID", clientID), zap.Error(err))
			return fmt.Errorf("failed to get dex public client '%s': %w", clientID, err)
		}
		logger.Info("Ensuring Dex public client exists/is updated", zap.String("clientID", clientID), zap.Strings("redirectURIs", validPublicUris))
		if clientResp != nil && clientResp.Client != nil {
			// Client exists, update it (e.g., to sync redirect URIs).
			req := dexApi.UpdateClientReq{Id: clientID, Name: clientName, RedirectUris: validPublicUris}
			_, err := dexClient.UpdateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to update dex public client", zap.String("clientID", clientID), zap.Error(err))
				return fmt.Errorf("failed to update dex public client '%s': %w", clientID, err)
			}
			logger.Info("Updated existing Dex public client.", zap.String("clientID", clientID))
		} else {
			// Client doesn't exist, create it.
			req := dexApi.CreateClientReq{Client: &dexApi.Client{Id: clientID, Name: clientName, RedirectUris: validPublicUris, Public: true}}
			_, err := dexClient.CreateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to create dex public client", zap.String("clientID", clientID), zap.Error(err))
				return fmt.Errorf("failed to create dex public client '%s': %w", clientID, err)
			}
			logger.Info("Created new Dex public client.", zap.String("clientID", clientID))
		}
	}

	// --- Ensure Private Client ---
	privateUrisList := strings.Split(strings.TrimSpace(dexPrivateClientRedirectUris), ",")
	var validPrivateUris []string
	for _, uri := range privateUrisList {
		if trimmed := strings.TrimSpace(uri); trimmed != "" {
			validPrivateUris = append(validPrivateUris, trimmed)
		}
	}

	if len(validPrivateUris) == 0 {
		logger.Warn("DEX_PRIVATE_CLIENT_REDIRECT_URIS is not set or empty, skipping private client setup.")
	} else {
		clientID := "private-client"
		clientName := "Private Client"
		clientResp, err := dexClient.GetClient(ctx, &dexApi.GetClientReq{Id: clientID})
		if err != nil && !strings.Contains(err.Error(), "not found") {
			logger.Error("Failed to get dex private client", zap.String("clientID", clientID), zap.Error(err))
			return fmt.Errorf("failed to get dex private client '%s': %w", clientID, err)
		}
		logger.Info("Ensuring Dex private client exists/is updated", zap.String("clientID", clientID), zap.Strings("redirectURIs", validPrivateUris))
		if clientResp != nil && clientResp.Client != nil {
			// Client exists, update it. Note: Dex API might not allow updating secrets via UpdateClient.
			req := dexApi.UpdateClientReq{Id: clientID, Name: clientName, RedirectUris: validPrivateUris}
			_, err := dexClient.UpdateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to update dex private client", zap.String("clientID", clientID), zap.Error(err))
				return fmt.Errorf("failed to update dex private client '%s': %w", clientID, err)
			}
			logger.Info("Updated existing Dex private client.", zap.String("clientID", clientID))
		} else {
			// Client doesn't exist, create it.
			// Production TODO: Load secret from secure source (env var, secrets manager).
			dexClientSecret := os.Getenv("DEX_PRIVATE_CLIENT_SECRET")
			if dexClientSecret == "" {
				dexClientSecret = "secret"
				logger.Warn("DEX_PRIVATE_CLIENT_SECRET not set, using insecure default secret for private client")
			}
			req := dexApi.CreateClientReq{Client: &dexApi.Client{Id: clientID, Name: clientName, RedirectUris: validPrivateUris, Secret: dexClientSecret}}
			_, err := dexClient.CreateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to create dex private client", zap.String("clientID", clientID), zap.Error(err))
				return fmt.Errorf("failed to create dex private client '%s': %w", clientID, err)
			}
			logger.Info("Created new Dex private client.", zap.String("clientID", clientID))
		}
	}
	return nil
}
