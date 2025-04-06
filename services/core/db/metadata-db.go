package db

import (
	"errors"
	"github.com/opengovern/opensecurity/services/core/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (db Database) upsertQueryParameter(queryParam models.PolicyParameterValues) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&queryParam).Error
}

func (db Database) upsertQueryParameters(queryParam []*models.PolicyParameterValues) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}, {Name: "control_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(queryParam).Error
}

func (db Database) SetQueryParameter(key string, value string) error {
	return db.upsertQueryParameter(models.PolicyParameterValues{
		Key:   key,
		Value: value,
	})
}

func (db Database) SetQueryParameters(queryParams []*models.PolicyParameterValues) error {
	return db.upsertQueryParameters(queryParams)
}

func (db Database) GetQueryParameter(key string) (*models.PolicyParameterValues, error) {
	var queryParam models.PolicyParameterValues
	err := db.orm.First(&queryParam, "key = ?", key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &queryParam, nil
}

func (db Database) GetQueryParametersValues(keyRegex *string) ([]models.PolicyParameterValues, error) {
	var queryParams []models.PolicyParameterValues
	tx := db.orm.Model(&models.PolicyParameterValues{})
	if keyRegex != nil {
		tx = tx.Where("key ~* ?", *keyRegex)
	}
	err := tx.Find(&queryParams).Error
	if err != nil {
		return nil, err
	}
	return queryParams, nil
}

func (db Database) GetQueryParametersByIds(ids []string) ([]models.PolicyParameterValues, error) {
	var queryParams []models.PolicyParameterValues
	err := db.orm.Where("key IN ?", ids).Find(&queryParams).Error
	if err != nil {
		return nil, err
	}
	return queryParams, nil
}

func (db Database) DeleteQueryParameter(key string) error {
	return db.orm.Unscoped().Delete(&models.PolicyParameterValues{}, "key = ?", key).Error
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
// get user layout
func (db Database) GetUserLayout(userID string) (*models.UserLayout, error) {
	var userLayout models.UserLayout
	err := db.orm.First(&userLayout, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userLayout, nil
}
// set user layout
func (db Database) SetUserLayout( layoutConfig models.UserLayout) error {
	err := db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"layout_config", "name", "is_private"}),
	}).Create(&layoutConfig).Error
	if err != nil {
		return err
	}
	return nil
	
}
// get public layouts 

func (db Database) GetPublicLayouts() ([]models.UserLayout, error) {
	var userLayouts []models.UserLayout
	err := db.orm.Where("is_private = ?", false).Find(&userLayouts).Error
	if err != nil {
		return nil, err
	}
	return userLayouts, nil
}
// change layout privacy
func (db Database) ChangeLayoutPrivacy(userID string, isPrivate bool) error {
	err := db.orm.Model(&models.UserLayout{}).Where("user_id = ?", userID).Update("is_private", isPrivate).Error
	if err != nil {
		return err
	}
	return nil
}

