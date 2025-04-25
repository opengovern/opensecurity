package utils

import (
	"context" // Added context import
	"time"

	// Ensure correct import path to where DatabaseInterface is defined
	"github.com/opengovern/opensecurity/services/auth/db"
)

// Service might hold context/configuration for utility functions.
// Updated to use the database interface.
type Service struct {
	domain       string
	clientID     string
	clientSecret string
	appClientID  string
	Connection   string
	InviteTTL    int

	database db.DatabaseInterface // Use interface type
}

// User represents the API/utility layer user structure.
type User struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	FullName      string    `json:"full_name"`
	LastLogin     time.Time `json:"last_login"`
	Username      string    `json:"username"`
	Role          string    `json:"role"` // Role as string in this layer
	IsActive      bool      `json:"is_active"`
	ConnectorId   string    `json:"connector_id"`
	ExternalId    string    `json:"external_id"`
}

// DbUserToApi converts a user object from the database layer to the utility/API layer.
// No changes needed here as it works on the DB type directly.
func DbUserToApi(u *db.User) (*User, error) {
	if u == nil {
		// Return nil, nil to represent "not found" consistently
		// if the underlying DB method returns nil, nil for not found.
		return nil, nil
	}

	return &User{
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt, // Make sure db.User has UpdatedAt from gorm.Model
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		FullName:      u.FullName,
		LastLogin:     u.LastLogin,
		Username:      u.Username,     // Username field seems to be string already in db.User
		Role:          string(u.Role), // Convert db Role (likely api.Role) to string
		ExternalId:    u.ExternalId,
		ID:            u.ID,
		IsActive:      u.IsActive,
		ConnectorId:   u.ConnectorId,
	}, nil
}

// New creates a new Service instance.
// Updated to accept the database interface.
func New(domain, appClientID, clientID, clientSecret, connection string, inviteTTL int, database db.DatabaseInterface) *Service {
	return &Service{
		domain:       domain,
		appClientID:  appClientID,
		clientID:     clientID,
		clientSecret: clientSecret,
		Connection:   connection,
		InviteTTL:    inviteTTL,
		database:     database, // Assign interface
	}
}

// GetUser retrieves user details by external ID using the provided database interface.
// Updated signature to accept context and interface.
func GetUser(ctx context.Context, id string, database db.DatabaseInterface) (*User, error) {
	// Pass context down to the database method call
	user, err := database.GetUserByExternalID(ctx, id) // Call via interface
	if err != nil {
		// Includes case where err is gorm.ErrRecordNotFound (handled by DbUserToApi returning nil, nil)
		// Or other DB errors
		return nil, err
	}

	// Convert db.User to utils.User
	resp, err := DbUserToApi(user)
	if err != nil {
		// This error likely only happens if DbUserToApi logic changes,
		// as conversion itself shouldn't fail here.
		return nil, err
	}

	return resp, nil // Returns nil, nil if user not found and DbUserToApi handles it
}

// UpdateUserLastLogin updates the user's last login time.
// Updated signature to accept context and interface.
func UpdateUserLastLogin(ctx context.Context, userId string, lastLogin time.Time, database db.DatabaseInterface) error {
	// Pass context down to the database method call
	err := database.UpdateUserLastLoginWithExternalID(ctx, userId, lastLogin) // Call via interface
	if err != nil {
		return err
	}
	return nil
}
