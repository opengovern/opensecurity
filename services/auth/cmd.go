// cmd.go
package auth

import (
	"context"
	"crypto" // Added
	"crypto/rand"
	"crypto/rsa" // Added
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	dexApi "github.com/dexidp/dex/api/v2"

	"github.com/go-jose/go-jose/v4" // Added
	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/services/auth/db"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
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
	// PLATFORM_KEY_ID env var reading is removed
)

// --- ADDED HELPER FUNCTION ---
// calculateKeyID computes the JWK thumbprint (SHA256, Base64URL encoded) for the given public key.
func calculateKeyID(pub *rsa.PublicKey) (string, error) {
	// Use the public key directly to create the JWK representation for thumbprint
	jwk := jose.JSONWebKey{Key: pub}
	thumbprintBytes, err := jwk.Thumbprint(crypto.SHA256) // Calculate SHA256 thumbprint
	if err != nil {
		return "", fmt.Errorf("failed to calculate JWK thumbprint: %w", err)
	}
	// The thumbprint needs to be base64url encoded for use as kid
	kid := base64.RawURLEncoding.EncodeToString(thumbprintBytes)
	return kid, nil
}

// --- END HELPER FUNCTION ---

func Command() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return start(cmd.Context())
		},
	}
}

type ServerConfig struct {
	PostgreSQL config2.Postgres
}

