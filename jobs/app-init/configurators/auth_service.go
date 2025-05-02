// /Users/anil/workspace/opensecurity/jobs/app-init/configurators/auth_service.go
package configurators

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io" // Added for io.Discard
	"log"
	"net/http"
	"net/url"
	"strings"

	// Dex v2 API client
	dexApi "github.com/dexidp/dex/api/v2"
	// Password hashing
	"golang.org/x/crypto/bcrypt"
	// gRPC for Dex connection
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"                // For checking gRPC error codes
	"google.golang.org/grpc/credentials/insecure" // For insecure gRPC connection (adjust if TLS needed)
	"google.golang.org/grpc/status"               // For checking gRPC status

	// Local types package (adjust import path if necessary)
	initTypes "github.com/opengovern/opensecurity/jobs/app-init/types"
	// PostgreSQL driver (using blank identifier for side effects)
	_ "github.com/lib/pq"
)

// Constants related to user creation logic
const (
	defaultAdminRole    = "admin" // Role assigned to the initial user
	localConnectorID    = "local" // Connector ID used for local Dex users
	dexPasswordHashCost = 10      // Explicit bcrypt cost factor to match Dex expectations
)

// AuthComponent handles all initialization steps.
type AuthComponent struct {
	// HTTP Auth Service fields
	healthURL string       // URL for the Auth service health check
	client    *http.Client // Shared HTTP client for health checks

	// PostgreSQL fields
	pgHost     string // DB Hostname/IP
	pgPort     string // DB Port
	pgUser     string // DB Username
	pgPassword string // DB Password
	pgDbName   string // Target database name
	pgSslMode  string // DB SSL mode

	// Dex fields
	dexGrpcAddr            string           // Dex gRPC service address (host:port)
	dexPublicUris          string           // Comma-separated public client redirect URIs
	dexPrivateUris         string           // Comma-separated private client redirect URIs (kept for potential future use/config validation)
	dexClient              dexApi.DexClient // Dex gRPC client instance
	dexPublicClientID      string           // Configurable ID for the public client
	dexPrivateClientID     string           // Configurable ID for the private client (kept for potential future use/config validation)
	dexPrivateClientSecret string           // Configurable secret for the private client (kept for potential future use/config validation)
	dexHTTPHealthURL       string           // Full URL for Dex HTTP health check

	// Default User fields
	defaultUserEmail    string // Email for the initial user
	defaultUserName     string // Username for the initial user
	defaultUserPassword string // Password for the initial user
}

// newDexGrpcClient establishes the Dex gRPC connection with a timeout.
func newDexGrpcClient(ctx context.Context, hostAndPort string) (dexApi.DexClient, error) {
	log.Printf("INFO: Attempting to establish Dex gRPC connection to %s...", hostAndPort)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	dialCtx, cancel := context.WithTimeout(ctx, initTypes.DexTimeout)
	defer cancel()
	conn, err := grpc.DialContext(dialCtx, hostAndPort, opts...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("ERROR: Timed out connecting to Dex gRPC at %s after %v.", hostAndPort, initTypes.DexTimeout)
			return nil, fmt.Errorf("timed out connecting to dex gRPC server at %s: %w", hostAndPort, err)
		}
		log.Printf("ERROR: Failed to dial Dex gRPC at %s: %v.", hostAndPort, err)
		return nil, fmt.Errorf("failed to dial dex gRPC server at %s: %w", hostAndPort, err)
	}
	log.Printf("INFO: Successfully established Dex gRPC connection to %s.", hostAndPort)
	return dexApi.NewDexClient(conn), nil
}

