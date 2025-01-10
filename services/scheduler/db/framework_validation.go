package db

import (
	"errors"
	"github.com/opengovern/opencomply/services/scheduler/db/model"
	"gorm.io/gorm"
)

func (db Database) CreateFrameworkValidation(validation *model.FrameworkValidation) error {
	tx := db.ORM.
		Model(&model.FrameworkValidation{}).
		Create(validation)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) GetFrameworkValidation(frameworkID string) (*model.FrameworkValidation, error) {
	var validations model.FrameworkValidation
	tx := db.ORM.Model(&model.FrameworkValidation{}).
		Where("framework_id = ?", frameworkID).
		First(&validations)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &validations, nil
}
