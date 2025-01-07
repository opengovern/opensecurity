package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/jobs/post-install-job/db"
	"github.com/opengovern/opencomply/services/core/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type Migration struct {
}

func (m Migration) AttachmentFolderPath() string {
	return "/core-migration"
}

func (m Migration) IsGitBased() bool {
	return false
}

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
	if err := CoreMigration(conf, logger, m.AttachmentFolderPath()+"/metadata.json"); err != nil {
		return err
	}
	return nil
}

func CoreMigration(conf config.MigratorConfig, logger *zap.Logger, metadataFilePath string) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "core",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{ORM: orm}

	content, err := os.ReadFile(metadataFilePath)
	if err != nil {
		return err
	}

	var metadata []models.ConfigMetadata
	err = json.Unmarshal(content, &metadata)
	if err != nil {
		return err
	}

	for _, obj := range metadata {
		err := dbm.ORM.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&obj).Error
		if err != nil {
			return err
		}
	}
	return nil
}
