package demo_importer

import (
	"context"
	"fmt"
	"os"

	"github.com/opengovern/og-util/pkg/config"
	"github.com/opengovern/og-util/pkg/postgres"
	db2 "github.com/opengovern/opensecurity/jobs/demo-importer-job/db"
	"github.com/opengovern/opensecurity/jobs/demo-importer-job/types"
	"github.com/opengovern/opensecurity/jobs/demo-importer-job/worker"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func Command() *cobra.Command {
	var (
		cnf types.DemoImporterConfig
	)
	config.ReadFromEnv(&cnf, nil)
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("running", zap.String("es_address", cnf.ElasticSearch.Address), zap.String("es_arn", cnf.ElasticSearch.AssumeRoleArn))

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			err = os.MkdirAll(types.DemoDataPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failure creating path: %w", err)
			}

		

			cmd.SilenceUsage = true

			// metadataClient := client.NewCoreServiceClient(cnf.Core.BaseURL)

			// s3Url := cnf.DemoDataS3URL
			// value, err := metadataClient.GetConfigMetadata(&httpclient.Context{
			// 	UserRole: api.AdminRole,
			// }, models.DemoDataS3URL)
			// if err == nil && len(value.GetValue().(string)) > 0 {
			// 	s3Url = value.GetValue().(string)
			// } else if err != nil {
			// 	logger.Error("failed to get demo data s3 url from metadata", zap.Error(err))
			// }

			// logger.Info("Downloading file", zap.String("address", s3Url))
			// filePath, err := fetch.DownloadS3Object(cnf.DemoDataS3URL)
			// if err != nil {
			// 	return err
			// }

			// logger.Info("File Downloaded", zap.String("file", filePath))

			// decryptedData, err := fetch.DecryptString(filePath, cnf.OpensslPassword)

			// logger.Info("Successfully decrypted", zap.String("file", filePath))

			// decryptedFile := types.DemoDecryptedDataFilePath
			// err = os.WriteFile(decryptedFile, decryptedData, 0644)
			// if err != nil {
			// 	return err
			// }

			// logger.Info("Successfully decrypted file written", zap.String("file", filePath))

			// file, err := os.Open(decryptedFile)
			// if err != nil {
			// 	return err
			// }
			// defer file.Close()

			// err = fetch.Unzip(file)
			// if err != nil {
			// 	return fmt.Errorf("failure while unzipping file: %w", err)
			// }

			// logger.Info("Successfully unzipped", zap.String("file", filePath))
			//err = worker.ImportSQLFiles(cnf, types.PostgresqlBackupPath)
			//if err != nil {
			//	return fmt.Errorf("failure while importing sql files to postgres: %w", err)
			//}
			psqlMigratorCfg := postgres.Config{
				Host:    cnf.PostgreSQL.Host,
				Port:    cnf.PostgreSQL.Port,
				User:    cnf.PostgreSQL.Username,
				Passwd:  cnf.PostgreSQL.Password,
				DB:      "migrator",
				SSLMode: cnf.PostgreSQL.SSLMode,
			}
			migratorOrm, err := postgres.NewClient(&psqlMigratorCfg, logger)
			if err != nil {
				return fmt.Errorf("new postgres client: %w", err)
			}

			migratorDb := db2.Database{
				ORM: migratorOrm,
			}
			err = migratorDb.Initialize()
			if err != nil {
				return fmt.Errorf("failure while initializing database: %w", err)
			}

			err = worker.ImportJob(ctx, logger, migratorDb)
			if err != nil {
				return fmt.Errorf("failure while importing es indices: %w", err)
			}

			return nil
		},
	}

	return cmd
}
