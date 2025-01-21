package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/opengovern/og-util/pkg/integration"
	"github.com/opengovern/opencomply/services/integration/models"
	"gorm.io/gorm/clause"
)

// CreateIntegration creates a new integration
func (db Database) CreateIntegration(s *models.Integration) error {
	tx := db.Orm.
		Model(&models.Integration{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(s)

	if tx.Error != nil {
		return tx.Error
	} else if tx.RowsAffected == 0 {
		return fmt.Errorf("create spn: didn't create spn due to id conflict")
	}

	return nil
}

// DeleteIntegration deletes a integration
func (db Database) DeleteIntegration(IntegrationID uuid.UUID) error {
	tx := db.Orm.
		Where("integration_id = ?", IntegrationID.String()).
		Unscoped().
		Delete(&models.Integration{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// DeleteSampleIntegrations deletes sample integrations
func (db Database) DeleteSampleIntegrations() error {
	tx := db.Orm.
		Where("state = ?", integration.IntegrationStateSample).
		Unscoped().
		Delete(&models.Integration{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// ListSampleIntegrations list sample integrations
func (db Database) ListSampleIntegrations() ([]models.Integration, error) {
	var integrations []models.Integration

	tx := db.Orm.
		Where("state = ?", integration.IntegrationStateSample).
		Find(&integrations)
	if tx.Error != nil {
		return integrations, tx.Error
	}

	return integrations, nil
}

// ListIntegration list Integration
func (db Database) ListIntegration(types []integration.Type) ([]models.Integration, error) {
	var integrations []models.Integration
	tx := db.Orm.
		Model(&models.Integration{})

	if len(types) > 0 {
		tx = tx.Where("integration_type IN ?", types)
	}

	tx = tx.Find(&integrations)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return integrations, nil
}

// ListIntegrationsByFilters list Integrations by filters
func (db Database) ListIntegrationsByFilters(IntegrationIDs []string, types []string, NameRegex, providerIDRegex *string) ([]models.Integration, error) {
	var integrations []models.Integration
	tx := db.Orm.
		Model(&models.Integration{})

	if len(IntegrationIDs) > 0 {
		tx = tx.Where("integration_id IN ?", IntegrationIDs)
	}
	if len(types) > 0 {
		tx = tx.Where("integration_type IN ?", types)
	}
	if NameRegex != nil {
		tx = tx.Where("integration_name ~* ?", *NameRegex)
	}
	if providerIDRegex != nil {
		tx = tx.Where("provider_id ~* ?", *providerIDRegex)
	}

	tx = tx.Find(&integrations)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return integrations, nil
}

// GetIntegration get a Integration
func (db Database) GetIntegration(tracker uuid.UUID) (*models.Integration, error) {
	var integration models.Integration
	tx := db.Orm.
		Model(&models.Integration{}).
		Where("integration_id = ?", tracker.String()).
		First(&integration)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &integration, nil
}

// UpdateIntegration deletes a integration
func (db Database) UpdateIntegration(integration *models.Integration) error {
	tx := db.Orm.
		Where("integration_id = ?", integration.IntegrationID.String()).
		Updates(integration)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// InactiveIntegrationType inactive integrations for an integration type
func (db Database) InactiveIntegrationType(it integration.Type) error {
	tx := db.Orm.
		Model(&models.Integration{}).
		Where("integration_type = ?", it).
		Update("state", integration.IntegrationStateInactive)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
