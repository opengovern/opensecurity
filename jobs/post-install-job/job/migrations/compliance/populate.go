package compliance

import (
	"context"
	"errors"
	"fmt"
	"github.com/goccy/go-yaml"
	authApi "github.com/opengovern/og-util/pkg/api"
	"github.com/opengovern/og-util/pkg/httpclient"
	"github.com/opengovern/og-util/pkg/model"
	"github.com/opengovern/opencomply/jobs/post-install-job/utils"
	coreClient "github.com/opengovern/opencomply/services/core/client"
	"github.com/opengovern/opencomply/services/core/db/models"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/opencomply/jobs/post-install-job/config"
	"github.com/opengovern/opencomply/services/compliance/db"
	"go.uber.org/zap"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var QueryParameters []models.PolicyParameterValues

type Migration struct {
}

func (m Migration) IsGitBased() bool {
	return true
}

func (m Migration) AttachmentFolderPath() string {
	return ""
}

func (m Migration) Run(ctx context.Context, conf config.MigratorConfig, logger *zap.Logger) error {
	orm, err := postgres.NewClient(&postgres.Config{
		Host:    conf.PostgreSQL.Host,
		Port:    conf.PostgreSQL.Port,
		User:    conf.PostgreSQL.Username,
		Passwd:  conf.PostgreSQL.Password,
		DB:      "compliance",
		SSLMode: conf.PostgreSQL.SSLMode,
	}, logger)
	if err != nil {
		return fmt.Errorf("new postgres client: %w", err)
	}
	dbm := db.Database{Orm: orm}

	ormCore, err := postgres.NewClient(&postgres.Config{
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
	dbCore := db.Database{Orm: ormCore}

	p := GitParser{
		logger:             logger,
		frameworksChildren: make(map[string][]string),
		controlsPolicies:   make(map[string]db.Policy),
		namedPolicies:      make(map[string]NamedQuery),
	}
	if err := p.ExtractCompliance(config.ComplianceGitPath, config.ControlEnrichmentGitPath); err != nil {
		logger.Error("failed to extract controls and benchmarks", zap.Error(err))
		return err
	}
	if err := p.ExtractQueryViews(config.QueryViewsGitPath); err != nil {
		logger.Error("failed to extract query views", zap.Error(err))
		return err
	}

	logger.Info("extracted controls, benchmarks and query views", zap.Int("controls", len(p.controls)), zap.Int("benchmarks", len(p.benchmarks)), zap.Int("query_views", len(p.policies)))

	loadedQueries := make(map[string]bool)
	err = dbm.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Model(&db.BenchmarkChild{}).Where("1=1").Unscoped().Delete(&db.BenchmarkChild{})
		tx.Model(&db.BenchmarkControls{}).Where("1=1").Unscoped().Delete(&db.BenchmarkControls{})
		tx.Model(&db.Benchmark{}).Where("1=1").Unscoped().Delete(&db.Benchmark{})
		tx.Model(&db.Control{}).Where("1=1").Unscoped().Delete(&db.Control{})
		tx.Model(&db.PolicyParameter{}).Where("1=1").Unscoped().Delete(&db.PolicyParameter{})
		tx.Model(&db.Policy{}).Where("1=1").Unscoped().Delete(&db.Policy{})

		for _, obj := range p.policies {
			obj.Controls = nil
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}}, // key column
				DoNothing: true,
			}).Create(&obj).Error
			if err != nil {
				return err
			}
			for _, param := range obj.Parameters {
				err = tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "key"}, {Name: "policy_id"}}, // key columns
					DoNothing: true,
				}).Create(&param).Error
				if err != nil {
					return fmt.Errorf("failure in query parameter insert: %v", err)
				}
			}
			loadedQueries[obj.ID] = true
		}
		return nil
	})
	if err != nil {
		logger.Error("failed to insert policies", zap.Error(err))
		return err
	}

	err = dbCore.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, obj := range p.policyParamValues {
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

	missingQueries := make(map[string]bool)
	err = dbm.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		for _, obj := range p.controls {
			obj.Benchmarks = nil
			if obj.PolicyID != nil && !loadedQueries[*obj.PolicyID] {
				missingQueries[*obj.PolicyID] = true
				logger.Info("query not found", zap.String("query_id", *obj.PolicyID))
				continue
			}
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).Create(&obj).Error
			if err != nil {
				return err
			}
			for _, tag := range obj.Tags {
				err = tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "key"}, {Name: "control_id"}}, // key columns
					DoNothing: true,
				}).Create(&tag).Error
				if err != nil {
					return fmt.Errorf("failure in control tag insert: %v", err)
				}
			}
		}

		for _, obj := range p.benchmarks {
			obj.Children = nil
			obj.Controls = nil
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}}, // key column
				DoNothing: true,
			}).Create(&obj).Error
			if err != nil {
				return err
			}
			for _, tag := range obj.Tags {
				err = tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "key"}, {Name: "benchmark_id"}}, // key columns
					DoNothing: true,
				}).Create(&tag).Error
				if err != nil {
					return fmt.Errorf("failure in benchmark tag insert: %v", err)
				}
			}
		}

		for _, obj := range p.benchmarks {
			for _, child := range obj.Children {
				err := tx.Clauses(clause.OnConflict{
					DoNothing: true,
				}).Create(&db.BenchmarkChild{
					BenchmarkID: obj.ID,
					ChildID:     child.ID,
				}).Error
				if err != nil {
					logger.Error("inserted controls and benchmarks", zap.Error(err))
					return err
				}
			}

			for _, control := range obj.Controls {
				if control.PolicyID != nil && !loadedQueries[*control.PolicyID] {
					continue
				}
				err := tx.Clauses(clause.OnConflict{
					DoNothing: true,
				}).Create(&db.BenchmarkControls{
					BenchmarkID: obj.ID,
					ControlID:   control.ID,
				}).Error
				if err != nil {
					logger.Info("inserted controls and benchmarks", zap.Error(err))
					return err
				}
			}
		}

		missingQueriesList := make([]string, 0, len(missingQueries))
		for query := range missingQueries {
			missingQueriesList = append(missingQueriesList, query)
		}
		if len(missingQueriesList) > 0 {
			logger.Warn("missing policies", zap.Strings("policies", missingQueriesList))
		}
		return nil
	})

	if err != nil {
		logger.Info("inserted controls and benchmarks", zap.Error(err))
		return err
	}

	loadedQueryViewsQueries := make(map[string]bool)
	missingQueryViewsQueries := make(map[string]bool)
	err = dbCore.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Model(&models.QueryView{}).Where("1=1").Unscoped().Delete(&models.QueryView{})
		tx.Model(&models.QueryParameter{}).Where("1=1").Unscoped().Delete(&models.QueryParameter{})
		tx.Model(&models.Query{}).Where("1=1").Unscoped().Delete(&models.Query{})
		for _, obj := range p.coreServiceQueries {
			obj.QueryViews = nil
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}}, // key column
				DoNothing: true,
			}).Create(&obj).Error
			if err != nil {
				return err
			}
			for _, param := range obj.Parameters {
				err = tx.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "key"}, {Name: "query_id"}}, // key columns
					DoNothing: true,
				}).Create(&param).Error
				if err != nil {
					return fmt.Errorf("failure in query parameter insert: %v", err)
				}
			}
			loadedQueryViewsQueries[obj.ID] = true
		}

		return nil
	})
	if err != nil {
		logger.Error("failed to insert query views", zap.Error(err))
		return err
	}

	err = dbCore.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, obj := range p.queryViews {
			if obj.QueryID != nil && !loadedQueryViewsQueries[*obj.QueryID] {
				missingQueryViewsQueries[*obj.QueryID] = true
				logger.Info("query not found", zap.String("query_id", *obj.QueryID))
				continue
			}
			err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).Create(&obj).Error
			if err != nil {
				logger.Error("error while inserting query view", zap.Error(err))
				return err
			}
			for _, tag := range obj.Tags {
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
	if err != nil {
		logger.Error("failed to insert query views", zap.Error(err))
		return err
	}

	err = populateQueries(logger, dbCore)
	if err != nil {
		return err
	}

	err = dbCore.Orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

	mClient := coreClient.NewCoreServiceClient(conf.Core.BaseURL)
	err = mClient.ReloadViews(&httpclient.Context{Ctx: ctx, UserRole: authApi.AdminRole})
	if err != nil {
		logger.Error("failed to reload views", zap.Error(err))
		return fmt.Errorf("failed to reload views: %s", err.Error())
	}

	return nil
}

func populateQueries(logger *zap.Logger, db db.Database) error {
	err := db.Orm.Transaction(func(tx *gorm.DB) error {

		tx.Model(&models.NamedQuery{}).Where("1=1").Unscoped().Delete(&models.NamedQuery{})
		tx.Model(&models.NamedQueryTag{}).Where("1=1").Unscoped().Delete(&models.NamedQueryTag{})

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

	dbMetric := models.NamedQuery{
		ID:               id,
		IntegrationTypes: integrationTypes,
		Title:            item.Title,
		Description:      item.Description,
		IsBookmarked:     isBookmarked,
		QueryID:          &id,
	}
	queryParams := []models.QueryParameter{}
	listOfTables, err := utils.ExtractTableRefsFromPolicy("sql", item.Query)
	if err != nil {
		logger.Error("failed to extract table refs from query", zap.String("query-id", dbMetric.ID), zap.Error(err))
	}
	query := models.Query{
		ID:             dbMetric.ID,
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
	}).Create(dbMetric).Error
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

	return nil
}
