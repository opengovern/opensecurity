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
// GetUserLayouts returns all dashboards for the specified user,
// preloading their widgets.
func (db Database) GetUserLayouts(userID string) ([]models.Dashboard, error) {
	var dashboards []models.Dashboard
	err := db.orm.
		Where("user_id = ?", userID).
		Preload("Widgets").
		Find(&dashboards).Error
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

// GetUserDefaultLayout returns the dashboard that is marked as default for the user.
// If no default is found, it returns nil.
func (db Database) GetUserDefaultLayout(userID string) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := db.orm.
		Where("user_id = ? AND is_default = ?", userID, true).
		Preload("Widgets").
		First(&dashboard).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dashboard, nil
}

// SetUserLayout upserts a dashboard (excluding its widgets), and then
// replaces the associated widgets. This function runs in a transaction.
func (db Database) SetUserLayout(layout models.Dashboard) error {
	return db.orm.Transaction(func(tx *gorm.DB) error {
		// Upsert the dashboard (excluding Widgets)
		err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "description", "is_default", "is_private", "updated_at", "user_id"}),
		}).Omit("Widgets").Create(&layout).Error
		if err != nil {
			return err
		}

		// Delete previous widgets for this dashboard
		err = tx.Where("user_id = ? AND dashboard_id = ?", layout.UserID, layout.ID).
			Delete(&models.Widget{}).Error
		if err != nil {
			return err
		}

		// Ensure each widget has the correct DashboardID and UserID
		for i := range layout.Widgets {
			layout.Widgets[i].DashboardID = layout.ID
			layout.Widgets[i].UserID = layout.UserID
		}

		// Insert the new widgets
		err = tx.Create(&layout.Widgets).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// GetPublicLayouts returns dashboards that are not private (public dashboards).
func (db Database) GetPublicLayouts() ([]models.Dashboard, error) {
	var dashboards []models.Dashboard
	err := db.orm.
		Where("is_private = ?", false).
		Preload("Widgets").
		Find(&dashboards).Error
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

// ChangeLayoutPrivacy updates the privacy status for all dashboards of a user.
func (db Database) ChangeLayoutPrivacy(userID string, isPrivate bool) error {
	err := db.orm.
		Model(&models.Dashboard{}).
		Where("user_id = ?", userID).
		Update("is_private", isPrivate).Error
	return err
}

// GetUserWidgets returns all widgets for the specified user.
func (db Database) GetUserWidgets(userID string) ([]models.Widget, error) {
	var widgets []models.Widget
	err := db.orm.
		Where("user_id = ?", userID).
		Find(&widgets).Error
	if err != nil {
		return nil, err
	}
	return widgets, nil
}

// GetWidget returns a single widget by its ID.
func (db Database) GetWidget(widgetID string) (*models.Widget, error) {
	var widget models.Widget
	err := db.orm.
		Where("id = ?", widgetID).
		First(&widget).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &widget, nil
}

// SetUserWidget upserts a widget using its ID as the conflict key.
func (db Database) SetUserWidget(widget models.Widget) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "description", "widget_type", "widget_props", "row_span", "column_span", "column_offset", "is_public", "user_id", "dashboard_id", "updated_at"}),
	}).Create(&widget).Error
}

// DeleteUserWidget deletes a widget by its ID.
func (db Database) DeleteUserWidget(widgetID string) error {
	return db.orm.Delete(&models.Widget{}, "id = ?", widgetID).Error
}

func (db Database) AddWidgetsToDashboard(dashboardID string, widgets []models.Widget) error {
	return db.orm.Transaction(func(tx *gorm.DB) error {
		// Set the DashboardID for each widget
		for i := range widgets {
			widgets[i].DashboardID = dashboardID
		}
		// Bulk insert all widgets
		if err := tx.Create(&widgets).Error; err != nil {
			return err
		}
		return nil
	})
}

func (db Database) AddWidgets(widgets []models.Widget) error {
    return db.orm.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&widgets).Error; err != nil {
            return err
        }
        return nil
    })
}
