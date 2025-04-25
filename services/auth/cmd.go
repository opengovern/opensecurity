package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log" // Use standard log for fatal startup errors before zap is ready
	"net/http"
	"os"
	"os/signal"
	"strconv" // Keep strconv only if needed elsewhere, removed for Port
	"strings"
	"syscall"
	"time"

	dexApi "github.com/dexidp/dex/api/v2"
	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/postgres"          // Used for postgres.Config
	"github.com/opengovern/opensecurity/services/auth/db" // Local DB package
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	// Use v4 as confirmed working by 'go get'
	jose "github.com/go-jose/go-jose/v4"

	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	// Keep environment variable reads for existing configs
	dexAuthDomain                = os.Getenv("DEX_AUTH_DOMAIN")
	dexAuthPublicClientID        = os.Getenv("DEX_AUTH_PUBLIC_CLIENT_ID")
	dexGrpcAddress               = os.Getenv("DEX_GRPC_ADDR")
	dexPublicClientRedirectUris  = os.Getenv("DEX_PUBLIC_CLIENT_REDIRECT_URIS")
	dexPrivateClientRedirectUris = os.Getenv("DEX_PRIVATE_CLIENT_REDIRECT_URIS")
	httpServerAddress            = os.Getenv("HTTP_ADDRESS")
	platformHost                 = os.Getenv("PLATFORM_HOST")
	platformKeyEnabledStr        = os.Getenv("PLATFORM_KEY_ENABLED")
	platformPublicKeyStr         = os.Getenv("PLATFORM_PUBLIC_KEY")
	platformPrivateKeyStr        = os.Getenv("PLATFORM_PRIVATE_KEY")
)

// calculateKeyID computes the JWK thumbprint (SHA256, Base64URL encoded).
func calculateKeyID(pub *rsa.PublicKey) (string, error) {
	jwk := jose.JSONWebKey{Key: pub}
	thumbprintBytes, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("failed to calculate JWK thumbprint: %w", err)
	}
	kid := base64.RawURLEncoding.EncodeToString(thumbprintBytes)
	return kid, nil
}

// Command creates the Cobra command for the auth service.
func Command() *cobra.Command {
	cmd := &cobra.Command{Use: "auth-service", Short: "Starts the OpenSecurity authentication and authorization service", SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			err := start(ctx)
			if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("ERROR: Service failed: %v\n", err)
				return err
			}
			if errors.Is(err, context.Canceled) {
				log.Println("Service shutdown requested via signal.")
			}
			log.Println("Service shutdown complete.")
			return nil
		},
	}
	return cmd
}

// ServerConfig holds configuration read from the environment.
type ServerConfig struct {
	PostgreSQL config2.Postgres // Assumes config2.Postgres has Port as string
}

