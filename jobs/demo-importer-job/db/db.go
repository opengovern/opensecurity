package db

import (
	"errors"

	"github.com/jackc/pgtype"
	"github.com/opengovern/opensecurity/jobs/demo-importer-job/db/model"
	"gorm.io/gorm"
)

type Database struct {
	ORM *gorm.DB
}

func (db Database) Initialize() error {
	err := db.ORM.AutoMigrate(
		&model.Migration{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) GetMigration(id string) (*model.Migration, error) {
	var mig model.Migration
	tx := db.ORM.Model(&model.Migration{}).Where("id = ?", id).First(&mig)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return &mig, nil
}

func (db Database) UpdateMigrationAdditionalInfo(id string, additionalInfo string) error {
	tx := db.ORM.Model(&model.Migration{}).Where("id = ?", id).Update("additional_info", additionalInfo)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) UpdateMigrationJob(id string, status string, JobsStatus pgtype.JSONB) error {
	tx := db.ORM.Model(&model.Migration{}).Where("id = ?", id).Update("status", status).Update("jobs_status", JobsStatus)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) CreateMigration(m *model.Migration) error {
	tx := db.ORM.Model(&model.Migration{}).Create(m)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
