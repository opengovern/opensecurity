package db

import (
	"context" // Added context import
	"errors"
	"fmt"
	"time"

	// Keep existing imports

	"github.com/opengovern/og-util/pkg/api"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Database struct remains the same
type Database struct {
	Orm *gorm.DB
}

// Note: Methods below now accept context.Context as the first argument.
// Consider changing the receiver to pointer '(db *Database)' for consistency,
// although '(db Database)' also works for interface satisfaction when passing &adb.

// Initialize remains the same, AutoMigrate usually doesn't need context.
func (db Database) Initialize() error {
	err := db.Orm.AutoMigrate(
		&ApiKey{},
		&User{},
		&Configuration{},
		&Connector{},
	)
	// Error handling remains the same
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}
	return nil
}

func (db Database) GetKeyPair(ctx context.Context) ([]Configuration, error) {
	var s []Configuration
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&Configuration{}).
		Where("key = 'private_key' or key = 'public_key'").Find(&s)
	if tx.Error != nil {
		// Add context to error message if desired
		return nil, fmt.Errorf("error getting key pair: %w", tx.Error)
	}
	return s, nil
}

func (db Database) AddConfiguration(ctx context.Context, c *Configuration) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Create(c)
	if tx.Error != nil {
		return fmt.Errorf("error adding configuration key=%s: %w", c.Key, tx.Error)
	}
	return nil
}

func (db Database) ListApiKeys(ctx context.Context) ([]ApiKey, error) {
	var s []ApiKey
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error listing all API keys: %w", tx.Error)
	}
	return s, nil
}

func (db Database) ListApiKeysForUser(ctx context.Context, userId string) ([]ApiKey, error) {
	var s []ApiKey
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Where("creator_user_id = ?", userId).
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error listing API keys for user %s: %w", userId, tx.Error)
	}
	return s, nil
}

func (db Database) AddApiKey(ctx context.Context, key *ApiKey) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Create(key)
	if tx.Error != nil {
		return fmt.Errorf("error adding API key '%s': %w", key.Name, tx.Error)
	}
	return nil
}

func (db Database) UpdateAPIKey(ctx context.Context, id string, is_active bool, role api.Role) error {
	// Use WithContext for the chain
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": is_active, "role": role}) // Use map for multiple updates

	// Check RowsAffected for not found condition
	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return specific error if nothing was updated
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating API key ID %s: %w", id, tx.Error)
	}
	return nil
}

func (db Database) DeleteAPIKey(ctx context.Context, id uint64) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ApiKey{}) // GORM needs the type for table name inference

	// Check RowsAffected for not found condition
	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error deleting API key ID %d: %w", id, tx.Error)
	}
	return nil
}

func (db Database) CreateUser(ctx context.Context, user *User) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // Ensure 'id' is the correct conflict column
		DoUpdates: clause.AssignmentColumns([]string{"updated_at", "email", "email_verified", "role", "connector_id", "external_id", "full_name", "last_login", "username", "is_active"}),
		// Note: 'created_at' and 'id' usually shouldn't be in DoUpdates unless you intend to overwrite them on conflict.
	}).Create(user)

	if tx.Error != nil {
		return fmt.Errorf("error creating/updating user email=%s: %w", user.Email, tx.Error)
	}
	return nil
}

func (db Database) UpdateUser(ctx context.Context, user *User) error {
	// Use WithContext. Updates handles non-zero fields of the struct.
	// Ensure user.ID is set correctly before calling.
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ?", user.ID).
		Updates(user) // Pass the whole struct

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // User with that ID not found
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating user ID %d: %w", user.ID, tx.Error)
	}
	return nil
}

func (db Database) DeleteUser(ctx context.Context, id uint) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).
		Where("id = ?", id).
		Delete(&User{}) // Provide type

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error deleting user ID %d: %w", id, tx.Error)
	}
	return nil
}

func (db Database) GetUsers(ctx context.Context) ([]User, error) {
	var s []User
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Order("updated_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error getting all users: %w", tx.Error)
	}
	return s, nil
}

func (db Database) GetUser(ctx context.Context, id string) (*User, error) {
	var s User
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ?", id).
		First(&s) // Use First to get ErrRecordNotFound
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil for not found
		}
		return nil, fmt.Errorf("error getting user by ID %s: %w", id, tx.Error)
	}
	return &s, nil
}

func (db Database) GetUserByExternalID(ctx context.Context, id string) (*User, error) {
	var s User
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("external_id = ?", id).
		First(&s) // Use First
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil for not found
		}
		return nil, fmt.Errorf("error getting user by external ID %s: %w", id, tx.Error)
	}
	return &s, nil
}

func (db Database) GetUsersCount(ctx context.Context) (int64, error) {
	var count int64
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).Count(&count)
	if tx.Error != nil {
		return 0, fmt.Errorf("error counting users: %w", tx.Error)
	}
	return count, nil
}