// start initializes and runs the auth service components.
func start(ctx context.Context) error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("CRITICAL: Failed to initialize Zap logger: %v", err)
		return err
	}
	defer func() { _ = logger.Sync() }()
	logger = logger.Named("auth-service")
	logger.Info("Auth service starting...")

	// --- Configuration Reading & Validation ---
	var conf ServerConfig
	// Fix #1: Call ReadFromEnv directly as it returns no error
	config2.ReadFromEnv(&conf, nil)
	// Note: The original ReadFromEnv panics on strconv errors inside, so no error check needed here.
	logger.Info("Configuration loaded from environment")

	// Validate essential configurations
	if dexAuthDomain == "" || dexAuthPublicClientID == "" || dexGrpcAddress == "" {
		return fmt.Errorf("required Dex configuration missing (DEX_AUTH_DOMAIN, DEX_AUTH_PUBLIC_CLIENT_ID, DEX_GRPC_ADDR)")
	}
	if httpServerAddress == "" {
		return fmt.Errorf("required HTTP server address missing (HTTP_ADDRESS)")
	}
	// Fix #2: Validate Port as string
	if conf.PostgreSQL.Host == "" || conf.PostgreSQL.Port == "" || conf.PostgreSQL.Username == "" || conf.PostgreSQL.DB == "" {
		return fmt.Errorf("incomplete PostgreSQL configuration provided (host, port, user, db required)")
	}
	// Fix #2: Remove Atoi conversion here. Use string directly below.
	// dbPortInt, err := strconv.Atoi(conf.PostgreSQL.Port)
	// if err != nil || dbPortInt <= 0 {
	// 	return fmt.Errorf("invalid PostgreSQL port number '%s': %w", conf.PostgreSQL.Port, err)
	// }
	logger.Info("Configuration validated")

	// --- OIDC Verifier Setup ---
	dexVerifier, err := newDexOidcVerifier(ctx, dexAuthDomain, dexAuthPublicClientID)
	if err != nil {
		return fmt.Errorf("failed to create OIDC dex verifier: %w", err)
	}
	logger.Info("Instantiated Open ID Connect verifier", zap.String("issuer", dexAuthDomain))

	// --- Database Setup ---
	pgCfg := postgres.Config{
		Host: conf.PostgreSQL.Host,
		// Fix #2: Use the original Port string, assuming postgres.Config expects string
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      conf.PostgreSQL.DB,
		SSLMode: conf.PostgreSQL.SSLMode,
	}
	orm, err := postgres.NewClient(&pgCfg, logger.Named("postgres"))
	if err != nil {
		return fmt.Errorf("failed to create postgres client: %w", err)
	}
	sqlDB, err := orm.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	adb := db.Database{Orm: orm} // Concrete DB struct
	logger.Info("Connected to the postgres database", zap.String("db_name", conf.PostgreSQL.DB))
	if err := adb.Initialize(); err != nil {
		return fmt.Errorf("database migration/initialization error: %w", err)
	}
	logger.Info("Database initialized successfully")

	// --- Platform Key Loading/Generation (Keep existing logic) ---
	if platformKeyEnabledStr == "" {
		platformKeyEnabledStr = "false"
	}
	platformKeyEnabled, err := strconv.ParseBool(platformKeyEnabledStr)
	if err != nil {
		return fmt.Errorf("invalid PLATFORM_KEY_ENABLED value [%s]: %w", platformKeyEnabledStr, err)
	}
	var platformPublicKey *rsa.PublicKey
	var platformPrivateKey *rsa.PrivateKey
	var platformKeyID string
	if platformKeyEnabled {
		logger.Info("Loading platform keys from environment variables.")
		if platformPublicKeyStr == "" || platformPrivateKeyStr == "" {
			return fmt.Errorf("PLATFORM_KEY_ENABLED=true but PLATFORM_PUBLIC_KEY or PLATFORM_PRIVATE_KEY is missing")
		}
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
		logger.Info("Attempting to load/generate platform keys from/to database.")
		keyPair, err := adb.GetKeyPair(ctx)
		if err != nil {
			return fmt.Errorf("failed to query key pair from db: %w", err)
		}
		if len(keyPair) == 0 {
			logger.Info("No keys found in database, generating new platform RSA key pair.")
			platformPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return fmt.Errorf("error generating RSA key: %w", err)
			}
			platformPublicKey = &platformPrivateKey.PublicKey
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
	dexClient, conn, err := newDexClient(dexGrpcAddress)
	if err != nil {
		logger.Error("Failed to create dex client", zap.Error(err))
		return err
	}
	defer func() {
		logger.Info("Closing Dex gRPC client connection...")
		if closeErr := conn.Close(); closeErr != nil {
			logger.Warn("Error closing Dex gRPC client connection", zap.Error(closeErr))
		}
	}()
	if err = ensureDexClients(ctx, logger, dexClient); err != nil {
		logger.Error("Failed to ensure dex clients", zap.Error(err))
		return err
	}
	logger.Info("Dex gRPC client connected and clients ensured")

	// --- Instantiate Servers ---
	// Pass pointer (&adb) which satisfies db.DatabaseInterface
	// !!! THIS WILL FAIL TO COMPILE until methods in db/db.go are updated !!!
	authServer := &Server{
		host: platformHost, platformPublicKey: platformPublicKey, platformKeyID: platformKeyID,
		dexVerifier: dexVerifier, dexClient: dexClient, logger: logger.Named("authServer"),
		db:          &adb, // Assign pointer to concrete struct to interface field
		updateLogin: make(chan User, 100000),
	}
	go authServer.UpdateLastLoginLoop()

	errorsChan := make(chan error, 1)

	go func() {
		// !!! THIS WILL FAIL TO COMPILE until methods in db/db.go are updated !!!
		httpRoutes := httpRoutes{
			logger: logger.Named("httpRoutes"), platformPrivateKey: platformPrivateKey,
			platformKeyID: platformKeyID, db: &adb, // Assign pointer to concrete struct to interface field
			authServer: authServer,
		}
		logger.Info("Starting HTTP server", zap.String("address", httpServerAddress))
		serverErr := httpserver.RegisterAndStart(ctx, logger, httpServerAddress, &httpRoutes)
		if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
			errorsChan <- fmt.Errorf("http server error: %w", serverErr)
		} else {
			logger.Info("HTTP server shut down.")
			close(errorsChan)
		}
	}()

	// --- Wait for Shutdown Signal or Error ---
	logger.Info("Auth service started successfully. Waiting for shutdown signal...")
	select {
	case err, ok := <-errorsChan:
		// Fix #1: Check error from channel correctly using errors.Is
		if ok && err != nil {
			logger.Error("Service failed", zap.Error(err))
			if errors.Is(err, http.ErrServerClosed) {
				logger.Info("HTTP server closed normally.")
				return nil
			}
			return err
		}
		logger.Info("Errors channel closed, service stopped gracefully.")
		return nil
	case <-ctx.Done():
		logger.Info("Service shutting down due to context cancellation signal...")
		return ctx.Err()
	}
}

// --- Helper Functions (Keep as before, ensure imports are correct) ---

// newServerCredentials (Example TLS setup)
func newServerCredentials(certPath string, keyPath string, caPath string) (credentials.TransportCredentials, error) {
	srv, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS key pair (%s, %s): %w", certPath, keyPath, err)
	}
	cp := x509.NewCertPool()
	if caPath != "" {
		caBytes, err := os.ReadFile(caPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate %s: %w", caPath, err)
		}
		if !cp.AppendCertsFromPEM(caBytes) {
			return nil, fmt.Errorf("failed to append CA certs from %s", caPath)
		}
	}
	return credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{srv}, RootCAs: cp, MinVersion: tls.VersionTLS12}), nil
}

