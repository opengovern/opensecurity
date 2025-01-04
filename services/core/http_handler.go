package core

import (
	"fmt"
	"time"

	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opencomply/services/core/config"
	integrationClient "github.com/opengovern/opencomply/services/integration/client"

	dexApi "github.com/dexidp/dex/api/v2"
	helmv2 "github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/og-util/pkg/steampipe"
	db2 "github.com/opengovern/opencomply/jobs/post-install-job/db"
	complianceClient "github.com/opengovern/opencomply/services/compliance/client"
	describeClient "github.com/opengovern/opencomply/services/describe/client"
	"go.uber.org/zap"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"github.com/opengovern/og-util/pkg/vault"

)

type HttpHandler struct {
	client            opengovernance.Client
	db                Database
	steampipeConn     *steampipe.Database
	schedulerClient   describeClient.SchedulerServiceClient
	integrationClient integrationClient.IntegrationServiceClient
	complianceClient  complianceClient.ComplianceServiceClient
	logger *zap.Logger
	viewCheckpoint time.Time
	cfg                config.Config
	kubeClient         client.Client
	vault              vault.VaultSourceConfig
	vaultSecretHandler vault.VaultSecretHandler
	dexClient          dexApi.DexClient
	migratorDb         *db2.Database

}

func InitializeHttpHandler(
	cfg config.Config,
	steampipeHost string, steampipePort string, steampipeDb string, steampipeUsername string, steampipePassword string,
	schedulerBaseUrl string, integrationBaseUrl string, complianceBaseUrl string,
	logger *zap.Logger,
) (h *HttpHandler, err error) {
	h = &HttpHandler{}

	fmt.Println("Initializing http handler")

	// setup postgres connection
	psqlCfg := postgres.Config{
		Host:    cfg.Postgres.Host,
		Port:    cfg.Postgres.Port,
		User:    cfg.Postgres.Username,
		Passwd:  cfg.Postgres.Password,
		DB:      cfg.Postgres.DB,
		SSLMode: cfg.Postgres.SSLMode,
	}
	orm, err := postgres.NewClient(&psqlCfg, logger)
	if err != nil {
		return nil, fmt.Errorf("new postgres client: %w", err)
	}

	h.db = Database{orm: orm}
	fmt.Println("Connected to the postgres database: ", postgresDb)

	err = h.db.Initialize()
	if err != nil {
		return nil, err
	}
	fmt.Println("Initialized postgres database: ", postgresDb)

	// setup steampipe connection
	steampipeConn, err := steampipe.NewSteampipeDatabase(steampipe.Option{
		Host: steampipeHost,
		Port: steampipePort,
		User: steampipeUsername,
		Pass: steampipePassword,
		Db:   steampipeDb,
	})
	h.steampipeConn = steampipeConn
	if err != nil {
		return nil, err
	}
	fmt.Println("Initialized steampipe database: ", steampipeConn)

	h.client, err = opengovernance.NewClient(opengovernance.ClientConfig{
		Addresses:     []string{cfg.ElasticSearch.Address},
		Username:     &cfg.ElasticSearch.Username,
		Password:     &cfg.ElasticSearch.Password,
		IsOnAks:      &cfg.ElasticSearch.IsOnAks,
		IsOpenSearch: &cfg.ElasticSearch.IsOpenSearch,
		AwsRegion:    &cfg.ElasticSearch.AwsRegion,
		AssumeRoleArn:&cfg.ElasticSearch.AssumeRoleArn,
	})
	if err != nil {
		return nil, err
	}
	h.schedulerClient = describeClient.NewSchedulerServiceClient(schedulerBaseUrl)

	h.integrationClient = integrationClient.NewIntegrationServiceClient(integrationBaseUrl)
	h.complianceClient = complianceClient.NewComplianceClient(complianceBaseUrl)

	h.logger = logger

	return h, nil
}


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
