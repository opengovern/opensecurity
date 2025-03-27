package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opensecurity/jobs/post-install-job/utils"
	"github.com/opengovern/opensecurity/pkg/types"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/opengovern/og-util/pkg/model"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opensecurity/jobs/post-install-job/config"
	"github.com/opengovern/opensecurity/jobs/post-install-job/db"
	"github.com/opengovern/opensecurity/services/core/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ResourceType struct {
	ResourceName         string
	Category             []string
	ResourceLabel        string
	ServiceName          string
	ListDescriber        string
	GetDescriber         string
	TerraformName        []string
	TerraformNameString  string `json:"-"`
	TerraformServiceName string
	Discovery            string
	IgnoreSummarize      bool
	SteampipeTable       string
	Model                string
}

type Migration struct {
}

func (m Migration) IsGitBased() bool {
	return false
}
func (m Migration) AttachmentFolderPath() string {
	return "/inventory-data-config"
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

	awsResourceTypesContent, err := os.ReadFile(path.Join(m.AttachmentFolderPath(), "aws-resource-types.json"))
	if err != nil {
		return err
	}
	azureResourceTypesContent, err := os.ReadFile(path.Join(m.AttachmentFolderPath(), "azure-resource-types.json"))
	if err != nil {
		return err
	}
	var awsResourceTypes []ResourceType
	var azureResourceTypes []ResourceType
	if err := json.Unmarshal(awsResourceTypesContent, &awsResourceTypes); err != nil {
		return err
	}
	if err := json.Unmarshal(azureResourceTypesContent, &azureResourceTypes); err != nil {
		return err
	}

	err = dbm.ORM.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.ResourceType{}).Where("integration_type = ?", integration.Type("aws_cloud_account")).Unscoped().Delete(&models.ResourceType{}).Error
		if err != nil {
			logger.Error("failed to delete aws resource types", zap.Error(err))
			return err
		}

		for _, resourceType := range awsResourceTypes {
			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&models.ResourceType{
				IntegrationType: "aws_cloud_account",
				ResourceType:    resourceType.ResourceName,
				ResourceLabel:   resourceType.ResourceLabel,
				ServiceName:     strings.ToLower(resourceType.ServiceName),
				DoSummarize:     !resourceType.IgnoreSummarize,
			}).Error
			if err != nil {
				logger.Error("failed to create aws resource type", zap.Error(err))
				return err
			}

			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&models.ResourceTypeTag{
				Tag: model.Tag{
					Key:   "category",
					Value: resourceType.Category,
				},
				ResourceType: resourceType.ResourceName,
			}).Error
			if err != nil {
				logger.Error("failed to create aws resource type tag", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failure in aws transaction: %v", err)
	}

	err = dbm.ORM.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.ResourceType{}).Where("integration_type = ?", integration.Type("azure_subscription")).Unscoped().Delete(&models.ResourceType{}).Error
		if err != nil {
			logger.Error("failed to delete azure resource types", zap.Error(err))
			return err
		}
		for _, resourceType := range azureResourceTypes {
			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&models.ResourceType{
				IntegrationType: "azure_subscription",
				ResourceType:    resourceType.ResourceName,
				ResourceLabel:   resourceType.ResourceLabel,
				ServiceName:     strings.ToLower(resourceType.ServiceName),
				DoSummarize:     !resourceType.IgnoreSummarize,
			}).Error
			if err != nil {
				logger.Error("failed to create azure resource type", zap.Error(err))
				return err
			}

			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&models.ResourceTypeTag{
				Tag: model.Tag{
					Key:   "category",
					Value: resourceType.Category,
				},
				ResourceType: resourceType.ResourceName,
			}).Error
			if err != nil {
				logger.Error("failed to create azure resource type tag", zap.Error(err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failure in azure transaction: %v", err)
	}

	err = ExtractQueryViews(ctx, logger, dbm, config.QueryViewsGitPath)

	return nil
}

func ExtractQueryViews(ctx context.Context, logger *zap.Logger, dbm db.Database, viewsPath string) error {
	var queries []models.Query
	var queryViews []models.QueryView
	err := filepath.WalkDir(viewsPath, func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			logger.Error("failed to read query view", zap.String("path", path), zap.Error(err))
			return err
		}

		var obj QueryView
		err = yaml.Unmarshal(content, &obj)
		if err != nil {
			logger.Error("failed to unmarshal query view", zap.String("path", path), zap.Error(err))
			return nil
		}

		qv := models.QueryView{
			ID:          obj.ID,
			Title:       obj.Title,
			Description: obj.Description,
		}

		listOfTables, err := utils.ExtractTableRefsFromPolicy(types.PolicyLanguageSQL, obj.Query)
		if err != nil {
			logger.Error("failed to extract table refs from query", zap.String("query-id", obj.ID), zap.Error(err))
		}

		q := models.Query{
			ID:             obj.ID,
			QueryToExecute: obj.Query,
			ListOfTables:   listOfTables,
			Engine:         "sql",
		}

		queries = append(queries, q)
		qv.QueryID = &obj.ID

		queryViews = append(queryViews, qv)

		return nil
	})
	if err != nil {
		return err
	}

	err = dbm.ORM.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Model(&models.QueryView{}).Where("1=1").Unscoped().Delete(&models.QueryView{})
		tx.Model(&models.Query{}).Where("1=1").Unscoped().Delete(&models.Query{})

		for _, q := range queries {
			q.QueryViews = nil
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}}, // key column
				DoNothing: true,
			}).Create(&q).Error
			if err != nil {
				return err
			}
		}

		for _, qv := range queryViews {
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}}, // key column
				DoNothing: true,
			}).Create(&qv).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
