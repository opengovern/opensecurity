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
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByExternalID(ctx context.Context, id string) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUsersCount(ctx context.Context) (int64, error)
	GetFirstUser(ctx context.Context) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uint) error
	UpdateUserLastLoginWithExternalID(ctx context.Context, id string, lastLogin time.Time) error
	UserPasswordUpdate(ctx context.Context, id uint) error
	FindIdByEmail(ctx context.Context, email string) (uint, error) // Added based on original db.go

	// --- API Key Operations ---
	AddApiKey(ctx context.Context, key *ApiKey) error
	CountApiKeysForUser(ctx context.Context, userID string) (int64, error)
	ListApiKeysForUser(ctx context.Context, userId string) ([]ApiKey, error)
	UpdateAPIKey(ctx context.Context, id string, isActive bool, role api.Role) error
	DeleteAPIKey(ctx context.Context, id uint64) error
	ListApiKeys(ctx context.Context) ([]ApiKey, error) // Added based on original db.go

	// --- Connector Operations ---
	GetConnectorByConnectorID(ctx context.Context, connectorID string) (*Connector, error)
	CreateConnector(ctx context.Context, connector *Connector) error
	UpdateConnector(ctx context.Context, connector *Connector) error
	DeleteConnector(ctx context.Context, connectorID string) error
	GetConnectors(ctx context.Context) ([]Connector, error)                                    // Added based on original db.go
	GetConnector(ctx context.Context, id string) (*Connector, error)                           // Added based on original db.go
	GetConnectorByConnectorType(ctx context.Context, connectorType string) (*Connector, error) // Added

	// --- Configuration Operations ---
	GetKeyPair(ctx context.Context) ([]Configuration, error)      // Added based on original db.go, added context
	AddConfiguration(ctx context.Context, c *Configuration) error // Added based on original db.go, added context

	// --- Initialization (Might not be needed in interface, but included if called externally) ---
	Initialize() error // Context usually not needed for AutoMigrate
}

// Optional: Compile-time check to ensure the real Database struct implements this interface.
// Place this below the struct definition in db.go or near the interface.
