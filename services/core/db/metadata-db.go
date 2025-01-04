package db

import (
	"github.com/opengovern/opencomply/services/core/db/models"
	"gorm.io/gorm/clause"
	"errors"
	"gorm.io/gorm"

)


func (db Database) upsertConfigMetadata(configMetadata models.ConfigMetadata) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "type"}),
	}).Create(&configMetadata).Error
}

func (db Database) SetConfigMetadata(cm models.ConfigMetadata) error {
	return db.upsertConfigMetadata(models.ConfigMetadata{
		Key:   cm.Key,
		Type:  cm.Type,
		Value: cm.Value,
	})
}

func (db Database) GetConfigMetadata(key string) (models.IConfigMetadata, error) {
	var configMetadata models.ConfigMetadata
	err := db.orm.First(&configMetadata, "key = ?", key).Error
	if err != nil {
		return nil, err
	}
	return configMetadata.ParseToType()
}

func (db Database) AddFilter(filter models.Filter) error {
	return db.orm.Model(&models.Filter{}).Create(filter).Error
}

func (db Database) ListFilters() ([]models.Filter, error) {
	var filters []models.Filter
	err := db.orm.Model(&models.Filter{}).First(&filters).Error
	if err != nil {
		return nil, err
	}
	return filters, nil
}

func (db Database) ListApp() ([]models.PlatformConfiguration, error) {
	var apps []models.PlatformConfiguration
	err := db.orm.Model(&models.PlatformConfiguration{}).Find(&apps).Error
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (db Database) CreateApp(app *models.PlatformConfiguration) error {
	return db.orm.Model(&models.PlatformConfiguration{}).Create(app).Error
}

func (db Database) AppConfigured(configured bool) error {
	return db.orm.Model(&models.PlatformConfiguration{}).Update("configured", configured).Error
}

func (db Database) GetAppConfiguration() (*models.PlatformConfiguration, error) {
	var appConfiguration models.PlatformConfiguration
	err := db.orm.Model(&models.PlatformConfiguration{}).First(&appConfiguration).Error
	if err != nil {
		return nil, err
	}
	return &appConfiguration, nil
}


func (db Database) upsertQueryParameter(queryParam models.QueryParameterValues) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&queryParam).Error
}

func (db Database) upsertQueryParameters(queryParam []*models.QueryParameterValues) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(queryParam).Error
}

func (db Database) SetQueryParameter(key string, value string) error {
	return db.upsertQueryParameter(models.QueryParameterValues{
		Key:   key,
		Value: value,
	})
}

func (db Database) SetQueryParameters(queryParams []*models.QueryParameterValues) error {
	return db.upsertQueryParameters(queryParams)
}

func (db Database) GetQueryParameter(key string) (*models.QueryParameterValues, error) {
	var queryParam models.QueryParameterValues
	err := db.orm.First(&queryParam, "key = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &queryParam, nil
}

func (db Database) GetQueryParametersValues() ([]models.QueryParameterValues, error) {
	var queryParams []models.QueryParameterValues
	err := db.orm.Find(&queryParams).Error
	if err != nil {
		return nil, err
	}
	return queryParams, nil
}

func (db Database) GetQueryParametersByIds(ids []string) ([]models.QueryParameterValues, error) {
	var queryParams []models.QueryParameterValues
	err := db.orm.Where("key IN ?", ids).Find(&queryParams).Error
	if err != nil {
		return nil, err
	}
	return queryParams, nil
}

func (db Database) DeleteQueryParameter(key string) error {
	return db.orm.Unscoped().Delete(&models.QueryParameterValues{}, "key = ?", key).Error
}

func (db Database) ListQueryViews() ([]models.QueryView, error) {
	var queryViews []models.QueryView
	err := db.orm.
		Model(&models.QueryView{}).
		Preload(clause.Associations).
		Find(&queryViews).Error
	if err != nil {
		return nil, err
	}
	return queryViews, nil
}
