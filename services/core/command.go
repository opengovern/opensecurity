package core

import (
	"context"
	"fmt"
	"os"

	"github.com/opengovern/og-util/pkg/httpserver"

	dexApi "github.com/dexidp/dex/api/v2"
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/vault"
	config2 "github.com/opengovern/opencomply/services/core/config"
	vault2 "github.com/opengovern/opencomply/services/core/vault"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strings"
)

var (
	SteampipeHost          = os.Getenv("STEAMPIPE_HOST")
	SteampipePort          = os.Getenv("STEAMPIPE_PORT")
	SteampipeDb            = os.Getenv("STEAMPIPE_DB")
	SteampipeUser          = os.Getenv("STEAMPIPE_USERNAME")
	SteampipePassword      = os.Getenv("STEAMPIPE_PASSWORD")
	PostgresPluginHost     = os.Getenv("POSTGRESPLUGIN_HOST")
	PostgresPluginPort     = os.Getenv("POSTGRESPLUGIN_PORT")
	PostgresPluginUsername = os.Getenv("POSTGRESPLUGIN_USERNAME")
	PostgresPluginPassword = os.Getenv("POSTGRESPLUGIN_PASSWORD")
	SchedulerBaseUrl       = os.Getenv("SCHEDULER_BASE_URL")
	IntegrationBaseUrl     = os.Getenv("INTEGRATION_BASE_URL")
	ComplianceBaseUrl      = os.Getenv("COMPLIANCE_BASE_URL")
)

func Command() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			var cnf config2.Config
			config.ReadFromEnv(&cnf, nil)

			return start(cmd.Context(), cnf)
		},
	}
}

func start(ctx context.Context, cnf config2.Config) error {
	cfg := koanf.Provide("core", config2.Config{})

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}
	if cfg.Vault.Provider == vault.HashiCorpVault {
		sealHandler, err := vault2.NewSealHandler(ctx, logger, cfg)
		if err != nil {
			return fmt.Errorf("new seal handler: %w", err)
		}
		// This blocks until vault is inited and unsealed
		sealHandler.Start(ctx)
	}

	dexClient, err := newDexClient(cfg.DexGrpcAddr)
	if err != nil {
		logger.Error("Auth Migrator: failed to create dex client", zap.Error(err))
		return err
	}

	publicUris := strings.Split(cfg.DexPublicClientRedirectUris, ",")

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

		_, err = dexClient.UpdateClient(ctx, &publicClientReq)
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

		_, err = dexClient.CreateClient(ctx, &publicClientReq)
		if err != nil {
			logger.Error("Auth Migrator: failed to create dex public client", zap.Error(err))
			return err
		}
	}

	privateUris := strings.Split(cfg.DexPrivateClientRedirectUris, ",")

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

		_, err = dexClient.UpdateClient(ctx, &privateClientReq)
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

		_, err = dexClient.CreateClient(ctx, &privateClientReq)
		if err != nil {
			logger.Error("Auth Migrator: failed to create dex private client", zap.Error(err))
			return err
		}
	}

	handler, err := InitializeHttpHandler(
		cfg,
		SteampipeHost, SteampipePort, SteampipeDb, SteampipeUser, SteampipePassword,
		SchedulerBaseUrl, IntegrationBaseUrl, ComplianceBaseUrl,
		logger, dexClient,
		cnf.ElasticSearch,
	)
	if err != nil {
		return fmt.Errorf("init http handler: %w", err)
	}

	return httpserver.RegisterAndStart(ctx, logger, cfg.Http.Address, handler)
}
