package cloudql_init_job

import (
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/opensecurity/services/integration/client"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os/exec"
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

			_, err = j.Run(cmd.Context())
			killCmd := exec.Command("steampipe", "service", "stop", "--force")
			err = killCmd.Run()
			if err != nil {
				j.logger.Error("first stop failed", zap.Error(err))
				return err
			}
			//NOTE: stop must be called twice. it's not a mistake
			killCmd = exec.Command("steampipe", "service", "stop", "--force")
			err = killCmd.Run()
			if err != nil {
				j.logger.Error("second stop failed", zap.Error(err))
				return err
			}
			return err
		},
	}

	return cmd
}
