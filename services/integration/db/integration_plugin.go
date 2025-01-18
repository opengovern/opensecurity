package db

import (
	"errors"
	"github.com/opengovern/opencomply/services/integration/models"
	"gorm.io/gorm"
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

func (db Database) CreatePlugin(plugin models.IntegrationPlugin) error {
	err := db.IntegrationTypeOrm.Create(&plugin).Error
	if err != nil {
		return err
	}
	return nil
}

func (db Database) UpdatePlugin(plugin models.IntegrationPlugin) error {
	err := db.IntegrationTypeOrm.Model(models.IntegrationPlugin{}).Where("plugin_id = ?", plugin.PluginID).Updates(&plugin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}
