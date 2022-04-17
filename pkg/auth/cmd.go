package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/spf13/cobra"
	"gitlab.com/keibiengine/keibi-engine/pkg/internal/httpserver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgreSQLHost     = os.Getenv("POSTGRESQL_HOST")
	postgreSQLPort     = os.Getenv("POSTGRESQL_PORT")
	postgreSQLDb       = os.Getenv("POSTGRESQL_DB")
	postgreSQLUser     = os.Getenv("POSTGRESQL_USERNAME")
	postgreSQLPassword = os.Getenv("POSTGRESQL_PASSWORD")

	azureAuthTenantName   = os.Getenv("AZURE_OAUTH_TENANT_NAME")
	azureAuthTenantID     = os.Getenv("AZURE_OAUTH_TENANT_ID")
	azureAuthClientID     = os.Getenv("AZURE_OAUTH_CLIENT_ID")
	azureAuthSignInPolicy = os.Getenv("AZURE_OAUTH_POLICY")

	httpServerAddress = os.Getenv("HTTP_ADDRESS")

	grpcServerAddress = os.Getenv("GRPC_ADDRESS")
	grpcTlsCertPath   = os.Getenv("GRPC_TLS_CERT_PATH")
	grpcTlsKeyPath    = os.Getenv("GRPC_TLS_KEY_PATH")
	grpcTlsCAPath     = os.Getenv("GRPC_TLS_CA_PATH")
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
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=GMT`,
		postgreSQLHost,
		postgreSQLPort,
		postgreSQLUser,
		postgreSQLPassword,
		postgreSQLDb,
	)

	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	logger.Info("Connected to the postgres database: ", zap.String("orm", "postgresDb"))

	verifier, err := newOidcVerifier(ctx, azureAuthTenantName, azureAuthTenantID, azureAuthClientID, azureAuthSignInPolicy)
	if err != nil {
		return fmt.Errorf("open id connect verifier: %w", err)
	}
	logger.Info("Instantiated a new Open ID Connect verifier")

	db := Database{orm: orm}

	err = db.Initialize()
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}

	authServer := Server{
		db:       db,
		verifier: verifier,
		authEcho: buildEchoRoutes(),
		logger:   logger,
	}

	creds, err := newServerCredentials(
		grpcTlsCertPath,
		grpcTlsKeyPath,
		grpcTlsCAPath,
	)
	if err != nil {
		return fmt.Errorf("grpc tls creds: %w", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	envoyauth.RegisterAuthorizationServer(grpcServer, authServer)

	lis, err := net.Listen("tcp", grpcServerAddress)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}

	errors := make(chan error, 1)
	go func() {
		errors <- fmt.Errorf("grpc server: %w", grpcServer.Serve(lis))
	}()

	go func() {
		errors <- fmt.Errorf("http server: %w", httpserver.RegisterAndStart(logger, httpServerAddress, httpRoutes{db: db}))
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
		ca, err := ioutil.ReadFile(caPath) //nolint(gosec)
		if err != nil {
			return nil, err
		}

		p.AppendCertsFromPEM(ca)
	}

	return credentials.NewTLS(&tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{srv},
		RootCAs:      p,
	}), nil
}
