package cloudql_init_job

import (
	"github.com/opengovern/og-util/pkg/config"
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

			integrationClient := client.NewIntegrationServiceClient(cnf.Integration.BaseURL)
			j := NewJob(logger, cnf, integrationClient)

			return j.Run(cmd.Context())
		},
	}

	return cmd
}
