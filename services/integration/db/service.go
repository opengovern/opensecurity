package db

import (
	"github.com/opengovern/opencomply/services/integration/models"
	"gorm.io/gorm"
)

type Database struct {
	Orm                *gorm.DB
	IntegrationTypeOrm *gorm.DB
}

func NewDatabase(orm *gorm.DB, integrationTypeOrm *gorm.DB) Database {
	return Database{
		Orm:                orm,
		IntegrationTypeOrm: integrationTypeOrm,
	}
}

func (db Database) Initialize() error {
	err := db.Orm.AutoMigrate(
		&models.Integration{},
		&models.Credential{},
		&models.IntegrationGroup{},
	)
	if err != nil {
		return err
	}

	err = db.IntegrationTypeOrm.AutoMigrate(
		&models.IntegrationPlugin{},
		&models.IntegrationPluginBinary{},
	)
	if err != nil {
		return err
	}

	return nil
}
