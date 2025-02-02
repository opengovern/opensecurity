package checkup

import (
	"errors"
	"os"

	config2 "github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/opencomply/jobs/checkup-job/config"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	IntegrationBaseUrl = os.Getenv("INTEGRATION_BASE_URL")
	AuthBaseUrl        = os.Getenv("AUTH_BASE_URL")
	CoreBaseUrl    = os.Getenv("CORE_BASE_URL")
	NATSAddress        = os.Getenv("NATS_URL")
)

func WorkerCommand() *cobra.Command {
	var id string
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
			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			cmd.SilenceUsage = true
			var cnf config.WorkerConfig
			config2.ReadFromEnv(&cnf, nil)

			w, err := NewWorker(
				id,
				NATSAddress,
				logger,
				IntegrationBaseUrl,
				AuthBaseUrl,
				CoreBaseUrl,
				cnf,
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
