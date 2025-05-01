package core

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	cloudql_init_job "github.com/opengovern/opensecurity/jobs/cloudql-init-job"
	authClient "github.com/opengovern/opensecurity/services/auth/client"
	complianceapi "github.com/opengovern/opensecurity/services/compliance/api"
	coreApi "github.com/opengovern/opensecurity/services/core/api"

	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/google/uuid"
	api6 "github.com/hashicorp/vault/api"
	config3 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/og-util/pkg/vault"
	db2 "github.com/opengovern/opensecurity/jobs/post-install-job/db"
	"github.com/opengovern/opensecurity/jobs/post-install-job/db/model"
	complianceClient "github.com/opengovern/opensecurity/services/compliance/client"
	"github.com/opengovern/opensecurity/services/core/config"
	"github.com/opengovern/opensecurity/services/core/db"
	"github.com/opengovern/opensecurity/services/core/db/models"
	integrationClient "github.com/opengovern/opensecurity/services/integration/client"
	describeClient "github.com/opengovern/opensecurity/services/scheduler/client"
	"go.uber.org/zap"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type HttpHandler struct {
	client             opengovernance.Client
	db                 db.Database
	steampipeConn      *steampipe.Database
	schedulerClient    describeClient.SchedulerServiceClient
	integrationClient  integrationClient.IntegrationServiceClient
	complianceClient   complianceClient.ComplianceServiceClient
	authClient         authClient.AuthServiceClient
	logger             *zap.Logger
	viewCheckpoint     time.Time
	cfg                config.Config
	kubeClient         client.Client
	vault              vault.VaultSourceConfig
	vaultSecretHandler vault.VaultSecretHandler
	vaultSealHandler   *vault.HashiCorpVaultSealHandler // <<< ADD THIS FIELD (use actual type/interface, e.g., *vault.HashiCorpVaultSealHandler)

	migratorDb *db2.Database

	queryParameters []coreApi.QueryParameter
	queryParamsMu   sync.RWMutex

	complianceEnabled bool

	PluginJob *cloudql_init_job.Job
}