// NewAuthComponent creates a new AuthComponent instance, validates input,
// and establishes the initial Dex gRPC connection.
func NewAuthComponent(
	// Auth HTTP
	healthCheckURL string,
	// Postgres
	pgHost, pgPort, pgUser, pgPassword, pgDbName, pgSslMode string,
	// Dex
	dexGrpcAddr string,
	dexPublicClientID string,
	dexPublicUris string,
	dexPrivateClientID string,
	dexPrivateUris string,
	dexPrivateClientSecret string,
	dexHTTPHealthURL string,
	// Default User
	defaultUserEmail, defaultUserName, defaultUserPassword string,
) (*AuthComponent, error) {
	log.Println("INFO: Creating new AuthComponent instance...")
	log.Println("INFO: Validating constructor parameters...")
	// --- Basic Input Validations ---
	if healthCheckURL == "" {
		return nil, errors.New("auth health check URL cannot be empty")
	}
	if _, err := url.ParseRequestURI(healthCheckURL); err != nil {
		return nil, fmt.Errorf("invalid auth health check URL format '%s': %w", healthCheckURL, err)
	}
	if pgHost == "" {
		return nil, errors.New("PGHOST cannot be empty")
	}
	if pgPort == "" {
		return nil, errors.New("PGPORT cannot be empty")
	}
	if pgUser == "" {
		return nil, errors.New("PGUSER cannot be empty")
	}
	if pgDbName == "" {
		return nil, errors.New("PGDATABASE cannot be empty")
	}
	if pgSslMode == "" {
		return nil, errors.New("PGSSLMODE cannot be empty")
	}
	if pgPassword == "" {
		log.Println("WARN: PGPASSWORD is empty.")
	}
	if dexGrpcAddr == "" {
		return nil, errors.New("DEX_GRPC_ADDR cannot be empty")
	}
	if dexPublicClientID == "" {
		return nil, errors.New("DEX_PUBLIC_CLIENT_ID cannot be empty")
	}
	if dexPublicUris == "" {
		return nil, errors.New("DEX_PUBLIC_CLIENT_REDIRECT_URIS cannot be empty")
	}
	if dexPrivateClientID == "" {
		return nil, errors.New("DEX_PRIVATE_CLIENT_ID cannot be empty")
	}
	if dexPrivateUris == "" {
		return nil, errors.New("DEX_PRIVATE_CLIENT_REDIRECT_URIS cannot be empty")
	}
	if dexPrivateClientSecret == "" {
		return nil, errors.New("DEX_PRIVATE_CLIENT_SECRET cannot be empty")
	}
	if dexHTTPHealthURL == "" {
		return nil, errors.New("DEX_HTTP_HEALTH_URL cannot be empty")
	}
	if _, err := url.ParseRequestURI(dexHTTPHealthURL); err != nil {
		return nil, fmt.Errorf("invalid DEX_HTTP_HEALTH_URL format '%s': %w", dexHTTPHealthURL, err)
	}
	if defaultUserEmail == "" {
		return nil, errors.New("DEFAULT_DEX_USER_EMAIL cannot be empty")
	}
	if !strings.Contains(defaultUserEmail, "@") {
		log.Printf("WARN: DEFAULT_DEX_USER_EMAIL ('%s') appears invalid.", defaultUserEmail)
	}
	if defaultUserName == "" {
		return nil, errors.New("DEFAULT_DEX_USER_NAME cannot be empty")
	}
	if defaultUserPassword == "" {
		return nil, errors.New("DEFAULT_DEX_USER_PASSWORD cannot be empty")
	}
	if len(defaultUserPassword) < 8 {
		log.Println("WARN: DEFAULT_DEX_USER_PASSWORD is less than 8 characters.")
	}
	log.Println("INFO: Constructor parameters validated.")

	// --- Create Dex gRPC Client ---
	log.Println("INFO: Establishing initial Dex gRPC client connection...")
	dexClient, err := newDexGrpcClient(context.Background(), dexGrpcAddr)
	if err != nil {
		log.Printf("ERROR: Component creation failed: Unable to connect to Dex gRPC at %s: %v", dexGrpcAddr, err)
		return nil, fmt.Errorf("failed to establish initial Dex gRPC connection during component creation: %w", err)
	}
	log.Println("INFO: Initial Dex gRPC client connection established.")

	httpClient := &http.Client{Timeout: initTypes.RequestTimeout}

	// --- Populate struct fields ---
	component := &AuthComponent{
		healthURL:              healthCheckURL,
		client:                 httpClient,
		pgHost:                 pgHost,
		pgPort:                 pgPort,
		pgUser:                 pgUser,
		pgPassword:             pgPassword,
		pgDbName:               pgDbName,
		pgSslMode:              pgSslMode,
		dexGrpcAddr:            dexGrpcAddr,
		dexPublicUris:          dexPublicUris,
		dexPrivateUris:         dexPrivateUris,
		dexClient:              dexClient,
		dexPublicClientID:      dexPublicClientID,
		dexPrivateClientID:     dexPrivateClientID,
		dexPrivateClientSecret: dexPrivateClientSecret,
		dexHTTPHealthURL:       dexHTTPHealthURL,
		defaultUserEmail:       defaultUserEmail,
		defaultUserName:        defaultUserName,
		defaultUserPassword:    defaultUserPassword,
	}
	log.Printf("INFO: AuthComponent instance created successfully.")
	return component, nil
}

