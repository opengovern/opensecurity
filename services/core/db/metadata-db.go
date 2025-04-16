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

// GetUserLayouts get user layout// Get user layouts (with widgets)
func (db Database) GetUserLayouts(userID string) ([]models.Dashboard, error) {
	var userLayouts []models.Dashboard
	err := db.orm.Preload("Widgets").Where("user_id = ?", userID).Find(&userLayouts).Error
	if err != nil {
		return nil, err
	}
	return userLayouts, nil
}

// GetUserDefaultLayout Get the default layout for a user
func (db Database) GetUserDefaultLayout(userID string) (*models.Dashboard, error) {
	var userLayout models.Dashboard
	err := db.orm.Preload("Widgets").Where("user_id = ? AND is_default = ?", userID, true).First(&userLayout).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userLayout, nil
}

// SetUserLayout Upsert dashboard and update associated widgets
func (db Database) SetUserLayout(layoutConfig models.Dashboard) error {
	return db.orm.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "description", "is_default", "is_private", "updated_at", "user_id",
			}),
		}).Omit("Widgets").Create(&layoutConfig).Error
		if err != nil {
			return err
		}

		// Replace dashboard-widgets association
		err = tx.Model(&layoutConfig).Association("Widgets").Replace(layoutConfig.Widgets)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetPublicLayouts Get all public dashboards
func (db Database) GetPublicLayouts() ([]models.Dashboard, error) {
	var publicLayouts []models.Dashboard
	err := db.orm.Preload("Widgets").Where("is_private = ?", false).Find(&publicLayouts).Error
	if err != nil {
		return nil, err
	}
	return publicLayouts, nil
}

// ChangeLayoutPrivacy Change privacy status of all layouts for a user
func (db Database) ChangeLayoutPrivacy(userID string, isPrivate bool) error {
	return db.orm.Model(&models.Dashboard{}).Where("user_id = ?", userID).Update("is_private", isPrivate).Error
}

// GetUserWidgets Get all widgets for a user
func (db Database) GetUserWidgets(userID string) ([]models.Widget, error) {
	var widgets []models.Widget
	err := db.orm.Where("user_id = ?", userID).Find(&widgets).Error
	if err != nil {
		return nil, err
	}
	return widgets, nil
}

// GetWidget Get a single widget by ID
func (db Database) GetWidget(widgetID string) (*models.Widget, error) {
	var widget models.Widget
	err := db.orm.Where("id = ?", widgetID).First(&widget).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &widget, nil
}

// AddWidgets Add widgets in bulk
func (db Database) AddWidgets(widgets []models.Widget) error {
	return db.orm.Create(&widgets).Error
}

// SetUserWidget Upsert a widget
func (db Database) SetUserWidget(widget models.Widget) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"title", "description", "widget_type", "widget_props",
			"row_span", "column_span", "column_offset", "is_public",
			"user_id", "updated_at",
		}),
	}).Create(&widget).Error
}

// DeleteWidgets Delete widgets by IDs
func (db Database) DeleteWidgets(widgetIDs []string) error {
	return db.orm.Where("id IN ?", widgetIDs).Delete(&models.Widget{}).Error
}

// UpdateWidget Update widget content (single widget)
func (db Database) UpdateWidget(widget models.Widget) error {
	return db.orm.Model(&models.Widget{}).Where("id = ?", widget.ID).Updates(&widget).Error
}

// UpdateDashboardWidgets Update widget-dashboard relationship (associate widgets to dashboard)
func (db Database) UpdateDashboardWidgets(dashboardID string, widgetIDs []string) error {
	var dashboard models.Dashboard
	dashboard.ID = dashboardID

	var widgets []models.Widget
	err := db.orm.Where("id IN ?", widgetIDs).Find(&widgets).Error
	if err != nil {
		return err
	}

	return db.orm.Model(&dashboard).Association("Widgets").Replace(widgets)
}

// DeleteUserWidget
func (db Database) DeleteUserWidget(widgetID string) error {
	return db.orm.Where("id = ?", widgetID).Delete(&models.Widget{}).Error
}

func (db Database) UpdateWidgetDashboards(widgetID string, dashboardIDs []string) error {
	var widget models.Widget
	// Fetch the widget by ID
	if err := db.orm.First(&widget, "id = ?", widgetID).Error; err != nil {
		return err
	}

	var dashboards []models.Dashboard
	// Fetch the dashboards by the provided dashboardIDs
	if err := db.orm.Where("id IN ?", dashboardIDs).Find(&dashboards).Error; err != nil {
		return err
	}

	// Replace the dashboards associated with the widget
	return db.orm.Model(&models.Widget{}).Association("Dashboards").Replace(dashboards)
}

// GetAllPublicWidgets get all public widgets
func (db Database) GetAllPublicWidgets() ([]models.Widget, error) {
	var widgets []models.Widget
	err := db.orm.
		Where("is_public = ?", true).
		Find(&widgets).Error
	if err != nil {
		return nil, err
	}
	return widgets, nil
}
