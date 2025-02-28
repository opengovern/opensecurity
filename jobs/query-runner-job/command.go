package query_runner

import (
	"errors"
	"github.com/opengovern/og-util/pkg/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func WorkerCommand() *cobra.Command {
	var (
		id  string
		cnf Config
	)
	config.ReadFromEnv(&cnf, nil)

	cmd := &cobra.Command{
		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case id == "":
				return errors.New("missing required flag 'id'")
			default:
				return nil
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			w, err := NewWorker(
				cnf,
				logger,
				cnf.PrometheusPushAddress,
				cmd.Context(),
			)
			if err != nil {
				return err
			}

			defer w.Stop()

			return w.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "The worker id")

	return cmd
}