// Name returns the component's descriptive name for logging.
func (a *AuthComponent) Name() string {
	return fmt.Sprintf("Auth Service Prerequisites (DB: %s, Dex Client: %s)",
		a.pgDbName, a.dexPublicClientID)
}

// getDSN constructs the PostgreSQL Data Source Name string.
func (a *AuthComponent) getDSN(targetDb string) string {
	dbToConnect := a.pgDbName
	if targetDb == "postgres" || targetDb == "" {
		dbToConnect = "postgres"
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(a.pgUser), url.QueryEscape(a.pgPassword), url.PathEscape(a.pgHost), url.PathEscape(a.pgPort), url.PathEscape(dbToConnect), url.QueryEscape(a.pgSslMode))
	maskedDSN := fmt.Sprintf("postgres://%s:[REDACTED]@%s:%s/%s?sslmode=%s",
		url.QueryEscape(a.pgUser), url.PathEscape(a.pgHost), url.PathEscape(a.pgPort), url.PathEscape(dbToConnect), url.QueryEscape(a.pgSslMode))
	log.Printf("DEBUG: [%s] Constructed DSN for DB '%s': %s", a.Name(), dbToConnect, maskedDSN)
	return dsn
}

// checkPostgresConnection attempts to connect and ping the specified PostgreSQL database.
func (a *AuthComponent) checkPostgresConnection(ctx context.Context, targetDb string) error {
	dbLogName := targetDb
	if dbLogName == "" || dbLogName == "postgres" {
		dbLogName = "(maintenance)"
	} else {
		dbLogName = a.pgDbName
	}
	log.Printf("DEBUG: [%s] Checking PostgreSQL connection to DB '%s'...", a.Name(), dbLogName)
	dsn := a.getDSN(targetDb)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("ERROR: [%s] Failed to prepare DB driver for '%s': %v", a.Name(), dbLogName, err)
		return fmt.Errorf("failed to open postgres driver for %s: %w", dbLogName, err)
	}
	defer db.Close()
	err = db.PingContext(ctx)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			log.Printf("WARN: [%s] Context error during ping for PostgreSQL DB '%s': %v", a.Name(), dbLogName, ctxErr)
			return fmt.Errorf("context error during ping for postgres database '%s': %w", dbLogName, ctxErr)
		}
		log.Printf("WARN: [%s] Failed to ping PostgreSQL DB '%s': %v.", a.Name(), dbLogName, err)
		return fmt.Errorf("failed to ping postgres database '%s': %w", dbLogName, err)
	}
	log.Printf("DEBUG: [%s] Successfully pinged PostgreSQL DB '%s'.", a.Name(), dbLogName)
	return nil
}

// checkHTTPGet performs a GET request to the specified URL, checking for a 2xx response.
func (a *AuthComponent) checkHTTPGet(ctx context.Context, urlToCheck string, purpose string) error {
	log.Printf("DEBUG: [%s] Performing HTTP GET check for %s at %s...", a.Name(), purpose, urlToCheck)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlToCheck, nil)
	if err != nil {
		log.Printf("ERROR: [%s] Failed to create HTTP request for %s (%s): %v", a.Name(), purpose, urlToCheck, err)
		return fmt.Errorf("failed create http request for %s (%s): %w", purpose, urlToCheck, err)
	}
	req.Header.Set("User-Agent", "opensecurity-app-init/1.0")
	resp, err := a.client.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			log.Printf("WARN: [%s] Context error during HTTP GET for %s (%s): %v", a.Name(), purpose, urlToCheck, ctxErr)
			return fmt.Errorf("context error during http get for %s (%s): %w", purpose, urlToCheck, ctxErr)
		}
		log.Printf("WARN: [%s] HTTP GET for %s (%s) failed: %v.", a.Name(), purpose, urlToCheck, err)
		return fmt.Errorf("http get for %s (%s) failed: %w", purpose, urlToCheck, err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("DEBUG: [%s] HTTP GET check for %s (%s) successful with status %s.", a.Name(), purpose, urlToCheck, resp.Status)
		return nil
	}
	log.Printf("WARN: [%s] Received non-2xx status %s from %s (%s).", a.Name(), resp.Status, purpose, urlToCheck)
	return fmt.Errorf("received non-2xx status from %s (%s): %s", purpose, urlToCheck, resp.Status)
}

