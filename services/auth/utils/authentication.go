// Package utils provides helper functions and shared logic for the auth service,
// including database interactions, connector configurations, and Kubernetes interactions.
// This specific file focuses on authentication-related utilities like user data mapping
// and potentially interacting with user data via the database interface.
package utils

import (
	"context" // Added context import
	"fmt"
	"time"

	// Ensure correct import path to where DatabaseInterface is defined
	"github.com/opengovern/opensecurity/services/auth/db"
	// Import other necessary packages, like API roles if used directly
	// "github.com/opengovern/og-util/pkg/api"
)

// Service might hold configuration or shared dependencies for utility functions
// within this package, such as domain info or client secrets if needed for
// operations beyond simple DB interactions (e.g., interacting with an external
// user management API).
// It now uses the db.DatabaseInterface for testability and abstraction.
type Service struct {
	domain       string
	clientID     string
	clientSecret string
	appClientID  string
	Connection   string
	InviteTTL    int

	database db.DatabaseInterface // Database access via interface.
}

// User represents the user data structure used within the utility package or
// potentially returned by API client functions. It may differ slightly from the
// database model (db.User), for example, by representing Role as a string.
type User struct {
	ID            uint      `json:"id"`             // Internal database primary key.
	CreatedAt     time.Time `json:"created_at"`     // Timestamp of user creation.
	UpdatedAt     time.Time `json:"updated_at"`     // Timestamp of last user update.
	Email         string    `json:"email"`          // User's email address.
	EmailVerified bool      `json:"email_verified"` // Whether the user's email has been verified.
	FullName      string    `json:"full_name"`      // User's full name.
	LastLogin     time.Time `json:"last_login"`     // Timestamp of the last recorded login.
	Username      string    `json:"username"`       // User's chosen username.
	Role          string    `json:"role"`           // User's application role (as a string).
	IsActive      bool      `json:"is_active"`      // Whether the user account is currently active.
	ConnectorId   string    `json:"connector_id"`   // Identifier of the connector used for authentication (e.g., "local", "oidc-google").
	ExternalId    string    `json:"external_id"`    // Canonical user identifier across the system (e.g., "connector|id").
}

// DbUserToApi converts a database user model (*db.User) to the utility/API layer
// user representation (*utils.User). It handles nil input gracefully, returning nil, nil.
// It also converts the Role type (likely an enum in db.User) to its string representation.
func DbUserToApi(u *db.User) (*User, error) {
	// If the input database user is nil (e.g., user not found), return nil, nil.
	if u == nil {
		return nil, nil
	}

	// Map fields from the database struct to the utility struct.
	return &User{
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt, // Assumes db.User embeds gorm.Model
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		FullName:      u.FullName,
		LastLogin:     u.LastLogin,
		Username:      u.Username,     // Assumes db.User.Username is string
		Role:          string(u.Role), // Convert db.User.Role (e.g., api.Role) to string
		ExternalId:    u.ExternalId,
		ID:            u.ID,
		IsActive:      u.IsActive,
		ConnectorId:   u.ConnectorId,
	}, nil
}

// New creates a new instance of the utility Service, injecting necessary
// configuration and the database interface dependency.
func New(domain, appClientID, clientID, clientSecret, connection string, inviteTTL int, database db.DatabaseInterface) *Service {
	return &Service{
		domain:       domain,
		appClientID:  appClientID,
		clientID:     clientID,
		clientSecret: clientSecret,
		Connection:   connection,
		InviteTTL:    inviteTTL,
		database:     database, // Store the provided database interface implementation.
	}
}

// GetUser is a helper function to retrieve user details mapped to the utils.User struct,
// identified by their external ID. It requires a context for cancellation propagation
// and uses the provided db.DatabaseInterface implementation to fetch the data.
// Returns nil, nil if the user is not found by the database layer.
func GetUser(ctx context.Context, id string, database db.DatabaseInterface) (*User, error) {
	// Call the database layer method via the interface, passing the context.
	user, err := database.GetUserByExternalID(ctx, id)
	if err != nil {
		// Propagate database errors (includes potential gorm.ErrRecordNotFound).
		return nil, err
	}

	// Convert the database model to the utility/API layer model.
	// DbUserToApi handles the case where 'user' is nil (not found).
	resp, err := DbUserToApi(user)
	if err != nil {
		// This error is unlikely if DbUserToApi only does field mapping.
		return nil, fmt.Errorf("failed to map db user to api user: %w", err)
	}

	// Return the mapped user or nil, nil if not found.
	return resp, nil
}

// UpdateUserLastLogin is a helper function to update the user's last login timestamp
// in the database, identified by their external user ID.
// It requires a context and uses the provided db.DatabaseInterface implementation.
func UpdateUserLastLogin(ctx context.Context, userId string, lastLogin time.Time, database db.DatabaseInterface) error {
	// Call the database layer method via the interface, passing the context.
	err := database.UpdateUserLastLoginWithExternalID(ctx, userId, lastLogin)
	if err != nil {
		// Propagate database errors.
		return err
	}
	return nil
}