// start runs both HTTP and GRPC server.
// GRPC server has Check method to ensure user is
// authenticated and authorized to perform an action.
// HTTP server has multiple endpoints to view and update
// the user roles.
func start(ctx context.Context) error {
	var conf ServerConfig
	config2.ReadFromEnv(&conf, nil)

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	logger = logger.Named("auth")

	dexVerifier, err := newDexOidcVerifier(ctx, dexAuthDomain, dexAuthPublicClientID)
	if err != nil {
		return fmt.Errorf("open id connect dex verifier: %w", err)
	}
	logger.Info("Instantiated a new Open ID Connect verifier")

	// setup postgres connection
	cfg := postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      conf.PostgreSQL.DB,
		SSLMode: conf.PostgreSQL.SSLMode,
	}
	orm, err := postgres.NewClient(&cfg, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	adb := db.Database{Orm: orm}
	fmt.Println("Connected to the postgres database: ", conf.PostgreSQL.DB)
	err = adb.Initialize()
	if err != nil {
		// Use a more specific error message than just "new postgres client"
		return fmt.Errorf("database initialization error: %w", err)
	}

	if platformKeyEnabledStr == "" {
		platformKeyEnabledStr = "false"
	}
	platformKeyEnabled, err := strconv.ParseBool(platformKeyEnabledStr)
	if err != nil {
		return fmt.Errorf("platformKeyEnabled [%s]: %w", platformKeyEnabledStr, err)
	}

	var platformPublicKey *rsa.PublicKey
	var platformPrivateKey *rsa.PrivateKey
	var platformKeyID string // Variable to store the calculated Key ID

	if platformKeyEnabled {
		// --- Load keys from Environment Variables ---
		logger.Info("Loading platform keys from environment variables.")
		b, err := base64.StdEncoding.DecodeString(platformPublicKeyStr)
		if err != nil {
			return fmt.Errorf("public key decode: %w", err)
		}
		block, _ := pem.Decode(b)
		if block == nil {
			return fmt.Errorf("failed to decode public key PEM block from env")
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse public key from env: %w", err)
		}
		var ok bool
		platformPublicKey, ok = pub.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("key parsed from env is not an RSA public key")
		}

		b, err = base64.StdEncoding.DecodeString(platformPrivateKeyStr)
		if err != nil {
			return fmt.Errorf("private key decode: %w", err)
		}
		block, _ = pem.Decode(b)
		if block == nil {
			// Use return instead of panic for better error handling during startup
			return fmt.Errorf("failed to decode private key PEM block from env")
		}
		pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			// Use return instead of panic
			return fmt.Errorf("failed to parse private key from env: %w", err)
		}
		platformPrivateKey, ok = pri.(*rsa.PrivateKey)
		if !ok {
			return fmt.Errorf("key parsed from env is not an RSA private key")
		}
		// --- END Load keys ---

		// --- Calculate Key ID from loaded public key ---
		platformKeyID, err = calculateKeyID(platformPublicKey)
		if err != nil {
			logger.Error("Failed to calculate Key ID from environment public key", zap.Error(err))
			return fmt.Errorf("failed to derive platform key ID from env key: %w", err)
		}
		logger.Info("Derived platform Key ID (kid) from environment key", zap.String("kid", platformKeyID))
		// --- END Calculate Key ID ---

	} else {
		// --- Load/Generate keys from/to Database ---
		logger.Info("Attempting to load platform keys from database.")
		keyPair, err := adb.GetKeyPair()
		// Use return instead of panic for DB errors during startup
		if err != nil {
			return fmt.Errorf("failed to query key pair from db: %w", err)
		}

		if len(keyPair) == 0 {
			// --- Generate New Keys ---
			logger.Info("No keys found in database, generating new platform RSA key pair.")
			platformPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				return fmt.Errorf("error generating RSA key: %w", err)
			}
			platformPublicKey = &platformPrivateKey.PublicKey

			// Store public key
			bPub, errPub := x509.MarshalPKIXPublicKey(platformPublicKey)
			if errPub != nil {
				return fmt.Errorf("failed to marshal generated public key: %w", errPub)
			}
			bpPub := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: bPub})
			strPub := base64.StdEncoding.EncodeToString(bpPub)
			errDbPub := adb.AddConfiguration(&db.Configuration{Key: "public_key", Value: strPub})
			if errDbPub != nil {
				return fmt.Errorf("failed to save generated public key to db: %w", errDbPub)
			}

			// Store private key
			bPri, errPri := x509.MarshalPKCS8PrivateKey(platformPrivateKey)
			if errPri != nil {
				return fmt.Errorf("failed to marshal generated private key: %w", errPri)
			}
			bpPri := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: bPri})
			strPri := base64.StdEncoding.EncodeToString(bpPri)
			errDbPri := adb.AddConfiguration(&db.Configuration{Key: "private_key", Value: strPri})
			if errDbPri != nil {
				return fmt.Errorf("failed to save generated private key to db: %w", errDbPri)
			}
			logger.Info("Saved generated key pair to database.")
			// --- END Generate New Keys ---

			// --- Calculate Key ID from generated public key ---
			platformKeyID, err = calculateKeyID(platformPublicKey)
			if err != nil {
				logger.Error("Failed to calculate Key ID from generated public key", zap.Error(err))
				return fmt.Errorf("failed to derive platform key ID from generated key: %w", err)
			}
			logger.Info("Derived platform Key ID (kid) from generated key", zap.String("kid", platformKeyID))
			// --- END Calculate Key ID ---

		} else {
			// --- Load Keys From DB ---
			var pubFound, privFound bool
			for _, k := range keyPair {
				if k.Key == "public_key" {
					b, err := base64.StdEncoding.DecodeString(k.Value)
					if err != nil {
						return fmt.Errorf("db public key decode error: %w", err)
					}
					block, _ := pem.Decode(b)
					if block == nil {
						return fmt.Errorf("failed to decode public key PEM block from db")
					}
					pub, err := x509.ParsePKIXPublicKey(block.Bytes)
					if err != nil {
						return fmt.Errorf("failed to parse public key from db: %w", err)
					}
					var ok bool
					platformPublicKey, ok = pub.(*rsa.PublicKey)
					if !ok {
						return fmt.Errorf("key from db is not an RSA public key")
					}
					pubFound = true
				} else if k.Key == "private_key" {
					b, err := base64.StdEncoding.DecodeString(k.Value)
					if err != nil {
						return fmt.Errorf("db private key decode error: %w", err)
					}
					block, _ := pem.Decode(b)
					if block == nil {
						// Use return instead of panic
						return fmt.Errorf("failed to decode private key PEM block from db")
					}
					pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
					if err != nil {
						// Use return instead of panic
						return fmt.Errorf("failed to parse private key from db: %w", err)
					}
					var ok bool
					platformPrivateKey, ok = pri.(*rsa.PrivateKey)
					if !ok {
						return fmt.Errorf("key from db is not an RSA private key")
					}
					privFound = true
				}
			}
			if !pubFound || !privFound {
				return fmt.Errorf("could not find both public and private keys in db configuration")
			}
			logger.Info("Loaded platform key pair from database.")
			// --- END Load Keys From DB ---

			// --- Calculate Key ID from loaded public key ---
			platformKeyID, err = calculateKeyID(platformPublicKey)
			if err != nil {
				logger.Error("Failed to calculate Key ID from database public key", zap.Error(err))
				return fmt.Errorf("failed to derive platform key ID from db key: %w", err)
			}
			logger.Info("Derived platform Key ID (kid) from database key", zap.String("kid", platformKeyID))
			// --- END Calculate Key ID ---
		}
	}

	// Final check after loading/generating
	if platformPrivateKey == nil || platformPublicKey == nil || platformKeyID == "" {
		return fmt.Errorf("platform key pair or key ID could not be initialized")
	}

	// --- Dex Client Setup ---
	dexClient, err := newDexClient(dexGrpcAddress)
	if err != nil {
		logger.Error("Failed to create dex client", zap.Error(err))
		return err // Return the error directly
	}
	err = ensureDexClients(ctx, logger, dexClient)
	if err != nil {
		logger.Error("Failed to ensure dex clients", zap.Error(err))
		return err // Return the error directly
	}
	// --- END Dex Client Setup ---

	// --- Instantiate Servers ---
	authServer := &Server{
		host:                platformHost,
		platformPublicKey:   platformPublicKey,
		platformKeyID:       platformKeyID, // Pass calculated Key ID
		dexVerifier:         dexVerifier,
		dexClient:           dexClient,
		logger:              logger,
		db:                  adb,
		updateLoginUserList: nil, // Initialize properly later if needed
		updateLogin:         make(chan User, 100000),
	}

	go authServer.UpdateLastLoginLoop()

	errors := make(chan error, 1)
	go func() {
		routes := httpRoutes{
			logger:             logger,
			platformPrivateKey: platformPrivateKey,
			platformKeyID:      platformKeyID, // Pass calculated Key ID
			db:                 adb,
			authServer:         authServer,
		}
		// Use the httpserver package's function correctly
		errors <- fmt.Errorf("http server: %w", httpserver.RegisterAndStart(ctx, logger, httpServerAddress, &routes))
	}()
	// --- END Instantiate Servers ---

	return <-errors
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
		RootCAs:      p, // For server validation of clients (mTLS), use ClientCAs
	}), nil
}