// checkDexGRPC performs a GetVersion call to check Dex gRPC reachability.
func (a *AuthComponent) checkDexGRPC(ctx context.Context) error {
	log.Printf("DEBUG: [%s] Performing Dex gRPC check (GetVersion) to %s...", a.Name(), a.dexGrpcAddr)
	_, err := a.dexClient.GetVersion(ctx, &dexApi.VersionReq{})
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			log.Printf("WARN: [%s] Context error during Dex gRPC check (GetVersion): %v", a.Name(), ctxErr)
			return fmt.Errorf("context error during dex grpc check: %w", ctxErr)
		}
		log.Printf("WARN: [%s] Dex gRPC GetVersion check failed: %v.", a.Name(), err)
		return fmt.Errorf("dex grpc check failed: %w", err)
	}
	log.Printf("DEBUG: [%s] Dex gRPC GetVersion check successful.", a.Name())
	return nil
}

// --- InitializableComponent Interface Implementation ---

// CheckAvailability verifies connectivity to ALL external dependencies using retries.
func (a *AuthComponent) CheckAvailability(ctx context.Context) error {
	log.Printf("INFO: [%s] CheckAvailability: Verifying connectivity to all dependencies...", a.Name())
	var combinedErr error
	pgCheckFunc := func(c context.Context) error { return a.checkPostgresConnection(c, "postgres") }
	authHTTPCheckFunc := func(c context.Context) error { return a.checkHTTPGet(c, a.healthURL, "Auth Service") }
	dexHTTPCheckFunc := func(c context.Context) error { return a.checkHTTPGet(c, a.dexHTTPHealthURL, "Dex HTTP Health") }
	dexGRPCCheckFunc := func(c context.Context) error { return a.checkDexGRPC(c) }

	log.Printf("INFO: [%s] CheckAvailability: Checking PostgreSQL server connectivity...", a.Name())
	err := initTypes.WaitForCondition(ctx, a.Name(), "PostgreSQL server connectivity", initTypes.MaxRetries, initTypes.RetryDelay, pgCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] CheckAvailability Failed: PostgreSQL server.", a.Name())
		combinedErr = errors.Join(combinedErr, fmt.Errorf("postgresql server connectivity check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] CheckAvailability OK: PostgreSQL server.", a.Name())
	}

	log.Printf("INFO: [%s] CheckAvailability: Checking Auth HTTP endpoint connectivity (%s)...", a.Name(), a.healthURL)
	err = initTypes.WaitForCondition(ctx, a.Name(), "Auth HTTP endpoint connectivity", initTypes.MaxRetries, initTypes.RetryDelay, authHTTPCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] CheckAvailability Failed: Auth HTTP endpoint.", a.Name())
		combinedErr = errors.Join(combinedErr, fmt.Errorf("auth http endpoint connectivity check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] CheckAvailability OK: Auth HTTP endpoint.", a.Name())
	}

	log.Printf("INFO: [%s] CheckAvailability: Checking Dex HTTP health connectivity (%s)...", a.Name(), a.dexHTTPHealthURL)
	err = initTypes.WaitForCondition(ctx, a.Name(), "Dex HTTP health connectivity", initTypes.MaxRetries, initTypes.RetryDelay, dexHTTPCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] CheckAvailability Failed: Dex HTTP health endpoint.", a.Name())
		combinedErr = errors.Join(combinedErr, fmt.Errorf("dex http health connectivity check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] CheckAvailability OK: Dex HTTP health endpoint.", a.Name())
	}

	log.Printf("INFO: [%s] CheckAvailability: Checking Dex gRPC connectivity (%s)...", a.Name(), a.dexGrpcAddr)
	err = initTypes.WaitForCondition(ctx, a.Name(), "Dex gRPC connectivity", initTypes.MaxRetries, initTypes.RetryDelay, dexGRPCCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] CheckAvailability Failed: Dex gRPC.", a.Name())
		combinedErr = errors.Join(combinedErr, fmt.Errorf("dex grpc connectivity check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] CheckAvailability OK: Dex gRPC.", a.Name())
	}

	if combinedErr != nil {
		log.Printf("ERROR: [%s] CheckAvailability finished with errors.", a.Name())
		return combinedErr
	}
	log.Printf("INFO: [%s] CheckAvailability: All dependencies seem reachable.", a.Name())
	return nil
}

