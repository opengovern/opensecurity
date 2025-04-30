package db

import (
	"errors"
	"fmt"
	"strconv" // Added import
	"time"

	"github.com/google/uuid"
	"github.com/opengovern/og-util/pkg/api"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Database struct definition
type Database struct {
	Orm    *gorm.DB
	Logger *zap.Logger // Logger for database operations
}

// --- Configuration Methods ---

func (db Database) GetKeyPair() ([]Configuration, error) {
	var s []Configuration
	// Use Model() for clarity, Where condition is fine
	tx := db.Orm.Model(&Configuration{}).
		Where("key = ? OR key = ?", "private_key", "public_key"). // Use placeholders
		Find(&s)
	if tx.Error != nil {
		db.Logger.Error("Failed to get key pair", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) AddConfiguration(c *Configuration) error {
	// Basic create operation
	tx := db.Orm.Create(c)
	if tx.Error != nil {
		db.Logger.Error("Failed to add configuration", zap.String("key", c.Key), zap.Error(tx.Error))
		return tx.Error
	}
	db.Logger.Info("Added configuration", zap.String("key", c.Key))
	return nil
}

// --- API Key Methods ---

func (db Database) ListApiKeys() ([]ApiKey, error) {
	var s []ApiKey
	// Order and find all keys
	tx := db.Orm.Model(&ApiKey{}).
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		db.Logger.Error("Failed to list all API keys", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) ListApiKeysForUser(userId string) ([]ApiKey, error) {
	var s []ApiKey
	// Filter by creator_user_id
	tx := db.Orm.Model(&ApiKey{}).
		Where("creator_user_id = ?", userId). // Use placeholder
		Order("created_at desc").
		Find(&s)
	if tx.Error != nil {
		db.Logger.Error("Failed to list API keys for user", zap.String("userId", userId), zap.Error(tx.Error))
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) AddApiKey(key *ApiKey) error {
	// Basic create operation
	tx := db.Orm.Create(key)
	if tx.Error != nil {
		db.Logger.Error("Failed to add API key", zap.String("name", key.Name), zap.String("userId", key.CreatorUserID), zap.Error(tx.Error))
		return tx.Error
	}
	db.Logger.Info("Added API key", zap.String("name", key.Name), zap.Uint("id", key.ID))
	return nil
}

func (db Database) UpdateAPIKey(id string, is_active bool, role api.Role) error {
	// Convert string ID to uint64
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		db.Logger.Error("Invalid API Key ID format for update", zap.String("id", id), zap.Error(err))
		return fmt.Errorf("invalid API Key ID format: %w", err)
	}

	// Use Updates with a map for clarity and specific field updates
	updateData := map[string]interface{}{
		"is_active": is_active,
		"role":      role,
	}
	tx := db.Orm.Model(&ApiKey{}).
		Where("id = ?", idUint).
		Updates(updateData)

	if tx.Error != nil {
		db.Logger.Error("Failed to update API key", zap.Uint64("id", idUint), zap.Error(tx.Error))
		return tx.Error
	}
	// Check if the record was actually found and updated
	if tx.RowsAffected == 0 {
		db.Logger.Warn("API key not found for update", zap.Uint64("id", idUint))
		return gorm.ErrRecordNotFound // Return specific error if not found
	}
	db.Logger.Info("Updated API key", zap.Uint64("id", idUint))
	return nil
}

func (db Database) DeleteAPIKey(id uint64) error {
	// Basic delete operation
	tx := db.Orm.Model(&ApiKey{}). // Specify model for clarity, although Where+Delete works
					Where("id = ?", id).
					Delete(&ApiKey{})
	if tx.Error != nil {
		db.Logger.Error("Failed to delete API key", zap.Uint64("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	// Check if the record was found and deleted
	if tx.RowsAffected == 0 {
		db.Logger.Warn("API key not found for deletion", zap.Uint64("id", id))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Deleted API key", zap.Uint64("id", id))
	return nil
}

func (db Database) CountApiKeysForUser(userID string) (int64, error) {
	var count int64
	// Combine WHERE clauses, use boolean true for is_active
	tx := db.Orm.Model(&ApiKey{}).
		Where("creator_user_id = ? AND is_active = ?", userID, true).
		Count(&count)
	if tx.Error != nil {
		db.Logger.Error("Failed to count active API keys for user", zap.String("userId", userID), zap.Error(tx.Error))
		return 0, tx.Error
	}
	return count, nil
}

// --- User Methods ---

func (db Database) CreateUser(user *User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	// Ensure ExternalId is provided if it's the conflict target
	if user.ExternalId == "" {
		db.Logger.Error("CreateUser attempted with empty ExternalId, which is used for OnConflict")
		return errors.New("external ID is required for user creation/upsert")
	}
	// Use external_id as the conflict target (ensure it has a unique index in DB)
	tx := db.Orm.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "external_id"}}, // Conflict on external_id
		DoUpdates: clause.AssignmentColumns([]string{ // Columns to update on conflict
			"updated_at", "email", "email_verified", "role", "connector_id",
			"full_name", "last_login", "username", "is_active",
		}),
	}).Create(user)

	if tx.Error != nil {
		db.Logger.Error("Failed to create/update user via OnConflict", zap.Error(tx.Error), zap.String("externalId", user.ExternalId), zap.String("email", user.Email))
		return tx.Error // Return the original GORM error
	}
	// Log success, potentially distinguishing between create and update if needed (check tx.RowsAffected maybe, though behavior varies)
	db.Logger.Info("Successfully created or updated user", zap.String("externalId", user.ExternalId), zap.Uint("id", user.ID))
	return nil
}

func (db Database) UpdateUser(user *User) error {
	// Ensure ID is present for targeting the update
	if user.ID == 0 {
		return errors.New("user ID must be provided for update")
	}
	// Use map for explicit updates
	updateData := map[string]interface{}{
		"email":          user.Email,
		"email_verified": user.EmailVerified,
		"full_name":      user.FullName,
		"role":           user.Role,
		"connector_id":   user.ConnectorId,
		"external_id":    user.ExternalId, // Be careful updating this if it's a key identifier
		"last_login":     user.LastLogin,
		"username":       user.Username,
		"is_active":      user.IsActive, // Explicitly include is_active
	}
	tx := db.Orm.Model(&User{}). // Target the User model
					Where("id = ?", user.ID).
					Updates(updateData)

	if tx.Error != nil {
		db.Logger.Error("Failed to update user", zap.Uint("id", user.ID), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("User not found for update", zap.Uint("id", user.ID))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Updated user", zap.Uint("id", user.ID))
	return nil
}

func (db Database) DeleteUser(id uint) error {
	// Add check for deleting the first user if necessary (business logic)
	// if id == 1 {
	// 	db.Logger.Warn("Attempted to delete the first user (ID 1)", zap.Uint("id", id))
	// 	return errors.New("cannot delete the first user")
	// }
	tx := db.Orm.
		Where("id = ?", id).
		Delete(&User{}) // Specify the model type
	if tx.Error != nil {
		db.Logger.Error("Failed to delete user", zap.Uint("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("User not found for deletion", zap.Uint("id", id))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Deleted user", zap.Uint("id", id))
	return nil
}

func (db Database) GetUsers() ([]User, error) {
	var s []User
	tx := db.Orm.Model(&User{}).
		Order("updated_at desc"). // Consider ordering by a more stable field like created_at or id if updates are frequent
		Find(&s)
	if tx.Error != nil {
		db.Logger.Error("Failed to get all users", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) GetUser(id string) (*User, error) {
	// Convert string ID to uint
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		db.Logger.Error("Invalid user ID format for GetUser", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	var s User
	// Use First for primary key lookup
	tx := db.Orm.First(&s, idUint)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			db.Logger.Debug("User not found by ID", zap.Uint64("id", idUint))
			return nil, nil // Return nil, nil for not found
		}
		db.Logger.Error("Failed to get user by ID", zap.Uint64("id", idUint), zap.Error(tx.Error))
		return nil, tx.Error // Return other GORM errors
	}
	return &s, nil
}

func (db Database) GetUserByExternalID(id string) (*User, error) {
	var s User
	// Use First for unique key lookup
	tx := db.Orm.Where("external_id = ?", id).First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			db.Logger.Debug("User not found by external ID", zap.String("externalId", id))
			return nil, nil
		}
		db.Logger.Error("Failed to get user by external ID", zap.String("externalId", id), zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &s, nil
}

func (db Database) GetUsersCount() (int64, error) {
	var count int64
	tx := db.Orm.Model(&User{}).Count(&count)
	if tx.Error != nil {
		db.Logger.Error("Failed to count users", zap.Error(tx.Error))
		return 0, tx.Error
	}
	return count, nil
}

func (db Database) GetFirstUser() (*User, error) {
	var user User
	// Order by ID for deterministic "first" user
	tx := db.Orm.Order("id asc").First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			db.Logger.Info("No users found in the database")
			return nil, nil // No users exist
		}
		db.Logger.Error("Failed to get first user", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &user, nil
}

func (db Database) UpdateUserLastLogin(id uuid.UUID, lastLogin *time.Time) error {
	// Validate input
	if lastLogin == nil {
		return errors.New("lastLogin time cannot be nil for UpdateUserLastLogin")
	}
	// Use Update for single column update
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id). // Assuming User.ID is uuid.UUID
		Update("last_login", lastLogin)

	if tx.Error != nil {
		db.Logger.Error("Failed to update last login by UUID", zap.Stringer("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		// This might happen if the user doesn't exist OR the time hasn't changed.
		// If distinguishing is critical, a prior SELECT might be needed.
		db.Logger.Warn("User not found or last login time unchanged (UUID)", zap.Stringer("id", id))
		// Return RecordNotFound only if user truly doesn't exist is desired.
		// For now, we don't return error if RowsAffected is 0.
	}
	return nil
}

func (db Database) UpdateUserLastLoginWithExternalID(id string, lastLogin time.Time) error {
	// Skip update if time is zero
	if lastLogin.IsZero() {
		db.Logger.Debug("Skipping last login update for zero time", zap.String("externalId", id))
		return nil
	}
	tx := db.Orm.Model(&User{}).
		Where("external_id = ?", id).
		Update("last_login", lastLogin) // Update single column

	if tx.Error != nil {
		db.Logger.Error("Failed to update last login by external ID", zap.String("externalId", id), zap.Error(tx.Error))
		return tx.Error
	}
	// Again, RowsAffected might be 0 if time hasn't changed.
	// if tx.RowsAffected == 0 {
	// 	db.Logger.Warn("User not found or last login time unchanged (External ID)", zap.String("externalId", id))
	// 	// return gorm.ErrRecordNotFound // Only if strict not found is required
	// }
	return nil
}

func (db Database) GetUserByEmail(email string) (*User, error) {
	var s User
	// Use First for unique email lookup
	tx := db.Orm.Where("email = ?", email).First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			db.Logger.Debug("User not found by email", zap.String("email", email))
			return nil, nil
		}
		db.Logger.Error("Failed to get user by email", zap.String("email", email), zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &s, nil
}

func (db Database) UserPasswordUpdate(id uint) error {
	// Update single column
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id).
		Update("require_password_change", false)

	if tx.Error != nil {
		db.Logger.Error("Failed to update require_password_change flag", zap.Uint("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("User not found for password update flag", zap.Uint("id", id))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Updated require_password_change flag for user", zap.Uint("id", id))
	return nil
}

// DisableUser sets is_active to false. Assumes User.ID is uuid.UUID.
func (db Database) DisableUser(id uuid.UUID) error {
	// Check if the User model actually has a 'disabled' field.
	// Based on other methods, it seems 'is_active' is used for this purpose.
	// Updating 'is_active' to false. If 'disabled' truly exists, change "is_active" to "disabled".
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id).
		Update("is_active", false) // Set user to inactive

	if tx.Error != nil {
		db.Logger.Error("Failed to disable user (set is_active=false)", zap.Stringer("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("User not found for disabling", zap.Stringer("id", id))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Disabled user (set is_active=false)", zap.Stringer("id", id))
	return nil
}

// EnableUser sets is_active to true. Assumes User.ID is uuid.UUID.
func (db Database) EnableUser(id uuid.UUID) error {
	// Corrected logic: set is_active to true
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id).
		Update("is_active", true) // Set user to active

	if tx.Error != nil {
		db.Logger.Error("Failed to enable user (set is_active=true)", zap.Stringer("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("User not found for enabling", zap.Stringer("id", id))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Enabled user (set is_active=true)", zap.Stringer("id", id))
	return nil
}

// DeActiveUser sets is_active to false. Assumes User.ID is uuid.UUID.
func (db Database) DeActiveUser(id uuid.UUID) error {
	// This is functionally identical to DisableUser based on current implementation
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id).
		Update("is_active", false)

	if tx.Error != nil {
		db.Logger.Error("Failed to deactivate user (set is_active=false)", zap.Stringer("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("User not found for deactivation", zap.Stringer("id", id))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Deactivated user (set is_active=false)", zap.Stringer("id", id))
	return nil
}

// ActiveUser sets is_active to true. Assumes User.ID is uuid.UUID.
func (db Database) ActiveUser(id uuid.UUID) error {
	// This is functionally identical to EnableUser
	tx := db.Orm.Model(&User{}).
		Where("id = ?", id).
		Update("is_active", true)

	if tx.Error != nil {
		db.Logger.Error("Failed to activate user (set is_active=true)", zap.Stringer("id", id), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("User not found for activation", zap.Stringer("id", id))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Activated user (set is_active=true)", zap.Stringer("id", id))
	return nil
}

func (db Database) FindIdByEmail(email string) (uint, error) {
	var s User // GORM needs struct to scan into, even if selecting one column
	// Select only ID for efficiency, use First for unique lookup
	tx := db.Orm.Model(&User{}).Select("id").Where("email = ?", email).First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			// Return 0 and the specific error for clarity
			return 0, gorm.ErrRecordNotFound
		}
		db.Logger.Error("Failed to find user ID by email", zap.String("email", email), zap.Error(tx.Error))
		return 0, tx.Error
	}
	return s.ID, nil
}

// --- Connector Methods ---

func (db Database) GetConnectors() ([]Connector, error) {
	var s []Connector
	tx := db.Orm.Model(&Connector{}).
		Order("last_update desc").
		Find(&s)
	if tx.Error != nil {
		db.Logger.Error("Failed to get connectors", zap.Error(tx.Error))
		return nil, tx.Error
	}
	return s, nil
}

func (db Database) GetConnector(id string) (*Connector, error) {
	// Convert string ID to uint
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		db.Logger.Error("Invalid connector ID format for GetConnector", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("invalid connector ID format: %w", err)
	}
	var s Connector
	// Use First for primary key lookup
	tx := db.Orm.First(&s, idUint)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			db.Logger.Debug("Connector not found by ID", zap.Uint64("id", idUint))
			return nil, nil
		}
		db.Logger.Error("Failed to get connector by ID", zap.Uint64("id", idUint), zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &s, nil
}

func (db Database) CreateConnector(connector *Connector) error {
	tx := db.Orm.Create(connector)
	if tx.Error != nil {
		db.Logger.Error("Failed to create connector", zap.String("connectorId", connector.ConnectorID), zap.String("type", connector.ConnectorType), zap.Error(tx.Error))
		return tx.Error
	}
	db.Logger.Info("Created connector", zap.Uint("id", connector.ID), zap.String("connectorId", connector.ConnectorID))
	return nil
}

func (db Database) UpdateConnector(connector *Connector) error {
	// Ensure ID is present
	if connector.ID == 0 {
		return errors.New("connector ID must be provided for update")
	}
	// Use Updates to update non-zero fields by default.
	// If specific fields need updating regardless of zero value, use a map.
	tx := db.Orm.Model(&Connector{}).
		Where("id = ?", connector.ID).
		Updates(connector) // Be aware this might skip zero values (like UserCount=0)
	// If you need to set UserCount to 0, use:
	// tx := db.Orm.Model(&Connector{}).Where("id = ?", connector.ID).Updates(map[string]interface{}{...fields...})

	if tx.Error != nil {
		db.Logger.Error("Failed to update connector", zap.Uint("id", connector.ID), zap.String("connectorId", connector.ConnectorID), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("Connector not found for update", zap.Uint("id", connector.ID))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Updated connector", zap.Uint("id", connector.ID), zap.String("connectorId", connector.ConnectorID))
	return nil
}

func (db Database) DeleteConnector(connectorID string) error {
	// Delete by the business key 'connector_id'
	tx := db.Orm.
		Where("connector_id = ?", connectorID).
		Delete(&Connector{})
	if tx.Error != nil {
		db.Logger.Error("Failed to delete connector by connector_id", zap.String("connectorId", connectorID), zap.Error(tx.Error))
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		db.Logger.Warn("Connector not found for deletion by connector_id", zap.String("connectorId", connectorID))
		return gorm.ErrRecordNotFound
	}
	db.Logger.Info("Deleted connector by connector_id", zap.String("connectorId", connectorID))
	return nil
}

func (db Database) GetConnectorByConnectorID(connectorID string) (*Connector, error) {
	var s Connector
	// Use First for unique key lookup
	tx := db.Orm.Where("connector_id = ?", connectorID).First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			db.Logger.Debug("Connector not found by connector_id", zap.String("connectorId", connectorID))
			return nil, nil
		}
		db.Logger.Error("Failed to get connector by connector_id", zap.String("connectorId", connectorID), zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &s, nil
}

func (db Database) GetConnectorByConnectorType(connectorType string) (*Connector, error) {
	var s Connector
	// Use First as we likely expect only one connector of a specific type?
	// If multiple are allowed, use Find(&s) instead.
	tx := db.Orm.Where("connector_type = ?", connectorType).First(&s)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			db.Logger.Debug("Connector not found by type", zap.String("type", connectorType))
			return nil, nil
		}
		db.Logger.Error("Failed to get connector by type", zap.String("type", connectorType), zap.Error(tx.Error))
		return nil, tx.Error
	}
	return &s, nil
}
