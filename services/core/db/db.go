// /Users/anil/workspace/opensecurity/services/core/db/db.go
package db

import (
	"context" // <<< Added import
	"errors"  // <<< Added import
	"fmt"     // <<< Added import

	"github.com/opengovern/opensecurity/services/core/db/models"
	"gorm.io/gorm"
	// Import database/sql only if needed directly, PingContext is on sql.DB
	// "database/sql"
)

type Database struct {
	orm *gorm.DB // Keep unexported for encapsulation
}

// NewDatabase creates a new Database instance.
func NewDatabase(orm *gorm.DB) Database {
	if orm == nil {
		// Optional: Handle nil ORM input, perhaps panic or return error?
		// Returning struct with nil orm will cause Ping to fail later.
		fmt.Println("Warning: NewDatabase called with nil gorm.DB") // Or use logger if available
	}
	return Database{orm: orm}
}

// Initialize performs database migrations.
func (db Database) Initialize() error {
	if db.orm == nil {
		return errors.New("cannot initialize database with nil connection")
	}
	err := db.orm.AutoMigrate(
		// shared
		&models.Query{},
		&models.QueryParameter{},
		// inventory
		&models.ResourceType{},
		//&models.NamedQuery{},
		//&models.NamedQueryTag{},
		//&models.NamedQueryHistory{},
		&models.ResourceTypeTag{},
		//&models.ResourceCollection{},
		//&models.ResourceCollectionTag{},
		&models.ResourceTypeV2{},
		// metadata
		&models.ConfigMetadata{},
		&models.PolicyParameterValues{},
		&models.QueryView{},
		&models.QueryViewTag{},
		&models.PlatformConfiguration{},
		&models.RunNamedQueryRunCache{},
		&models.Dashboard{},
		&models.Widget{},
		&models.ChatbotSecret{},
		&models.Session{},
		&models.Chat{},
		&models.ChatSuggestion{},
		&models.ChatClarification{},
		// Add any other models here
	)
	if err != nil {
		// Wrap error for context
		return fmt.Errorf("database auto-migration failed: %w", err)
	}
	return nil
}

// Ping checks the database connectivity using the underlying connection pool.
func (db Database) Ping(ctx context.Context) error {
	if db.orm == nil {
		return errors.New("database connection (orm) is nil, cannot ping")
	}
	// Get the underlying *sql.DB connection pool from GORM
	sqlDB, err := db.orm.DB()
	if err != nil {
		// Error retrieving the sql.DB handle from GORM
		return fmt.Errorf("failed to get underlying *sql.DB from gorm: %w", err)
	}
	if sqlDB == nil {
		// GORM returned a nil handle, which shouldn't happen if orm was non-nil and retrieving didn't error. Defensive check.
		return errors.New("gorm returned a nil *sql.DB connection")
	}

	// Ping the database using the provided context (for deadlines/cancellation)
	if err := sqlDB.PingContext(ctx); err != nil {
		// Wrap ping error for context
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Ping successful
	return nil
}

// Add other Database methods ( AddFilter, ListFilters, etc.) here...
