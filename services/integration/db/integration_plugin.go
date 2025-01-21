package db

import (
	"errors"
	"github.com/opengovern/opencomply/services/integration/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db Database) ListPlugins() ([]models.IntegrationPlugin, error) {
	var plugin []models.IntegrationPlugin
	err := db.IntegrationTypeOrm.Model(models.IntegrationPlugin{}).Find(&plugin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return plugin, nil
}

func (db Database) GetPluginByID(pluginID string) (*models.IntegrationPlugin, error) {
	var plugin models.IntegrationPlugin
	err := db.IntegrationTypeOrm.Model(models.IntegrationPlugin{}).Where("plugin_id = ?", pluginID).Find(&plugin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plugin, nil
}

func (db Database) GetPluginBinaryByID(pluginID string) (*models.IntegrationPluginBinary, error) {
	var pluginBinary models.IntegrationPluginBinary
	err := db.IntegrationTypeOrm.Model(models.IntegrationPluginBinary{}).Where("plugin_id = ?", pluginID).Find(&pluginBinary).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &pluginBinary, nil
}

func (db Database) GetPluginByIntegrationType(pluginID string) (*models.IntegrationPlugin, error) {
	var plugin models.IntegrationPlugin
	err := db.IntegrationTypeOrm.Model(models.IntegrationPlugin{}).Where("plugin_id = ?", pluginID).Find(&plugin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plugin, nil
}

func (db Database) GetPluginByURL(url string) (*models.IntegrationPlugin, error) {
	var plugin models.IntegrationPlugin
	err := db.IntegrationTypeOrm.Model(models.IntegrationPlugin{}).Where("url = ?", url).Find(&plugin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plugin, nil
}

func (db Database) CreatePlugin(plugin *models.IntegrationPlugin) error {
	err := db.IntegrationTypeOrm.Create(plugin).Error
	if err != nil {
		return err
	}
	return nil
}

func (db Database) CreatePluginBinary(pluginBinary *models.IntegrationPluginBinary) error {
	// Use GORM's upsert functionality
	err := db.IntegrationTypeOrm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_id"}},                                        // Define the unique key or conflict target
		DoUpdates: clause.AssignmentColumns([]string{"integration_plugin", "cloud_ql_plugin"}), // Columns to update
	}).Create(pluginBinary).Error

	if err != nil {
		return err
	}
	return nil
}

func (db Database) UpdatePlugin(plugin *models.IntegrationPlugin) error {
	err := db.IntegrationTypeOrm.Model(models.IntegrationPlugin{}).Where("plugin_id = ?", plugin.PluginID).Updates(plugin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}
