package es_sink

import (
	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/jq"
	"github.com/opengovern/og-util/pkg/koanf"
	es "github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/opensecurity/services/es-sink/api"
	"github.com/opengovern/opensecurity/services/es-sink/config"
	"github.com/opengovern/opensecurity/services/es-sink/grpcApi"
	"github.com/opengovern/opensecurity/services/es-sink/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Command() *cobra.Command {
	cnf := koanf.Provide("essink", config.EsSinkConfig{})

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			logger = logger.Named("es-sink")

			cmd.SilenceUsage = true

			esClient, err := es.NewClient(es.ClientConfig{
				Addresses:     []string{cnf.ElasticSearch.Address},
				Username:      &cnf.ElasticSearch.Username,
				Password:      &cnf.ElasticSearch.Password,
				IsOpenSearch:  &cnf.ElasticSearch.IsOpenSearch,
				IsOnAks:       &cnf.ElasticSearch.IsOnAks,
				AwsRegion:     &cnf.ElasticSearch.AWSRegion,
				AssumeRoleArn: &cnf.ElasticSearch.AssumeRoleARN,
				ExternalID:    &cnf.ElasticSearch.ExternalID,
			})

			if err != nil {
				logger.Error("failed to create es client", zap.Error(err))
				return err
			}

			nats, err := jq.New(cnf.NATS.URL, logger)
			if err != nil {
				logger.Error("failed to create nats client", zap.Error(err))
				return err
			}

			sinkService, err := service.NewEsSinkService(ctx, logger, esClient, nats)
			if err != nil {
				logger.Error("failed to create es sink service", zap.Error(err))
				return err
			}

			go sinkService.Start(ctx)

			grpcServer, err := grpcApi.NewGRPCSinkServer(logger, sinkService, cnf.Grpc.Address)
			if err != nil {
				logger.Error("failed to create grpc server", zap.Error(err))
				return err
			}

			go grpcServer.Start()

			return httpserver.RegisterAndStart(
				ctx,
				logger,
				cnf.Http.Address,
				api.New(logger, sinkService),
			)
		},
	}

	return cmd
}
