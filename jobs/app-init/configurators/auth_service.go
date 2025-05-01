// jobs/app-init/auth_service.go
package configurators // Package name should match the directory name usually

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	// Required for context timeout
)

// AuthComponent checks Auth service availability via HTTP.
type AuthComponent struct {
	healthURL string
	client    *http.Client
}

// NewAuthComponent creates a new component for Auth service checks.
func NewAuthComponent(healthCheckURL string) (*AuthComponent, error) {
	if healthCheckURL == "" {
		return nil, errors.New("auth health check URL cannot be empty")
	}
	_, err := url.ParseRequestURI(healthCheckURL) // Validate format
	if err != nil {
		return nil, fmt.Errorf("invalid auth health check URL format '%s': %w", healthCheckURL, err)
	}
	// Client timeout is handled by the context passed to Do in checkHTTPGet
	client := &http.Client{}
	return &AuthComponent{healthURL: healthCheckURL, client: client}, nil
}

// Name returns the component name.
func (a *AuthComponent) Name() string {
	return "Auth Service"
}

// checkHTTPGet performs the actual HTTP GET request.
func (a *AuthComponent) checkHTTPGet(ctx context.Context) error {
	// The context passed here already has the attempt timeout (requestTimeout)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.healthURL, nil)
	if err != nil {
		// Error creating request object (should be rare with validated URL)
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		// Error performing request (DNS, connection, timeout)
		return fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close() // Ensure body is always closed

	// Check for successful status code (2xx)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil // Success
	}

	// Return error for non-2xx status codes
	return fmt.Errorf("received non-2xx status: %s", resp.Status)
}

// CheckAvailability uses the waitForCondition helper to check the HTTP endpoint with retries.
func (a *AuthComponent) CheckAvailability(ctx context.Context) error {
	// Assuming maxRetries and retryDelay are defined constants in the package
	return waitForCondition(ctx, a.Name(), "availability", maxRetries, retryDelay, a.checkHTTPGet)
}

// Configure is a placeholder for future Auth configuration steps.
func (a *AuthComponent) Configure(ctx context.Context) error {
	log.Printf("INFO: Placeholder - Configure step for %s.", a.Name())
	// Add Auth configuration logic here if needed
	return nil
}

// CheckHealth is a placeholder for future post-configuration checks.
func (a *AuthComponent) CheckHealth(ctx context.Context) error {
	log.Printf("INFO: Placeholder - Health check for %s (using availability check).", a.Name())
	// For now, reuse availability check. Could be more specific later.
	return a.CheckAvailability(ctx)
}

// Compile-time check to ensure AuthComponent implements the full (future) interface
// var _ InitializableComponent = (*AuthComponent)(nil)
// Compile-time check to ensure AuthComponent implements AvailabilityChecker
var _ AvailabilityChecker = (*AuthComponent)(nil)
