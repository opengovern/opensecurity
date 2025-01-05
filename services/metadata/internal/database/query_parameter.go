package database

import (
	"errors"

	"github.com/opengovern/opencomply/services/metadata/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (db Database) upsertQueryParameter(queryParam models.PolicyParameterValues) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&queryParam).Error
}

func (db Database) upsertQueryParameters(queryParam []*models.PolicyParameterValues) error {
	return db.orm.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
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

func (db Database) GetQueryParameters() ([]models.PolicyParameterValues, error) {
	var queryParams []models.PolicyParameterValues
	err := db.orm.Find(&queryParams).Error
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