// newDexOidcVerifier creates a verifier for tokens issued by Dex.
func newDexOidcVerifier(ctx context.Context, domain, clientId string) (*oidc.IDTokenVerifier, error) {
	transport := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		MaxIdleConnsPerHost: 10,
	}
	httpClient := &http.Client{
		Transport: transport,
	}

	// Use context with client for provider creation
	providerCtx := oidc.ClientContext(ctx, httpClient)

	// Allow insecure issuer URL for http endpoints (like localhost testing)
	// Use oidc.Provider for production with https
	provider, err := oidc.NewProvider(
		oidc.InsecureIssuerURLContext(providerCtx, domain),
		domain,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	return provider.Verifier(&oidc.Config{
		ClientID: clientId,
		// Skipping checks might be necessary if Dex issuer/client_id setup is unusual,
		// but it's generally safer to validate them if possible.
		SkipClientIDCheck: true,
		SkipIssuerCheck:   true,
	}), nil
}

// newDexClient creates a gRPC client connection to Dex.
func newDexClient(hostAndPort string) (dexApi.DexClient, error) {
	// Use grpc.WithInsecure() for non-TLS connections (e.g., localhost testing)
	// For production, configure TLS credentials:
	// creds, err := credentials.NewClientTLSFromFile(caPath, serverNameOverride)
	// conn, err := grpc.NewClient(hostAndPort, grpc.WithTransportCredentials(creds))
	conn, err := grpc.NewClient(hostAndPort, grpc.WithInsecure()) // Assuming insecure for now
	if err != nil {
		return nil, fmt.Errorf("failed to dial dex grpc server: %w", err)
	}
	return dexApi.NewDexClient(conn), nil
}

