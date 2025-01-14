package query_runner

import (
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/services/integration/client"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func JobCommand() *cobra.Command {
	var cnf Config

	config.ReadFromEnv(&cnf, nil)

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			db, err := postgres.NewClient(&postgres.Config{
				Host:    cnf.Postgres.Host,
				Port:    cnf.Postgres.Port,
				User:    cnf.Postgres.Username,
				Passwd:  cnf.Postgres.Password,
				DB:      "integration_types",
				SSLMode: cnf.Postgres.SSLMode,
			}, logger.Named("postgres"))
			if err != nil {
				logger.Error("failed to create postgres client", zap.Error(err))
				return err
			}

			integrationClient := client.NewIntegrationServiceClient(cnf.Integration.BaseURL)
			j := NewJob(logger, cnf, db, integrationClient)

			return j.Run(cmd.Context())
		},
	}

	return cmd
}