// newDexOidcVerifier creates a verifier for tokens issued by Dex.
func newDexOidcVerifier(ctx context.Context, domain, clientId string) (*oidc.IDTokenVerifier, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	providerCtx := oidc.ClientContext(ctx, httpClient)
	// IMPORTANT: Production Dex should use HTTPS, remove InsecureIssuerURLContext wrapper unless required for dev
	provider, err := oidc.NewProvider(oidc.InsecureIssuerURLContext(providerCtx, domain), domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider for %s: %w", domain, err)
	}
	return provider.Verifier(&oidc.Config{ClientID: clientId}), nil // Removed Skip checks
}

// newDexClient creates a gRPC client connection to Dex. Returns client, connection, error.
func newDexClient(hostAndPort string) (dexApi.DexClient, *grpc.ClientConn, error) {
	// Production TODO: Replace WithTransportCredentials(insecure.NewCredentials()) with TLS credentials
	log.Printf("WARNING: Connecting to Dex gRPC at %s using insecure credentials. THIS IS NOT PRODUCTION SAFE.", hostAndPort)
	conn, err := grpc.NewClient(hostAndPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial dex grpc server at %s: %w", hostAndPort, err)
	}
	return dexApi.NewDexClient(conn), conn, nil
}

// ensureDexClients ensures the required public and private OAuth clients exist in Dex.
func ensureDexClients(ctx context.Context, logger *zap.Logger, dexClient dexApi.DexClient) error {
	// Public Client
	publicUris := strings.Split(strings.TrimSpace(dexPublicClientRedirectUris), ",")
	if len(publicUris) == 0 || publicUris[0] == "" {
		logger.Warn("DEX_PUBLIC_CLIENT_REDIRECT_URIS is not set or empty, skipping public client setup.")
	} else {
		clientResp, err := dexClient.GetClient(ctx, &dexApi.GetClientReq{Id: "public-client"})
		if err != nil && !strings.Contains(err.Error(), "not found") {
			logger.Error("Failed to get dex public client", zap.Error(err))
			return fmt.Errorf("failed to get dex public client: %w", err)
		}
		logger.Info("Ensuring Dex public client exists/is updated", zap.Any("redirectURIs", publicUris))
		if clientResp != nil && clientResp.Client != nil {
			req := dexApi.UpdateClientReq{Id: "public-client", Name: "Public Client", RedirectUris: publicUris}
			_, err := dexClient.UpdateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to update dex public client", zap.Error(err))
				return fmt.Errorf("failed to update dex public client: %w", err)
			}
			logger.Info("Updated existing Dex public client.")
		} else {
			req := dexApi.CreateClientReq{Client: &dexApi.Client{Id: "public-client", Name: "Public Client", RedirectUris: publicUris, Public: true}}
			_, err := dexClient.CreateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to create dex public client", zap.Error(err))
				return fmt.Errorf("failed to create dex public client: %w", err)
			}
			logger.Info("Created new Dex public client.")
		}
	}
	// Private Client
	privateUris := strings.Split(strings.TrimSpace(dexPrivateClientRedirectUris), ",")
	if len(privateUris) == 0 || privateUris[0] == "" {
		logger.Warn("DEX_PRIVATE_CLIENT_REDIRECT_URIS is not set or empty, skipping private client setup.")
	} else {
		clientResp, err := dexClient.GetClient(ctx, &dexApi.GetClientReq{Id: "private-client"})
		if err != nil && !strings.Contains(err.Error(), "not found") {
			logger.Error("Failed to get dex private client", zap.Error(err))
			return fmt.Errorf("failed to get dex private client: %w", err)
		}
		logger.Info("Ensuring Dex private client exists/is updated", zap.Any("redirectURIs", privateUris))
		if clientResp != nil && clientResp.Client != nil {
			req := dexApi.UpdateClientReq{Id: "private-client", Name: "Private Client", RedirectUris: privateUris}
			_, err := dexClient.UpdateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to update dex private client", zap.Error(err))
				return fmt.Errorf("failed to update dex private client: %w", err)
			}
			logger.Info("Updated existing Dex private client.")
		} else {
			dexClientSecret := os.Getenv("DEX_PRIVATE_CLIENT_SECRET")
			if dexClientSecret == "" {
				dexClientSecret = "secret"
				logger.Warn("DEX_PRIVATE_CLIENT_SECRET not set, using insecure default secret for private client")
			}
			req := dexApi.CreateClientReq{Client: &dexApi.Client{Id: "private-client", Name: "Private Client", RedirectUris: privateUris, Secret: dexClientSecret}}
			_, err := dexClient.CreateClient(ctx, &req)
			if err != nil {
				logger.Error("Failed to create dex private client", zap.Error(err))
				return fmt.Errorf("failed to create dex private client: %w", err)
			}
			logger.Info("Created new Dex private client.")
		}
	}
	return nil
}
