package api

import (
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
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
	elastic         opengovernance.Client
	kubeClient      client.Client
	typeManager     *integration_type.IntegrationTypeManager
	vaultKeyId      string
	masterAccessKey string
	masterSecretKey string
	vault           vault.VaultSourceConfig

	steampipeOption *steampipe.Option
}

func New(
	logger *zap.Logger,
	db db.Database,
	vault vault.VaultSourceConfig,
	steampipeOption *steampipe.Option,
	kubeClient client.Client,
	typeManager *integration_type.IntegrationTypeManager,
	elastic opengovernance.Client,
) *API {
	return &API{
		logger:          logger.Named("api"),
		database:        db,
		vault:           vault,
		steampipeOption: steampipeOption,
		kubeClient:      kubeClient,
		typeManager:     typeManager,
		elastic:         elastic,
	}
}

func (api *API) Register(e *echo.Echo) {
	integrationsApi := integrations.New(api.vault, api.database, api.logger, api.steampipeOption, api.kubeClient, api.typeManager)
	cred := credentials.New(api.vault, api.database, api.logger)
	integrationType := integration_type2.New(api.typeManager, api.database, api.logger, api.elastic)

	integrationsApi.Register(e.Group("/api/v1/integrations"))
	cred.Register(e.Group("/api/v1/credentials"))
	integrationType.Register(e.Group("/api/v1/integration-types"))
}
