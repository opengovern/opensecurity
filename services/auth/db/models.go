// Package db provides the data persistence layer for the auth service,
// handling interactions with the PostgreSQL database using GORM.
package db

import (
	// Use the correct import path for your project's api.Role
	"time"

	"github.com/opengovern/og-util/pkg/api"
	"gorm.io/gorm"
)

// Configuration represents a generic key-value configuration setting stored in the database.
// It is currently used to store the platform's RSA key pair if not provided via environment variables.
type Configuration struct {
	gorm.Model        // Embeds standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).
	Key        string `gorm:"uniqueIndex"` // The unique key identifying the configuration setting (e.g., "public_key", "private_key").
	Value      string // The value associated with the key (e.g., base64 encoded PEM key).
}

// ApiKey represents a platform-generated API key stored in the database.
// These keys are JWTs signed by the platform's private key and grant specific roles.
type ApiKey struct {
	gorm.Model
	Name          string   `gorm:"index"` // User-defined name for the API key for easy identification.
	Role          api.Role // The role granted by this API key (e.g., Admin, Editor, Viewer).
	CreatorUserID string   `gorm:"index"`       // The ExternalId of the user who created this API key.
	IsActive      bool     `gorm:"index"`       // Flag indicating if the API key is currently active and usable.
	KeyHash       string   `gorm:"uniqueIndex"` // A cryptographic hash (e.g., SHA512) of the full JWT token, used for potential revocation/lookup (though direct verification uses public key).
	MaskedKey     string   // A partially obscured version of the JWT token for display purposes (e.g., "abc...xyz").
}

// Connector represents metadata stored locally about an identity provider connector configured in Dex.
// This complements the configuration stored within Dex itself.
type Connector struct {
	gorm.Model
	UserCount        uint      `gorm:"default:0"`       // A count of users associated with this connector (usage/update mechanism not shown).
	ConnectorID      string    `gorm:"uniqueIndex"`     // The unique identifier used by Dex for this connector (e.g., "oidc-google", "auth0").
	ConnectorType    string    `gorm:"index"`           // The type of the connector (e.g., "oidc", "saml").
	ConnectorSubType string    `gorm:"index,omitempty"` // The specific subtype, if applicable (e.g., "general", "entraid", "google-workspace").
	LastUpdate       time.Time // Timestamp of the last update to this connector's metadata in the local DB.
}

// User represents a user account within the platform.
// It stores core identity information, application-specific role, status,
// and links to the authentication method (connector/external ID).
type User struct {
	gorm.Model
	Email                 string    `gorm:"uniqueIndex"` // User's unique email address, used as a primary identifier for lookups.
	EmailVerified         bool      // Flag indicating if the user's email address has been verified (usually via the IdP).
	FullName              string    // User's full name.
	Role                  api.Role  `gorm:"index"`       // The user's assigned role within the application (e.g., Admin, Editor, Viewer).
	ConnectorId           string    `gorm:"index"`       // Identifier of the connector used for authentication (e.g., "local", "oidc-google").
	ExternalId            string    `gorm:"uniqueIndex"` // The canonical, unique identifier for the user across the system, typically "connectorId|userIdFromProvider".
	LastLogin             time.Time // Timestamp of the last recorded successful login or verified activity.
	Username              string    `gorm:"index"`        // User's username.
	RequirePasswordChange bool      `gorm:"default:true"` // Flag indicating if a local user must change their password on next login.
	IsActive              bool      `gorm:"default:true"` // Flag indicating if the user account is active and allowed to log in.
}
