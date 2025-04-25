// Package db provides the data persistence layer for the auth service,
// handling interactions with the PostgreSQL database using GORM.
package db

import (
	"context" // Added context import
	"errors"
	"fmt"
	"time"

	// Keep existing imports
	// Used only in commented-out methods
	"github.com/opengovern/og-util/pkg/api"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Database struct holds the GORM database connection pool (*gorm.DB).
// Methods on this struct implement the DatabaseInterface.
type Database struct {
	Orm *gorm.DB
}

// Initialize performs GORM auto-migration for all necessary database models.
// It ensures the tables corresponding to ApiKey, User, Configuration, and Connector
// exist and have the correct schema. Context is typically not required for AutoMigrate.
func (db *Database) Initialize() error { // Changed receiver to pointer for consistency
	err := db.Orm.AutoMigrate(
		&ApiKey{},
		&User{},
		&Configuration{},
		&Connector{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}
	return nil
}

// GetKeyPair retrieves Configuration entries matching specific keys,
// typically used to load the platform's RSA public and private keys.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetKeyPair(ctx context.Context) ([]Configuration, error) {
	var s []Configuration
	tx := db.Orm.WithContext(ctx).Model(&Configuration{}).
		Where("key = ? OR key = ?", "private_key", "public_key").Find(&s) // Use placeholders
	if tx.Error != nil {
		return nil, fmt.Errorf("error getting key pair: %w", tx.Error)
	}
	return s, nil
}

// AddConfiguration creates a new key-value configuration record in the database.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) AddConfiguration(ctx context.Context, c *Configuration) error {
	tx := db.Orm.WithContext(ctx).Create(c)
	if tx.Error != nil {
		return fmt.Errorf("error adding configuration key=%s: %w", c.Key, tx.Error)
	}
	return nil
}

// ListApiKeys retrieves all API key records from the database, ordered by creation time descending.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) ListApiKeys(ctx context.Context) ([]ApiKey, error) {
	var s []ApiKey
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error listing all API keys: %w", tx.Error)
	}
	return s, nil
}

// ListApiKeysForUser retrieves all API key records created by a specific user,
// identified by their external user ID (userId). Ordered by creation time descending.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) ListApiKeysForUser(ctx context.Context, userId string) ([]ApiKey, error) {
	var s []ApiKey
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Where("creator_user_id = ?", userId).
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error listing API keys for user %s: %w", userId, tx.Error)
	}
	return s, nil
}

// AddApiKey creates a new API key record in the database.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) AddApiKey(ctx context.Context, key *ApiKey) error {
	tx := db.Orm.WithContext(ctx).Create(key)
	if tx.Error != nil {
		return fmt.Errorf("error adding API key '%s': %w", key.Name, tx.Error)
	}
	return nil
}

// UpdateAPIKey updates the 'is_active' status and 'role' for an existing API key,
// identified by its stringified internal database ID.
// Returns gorm.ErrRecordNotFound if no key with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) UpdateAPIKey(ctx context.Context, id string, is_active bool, role api.Role) error {
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"is_active": is_active, "role": role})

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that the record was not found
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating API key ID %s: %w", id, tx.Error)
	}
	return nil
}

// DeleteAPIKey deletes an API key record by its internal database ID (uint64).
// Returns gorm.ErrRecordNotFound if no key with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) DeleteAPIKey(ctx context.Context, id uint64) error {
	tx := db.Orm.WithContext(ctx).
		Where("id = ?", id).
		Delete(&ApiKey{})

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error deleting API key ID %d: %w", id, tx.Error)
	}
	return nil
}

// CreateUser creates a new user record. It uses clause.OnConflict to perform
// an "upsert" based on the 'id' column, updating specified columns if a conflict occurs.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) CreateUser(ctx context.Context, user *User) error {
	tx := db.Orm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // Assuming 'id' is the primary key for conflict resolution
		DoUpdates: clause.AssignmentColumns([]string{"updated_at", "email", "email_verified", "role", "connector_id", "external_id", "full_name", "last_login", "username", "is_active"}),
	}).Create(user)

	if tx.Error != nil {
		return fmt.Errorf("error creating/updating user email=%s: %w", user.Email, tx.Error)
	}
	return nil
}

