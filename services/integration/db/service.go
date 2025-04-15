package db

import (
	"github.com/opengovern/opensecurity/services/integration/models"
	taskModels "github.com/opengovern/opensecurity/services/tasks/db/models"
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
		&models.IntegrationResourcetypes{},
	)
	if err != nil {
		return err
	}

	err = db.IntegrationTypeOrm.AutoMigrate(
		&models.IntegrationPlugin{},
		&models.IntegrationPluginBinary{},
		&taskModels.TaskBinary{},
	)
	if err != nil {
		return err
	}

	return nil
}
