package db

import (
	"github.com/opengovern/og-util/pkg/api"
	"gorm.io/gorm"
	"time"
)

type Configuration struct {
	gorm.Model
	Key   string
	Value string
}

type ApiKey struct {
	gorm.Model
	Name          string
	Role          api.Role
	CreatorUserID string
	IsActive      bool
	KeyHash       string
	MaskedKey     string
}

type Connector struct {
	gorm.Model
	UserCount        uint `gorm:"default:0"`
	ConnectorID      string
	ConnectorType    string
	ConnectorSubType string
	LastUpdate       time.Time
}

type User struct {
	gorm.Model
	Email                 string    `gorm:"not null"` // Enforce NOT NULL
	EmailVerified         bool	`gorm:"default:false;not null"`
	FullName              string
	Role                  api.Role  `gorm:"not null"` // Enforce NOT NULL
	ConnectorId           string    `gorm:"not null"` // Enforce NOT NULL
	ExternalId            string    `gorm:"unique"`
	LastLogin             time.Time
	Username              string
	RequirePasswordChange bool      `gorm:"default:true;not null"`
	IsActive              bool      `gorm:"default:true;not null"`
}