// UpdateUser updates an existing user record based on the provided user struct.
// It identifies the user by user.ID and updates non-zero fields using gorm.Updates.
// Returns gorm.ErrRecordNotFound if no user with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) UpdateUser(ctx context.Context, user *User) error {
	// Ensure user.ID is set before calling this method.
	if user.ID == 0 {
		return errors.New("cannot update user with zero ID")
	}
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ?", user.ID).
		Updates(user) // Updates non-zero fields

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating user ID %d: %w", user.ID, tx.Error)
	}
	return nil
}

// DeleteUser deletes a user record by its internal database ID (uint).
// Returns gorm.ErrRecordNotFound if no user with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) DeleteUser(ctx context.Context, id uint) error {
	tx := db.Orm.WithContext(ctx).
		Where("id = ?", id).
		Delete(&User{})

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error deleting user ID %d: %w", id, tx.Error)
	}
	return nil
}

// GetUsers retrieves all user records from the database, ordered by last update time descending.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetUsers(ctx context.Context) ([]User, error) {
	var s []User
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Order("updated_at desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error getting all users: %w", tx.Error)
	}
	return s, nil
}

// GetUser retrieves a single user record by its internal database ID (stringified uint).
// Returns nil, nil if the user is not found.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetUser(ctx context.Context, id string) (*User, error) {
	var s User
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ?", id).
		First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Convention for not found
		}
		return nil, fmt.Errorf("error getting user by ID %s: %w", id, tx.Error)
	}
	return &s, nil
}

// GetUserByExternalID retrieves a single user record by their external identifier string.
// Returns nil, nil if the user is not found.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetUserByExternalID(ctx context.Context, id string) (*User, error) {
	var s User
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("external_id = ?", id).
		First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Convention for not found
		}
		return nil, fmt.Errorf("error getting user by external ID %s: %w", id, tx.Error)
	}
	return &s, nil
}

// GetUsersCount returns the total number of user records in the database.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetUsersCount(ctx context.Context) (int64, error) {
	var count int64
	tx := db.Orm.WithContext(ctx).Model(&User{}).Count(&count)
	if tx.Error != nil {
		return 0, fmt.Errorf("error counting users: %w", tx.Error)
	}
	return count, nil
}

// GetFirstUser retrieves the first user record based on ascending ID order.
// Typically used for admin bootstrap scenarios. Returns nil, nil if no users exist.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetFirstUser(ctx context.Context) (*User, error) {
	var user User
	tx := db.Orm.WithContext(ctx).Model(&User{}).Order("id asc").First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // No users exist
		}
		return nil, fmt.Errorf("error getting first user: %w", tx.Error)
	}
	return &user, nil
}

// UpdateUserLastLogin updates the last_login timestamp for a user identified by internal ID (uuid.UUID).
// Deprecated or needs adjustment if internal ID is consistently uint.
// func (db *Database) UpdateUserLastLogin(ctx context.Context, id uuid.UUID, lastLogin *time.Time) error {
// 	if lastLogin == nil || lastLogin.IsZero() { return nil } // Avoid zero time updates
// 	tx := db.Orm.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("last_login", lastLogin)
// 	if tx.Error == nil && tx.RowsAffected == 0 { return gorm.ErrRecordNotFound }
// 	if tx.Error != nil { return fmt.Errorf("error updating last login for user ID %s: %w", id.String(), tx.Error) }
// 	return nil
// }

// UpdateUserLastLoginWithExternalID updates the last_login timestamp for a user identified by their external ID string.
// Skips update if lastLogin time is zero. Returns nil if user not found (non-critical).
// Accepts a context for cancellation/timeout propagation.
func (db *Database) UpdateUserLastLoginWithExternalID(ctx context.Context, id string, lastLogin time.Time) error {
	if lastLogin.IsZero() {
		return nil
	} // Avoid zero time updates
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("external_id = ?", id).
		Update("last_login", lastLogin)

	if tx.Error == nil && tx.RowsAffected == 0 {
		// User not found is not treated as an error for this non-critical update.
		return nil
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating last login for external ID %s: %w", id, tx.Error)
	}
	return nil
}

// GetUserByEmail retrieves a single user record by their email address.
// Returns nil, nil if the user is not found.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var s User
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("email = ?", email).
		First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Convention for not found
		}
		return nil, fmt.Errorf("error getting user by email %s: %w", email, tx.Error)
	}
	return &s, nil
}

