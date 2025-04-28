package auth

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"

	dexApi "github.com/dexidp/dex/api/v2"

	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/services/auth/db"

	"crypto/rand"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
)

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
	//m := email.NewSendGridClient(mailApiKey, mailSender, mailSenderName, logger)

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
		return fmt.Errorf("new postgres client: %w", err)
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
	if platformKeyEnabled {
		b, err := base64.StdEncoding.DecodeString(platformPublicKeyStr)
		if err != nil {
			return fmt.Errorf("public key decode: %w", err)
		}
		block, _ := pem.Decode(b)
		if block == nil {
			return fmt.Errorf("failed to decode my private key")
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return err
		}
		platformPublicKey = pub.(*rsa.PublicKey)

		b, err = base64.StdEncoding.DecodeString(platformPrivateKeyStr)
		if err != nil {
			return fmt.Errorf("private key decode: %w", err)
		}
		block, _ = pem.Decode(b)
		if block == nil {
			panic("failed to decode private key")
		}
		pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			panic(err)
		}
		platformPrivateKey = pri.(*rsa.PrivateKey)
	} else {
		keyPair, err := adb.GetKeyPair()
		if err != nil {
			panic(err)
		}

		if len(keyPair) == 0 {
			platformPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
			if err != nil {
				panic(fmt.Sprintf("Error generating RSA key: %v", err))
			}
			platformPublicKey = &platformPrivateKey.PublicKey

			b, err := x509.MarshalPKIXPublicKey(platformPublicKey)
			if err != nil {
				panic(err)
			}
			bp := pem.EncodeToMemory(&pem.Block{
				Type:    "RSA PUBLIC KEY",
				Headers: nil,
				Bytes:   b,
			})
			str := base64.StdEncoding.EncodeToString(bp)
			err = adb.AddConfiguration(&db.Configuration{
				Key:   "public_key",
				Value: str,
			})
			if err != nil {
				panic(err)
			}

			b, err = x509.MarshalPKCS8PrivateKey(platformPrivateKey)
			if err != nil {
				panic(err)
			}
			bp = pem.EncodeToMemory(&pem.Block{
				Type:    "RSA PRIVATE KEY",
				Headers: nil,
				Bytes:   b,
			})
			str = base64.StdEncoding.EncodeToString(bp)
			err = adb.AddConfiguration(&db.Configuration{
				Key:   "private_key",
				Value: str,
			})
			if err != nil {
				panic(err)
			}

		} else {
			for _, k := range keyPair {
				if k.Key == "public_key" {
					b, err := base64.StdEncoding.DecodeString(k.Value)
					if err != nil {
						return fmt.Errorf("public key decode: %w", err)
					}
					block, _ := pem.Decode(b)
					if block == nil {
						return fmt.Errorf("failed to decode my private key")
					}
					pub, err := x509.ParsePKIXPublicKey(block.Bytes)
					if err != nil {
						return err
					}
					platformPublicKey = pub.(*rsa.PublicKey)
				} else if k.Key == "private_key" {
					b, err := base64.StdEncoding.DecodeString(k.Value)
					if err != nil {
						return fmt.Errorf("private key decode: %w", err)
					}
					block, _ := pem.Decode(b)
					if block == nil {
						panic("failed to decode private key")
					}
					pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
					if err != nil {
						panic(err)
					}
					platformPrivateKey = pri.(*rsa.PrivateKey)
				}
			}
		}
	}

	dexClient, err := newDexClient(dexGrpcAddress)
	if err != nil {
		logger.Error("Auth Migrator: failed to create dex client", zap.Error(err))
		return err
	}

	err = ensureDexClients(ctx, logger, dexClient)
	if err != nil {
		logger.Error("Auth Migrator: failed to ensure dex clients", zap.Error(err))
		return err
	}

	authServer := &Server{
		host:                platformHost,
		platformPublicKey:   platformPublicKey,
		dexVerifier:         dexVerifier,
		dexClient:           dexClient,
		logger:              logger,
		db:                  adb,
		updateLoginUserList: nil,
		updateLogin:         make(chan User, 100000),
	}

	go authServer.UpdateLastLoginLoop()

	errors := make(chan error, 1)
	go func() {
		routes := httpRoutes{
			logger:             logger,
			platformPrivateKey: platformPrivateKey,
			db:                 adb,
			authServer:         authServer,
		}
		errors <- fmt.Errorf("http server: %w", httpserver.RegisterAndStart(ctx, logger, httpServerAddress, &routes))
	}()

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