// InitializeHttpHandler initializes the main HttpHandler.
func InitializeHttpHandler(
	cfg config.Config, // core/config type
	schedulerBaseUrl string, integrationBaseUrl string, complianceBaseUrl string, authBaseUrl string,
	sealHandler *vault.HashiCorpVaultSealHandler, // Parameter is pointer type
	logger *zap.Logger, esConf config3.ElasticSearch, // Use og-util/pkg/config type
	complianceEnabled string,
) (h *HttpHandler, err error) {
	h = &HttpHandler{
		queryParamsMu: sync.RWMutex{},
		cfg:           cfg,
		logger:        logger,
		// *** Assign the pointer directly ***
		vaultSealHandler: sealHandler, // <<< CORRECTED assignment (no dereference *)
		// Initialize other fields... (Need to add based on struct definition)
		// Example: ensure these are initialized if needed later, even if nil initially
		// vault:             nil, // Or initialize properly if needed here
		// vaultSecretHandler: nil, // Or initialize properly
	}
	ctx := context.Background()

	logger.Info("Initializing http handler...")

	// --- Initialize Database Connections ---
	logger.Debug("Connecting to core database...")
	psqlCfg := postgres.Config{Host: cfg.Postgres.Host, Port: cfg.Postgres.Port, User: cfg.Postgres.Username, Passwd: cfg.Postgres.Password, DB: cfg.Postgres.DB, SSLMode: cfg.Postgres.SSLMode}
	orm, err := postgres.NewClient(&psqlCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("new core postgres client: %w", err)
	}
	h.db = db.NewDatabase(orm)
	if err = h.db.Initialize(); err != nil {
		return nil, fmt.Errorf("initialize core db: %w", err)
	}
	logger.Info("Connected and initialized core database", zap.String("db", cfg.Postgres.DB))

	logger.Debug("Connecting to migrator database...")
	migratorDbCfg := postgres.Config{Host: cfg.Postgres.Host, Port: cfg.Postgres.Port, User: cfg.Postgres.Username, Passwd: cfg.Postgres.Password, DB: "migrator", SSLMode: cfg.Postgres.SSLMode}
	migratorOrm, err := postgres.NewClient(&migratorDbCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("new migrator postgres client: %w", err)
	}
	if err := migratorOrm.AutoMigrate(&model.Migration{}); err != nil {
		return nil, fmt.Errorf("gorm migrate migrator db: %w", err)
	}
	h.migratorDb = &db2.Database{ORM: migratorOrm}
	logger.Info("Connected and initialized migrator database", zap.String("db", "migrator"))

	// --- Initialize App Config (if needed) ---
	apps, err := h.db.ListApp()
	if err != nil {
		return nil, fmt.Errorf("list app config: %w", err)
	}
	if len(apps) == 0 {
		logger.Info("No platform configuration found, creating initial entry.")
		err = h.db.CreateApp(&models.PlatformConfiguration{InstallID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()})
		if err != nil {
			return nil, fmt.Errorf("create app config: %w", err)
		}
	}

	// --- Initialize Kubernetes Client ---
	logger.Debug("Initializing Kubernetes client...")
	kubeClient, err := NewKubeClient()
	if err != nil {
		return nil, fmt.Errorf("new kube client: %w", err)
	}
	h.kubeClient = kubeClient
	logger.Debug("Kubernetes client initialized.")

	// --- Initialize Vault Clients (Secret Handler and Source Config) ---
	if cfg.Vault.Provider == vault.HashiCorpVault {
		logger.Debug("Initializing HashiCorp Vault Secret Handler...")
		h.vaultSecretHandler, err = vault.NewHashiCorpVaultSecretHandler(ctx, logger, cfg.Vault.HashiCorp)
		if err != nil {
			logger.Error("new hashicorp vault secret handler failed", zap.Error(err))
			return nil, fmt.Errorf("new hashicorp vault secret handler: %w", err)
		}

		logger.Debug("Initializing HashiCorp Vault Source Config...")
		h.vault, err = vault.NewHashiCorpVaultClient(ctx, logger, cfg.Vault.HashiCorp, cfg.Vault.KeyId)
		if err != nil {
			if h.vaultSecretHandler != nil && strings.Contains(err.Error(), api6.ErrSecretNotFound.Error()) {
				logger.Warn("Vault AES key secret not found, attempting to create.", zap.String("keyID", cfg.Vault.KeyId))
				b := make([]byte, 32)
				_, randErr := rand.Read(b)
				if randErr != nil {
					return nil, fmt.Errorf("generate random bytes for vault key: %w", randErr)
				}
				_, setErr := h.vaultSecretHandler.SetSecret(ctx, cfg.Vault.KeyId, b)
				if setErr != nil {
					return nil, fmt.Errorf("set initial vault secret for key %s: %w", cfg.Vault.KeyId, setErr)
				}
				logger.Info("Created new Vault AES key secret.", zap.String("keyID", cfg.Vault.KeyId))
				h.vault, err = vault.NewHashiCorpVaultClient(ctx, logger, cfg.Vault.HashiCorp, cfg.Vault.KeyId)
				if err != nil {
					logger.Error("new hashicorp vault source config after setSecret failed", zap.Error(err))
					return nil, fmt.Errorf("new hashicorp vault source config after setSecret: %w", err)
				}
			} else {
				logger.Error("new hashicorp vault source config failed", zap.Error(err))
				return nil, fmt.Errorf("new hashicorp vault source config: %w", err)
			}
		}
		logger.Info("HashiCorp Vault Secret Handler and Source Config initialized.")
	} else if cfg.Vault.Provider != "" {
		logger.Error("Unsupported vault provider specified in config", zap.String("provider", string(cfg.Vault.Provider)))
		return nil, fmt.Errorf("unsupported vault provider: %s", cfg.Vault.Provider)
	} else {
		logger.Info("No Vault provider configured in core config.")
	}

	// --- Initialize ES Client ---
	logger.Debug("Initializing Elasticsearch client...")
	h.client, err = opengovernance.NewClient(opengovernance.ClientConfig{Addresses: []string{esConf.Address}, Username: &esConf.Username, Password: &esConf.Password, IsOnAks: &esConf.IsOnAks, IsOpenSearch: &esConf.IsOpenSearch, AwsRegion: &esConf.AwsRegion, AssumeRoleArn: &esConf.AssumeRoleArn})
	if err != nil {
		logger.Error("Failed to create Elasticsearch client", zap.Error(err))
		return nil, fmt.Errorf("new opengovernance es client: %w", err)
	}
	logger.Info("Elasticsearch client initialized.")

	// --- Initialize Service Clients ---
	logger.Debug("Initializing service clients...")
	h.schedulerClient = describeClient.NewSchedulerServiceClient(schedulerBaseUrl)
	h.integrationClient = integrationClient.NewIntegrationServiceClient(integrationBaseUrl)
	h.authClient = authClient.NewAuthClient(authBaseUrl)
	logger.Debug("Initialized service clients", zap.Strings("urls", []string{schedulerBaseUrl, integrationBaseUrl, authBaseUrl}))

	// --- Initialize Compliance Client (optional) ---
	if lcComplianceEnabled := strings.ToLower(complianceEnabled); lcComplianceEnabled == "true" {
		h.complianceClient = complianceClient.NewComplianceClient(complianceBaseUrl)
		h.complianceEnabled = true
		logger.Info("Compliance client initialized.", zap.String("url", complianceBaseUrl))
	} else {
		h.complianceEnabled = false
		if lcComplianceEnabled != "false" && lcComplianceEnabled != "" {
			logger.Warn("Unsupported value for COMPLIANCE_ENABLED env var, disabling compliance client.", zap.String("value", complianceEnabled))
		} else {
			logger.Info("Compliance client disabled.")
		}
	}

	// --- Initialize Plugin Job Runner ---
	logger.Debug("Initializing Plugin Job Runner...")
	pluginJob := cloudql_init_job.NewJob(h.logger, cloudql_init_job.Config{
		Postgres:      config3.Postgres{Host: PostgresPluginHost, Port: PostgresPluginPort, Username: PostgresPluginUsername, Password: PostgresPluginPassword}, // Use package vars
		ElasticSearch: esConf, Steampipe: config3.Postgres{}, Integration: config3.OpenGovernanceService{BaseURL: integrationBaseUrl},
	}, h.integrationClient)
	h.PluginJob = pluginJob
	h.initializeSteampipePluginsWithRetry(ctx, 5, 2*time.Second) // Run Steampipe init
	logger.Debug("Plugin Job Runner initialized.")

	// Start background parameter fetch
	go h.fetchParameters()
	logger.Debug("Started background parameter fetch.")

	logger.Info("HttpHandler initialization complete.")
	return h, nil
} // --- END InitializeHttpHandler ---

func NewKubeClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := helmv2.AddToScheme(scheme); err != nil {
		return nil, err
	}
	if err := v1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	kubeClient, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

func (h *HttpHandler) fetchParameters() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	h.logger.Info("fetching parameters values")
	queryParams, err := h.listQueryParametersInternal()
	if err != nil {
		h.logger.Error("failed to get query parameters", zap.Error(err))
	} else {
		h.queryParamsMu.Lock()
		h.queryParameters = queryParams.Items
		h.queryParamsMu.Unlock()
	}

	for {
		select {
		case <-ticker.C:
			h.logger.Info("fetching parameters values")
			queryParams, err = h.listQueryParametersInternal()
			if err != nil {
				h.logger.Error("failed to get query parameters", zap.Error(err))
			} else {
				h.queryParamsMu.Lock()
				h.queryParameters = queryParams.Items
				h.queryParamsMu.Unlock()
			}
		}
	}
}

func (h *HttpHandler) listQueryParametersInternal() (coreApi.ListQueryParametersResponse, error) {
	clientCtx := &httpclient.Context{UserRole: authApi.AdminRole}
	var resp coreApi.ListQueryParametersResponse
	var err error

	var controls []complianceapi.Control
	if h.complianceEnabled {
		controls, err = h.complianceClient.ListControl(clientCtx, nil, nil)
		if err != nil {
			h.logger.Error("error listing controls", zap.Error(err))
			return resp, echo.NewHTTPError(http.StatusInternalServerError, "error listing controls")
		}
	}
	namedQueries, err := h.ListQueriesV2Internal(coreApi.ListQueryV2Request{})
	if err != nil {
		h.logger.Error("error listing queries", zap.Error(err))
		return resp, echo.NewHTTPError(http.StatusInternalServerError, "error listing queries")
	}

	var filteredQueryParams []string

	var queryParams []models.PolicyParameterValues
	if len(filteredQueryParams) > 0 {
		queryParams, err = h.db.GetQueryParametersByIds(filteredQueryParams)
		if err != nil {
			h.logger.Error("error getting query parameters", zap.Error(err))
			return resp, err
		}
	} else {
		queryParams, err = h.db.GetQueryParametersValues(nil)
		if err != nil {
			h.logger.Error("error getting query parameters", zap.Error(err))
			return resp, err
		}
	}

	parametersMap := make(map[string]*coreApi.QueryParameter)
	for _, dbParam := range queryParams {
		apiParam := coreApi.QueryParameter{
			Key:       dbParam.Key,
			ControlID: dbParam.ControlID,
			Value:     dbParam.Value,
		}
		parametersMap[apiParam.Key] = &apiParam
	}

	for _, c := range controls {
		for _, p := range c.Policy.Parameters {
			if _, ok := parametersMap[p.Key]; ok {
				parametersMap[p.Key].ControlsCount += 1
			}
		}
	}
	for _, q := range namedQueries.Items {
		for _, p := range q.Query.Parameters {
			if _, ok := parametersMap[p.Key]; ok {
				parametersMap[p.Key].QueriesCount += 1
			}
		}
	}

	var items []coreApi.QueryParameter
	for _, i := range parametersMap {
		items = append(items, *i)
	}

	totalCount := len(items)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})

	return coreApi.ListQueryParametersResponse{
		TotalCount: totalCount,
		Items:      items,
	}, nil
}

func (h *HttpHandler) initializeSteampipePlugins(ctx context.Context) {
	h.logger.Info("running plugin job to initialize integrations in cloudql")
	steampipeConn, err := h.PluginJob.Run(ctx)
	if err != nil {
		h.logger.Error("failed to run plugin job", zap.Error(err))
	}

	h.steampipeConn = steampipeConn
	fmt.Println("Initialized steampipe database: ", steampipeConn)
}

func (h *HttpHandler) initializeSteampipePluginsWithRetry(ctx context.Context, maxRetries int, initialBackoff time.Duration) {
	retries := 0
	backoff := initialBackoff

	for {
		h.logger.Info("Initializing Steampipe plugins. Attempt:", zap.Int("retry", retries+1))
		h.initializeSteampipePlugins(ctx)

		if h.steampipeConn != nil {
			h.logger.Info("Successfully initialized Steampipe plugins")
			return
		}

		if retries >= maxRetries {
			h.logger.Error("Max retries reached. Failed to initialize Steampipe plugins.")
			return
		}

		retries++
		h.logger.Warn("Retrying initialization after backoff...", zap.Duration("backoff", backoff))
		time.Sleep(backoff)
		backoff *= 2
	}
}
