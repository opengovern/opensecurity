package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/jobs/config-manager/config"
	"github.com/opengovern/opensecurity/jobs/config-manager/db"
	"github.com/opengovern/opensecurity/services/core/db/models"
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
	// Create three default widgets
	var widgets []models.Widget

	widgets = append(widgets, models.Widget{

		ID:          "integration",
		UserID:      "system",
		Title:       "Integrations",
		Description: "",
		WidgetType:  "integration",
		WidgetProps: func() pgtype.JSONB {
			var jsonb pgtype.JSONB
			// Ensure WidgetProps is never undefined
			_ = jsonb.Set(map[string]interface{}{})
			return jsonb
		}(),
		RowSpan:      8,
		ColumnSpan:   1,
		ColumnOffset: 3,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsPublic:     true,
	})
	widgets = append(widgets, models.Widget{

		ID:          "shortcut",
		UserID:      "system",
		Title:       "Shortcuts",
		Description: "",
		WidgetType:  "shortcut",
		WidgetProps: func() pgtype.JSONB {
			var jsonb pgtype.JSONB
			// Ensure WidgetProps is never undefined
			_ = jsonb.Set(map[string]interface{}{})
			return jsonb
		}(),
		RowSpan:      2,
		ColumnSpan:   3,
		ColumnOffset: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsPublic:     true,
	})
	widgets = append(widgets, models.Widget{

		ID:          "sre",
		UserID:      "system",
		Title:       "SRE",
		Description: "",
		WidgetType:  "sre",
		WidgetProps: func() pgtype.JSONB {
			var jsonb pgtype.JSONB
			// Ensure WidgetProps is never undefined
			_ = jsonb.Set(map[string]interface{}{})
			return jsonb
		}(),
		RowSpan:      2,
		ColumnSpan:   3,
		ColumnOffset: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsPublic:     true,
	})
	err = dbm.ORM.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"title", "description", "widget_type", "widget_props",
			"row_span", "column_span", "column_offset", "is_public",
			"user_id", "updated_at",
		}),
	}).Create(&widgets).Error
	if err != nil {
		return err
	}

	return nil
}
