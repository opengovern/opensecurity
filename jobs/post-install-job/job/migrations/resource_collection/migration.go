package resource_collection

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"github.com/opengovern/og-util/pkg/model"
	"github.com/opengovern/og-util/pkg/opengovernance-es-sdk"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/jobs/post-install-job/db"
	"github.com/opengovern/opencomply/services/core/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResourceCollection struct {
	ID          string                                    `json:"id" yaml:"id"`
	Name        string                                    `json:"name" yaml:"name"`
	Tags        map[string][]string                       `json:"tags" yaml:"tags"`
	Filters     []opengovernance.ResourceCollectionFilter `json:"filters" yaml:"filters"`
	Description string                                    `json:"description" yaml:"description"`
	Status      models.ResourceCollectionStatus        `json:"status" yaml:"status"`
}

type Migration struct {
}

func (m Migration) AttachmentFolderPath() string {
	return config.ResourceCollectionGitPath
}
func (m Migration) IsGitBased() bool {
	return true
}

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
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

	resourceCollections, err := ExtractResourceCollections(m.AttachmentFolderPath())
	if err != nil {
		logger.Error("failed to extract resource collections", zap.Error(err))
		return err
	}

	err = dbm.ORM.Transaction(func(tx *gorm.DB) error {
		currentRCs := make([]models.ResourceCollection, 0)
		err := tx.Model(&models.ResourceCollection{}).Find(&currentRCs).Error
		if err != nil {
			logger.Error("failed to get current resource collections", zap.Error(err))
			return err
		}
		currentRcMap := make(map[string]models.ResourceCollection)
		for _, rc := range currentRCs {
			currentRcMap[rc.ID] = rc
		}

		tx.Model(&models.ResourceCollection{}).Where("1=1").Unscoped().Delete(&models.ResourceCollection{})
		tx.Model(&models.ResourceCollectionTag{}).Where("1=1").Unscoped().Delete(&models.ResourceCollectionTag{})
		for _, resourceCollection := range resourceCollections {
			filtersJson, err := json.Marshal(resourceCollection.Filters)
			if err != nil {
				logger.Error("failed to marshal filters", zap.Error(err))
				return err
			}

			jsonb := pgtype.JSONB{}
			err = jsonb.Set(filtersJson)
			if err != nil {
				logger.Error("failed to set jsonb", zap.Error(err))
				return err
			}

			createdAt := time.Now()
			if currentRc, ok := currentRcMap[resourceCollection.ID]; ok {
				createdAt = currentRc.Created
				if createdAt.IsZero() || createdAt.Year() == 1 {
					createdAt = time.Now()
				}
			}
			if resourceCollection.Status == "" {
				resourceCollection.Status = models.ResourceCollectionStatusActive
			}

			dbResourceCollection := models.ResourceCollection{
				ID:          resourceCollection.ID,
				Name:        resourceCollection.Name,
				FiltersJson: jsonb,
				Description: resourceCollection.Description,
				Status:      resourceCollection.Status,
				Created:     createdAt,
			}
			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&dbResourceCollection).Error
			if err != nil {
				logger.Error("failed to create resource collection", zap.Error(err))
				return err
			}

			for key, values := range resourceCollection.Tags {
				err = tx.Clauses(clause.OnConflict{
					DoNothing: true,
				}).Create(&models.ResourceCollectionTag{
					Tag: model.Tag{
						Key:   key,
						Value: values,
					},
					ResourceCollectionID: resourceCollection.ID,
				}).Error
				if err != nil {
					logger.Error("failed to create resource collection tag", zap.Error(err))
					return err
				}
			}
		}
		return nil
	})

	return nil
}
