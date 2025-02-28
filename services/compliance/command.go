package compliance

import (
	"context"
	"fmt"
	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/httpserver"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
// S3AccessKey    = os.Getenv("S3_ACCESS_KEY")
// S3AccessSecret = os.Getenv("S3_ACCESS_SECRET")
// S3Region       = os.Getenv("S3_REGION")
)

type OpenAI struct {
	Token string
}

type ServerConfig struct {
	ElasticSearch config.ElasticSearch
	PostgreSQL    config.Postgres
	Scheduler     config.OpenGovernanceService
	Integration   config.OpenGovernanceService
	Core          config.OpenGovernanceService
	OpenAI        OpenAI
	Http          config.HttpServer

	MigratorJobQueueName string `yaml:"migrator_job_queue_name"`
}

func ServerCommand() *cobra.Command {
	return &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return startHttpServer(cmd.Context())
		},
	}
}

func startHttpServer(ctx context.Context) error {
	var conf ServerConfig
	config.ReadFromEnv(&conf, nil)

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new logger: %w", err)
	}

	handler, err := InitializeHttpHandler(ctx, conf,
		//S3Region, S3AccessKey, S3AccessSecret,
		logger)
	if err != nil {
		return fmt.Errorf("init http handler: %w", err)
	}

	return httpserver.RegisterAndStart(ctx, logger, conf.Http.Address, handler)
}