// UserPasswordUpdate sets the 'require_password_change' flag to false for a user identified by internal ID (uint).
// Returns gorm.ErrRecordNotFound if no user with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) UserPasswordUpdate(ctx context.Context, id uint) error {
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

// FindIdByEmail retrieves the internal database ID (uint) for a user by their email address.
// Returns 0 and an error if the user is not found or if another error occurs.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) FindIdByEmail(ctx context.Context, email string) (uint, error) {
	var s User
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("email = ?", email).
		Select("id"). // Only select the ID field
		First(&s)
	if tx.Error != nil {
		// Return 0 for ID on any error, including not found.
		return 0, fmt.Errorf("error finding ID for email %s: %w", email, tx.Error)
	}
	return s.ID, nil
}

// CountApiKeysForUser counts the number of *active* API keys associated with a given external user ID.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) CountApiKeysForUser(ctx context.Context, userID string) (int64, error) {
	var s int64
	tx := db.Orm.WithContext(ctx).Model(&ApiKey{}).
		Where("creator_user_id = ?", userID).
		Where("is_active = ?", true). // Use boolean true
		Count(&s)
	if tx.Error != nil {
		return 0, fmt.Errorf("error counting API keys for user %s: %w", userID, tx.Error)
	}
	return s, nil
}

// GetConnectors retrieves all connector metadata records from the local database, ordered by last update time descending.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetConnectors(ctx context.Context) ([]Connector, error) {
	var s []Connector
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Order("last_update desc").
		Find(&s)
	if tx.Error != nil {
		return nil, fmt.Errorf("error getting all connectors: %w", tx.Error)
	}
	return s, nil
}

// GetConnector retrieves a single connector metadata record by its internal database ID (stringified uint).
// Returns nil, nil if not found.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetConnector(ctx context.Context, id string) (*Connector, error) {
	var s Connector
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("id = ?", id).
		First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting connector by local ID %s: %w", id, tx.Error)
	}
	return &s, nil
}

// CreateConnector creates a new connector metadata record in the local database.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) CreateConnector(ctx context.Context, connector *Connector) error {
	tx := db.Orm.WithContext(ctx).Create(connector)
	if tx.Error != nil {
		return fmt.Errorf("error creating connector record ID %s: %w", connector.ConnectorID, tx.Error)
	}
	return nil
}

// UpdateConnector updates an existing connector metadata record identified by its local database ID (connector.ID).
// It typically updates non-zero fields of the provided struct.
// Returns gorm.ErrRecordNotFound if no record with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) UpdateConnector(ctx context.Context, connector *Connector) error {
	if connector.ID == 0 {
		return errors.New("cannot update connector with zero ID")
	}
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("id = ?", connector.ID).
		Updates(connector) // Updates non-zero fields
	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error updating connector local ID %d: %w", connector.ID, tx.Error)
	}
	return nil
}

// DeleteConnector deletes a connector metadata record identified by its Dex Connector ID string.
// Returns gorm.ErrRecordNotFound if no record with the given ConnectorID exists.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) DeleteConnector(ctx context.Context, connectorID string) error {
	tx := db.Orm.WithContext(ctx).
		Where("connector_id = ?", connectorID).
		Delete(&Connector{})

	if tx.Error == nil && tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if tx.Error != nil {
		return fmt.Errorf("error deleting connector record ID %s: %w", connectorID, tx.Error)
	}
	return nil
}

// GetConnectorByConnectorID retrieves a single connector metadata record by its Dex Connector ID string.
// Returns nil, nil if not found.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetConnectorByConnectorID(ctx context.Context, connectorID string) (*Connector, error) {
	var s Connector
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("connector_id = ?", connectorID).
		First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting connector by connector ID %s: %w", connectorID, tx.Error)
	}
	return &s, nil
}

// GetConnectorByConnectorType retrieves the first connector metadata record matching the given type string.
// Use with caution if multiple connectors of the same type can exist.
// Returns nil, nil if not found.
// Accepts a context for cancellation/timeout propagation.
func (db *Database) GetConnectorByConnectorType(ctx context.Context, connectorType string) (*Connector, error) {
	var s Connector
	tx := db.Orm.WithContext(ctx).Model(&Connector{}).
		Where("connector_type = ?", connectorType).
		First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting connector by type %s: %w", connectorType, tx.Error)
	}
	return &s, nil
}