// CheckIfInitializationIsRequired checks if the target PostgreSQL database exists.
func (a *AuthComponent) CheckIfInitializationIsRequired(ctx context.Context) error {
	log.Printf("INFO: [%s] CheckIfInitializationRequired: Checking if target DB '%s' exists...", a.Name(), a.pgDbName)
	dsn := a.getDSN("postgres")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("ERROR: [%s] CheckIfInitRequired: Failed to open connection to maintenance DB: %v", a.Name(), err)
		return fmt.Errorf("checkIfInit: failed open conn to maintenance db: %w", err)
	}
	defer db.Close()
	queryCtx, cancel := context.WithTimeout(ctx, initTypes.DBTimeout)
	defer cancel()
	var exists int
	query := "SELECT 1 FROM pg_database WHERE datname = $1"
	log.Printf("DEBUG: [%s] CheckIfInitRequired: Executing query: %s with arg: %s", a.Name(), query, a.pgDbName)
	err = db.QueryRowContext(queryCtx, query, a.pgDbName).Scan(&exists)
	switch {
	case err == nil:
		log.Printf("INFO: [%s] CheckIfInitializationRequired: Prerequisite DB '%s' exists. Configuration steps will run.", a.Name(), a.pgDbName)
		return nil
	case errors.Is(err, sql.ErrNoRows):
		log.Printf("ERROR: [%s] CheckIfInitializationRequired: Prerequisite DB '%s' does NOT exist. Halting.", a.Name(), a.pgDbName)
		return fmt.Errorf("prerequisite database '%s' not found", a.pgDbName)
	default:
		if ctxErr := queryCtx.Err(); ctxErr != nil {
			log.Printf("ERROR: [%s] CheckIfInitRequired: Context error checking DB existence for '%s': %v", a.Name(), a.pgDbName, ctxErr)
			return fmt.Errorf("checkIfInit: context error checking db existence for '%s': %w", a.pgDbName, ctxErr)
		}
		log.Printf("ERROR: [%s] CheckIfInitRequired: Failed to query DB existence for '%s': %v", a.Name(), a.pgDbName, err)
		return fmt.Errorf("checkIfInit: failed to query db existence for '%s': %w", a.pgDbName, err)
	}
}

// configureDexClients ensures only the public Dex OAuth client exists.
func (a *AuthComponent) configureDexClients(ctx context.Context) error {
	log.Printf("INFO: [%s] Configure: Ensuring Dex public client (%s) exists...", a.Name(), a.dexPublicClientID)
	publicUris := strings.Split(a.dexPublicUris, ",")

	// --- Public Client ---
	publicClientID := a.dexPublicClientID
	log.Printf("INFO: [%s] Configure: Checking Dex client ID=%s", a.Name(), publicClientID)
	getCtxPub, getCancelPub := context.WithTimeout(ctx, initTypes.DexTimeout)
	_, err := a.dexClient.GetClient(getCtxPub, &dexApi.GetClientReq{Id: publicClientID})
	getCancelPub()

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			log.Printf("INFO: [%s] Configure: Public client %s not found, attempting creation...", a.Name(), publicClientID)
			createReq := dexApi.CreateClientReq{Client: &dexApi.Client{Id: publicClientID, Name: "Public Client", RedirectUris: publicUris, Public: true}}
			createCtx, createCancel := context.WithTimeout(ctx, initTypes.DexTimeout)
			_, createErr := a.dexClient.CreateClient(createCtx, &createReq)
			createCancel()
			if createErr != nil {
				log.Printf("ERROR: [%s] Configure: Failed to create dex public client %s: %v", a.Name(), publicClientID, createErr)
				return fmt.Errorf("failed to create dex public client %s: %w", publicClientID, createErr)
			}
			log.Printf("INFO: [%s] Configure: Successfully created Dex client: %s", a.Name(), publicClientID)
		} else {
			log.Printf("ERROR: [%s] Configure: Failed to get dex public client %s: %v", a.Name(), publicClientID, err)
			return fmt.Errorf("failed to get dex public client %s: %w", publicClientID, err)
		}
	} else {
		log.Printf("INFO: [%s] Configure: Public client %s already exists. No action needed.", a.Name(), publicClientID)
	}

	log.Printf("INFO: [%s] Configure: Skipping private client check/creation as requested.", a.Name())
	log.Printf("INFO: [%s] Configure: Dex client configuration check completed.", a.Name())
	return nil
}

