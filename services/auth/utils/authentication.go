package utils

import (
	"time"

	"github.com/opengovern/opensecurity/services/auth/db"
)

type Service struct {
	domain       string
	clientID     string
	clientSecret string
	appClientID  string
	Connection   string
	InviteTTL    int

	database db.Database
}
type User struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	FullName      string    `json:"full_name"`
	LastLogin     time.Time `json:"last_login"`
	Username      string    `json:"username"`
	Role          string    `json:"role"`
	IsActive      bool      `json:"is_active"`
	ConnectorId   string    `json:"connector_id"`
	ExternalId    string    `json:"external_id"`
}

func DbUserToApi(u *db.User) (*User, error) {
	if u == nil {
		return nil, nil
	}

	return &User{
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		FullName:      u.FullName,
		LastLogin:     u.LastLogin,
		Username:      string(u.Username),
		Role:          string(u.Role),
		ExternalId:    u.ExternalId,
		ID:            u.ID,
		IsActive:      u.IsActive,
		ConnectorId:   u.ConnectorId,
	}, nil
}


func New(domain, appClientID, clientID, clientSecret, connection string, inviteTTL int, database db.Database) *Service {
	return &Service{
		domain:       domain,
		appClientID:  appClientID,
		clientID:     clientID,
		clientSecret: clientSecret,
		Connection:   connection,
		InviteTTL:    inviteTTL,
		database:     database,
	}
}

func GetUser(id string, database db.Database) (*User, error) {
	user, err := database.GetUserByExternalID(id)
	if err != nil {
		return nil, err
	}

	resp, err := DbUserToApi(user)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func UpdateUserLastLogin(userId string, lastLogin time.Time, database db.Database) error {

	err := database.UpdateUserLastLoginWithExternalID(userId, lastLogin)

	if err != nil {
		return err
	}

	return nil
}