// DisableUser marks a user account as inactive by setting its 'is_active' field to false.
// It identifies the user by their internal database ID (uint).
// Returns gorm.ErrRecordNotFound if no user with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
// Note: Role-based authorization checks should be performed in the calling layer (e.g., HTTP handler)
// *before* invoking this function.
func (db *Database) DisableUser(ctx context.Context, id uint) error { // Changed id type to uint
	// Use WithContext and update the 'is_active' column to false.
	tx := db.Orm.WithContext(ctx).Model(&User{}). // Specify the model for the table name
							Where("id = ?", id).       // Find user by primary key ID
							Update("is_active", false) // Set the active flag to false

	// Check for errors during the update operation.
	if tx.Error != nil {
		return fmt.Errorf("error disabling user ID %d: %w", id, tx.Error)
	}

	// Check if any row was actually updated. If not, the user was not found.
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return standard "not found" error
	}

	// Return nil on successful update.
	return nil
}

// --- Commented out methods using uuid.UUID ---
// These methods were present in the original code but used uuid.UUID for ID,
// which seems inconsistent with other methods using uint or string.
// They are commented out here but kept for reference. If needed, the interface
// and callers need to be adjusted to handle the correct ID type consistently.

/*

func (db Database) DisableUser(ctx context.Context, id uuid.UUID) error {
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ? ", id). // GORM might handle UUID automatically, check documentation
		Update("is_active", false) // Assuming 'disabled' was a typo for '!is_active'

	if tx.Error == nil && tx.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	if tx.Error != nil { return fmt.Errorf("error disabling user ID %s: %w", id.String(), tx.Error) }
	return nil
}

func (db Database) EnableUser(ctx context.Context, id uuid.UUID) error {
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ? ", id).
		Update("is_active", true) // Assuming 'Enable' means setting IsActive to true

	if tx.Error == nil && tx.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	if tx.Error != nil { return fmt.Errorf("error enabling user ID %s: %w", id.String(), tx.Error) }
	return nil
}

// DeActiveUser seems redundant with DisableUser if it just sets is_active=false
func (db Database) DeActiveUser(ctx context.Context, id uuid.UUID) error {
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ? ", id).
		Update("is_active", false)

	if tx.Error == nil && tx.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	if tx.Error != nil { return fmt.Errorf("error deactivating user ID %s: %w", id.String(), tx.Error) }
	return nil
}

// ActiveUser seems redundant with EnableUser if it just sets is_active=true
func (db Database) ActiveUser(ctx context.Context, id uuid.UUID) error {
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ? ", id).
		Update("is_active", true)

	if tx.Error == nil && tx.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	if tx.Error != nil { return fmt.Errorf("error activating user ID %s: %w", id.String(), tx.Error) }
	return nil
}
*/
// EnableUser marks a user account as active by setting its 'is_active' field to true.
// It identifies the user by their internal database ID (uint).
// Returns gorm.ErrRecordNotFound if no user with the given ID exists.
// Accepts a context for cancellation/timeout propagation.
// Note: Role-based authorization checks should be performed in the calling layer (e.g., HTTP handler)
// *before* invoking this function.
func (db *Database) EnableUser(ctx context.Context, id uint) error { // Using pointer receiver
	// Check if the ORM connection is initialized
	if db.Orm == nil {
		return errors.New("database connection not initialized")
	}

	// Use WithContext and update the 'is_active' column to true.
	// Specify the model to ensure GORM targets the correct table.
	tx := db.Orm.WithContext(ctx).Model(&User{}).
		Where("id = ?", id).      // Find user by primary key ID
		Update("is_active", true) // Set the active flag to true

	// Check for errors during the update operation.
	if tx.Error != nil {
		// Check if the error is specific to the record not being found,
		// although RowsAffected check below is more reliable for updates.
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return fmt.Errorf("error enabling user ID %d: %w", id, tx.Error)
	}

	// Check if any row was actually updated. If RowsAffected is 0 after a nil error,
	// it means the WHERE clause (id = ?) didn't match any rows.
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return standard "not found" error
	}

	// Return nil on successful update.
	return nil
}