func (db Database) GetFirstUser(ctx context.Context) (*User, error) {
	var user User
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).Order("id asc").First(&user) // Explicit order
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // No users exist
		}
		return nil, fmt.Errorf("error getting first user: %w", tx.Error)
	}
	return &user, nil
}

// This method signature from original code used uuid.UUID, interface uses uint.
// Let's keep uint from interface. Caller (ResetUserPassword) needs user.ID (uint).
func (db Database) UserPasswordUpdate(ctx context.Context, id uint) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ?", id).
		Update("require_password_change", false)

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating password status for user ID %d: %w", id, tx.Error)
	}
	return nil
}

// UpdateUserLastLoginWithExternalID accepts context now
func (db Database) UpdateUserLastLoginWithExternalID(ctx context.Context, id string, lastLogin time.Time) error {
	if lastLogin.IsZero() { // Avoid zero time updates
		return nil
	}
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("external_id = ?", id).
		Update("last_login", lastLogin)

	if tx.Error == nil && tx.RowsAffected == 0 {
		// Don't return error if user not found, maybe just log?
		// Or return ErrRecordNotFound if strict checking needed.
		// For last login update, maybe non-critical if user disappeared.
		return nil // Or return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating last login for external ID %s: %w", id, tx.Error)
	}
	return nil
}

func (db Database) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var s User
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("email = ?", email).
		First(&s) // Use First

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil for not found
		}
		return nil, fmt.Errorf("error getting user by email %s: %w", email, tx.Error)
	}
	return &s, nil
}

// FindIdByEmail needed for interface? Only used internally in original code?
// Adding it to implementation with context just in case.
func (db Database) FindIdByEmail(ctx context.Context, email string) (uint, error) {
	var s User
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("email = ?", email).
		Select("id"). // Only select ID
		First(&s)
	if tx.Error != nil {
		// Return 0 for ID on error, including not found
		return 0, fmt.Errorf("error finding ID for email %s: %w", email, tx.Error)
	}
	return s.ID, nil
}

func (db Database) CountApiKeysForUser(ctx context.Context, userID string) (int64, error) {
	var s int64
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Where("creator_user_id = ?", userID).
		Where("is_active = ?", true). // Use boolean true
		Count(&s)
	if tx.Error != nil {
		return 0, fmt.Errorf("error counting API keys for user %s: %w", userID, tx.Error)
	}
	return s, nil
}

func (db Database) GetConnectors(ctx context.Context) ([]Connector, error) {
	var s []Connector
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Order("last_update desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error getting all connectors: %w", tx.Error)
	}
	return s, nil
}

// GetConnector gets by *local* DB ID (uint)
func (db Database) GetConnector(ctx context.Context, id string) (*Connector, error) {
	var s Connector
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("id = ?", id).
		First(&s) // Use First
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting connector by local ID %s: %w", id, tx.Error)
	}
	return &s, nil
}

func (db Database) CreateConnector(ctx context.Context, connector *Connector) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Create(connector)
	if tx.Error != nil {
		return fmt.Errorf("error creating connector record ID %s: %w", connector.ConnectorID, tx.Error)
	}
	return nil
}

func (db Database) UpdateConnector(ctx context.Context, connector *Connector) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("id = ?", connector.ID). // Assuming update is by local DB ID
		Updates(connector)             // Updates non-zero fields
	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating connector local ID %d: %w", connector.ID, tx.Error)
	}
	return nil
}

// DeleteConnector deletes by Dex Connector ID (string)
func (db Database) DeleteConnector(ctx context.Context, connectorID string) error {
	// Use WithContext
	tx := db.Orm.WithContext(ctx).
		Where("connector_id = ?", connectorID).
		Delete(&Connector{})

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Or maybe return nil? If Dex delete succeeded, local missing isn't fatal.
	}
	if tx.Error != nil {
		return fmt.Errorf("error deleting connector record ID %s: %w", connectorID, tx.Error)
	}
	return nil
}

func (db Database) GetConnectorByConnectorID(ctx context.Context, connectorID string) (*Connector, error) {
	var s Connector
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("connector_id = ?", connectorID).
		First(&s) // Use First
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting connector by connector ID %s: %w", connectorID, tx.Error)
	}
	return &s, nil
}

func (db Database) GetConnectorByConnectorType(ctx context.Context, connectorType string) (*Connector, error) {
	var s Connector
	// Use WithContext
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("connector_type = ?", connectorType).
		First(&s) // Use First, assuming only one per type? Might need Find.
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting connector by type %s: %w", connectorType, tx.Error)
	}
	return &s, nil
}

// Methods below were in original db.go but maybe not used by http handlers / not in interface:
// DisableUser, EnableUser, DeActiveUser, ActiveUser - These used uuid.UUID which differs
// from the uint ID used elsewhere. They are omitted from the interface for now unless needed.
// If needed, reconcile ID types or add specific methods to interface.