// configureInitialUser checks and potentially creates the default user in Dex and the local DB.
func (a *AuthComponent) configureInitialUser(ctx context.Context, db *sql.DB) error {
	log.Printf("INFO: [%s] Configure: Checking if initial user '%s' needs setup...", a.Name(), a.defaultUserEmail)

	// 1. Check if user already exists in local DB
	checkCtx, checkCancel := context.WithTimeout(ctx, initTypes.DBTimeout)
	defer checkCancel()
	var existingUserID int64
	checkQuery := "SELECT id FROM users WHERE email = $1 LIMIT 1"
	log.Printf("DEBUG: [%s] Configure: Checking if user '%s' exists in local DB '%s'...", a.Name(), a.defaultUserEmail, a.pgDbName)
	err := db.QueryRowContext(checkCtx, checkQuery, a.defaultUserEmail).Scan(&existingUserID)

	switch {
	case err == nil:
		log.Printf("INFO: [%s] Configure: User '%s' (ID: %d) already exists in local database. Skipping initial user setup.", a.Name(), a.defaultUserEmail, existingUserID)
		return nil
	case errors.Is(err, sql.ErrNoRows):
		log.Printf("INFO: [%s] Configure: User '%s' not found in local DB. Proceeding with creation...", a.Name(), a.defaultUserEmail)
	default:
		if ctxErr := checkCtx.Err(); ctxErr != nil {
			log.Printf("ERROR: [%s] Configure: Context error checking user existence for '%s': %v", a.Name(), a.defaultUserEmail, ctxErr)
			return fmt.Errorf("initialUser: context error checking user existence: %w", ctxErr)
		}
		log.Printf("ERROR: [%s] Configure: Failed to query user existence for '%s': %v", a.Name(), a.defaultUserEmail, err)
		return fmt.Errorf("initialUser: failed to query user existence: %w", err)
	}

	// 2. User does not exist locally, proceed with creation in Dex and DB.
	log.Printf("INFO: [%s] Configure: Creating default user '%s' (Username: '%s')...", a.Name(), a.defaultUserEmail, a.defaultUserName)

	// --- DEBUG: Log password before hashing ---
	// WARNING: Logging raw passwords is a security risk. Remove this line before production.
	log.Printf("DEBUG: [%s] Configure: Raw password for default user '%s' before hashing: '%s' <<< REMOVE THIS LOG IN PRODUCTION >>>", a.Name(), a.defaultUserEmail, a.defaultUserPassword)
	// --- END DEBUG ---

	log.Printf("DEBUG: [%s] Configure: Hashing password for default user using cost %d...", a.Name(), dexPasswordHashCost)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.defaultUserPassword), dexPasswordHashCost)
	if err != nil {
		log.Printf("ERROR: [%s] Configure: Failed to hash default user password: %v", a.Name(), err)
		return fmt.Errorf("initialUser: failed to hash password: %w", err)
	}
	log.Printf("DEBUG: [%s] Configure: Password hashed successfully.", a.Name())

	// 3. Create Dex password entry
	dexUserID := fmt.Sprintf("%s|%s", localConnectorID, a.defaultUserEmail)
	log.Printf("INFO: [%s] Configure: Creating password entry in Dex (UserID: %s) for user '%s'...", a.Name(), dexUserID, a.defaultUserEmail)
	dexReq := dexApi.CreatePasswordReq{Password: &dexApi.Password{Email: a.defaultUserEmail, Username: a.defaultUserName, UserId: dexUserID, Hash: hashedPassword}}
	dexCtx, dexCancel := context.WithTimeout(ctx, initTypes.DexTimeout)
	defer dexCancel()
	_, err = a.dexClient.CreatePassword(dexCtx, &dexReq)
	if err != nil {
		if ctxErr := dexCtx.Err(); ctxErr != nil {
			log.Printf("ERROR: [%s] Configure: Context error creating password in Dex for user '%s': %v", a.Name(), a.defaultUserEmail, ctxErr)
			return fmt.Errorf("initialUser: context error creating dex password: %w", ctxErr)
		}
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Printf("WARN: [%s] Configure: Password entry already exists in Dex for user '%s'. Proceeding.", a.Name(), a.defaultUserEmail)
		} else {
			log.Printf("ERROR: [%s] Configure: Failed to create password in Dex for user '%s': %v", a.Name(), a.defaultUserEmail, err)
			return fmt.Errorf("initialUser: failed to create dex password: %w", err)
		}
	} else {
		log.Printf("INFO: [%s] Configure: Successfully created password entry in Dex for user '%s'.", a.Name(), a.defaultUserEmail)
	}

	// 4. Create local DB user entry
	log.Printf("INFO: [%s] Configure: Creating user entry in local database for '%s'...", a.Name(), a.defaultUserEmail)
	dbCtx, dbCancel := context.WithTimeout(ctx, initTypes.DBTimeout)
	defer dbCancel()
	insertQuery := `INSERT INTO users (email, username, full_name, role, connector_id, external_id, is_active, require_password_change, email_verified, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`
	log.Printf("DEBUG: [%s] Configure: Executing DB insert for user '%s'...", a.Name(), a.defaultUserEmail)
	_, insertErr := db.ExecContext(dbCtx, insertQuery, a.defaultUserEmail, a.defaultUserName, a.defaultUserName, defaultAdminRole, localConnectorID, dexUserID, true, true, false)
	if insertErr != nil {
		if ctxErr := dbCtx.Err(); ctxErr != nil {
			log.Printf("ERROR: [%s] Configure: Context error creating user '%s' in database: %v", a.Name(), a.defaultUserEmail, ctxErr)
			return fmt.Errorf("initialUser: context error creating db user: %w", ctxErr)
		}
		log.Printf("ERROR: [%s] Configure: Failed to create user '%s' in database: %v", a.Name(), a.defaultUserEmail, insertErr)
		return fmt.Errorf("initialUser: failed to insert user into database: %w", insertErr)
	}
	log.Printf("INFO: [%s] Configure: Successfully created initial user '%s' in database.", a.Name(), a.defaultUserEmail)
	return nil
}

