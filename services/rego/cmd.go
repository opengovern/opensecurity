package rego

import (
	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/steampipe"
	cloudql_init_job "github.com/opengovern/opencomply/jobs/cloudql-init-job"
	"github.com/opengovern/opencomply/services/integration/client"
	"github.com/opengovern/opencomply/services/rego/api"
	"github.com/opengovern/opencomply/services/rego/config"
	"github.com/opengovern/opencomply/services/rego/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"time"
)

func Command() *cobra.Command {
	var cnf config.RegoConfig
	config2.ReadFromEnv(&cnf, nil)

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			logger = logger.Named("rego")

			integrationClient := client.NewIntegrationServiceClient(cnf.Integration.BaseURL)

			pluginJob := cloudql_init_job.NewJob(logger, cloudql_init_job.Config{
				Postgres:      cnf.PostgresPlugin,
				ElasticSearch: cnf.ElasticSearch,
				Steampipe:     cnf.Steampipe,
			}, integrationClient)
			err = pluginJob.Run(ctx)
			if err != nil {
				logger.Error("failed to run plugin job", zap.Error(err))
				return err
			}

			time.Sleep(2 * time.Minute)

			steampipeConn, err := steampipe.StartSteampipeServiceAndGetConnection(logger)
			if err != nil {
				return err
			}

			regoEngine, err := service.NewRegoEngine(ctx, logger, steampipeConn)
			if err != nil {
				return err
			}

			return httpserver.RegisterAndStart(
				ctx,
				logger,
				cnf.Http.Address,
				api.New(logger, regoEngine),
			)
		},
	}

	return cmd
}
