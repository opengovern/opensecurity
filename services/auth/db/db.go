package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/opengovern/og-util/pkg/api"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	Orm *gorm.DB
}

func (db Database) Initialize() error {
	err := db.Orm.AutoMigrate(
		&ApiKey{},
		&User{},
		&Configuration{},
		&Connector{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) GetKeyPair() ([]Configuration, error) {
	var s []Configuration
	tx := db.Orm.Model(&Configuration{}).
		Where("key = 'private_key' or key = 'public_key'").Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) AddConfiguration(c *Configuration) error {
	tx := db.Orm.Create(c)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) ListApiKeys() ([]ApiKey, error) {
	var s []ApiKey
	tx := db.Orm.Model(&ApiKey{}).	
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) ListApiKeysForUser(userId string) ([]ApiKey, error) {
	var s []ApiKey
	tx := db.Orm.Model(&ApiKey{}).
		Where("creator_user_id", userId).
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) AddApiKey(key *ApiKey) error {
	tx := db.Orm.Create(key)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateAPIKey( id string, is_active bool,role api.Role) error {
	tx := db.Orm.Model(&ApiKey{}).
		Where("id = ?", id).
		Update("is_active", is_active).
		Update("role", role)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}


func (db Database) DeleteAPIKey(id uint64) error {
	tx := db.Orm.Model(&ApiKey{}).
		Where("id = ?", id).
		Delete(&ApiKey{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) CreateUser(user *User) error {
	tx := db.Orm.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"id", "created_at", "updated_at", "email", "email_verified",
			 "role", "connector_id", "external_id",
			"full_name", "last_login", "username", "is_active"}),
	}).Create(user)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateUser(user *User) error {
	fmt.Println("User", user)
	tx := db.Orm.Model(&User{}).
		Where("id = ?", user.ID).
		Updates(user).
		Update("is_active", user.IsActive)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) DeleteUser(id uint) error {
	tx := db.Orm.
		Where("id = ?", id).
		Delete(&User{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}


// Get all users

func (db Database) GetUsers() ([]User, error) {
	var s []User
	tx := db.Orm.Model(&User{}).
		Order("updated_at desc").
		Find(&s)
		
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}



func (db Database) GetUser(id string) (*User, error) {
	var s User
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id).
		Find(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &s, nil
}
func (db Database) GetUserByExternalID(id string) (*User, error) {
	var s User
	tx := db.Orm.Model(&User{}).
		Where("external_id = ?", id).
		Find(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &s, nil
}

func (db Database) GetUsersCount() (int64, error) {
	var count int64
	tx := db.Orm.Model(&User{}).
		Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}

func (db Database) GetFirstUser() (*User, error) {
	var user User
	tx := db.Orm.Model(&User{}).
		First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (db Database) UpdateUserLastLogin(id uuid.UUID, lastLogin *time.Time) error {
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id)

	if lastLogin != nil {
		tx = tx.Update("last_login", lastLogin)
	}

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db Database) UpdateUserLastLoginWithExternalID(id string, lastLogin time.Time) error {
	tx := db.Orm.Model(&User{}).
		Where("external_id = ?", id)

	if !lastLogin.IsZero(){
		tx = tx.Update("last_login", lastLogin)
	}

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}


func (db Database) GetUserByEmail(email string) (*User, error) {
	var s User
	tx := db.Orm.Model(&User{}).
		Where("email = ? ", email).
		Find(&s)
	
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if tx.Error != nil {
		return nil, tx.Error
	}
	return &s, nil
}




func (db Database) UserPasswordUpdate(id uint) error {
	tx := db.Orm.Model(&User{}).
		Where("id = ? ", id).
		Update("require_password_change", false)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}


func (db Database) DisableUser(id uuid.UUID) error {
	tx := db.Orm.Model(&User{}).
		Where("id = ? ", id).
		Update("disabled", true)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (db Database) EnableUser(id uuid.UUID) error {
	tx := db.Orm.Model(&User{}).
		Where("id = ? ", id).
		Update("is_active", false)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//  De Active User
func (db Database) DeActiveUser(id uuid.UUID) error {
	tx := db.Orm.Model(&User{}).
		Where("id = ? ", id).
		Update("is_active", false)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//  Active User
func (db Database) ActiveUser(id uuid.UUID) error {
	tx := db.Orm.Model(&User{}).
		Where("id = ? ", id).
		Update("is_active", true)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// find id by email

func (db Database) FindIdByEmail(email string) (uint, error) {
	var s User
	tx := db.Orm.Model(&User{}).
		Where("email = ? ", email).
		First(&s)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return s.ID, nil
}

func (db Database) CountApiKeysForUser(userID string) (int64, error) {
	var s int64
	tx := db.Orm.Model(&ApiKey{}).
		Where("creator_user_id", userID).
		Where("is_active", "true").
		Count(&s)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return s, nil
}

// Get all connectors

func (db Database) GetConnectors() ([]Connector, error) {
	var s []Connector
	tx := db.Orm.Model(&Connector{}).
		Order("last_update desc").
		Find(&s)
		
	if tx.Error != nil {
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) GetConnector(id string) (*Connector, error) {
	var s Connector
	tx := db.Orm.Model(&Connector{}).
		Where("id = ?", id).
		Find(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &s, nil
}

func (db Database) CreateConnector(connector *Connector) error {
	tx := db.Orm.Create(connector)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (db Database) UpdateConnector(connector *Connector) error {
	tx := db.Orm.Model(&Connector{}).
		Where("id = ?", connector.ID).
		Updates(connector)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// delete connector by connector id
func (db Database) DeleteConnector(connectorID string) error {
	tx := db.Orm.Model(&Connector{}).
		Where("connector_id = ?", connectorID).
		Delete(&Connector{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}


func (db Database) GetConnectorByConnectorID(connectorID string) (*Connector, error) {
	var s Connector
	tx := db.Orm.Model(&Connector{}).
		Where("connector_id = ?", connectorID).
		Find(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &s, nil
}

// get Connector by connector type 

func (db Database) GetConnectorByConnectorType(connectorType string) (*Connector, error) {
	var s Connector
	tx := db.Orm.Model(&Connector{}).
		Where("connector_type = ?", connectorType).
		Find(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &s, nil
}