// ensureDexClients ensures the required public and private OAuth clients exist in Dex.
func ensureDexClients(ctx context.Context, logger *zap.Logger, dexClient dexApi.DexClient) error {
	// Public Client
	publicUris := strings.Split(dexPublicClientRedirectUris, ",")
	if len(publicUris) == 0 || (len(publicUris) == 1 && publicUris[0] == "") {
		logger.Warn("DEX_PUBLIC_CLIENT_REDIRECT_URIS is not set or empty, skipping public client setup.")
	} else {
		publicClientResp, _ := dexClient.GetClient(ctx, &dexApi.GetClientReq{Id: "public-client"})
		logger.Info("Checking Dex public client", zap.Any("redirectURIs", publicUris))
		if publicClientResp != nil && publicClientResp.Client != nil {
			// Update existing client if needed (e.g., to update redirect URIs)
			// You might want to compare existing URIs before updating
			publicClientReq := dexApi.UpdateClientReq{
				Id:           "public-client",
				Name:         "Public Client", // Keep name consistent
				RedirectUris: publicUris,
			}
			_, err := dexClient.UpdateClient(ctx, &publicClientReq)
			if err != nil {
				logger.Error("Failed to update dex public client", zap.Error(err))
				return fmt.Errorf("failed to update dex public client: %w", err)
			}
			logger.Info("Updated existing Dex public client.")
		} else {
			// Create new client
			publicClientReq := dexApi.CreateClientReq{
				Client: &dexApi.Client{
					Id:           "public-client",
					Name:         "Public Client",
					RedirectUris: publicUris,
					Public:       true, // Mark as public
				},
			}
			_, err := dexClient.CreateClient(ctx, &publicClientReq)
			if err != nil {
				logger.Error("Failed to create dex public client", zap.Error(err))
				return fmt.Errorf("failed to create dex public client: %w", err)
			}
			logger.Info("Created new Dex public client.")
		}
	}

	// Private Client
	privateUris := strings.Split(dexPrivateClientRedirectUris, ",")
	if len(privateUris) == 0 || (len(privateUris) == 1 && privateUris[0] == "") {
		logger.Warn("DEX_PRIVATE_CLIENT_REDIRECT_URIS is not set or empty, skipping private client setup.")
	} else {
		privateClientResp, _ := dexClient.GetClient(ctx, &dexApi.GetClientReq{Id: "private-client"})
		logger.Info("Checking Dex private client", zap.Any("redirectURIs", privateUris))
		if privateClientResp != nil && privateClientResp.Client != nil {
			// Update existing client
			privateClientReq := dexApi.UpdateClientReq{
				Id:           "private-client",
				Name:         "Private Client", // Keep name consistent
				RedirectUris: privateUris,
				// Secret cannot be updated via API, it seems
			}
			_, err := dexClient.UpdateClient(ctx, &privateClientReq)
			if err != nil {
				logger.Error("Failed to update dex private client", zap.Error(err))
				return fmt.Errorf("failed to update dex private client: %w", err)
			}
			logger.Info("Updated existing Dex private client.")
		} else {
			// Create new client
			privateClientReq := dexApi.CreateClientReq{
				Client: &dexApi.Client{
					Id:           "private-client",
					Name:         "Private Client",
					RedirectUris: privateUris,
					Secret:       "secret", // Use a configurable secret in production
				},
			}
			_, err := dexClient.CreateClient(ctx, &privateClientReq)
			if err != nil {
				logger.Error("Failed to create dex private client", zap.Error(err))
				return fmt.Errorf("failed to create dex private client: %w", err)
			}
			logger.Info("Created new Dex private client.")
		}
	}
	return nil
}
