// jobs/app-init/postgres_component.go
package configurators // Package name should match the directory name usually

import (
	"context"
	"database/sql" // Required for DB interactions
	"errors"
	"fmt"
	"log"
	"net/url" // For potential URL parsing/validation

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver (ensure this is in your go.mod)
	// Add import for zap logger if you use it instead of standard log
)

// PostgresComponent holds configuration and state for checking/configuring PostgreSQL.
type PostgresComponent struct {
	dbURL string // Connection string (e.g., "postgresql://user:pass@host:port/db?sslmode=disable")
	// logger *zap.Logger // Optional: Pass in a structured logger if preferred
}

// NewPostgresComponent creates a new component for PostgreSQL checks.
func NewPostgresComponent(dbURL string) (*PostgresComponent, error) {
	if dbURL == "" {
		return nil, errors.New("database URL cannot be empty for PostgresComponent")
	}
	// Optional: More sophisticated validation of the dbURL format
	_, err := url.Parse(dbURL) // Basic check if it parses
	if err != nil {
		return nil, fmt.Errorf("invalid database URL format: %w", err)
	}

	return &PostgresComponent{
		dbURL: dbURL,
		// logger: logger, // Assign if using zap logger
	}, nil
}

// Name returns the human-readable name of the component.
func (p *PostgresComponent) Name() string {
	return "PostgreSQL Database"
}

// pingDB performs the actual database ping operation.
func (p *PostgresComponent) pingDB(ctx context.Context) error {
	// Open potentially opens a connection pool, Ping checks connectivity.
	// Ensure driver is imported with blank identifier: _ "github.com/jackc/pgx/v4/stdlib"
	db, err := sql.Open("pgx", p.dbURL)
	if err != nil {
		// This error happens if the DSN is invalid or driver unavailable
		return fmt.Errorf("failed to open DB connection using driver: %w", err)
	}
	// Close should be called to release resources, although PingContext usually doesn't keep it open long.
	defer db.Close()

	// Use a timeout specific to the Ping operation itself, within the overall check context.
	// requestTimeout should be defined (e.g., in interface.go or config)
	pingCtx, cancel := context.WithTimeout(ctx, consts.requestTimeout) // Assuming requestTimeout is defined
	defer cancel()

	err = db.PingContext(pingCtx)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err) // Wrap error for context
	}
	return nil // Ping successful
}

// CheckAvailability uses the waitForCondition helper to ping the database with retries.
func (p *PostgresComponent) CheckAvailability(ctx context.Context) error {
	// Assuming maxRetries and retryDelay are defined constants in the package
	return waitForCondition(ctx, p.Name(), "availability", maxRetries, retryDelay, p.pingDB)
}

// Configure performs any necessary setup for PostgreSQL (e.g., migrations).
// For now, it's a placeholder.
func (p *PostgresComponent) Configure(ctx context.Context) error {
	// TODO: Implement database migration or other configuration logic here if needed.
	log.Printf("INFO: Placeholder - Configuration step for %s executed.", p.Name())
	return nil
}

// CheckHealth performs a health check after configuration.
// For now, it simply re-uses the availability check.
func (p *PostgresComponent) CheckHealth(ctx context.Context) error {
	// For PostgreSQL, a successful ping often indicates basic health.
	// More complex checks could involve running a simple query like "SELECT 1".
	log.Printf("INFO: Performing health check for %s (using availability ping).", p.Name())
	return p.CheckAvailability(ctx) // Reuse availability check for basic health
}
