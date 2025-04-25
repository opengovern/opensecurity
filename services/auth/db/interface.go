// Package db provides the data persistence layer for the auth service,
// handling interactions with the PostgreSQL database using GORM.
package db

import (
	"context" // Added context
	"time"    // For time parameters/returns

	// Use the correct import path for your project's api.Role
	"github.com/opengovern/og-util/pkg/api"
	// Import uuid if needed by signatures (it's used in implementation but not signatures here)
	// "github.com/google/uuid"
)

// DatabaseInterface defines the set of database operations required by the auth service's
// core logic and HTTP handlers. Using an interface allows for easier testing (mocking)
// and potentially swapping database implementations in the future.
type DatabaseInterface interface {
	// --- User Operations ---

	// GetUserByEmail retrieves a user by their email address.
	// Returns nil, nil if the user is not found (following gorm.ErrRecordNotFound convention).
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// GetUserByExternalID retrieves a user by their external identifier (e.g., "connector|id").
	// Returns nil, nil if the user is not found.
	GetUserByExternalID(ctx context.Context, id string) (*User, error)

	// GetUser retrieves a user by their internal database ID (primary key, provided as string).
	// Returns nil, nil if the user is not found.
	GetUser(ctx context.Context, id string) (*User, error)

	// GetUsers retrieves a list of all users, typically ordered by last update time.
	GetUsers(ctx context.Context) ([]User, error)

	// GetUsersCount returns the total number of users in the database.
	GetUsersCount(ctx context.Context) (int64, error)

	// GetFirstUser retrieves the user with the lowest ID (often the initial admin).
	// Returns nil, nil if no users exist.
	GetFirstUser(ctx context.Context) (*User, error)

	// CreateUser creates a new user record or updates an existing one based on conflict resolution (e.g., OnConflict).
	CreateUser(ctx context.Context, user *User) error
	// EnableUser marks a user account as active (is_active=true).
	// Takes the internal database ID (uint).
	EnableUser(ctx context.Context, id uint) error // Added for handler

	// UpdateUser updates specific fields of an existing user identified by user.ID.
	// Implementations typically update non-zero fields of the provided user struct.
	UpdateUser(ctx context.Context, user *User) error
	// DisableUser marks a user account as inactive (is_active=false).
	// Takes the internal database ID (uint).
	DisableUser(ctx context.Context, id uint) error // <-- It's defined here!

	// DeleteUser deletes a user by their internal database ID (uint).
	DeleteUser(ctx context.Context, id uint) error

	// UpdateUserLastLoginWithExternalID updates the last_login timestamp for a user identified by external ID.
	UpdateUserLastLoginWithExternalID(ctx context.Context, id string, lastLogin time.Time) error

	// UserPasswordUpdate marks the user's password change requirement as false, typically after a password reset or update.
	UserPasswordUpdate(ctx context.Context, id uint) error

	// FindIdByEmail retrieves the internal database ID (uint) for a user by their email address.
	// Returns 0, error if not found or on other errors.
	FindIdByEmail(ctx context.Context, email string) (uint, error)

	// --- API Key Operations ---

	// AddApiKey creates a new API key record in the database.
	AddApiKey(ctx context.Context, key *ApiKey) error

	// CountApiKeysForUser returns the number of *active* API keys for a specific user (identified by External ID).
	CountApiKeysForUser(ctx context.Context, userID string) (int64, error)

	// ListApiKeysForUser retrieves all API keys created by a specific user (identified by External ID).
	ListApiKeysForUser(ctx context.Context, userId string) ([]ApiKey, error)

	// UpdateAPIKey updates the status (active/inactive) and role of an API key identified by its stringified internal ID.
	UpdateAPIKey(ctx context.Context, id string, isActive bool, role api.Role) error

	// DeleteAPIKey deletes an API key by its internal database ID (uint64).
	DeleteAPIKey(ctx context.Context, id uint64) error

	// ListApiKeys retrieves all API keys stored in the database.
	ListApiKeys(ctx context.Context) ([]ApiKey, error)

	// --- Connector Operations ---

	// GetConnectorByConnectorID retrieves connector metadata using the Dex connector ID string (e.g., "oidc-google").
	// Returns nil, nil if not found.
	GetConnectorByConnectorID(ctx context.Context, connectorID string) (*Connector, error)

	// CreateConnector creates a new connector metadata record in the local database.
	CreateConnector(ctx context.Context, connector *Connector) error

	// UpdateConnector updates an existing connector metadata record identified by its local database ID (connector.ID).
	UpdateConnector(ctx context.Context, connector *Connector) error

	// DeleteConnector deletes a connector metadata record using the Dex connector ID string.
	DeleteConnector(ctx context.Context, connectorID string) error

	// GetConnectors retrieves all connector metadata records from the local database.
	GetConnectors(ctx context.Context) ([]Connector, error)

	// GetConnector retrieves connector metadata by its local database ID (stringified uint).
	// Returns nil, nil if not found.
	GetConnector(ctx context.Context, id string) (*Connector, error)

	// GetConnectorByConnectorType retrieves connector metadata by its type string (e.g., "oidc").
	// Returns nil, nil if not found. Note: May return multiple if not unique.
	GetConnectorByConnectorType(ctx context.Context, connectorType string) (*Connector, error)

	// --- Configuration Operations ---

	// GetKeyPair retrieves the platform's RSA public and private key configurations from the database.
	GetKeyPair(ctx context.Context) ([]Configuration, error)

	// AddConfiguration adds a new key-value configuration entry to the database.
	AddConfiguration(ctx context.Context, c *Configuration) error

	// --- Initialization ---

	// Initialize performs necessary database setup, such as auto-migration of tables.
	// Context is usually not required for AutoMigrate.
	Initialize() error
}

// Optional: Compile-time check to ensure the real Database struct implements this interface.
// Place this below the struct definition in db.go or near the interface.
// var _ DatabaseInterface = (*Database)(nil) // Check against pointer receiver if methods use pointer receiver
