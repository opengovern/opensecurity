package describe

import (
	"errors"
	"os"

	"github.com/opengovern/og-util/pkg/config"
	config2 "github.com/opengovern/opencomply/services/describe/config"
	"github.com/spf13/cobra"
)

const (
	DescribeResultsQueueName = "opengovernance-describe-results-queue"
	DescribeStreamName       = "describe"
)

var (
	PostgreSQLHost     = os.Getenv("POSTGRESQL_HOST")
	PostgreSQLPort     = os.Getenv("POSTGRESQL_PORT")
	PostgreSQLDb       = os.Getenv("POSTGRESQL_DB")
	PostgreSQLUser     = os.Getenv("POSTGRESQL_USERNAME")
	PostgreSQLPassword = os.Getenv("POSTGRESQL_PASSWORD")
	PostgreSQLSSLMode  = os.Getenv("POSTGRESQL_SSLMODE")

	HttpServerAddress = os.Getenv("HTTP_ADDRESS")
	GRPCServerAddress = os.Getenv("GRPC_ADDRESS")

	DescribeIntervalHours      = os.Getenv("DESCRIBE_INTERVAL_HOURS")
	CostDiscoveryIntervalHours = os.Getenv("COST_DISCOVERY_INTERVAL_HOURS")
	DescribeTimeoutHours       = os.Getenv("DESCRIBE_TIMEOUT_HOURS")
	CheckupIntervalHours       = os.Getenv("CHECKUP_INTERVAL_HOURS")
	MustSummarizeIntervalHours = os.Getenv("MUST_SUMMARIZE_INTERVAL_HOURS")
	CoreBaseURL            = os.Getenv("Core_BASE_URL")
	ComplianceBaseURL          = os.Getenv("COMPLIANCE_BASE_URL")
	IntegrationBaseURL         = os.Getenv("INTEGRATION_BASE_URL")
	
	EsSinkBaseURL              = os.Getenv("ESSINK_BASEURL")
	AuthGRPCURI                = os.Getenv("AUTH_GRPC_URI")

	KeyARN                       = os.Getenv("VAULT_KEY_ID")
	KeyRegion                    = os.Getenv("KMS_ACCOUNT_REGION")
	DescribeLocalJobEndpoint     = os.Getenv("DESCRIBE_JOB_ENDPOINT_LOCAL")
	DescribeLocalDeliverEndpoint = os.Getenv("DESCRIBE_DELIVER_ENDPOINT_LOCAL")
	DescribeExternalEndpoint     = os.Getenv("DESCRIBE_DELIVER_ENDPOINT_EXTERNAL")

	DoDeleteOldResources  = os.Getenv("DO_DELETE_OLD_RESOURCES")
	OperationModeConfig   = os.Getenv("OPERATION_MODE_CONFIG")
	DoProcessReceivedMsgs = os.Getenv("DO_PROCESS_RECEIVED_MSGS")

	MaxConcurrentCall = os.Getenv("MAX_CONCURRENT_CALL")
)

func SchedulerCommand() *cobra.Command {
	var id string
	var conf config2.SchedulerConfig
	config.ReadFromEnv(&conf, nil)

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
			s, err := InitializeScheduler(
				id,
				conf,
				PostgreSQLUser,
				PostgreSQLPassword,
				PostgreSQLHost,
				PostgreSQLPort,
				PostgreSQLDb,
				PostgreSQLSSLMode,
				HttpServerAddress,
				DescribeTimeoutHours,
				CheckupIntervalHours,
				MustSummarizeIntervalHours,
				cmd.Context(),
			)
			if err != nil {
				return err
			}

			defer s.Stop()

			return s.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&id, "id", "", "The scheduler id")

	return cmd
}
