package inventory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opensecurity/jobs/post-install-job/utils"
	"github.com/opengovern/opensecurity/pkg/types"
	coreClient "github.com/opengovern/opensecurity/services/core/client"
	integrationClient "github.com/opengovern/opensecurity/services/integration/client"
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

	err = ExtractQueryViews(ctx, logger, dbm, conf, config.QueryViewsGitPath)
	if err != nil {
		return err
	}
	err = populateQueries(ctx, logger, dbm, conf)
	if err != nil {
		return err
	}

	return nil
}

func ExtractQueryViews(ctx context.Context, logger *zap.Logger, dbm db.Database, conf config.MigratorConfig, viewsPath string) error {
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
		tx.Model(&models.QueryParameter{}).Where("1=1").Unscoped().Delete(&models.QueryParameter{})
		tx.Model(&models.NamedQuery{}).Where("1=1").Unscoped().Delete(&models.NamedQuery{})
		tx.Model(&models.NamedQueryTag{}).Where("1=1").Unscoped().Delete(&models.NamedQueryTag{})
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
			for _, tag := range qv.Tags {
				err = tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "key"}, {Name: "query_view_id"}}, // key columns
					DoNothing: true,
				}).Create(&tag).Error
				if err != nil {
					return fmt.Errorf("failure in control tag insert: %v", err)
				}
			}
		}
		return nil
	})

	mClient := coreClient.NewCoreServiceClient(conf.Core.BaseURL)
	err = mClient.ReloadViews(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole})
	if err != nil {
		logger.Error("failed to reload views", zap.Error(err))
		return fmt.Errorf("failed to reload views: %s", err.Error())
	}

	return err
}

func populateQueries(ctx context.Context, logger *zap.Logger, dbm db.Database, conf config.MigratorConfig) error {
	iClient := integrationClient.NewIntegrationServiceClient(conf.Integration.BaseURL)
	pluginTables, err := iClient.GetPluginsTables(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole})
	if err != nil {
		logger.Error("failed to get plugin tables", zap.Error(err))
		return nil
	}
	tablesPluginMap := make(map[string]string)
	for _, p := range pluginTables {
		for _, t := range p.Tables {
			tablesPluginMap[t] = p.PluginID
		}
	}

	err = dbm.ORM.Transaction(func(tx *gorm.DB) error {
		err := filepath.Walk(config.QueriesGitPath, func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
				return populateFinderItem(logger, tx, path, info, tablesPluginMap)
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

func populateFinderItem(logger *zap.Logger, tx *gorm.DB, path string, info fs.FileInfo, tablesPluginMap map[string]string) error {
	id := strings.TrimSuffix(info.Name(), ".yaml")

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var item NamedQuery
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

	listOfTables, err := utils.ExtractTableRefsFromPolicy("sql", item.Query)
	if err != nil {
		logger.Error("failed to extract table refs from query", zap.String("query-id", id), zap.Error(err))
	}
	if len(integrationTypes) == 0 {
		integrationTypesMap := make(map[string]bool)
		for _, t := range listOfTables {
			if v, ok := tablesPluginMap[t]; ok {
				integrationTypesMap[v] = true
			}
		}
		for it := range integrationTypesMap {
			integrationTypes = append(integrationTypes, it)
		}
	}

	namedQuery := models.NamedQuery{
		ID:               id,
		IntegrationTypes: integrationTypes,
		Title:            item.Title,
		Description:      item.Description,
		IsBookmarked:     isBookmarked,
		QueryID:          &id,
	}

	parameters, err := utils.ExtractParameters("sql", item.Query)
	if err != nil {
		logger.Error("extract control failed: failed to extract parameters from query", zap.String("control-id", namedQuery.ID), zap.Error(err))
		return nil
	}
	queryParams := []models.QueryParameter{}
	for _, p := range parameters {
		queryParams = append(queryParams, models.QueryParameter{
			QueryID: namedQuery.ID,
			Key:     p,
		})
	}

	query := models.Query{
		ID:             namedQuery.ID,
		QueryToExecute: item.Query,
		ListOfTables:   listOfTables,
		Engine:         "sql",
		Parameters:     queryParams,
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
	}).Create(namedQuery).Error
	if err != nil {
		logger.Error("failure in insert query", zap.Error(err))
		return err
	}

	if len(tags) > 0 {
		for _, tag := range tags {
			err = tx.Model(&models.NamedQueryTag{}).Create(&tag).Error
			if err != nil {
				logger.Error("failure in insert tags", zap.Error(err))
				return err
			}
		}
	}

	for _, p := range item.Parameters {
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}, {Name: "control_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"value": gorm.Expr("CASE WHEN policy_parameter_values.value = '' THEN ? ELSE policy_parameter_values.value END", p.Value),
			}),
		}).Create(&models.PolicyParameterValues{
			Key:       p.Key,
			ControlID: "",
			Value:     p.Value,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}