// Configure orchestrates Dex client and initial user setup if the DB exists.
func (a *AuthComponent) Configure(ctx context.Context) error {
	log.Printf("INFO: [%s] Configure step starting: Ensure Dex client and initial user...", a.Name())

	// 1. Configure Dex Client (Now only public client)
	dexConfigureCtx, dexCancel := context.WithTimeout(ctx, initTypes.DexTimeout*2)
	defer dexCancel()
	log.Printf("INFO: [%s] Configure: Calling configureDexClients...", a.Name())
	err := a.configureDexClients(dexConfigureCtx)
	if err != nil {
		return fmt.Errorf("failed during Dex client configuration: %w", err)
	}
	log.Printf("INFO: [%s] Configure: Dex client configuration complete.", a.Name())

	// 2. Configure Initial User
	log.Printf("INFO: [%s] Configure: Connecting to target database '%s' for initial user setup...", a.Name(), a.pgDbName)
	dsn := a.getDSN(a.pgDbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("ERROR: [%s] Configure: Failed to open connection to target DB '%s': %v", a.Name(), a.pgDbName, err)
		return fmt.Errorf("configure: failed open conn target db '%s': %w", a.pgDbName, err)
	}
	defer db.Close()

	userConfigureCtx, userCancel := context.WithTimeout(ctx, initTypes.DBTimeout+initTypes.DexTimeout)
	defer userCancel()
	log.Printf("INFO: [%s] Configure: Calling configureInitialUser...", a.Name())
	err = a.configureInitialUser(userConfigureCtx, db)
	if err != nil {
		return fmt.Errorf("failed during initial user configuration: %w", err)
	}
	log.Printf("INFO: [%s] Configure: Initial user configuration complete.", a.Name())

	log.Printf("INFO: [%s] Configure step finished successfully.", a.Name())
	return nil
}

