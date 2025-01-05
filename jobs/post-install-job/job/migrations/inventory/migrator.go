package inventory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	utils "github.com/opengovern/opencomply/jobs/post-install-job/utils"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/opengovern/og-util/pkg/model"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/jobs/post-install-job/db"
	integration_type "github.com/opengovern/opencomply/services/integration/integration-type"
	"github.com/opengovern/opencomply/services/core/db/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var QueryParameters []models.PolicyParameterValues

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
		err := tx.Model(&models.ResourceType{}).Where("integration_type = ?", integration_type.IntegrationTypeAWSAccount).Unscoped().Delete(&models.ResourceType{}).Error
		if err != nil {
			logger.Error("failed to delete aws resource types", zap.Error(err))
			return err
		}

		for _, resourceType := range awsResourceTypes {
			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&models.ResourceType{
				IntegrationType: integration_type.IntegrationTypeAWSAccount,
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
		err := tx.Model(&models.ResourceType{}).Where("integration_type = ?", integration_type.IntegrationTypeAzureSubscription).Unscoped().Delete(&models.ResourceType{}).Error
		if err != nil {
			logger.Error("failed to delete azure resource types", zap.Error(err))
			return err
		}
		for _, resourceType := range azureResourceTypes {
			err = tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&models.ResourceType{
				IntegrationType: integration_type.IntegrationTypeAzureSubscription,
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

	err = populateQueries(logger, dbm)
	if err != nil {
		return err
	}

	err = dbm.ORM.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, obj := range QueryParameters {
			err := tx.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&obj).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.Error("failed to insert query params", zap.Error(err))
		return err
	}

	return nil
}

func populateQueries(logger *zap.Logger, db db.Database) error {
	err := db.ORM.Transaction(func(tx *gorm.DB) error {

		tx.Model(&models.NamedQuery{}).Where("1=1").Unscoped().Delete(&models.NamedQuery{})
		tx.Model(&models.NamedQueryTag{}).Where("1=1").Unscoped().Delete(&models.NamedQueryTag{})
		tx.Model(&models.QueryParameter{}).Where("1=1").Unscoped().Delete(&models.QueryParameter{})
		tx.Model(&models.Query{}).Where("1=1").Unscoped().Delete(&models.Query{})

		err := filepath.Walk(config.QueriesGitPath, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
				return populateFinderItem(logger, tx, path, info)
			}
			return nil
		})
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			logger.Error("failed to get queries", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func populateFinderItem(logger *zap.Logger, tx *gorm.DB, path string, info fs.FileInfo) error {
	id := strings.TrimSuffix(info.Name(), ".yaml")

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var item NamedPolicy
	err = yaml.Unmarshal(content, &item)
	if err != nil {
		logger.Error("failure in unmarshal", zap.String("path", path), zap.Error(err))
		return err
	}

	if item.ID != "" {
		id = item.ID
	}

	var integrationTypes []string
	for _, c := range item.IntegrationTypes {
		integrationTypes = append(integrationTypes, string(c))
	}

	isBookmarked := false
	tags := make([]models.NamedQueryTag, 0, len(item.Tags))
	for k, v := range item.Tags {
		if k == "platform_queries_bookmark" {
			isBookmarked = true
		}
		tag := models.NamedQueryTag{
			NamedQueryID: id,
			Tag: model.Tag{
				Key:   k,
				Value: v,
			},
		}
		tags = append(tags, tag)
	}

	dbMetric := models.NamedQuery{
		ID:               id,
		IntegrationTypes: integrationTypes,
		Title:            item.Title,
		Description:      item.Description,
		IsBookmarked:     isBookmarked,
		QueryID:          &id,
	}
	queryParams := []models.QueryParameter{}
	for _, qp := range item.Policy.Parameters {
		queryParams = append(queryParams, models.QueryParameter{
			Key:      qp.Key,
			Required: qp.Required,
			QueryID:  dbMetric.ID,
		})
		if qp.DefaultValue != "" {
			queryParamObj := models.PolicyParameterValues{
				Key:   qp.Key,
				Value: qp.DefaultValue,
			}
			QueryParameters = append(QueryParameters, queryParamObj)
		}
	}
	listOfTables, err := utils.ExtractTableRefsFromPolicy("sql", item.Policy.QueryToExecute)
	if err != nil {
		logger.Error("failed to extract table refs from query", zap.String("query-id", dbMetric.ID), zap.Error(err))
		listOfTables = item.Policy.ListOfTables
	}
	query := models.Query{
		ID:             dbMetric.ID,
		QueryToExecute: item.Policy.QueryToExecute,
		PrimaryTable:   item.Policy.PrimaryTable,
		ListOfTables:   listOfTables,
		Engine:         item.Policy.Engine,
		Parameters:     queryParams,
		Global:         item.Policy.Global,
	}
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // key column
		DoNothing: true,
	}).Create(&query).Error
	if err != nil {
		logger.Error("failure in Creating Policy", zap.String("query_id", id), zap.Error(err))
		return err
	}
	for _, param := range query.Parameters {
		err = tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}, {Name: "query_id"}}, // key columns
			DoNothing: true,
		}).Create(&param).Error
		if err != nil {
			return fmt.Errorf("failure in query parameter insert: %v", err)
		}
	}

	err = tx.Model(&models.NamedQuery{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // key column
		DoNothing: true,                          // column needed to be updated
	}).Create(dbMetric).Error
	if err != nil {
		logger.Error("failure in insert query", zap.Error(err))
		return err
	}

	// logger.Info("parsed the tags", zap.String("id", id), zap.Any("tags", tags))

	if len(tags) > 0 {
		for _, tag := range tags {
			err = tx.Model(&models.NamedQueryTag{}).Create(&tag).Error
			if err != nil {
				logger.Error("failure in insert tags", zap.Error(err))
				return err
			}
		}
	}

	return nil
}
