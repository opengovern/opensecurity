package api

import (
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/steampipe"
	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opencomply/services/integration/api/credentials"
	integration_type2 "github.com/opengovern/opencomply/services/integration/api/integration-types"
	"github.com/opengovern/opencomply/services/integration/api/integrations"
	"github.com/opengovern/opencomply/services/integration/db"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type API struct {
	logger          *zap.Logger
	database        db.Database
	steampipeConn   *steampipe.Database
	vault           vault.VaultSourceConfig
	kubeClient      client.Client
	typeManager     *integration_type.IntegrationTypeManager
	vaultKeyId      string
	masterAccessKey string
	masterSecretKey string
}

func New(
	logger *zap.Logger,
	db db.Database,
	vault vault.VaultSourceConfig,
	steampipeConn *steampipe.Database,
	kubeClient client.Client,
	typeManager *integration_type.IntegrationTypeManager,
) *API {
	return &API{
		logger:        logger.Named("api"),
		database:      db,
		vault:         vault,
		steampipeConn: steampipeConn,
		kubeClient:    kubeClient,
		typeManager:   typeManager,
	}
}

func (api *API) Register(e *echo.Echo) {
	integrationsApi := integrations.New(api.vault, api.database, api.logger, api.steampipeConn, api.kubeClient, api.typeManager)
	cred := credentials.New(api.vault, api.database, api.logger)
	integrationType := integration_type2.New(api.logger, api.typeManager)

	integrationsApi.Register(e.Group("/api/v1/integrations"))
	cred.Register(e.Group("/api/v1/credentials"))
	integrationType.Register(e.Group("/api/v1/integration-types"))
}
