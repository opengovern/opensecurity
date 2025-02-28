package rego

import (
	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"
	cloudql_init_job "github.com/opengovern/opensecurity/jobs/cloudql-init-job"
	"github.com/opengovern/opensecurity/services/integration/client"
	"github.com/opengovern/opensecurity/services/rego/api"
	"github.com/opengovern/opensecurity/services/rego/config"
	"github.com/opengovern/opensecurity/services/rego/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
			steampipeConn, err := pluginJob.Run(ctx)
			if err != nil {
				logger.Error("failed to run plugin job", zap.Error(err))
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