// CheckHealth verifies the operational status of all dependencies after configuration.
func (a *AuthComponent) CheckHealth(ctx context.Context) error {
	log.Printf("INFO: [%s] Performing health check on all dependencies post-configuration...", a.Name())
	var combinedErr error
	pgCheckFunc := func(c context.Context) error { return a.checkPostgresConnection(c, a.pgDbName) }
	authHTTPCheckFunc := func(c context.Context) error { return a.checkHTTPGet(c, a.healthURL, "Auth Service") }
	dexHTTPCheckFunc := func(c context.Context) error { return a.checkHTTPGet(c, a.dexHTTPHealthURL, "Dex HTTP Health") }
	dexGRPCCheckFunc := func(c context.Context) error { return a.checkDexGRPC(c) }

	log.Printf("INFO: [%s] Health Check: Verifying connection to target DB '%s'...", a.Name(), a.pgDbName)
	err := initTypes.WaitForCondition(ctx, a.Name(), fmt.Sprintf("PostgreSQL target DB '%s' connection", a.pgDbName), initTypes.MaxRetries, initTypes.RetryDelay, pgCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] Health Check Failed: PostgreSQL target DB '%s'.", a.Name(), a.pgDbName)
		combinedErr = errors.Join(combinedErr, fmt.Errorf("postgresql target database health check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] Health Check OK: PostgreSQL target DB '%s'.", a.Name(), a.pgDbName)
	}

	log.Printf("INFO: [%s] Health Check: Verifying Auth HTTP endpoint '%s'...", a.Name(), a.healthURL)
	err = initTypes.WaitForCondition(ctx, a.Name(), "Auth HTTP service health", initTypes.MaxRetries, initTypes.RetryDelay, authHTTPCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] Health Check Failed: Auth HTTP endpoint '%s'.", a.Name(), a.healthURL)
		combinedErr = errors.Join(combinedErr, fmt.Errorf("auth http service health check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] Health Check OK: Auth HTTP endpoint '%s'.", a.Name(), a.healthURL)
	}

	log.Printf("INFO: [%s] Health Check: Verifying Dex HTTP endpoint '%s'...", a.Name(), a.dexHTTPHealthURL)
	err = initTypes.WaitForCondition(ctx, a.Name(), "Dex HTTP health", initTypes.MaxRetries, initTypes.RetryDelay, dexHTTPCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] Health Check Failed: Dex HTTP endpoint '%s'.", a.Name(), a.dexHTTPHealthURL)
		combinedErr = errors.Join(combinedErr, fmt.Errorf("dex http health check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] Health Check OK: Dex HTTP endpoint '%s'.", a.Name(), a.dexHTTPHealthURL)
	}

	log.Printf("INFO: [%s] Health Check: Verifying Dex gRPC endpoint '%s'...", a.Name(), a.dexGrpcAddr)
	err = initTypes.WaitForCondition(ctx, a.Name(), "Dex gRPC health", initTypes.MaxRetries, initTypes.RetryDelay, dexGRPCCheckFunc)
	if err != nil {
		log.Printf("ERROR: [%s] Health Check Failed: Dex gRPC endpoint '%s'.", a.Name(), a.dexGrpcAddr)
		combinedErr = errors.Join(combinedErr, fmt.Errorf("dex grpc health check failed: %w", err))
	} else {
		log.Printf("INFO: [%s] Health Check OK: Dex gRPC endpoint '%s'.", a.Name(), a.dexGrpcAddr)
	}

	if combinedErr != nil {
		log.Printf("ERROR: [%s] Health check finished with errors.", a.Name())
		return combinedErr
	}
	log.Printf("INFO: [%s] Health check completed successfully. All dependencies appear healthy.", a.Name())
	return nil
}

// Compile-time check to ensure AuthComponent correctly implements the InitializableComponent interface.
var _ initTypes.InitializableComponent = (*AuthComponent)(nil)
