package db

import (
	"github.com/opengovern/opensecurity/services/integration/models"
)

// DeleteIntegrationGroup deletes an Integration Group
func (db Database) DeleteIntegrationGroup(name string) error {
	tx := db.Orm.
		Where("name = ?", name).
		Unscoped().
		Delete(&models.IntegrationGroup{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// ListIntegrationGroups list Integration Groups
func (db Database) ListIntegrationGroups() ([]models.IntegrationGroup, error) {
	var integrationGroups []models.IntegrationGroup
	tx := db.Orm.
		Model(&models.IntegrationGroup{}).
		Find(&integrationGroups)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return integrationGroups, nil
}

// GetIntegrationGroup get an integration group
func (db Database) GetIntegrationGroup(name string) (*models.IntegrationGroup, error) {
	var integrationGroup models.IntegrationGroup
	tx := db.Orm.
		Model(&models.IntegrationGroup{}).
		Where("name = ?", name).
		First(&integrationGroup)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &integrationGroup, nil
}
