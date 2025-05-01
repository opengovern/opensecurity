// /Users/anil/workspace/opensecurity/jobs/app-init/configurators/auth_component.go
package configurators

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	// Required for RequestTimeout constant usage in error messages
	// Import the types package from the parent directory's sibling 'types'
	// <<< ADJUST IMPORT PATH to your actual 'types' package location >>>
	initTypes "github.com/opengovern/opensecurity/jobs/app-init/types"
)

// AuthComponent checks Auth service availability via HTTP and implements the InitializableComponent interface.
type AuthComponent struct {
	healthURL string
	client    *http.Client // Use a shared client for potential connection reuse
	// logger    *zap.Logger // Optional: Pass in if zap logger is preferred
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
	// Create an HTTP client instance for this component.
	// Timeout is handled by the context passed to Do().
	client := &http.Client{}
	return &AuthComponent{healthURL: healthCheckURL, client: client}, nil
}

// Name returns the component name.
func (a *AuthComponent) Name() string {
	return "Auth Service"
}

// checkHTTPGet performs the actual HTTP GET request within the given context's deadline.
// This is used by both CheckAvailability and CheckHealth via WaitForCondition.
func (a *AuthComponent) checkHTTPGet(ctx context.Context) error {
	// The context passed here from WaitForCondition controls the timeout for this specific attempt.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	// Optional: Add a user-agent
	// req.Header.Set("User-Agent", "opensecurity-app-init-job/1.0")

	resp, err := a.client.Do(req)
	if err != nil {
		// Check specifically for context deadline exceeded to clarify timeout
		if errors.Is(err, context.DeadlineExceeded) {
			// Use the imported constant in the error message
			return fmt.Errorf("http get timed out after %v: %w", initTypes.RequestTimeout, err)
		}
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

// CheckAvailability uses the WaitForCondition helper from the 'types' package
// to check the HTTP endpoint with retries, using constants from 'types'.
func (a *AuthComponent) CheckAvailability(ctx context.Context) error {
	// Call the exported helper and use exported constants from the imported 'types' package
	return initTypes.WaitForCondition(ctx, a.Name(), "availability", initTypes.MaxRetries, initTypes.RetryDelay, a.checkHTTPGet)
}

// CheckIfInitializationIsRequired checks if the Auth service needs configuration by this job.
// Returns ErrAlreadyInitialized indicating configuration should be skipped.
func (a *AuthComponent) CheckIfInitializationIsRequired(ctx context.Context) error {
	log.Printf("INFO: [%s] Checking if configuration is required...", a.Name())
	// TODO: Implement actual logic if the Auth service *can* be configured by this job
	// and needs a check (e.g., check if a specific default user/role exists via API).
	// For now, assume configuration is handled elsewhere or not needed once available.
	log.Printf("INFO: [%s] Assuming configuration is not required by this job.", a.Name())
	return initTypes.ErrAlreadyInitialized // Signal configuration not needed
}

// Configure is a placeholder for future Auth configuration steps.
func (a *AuthComponent) Configure(ctx context.Context) error {
	// This would only be called by the orchestrator in command.go if
	// CheckIfInitializationIsRequired returned nil.
	// TODO: Implement Auth service configuration logic if required by this job.
	log.Printf("INFO: [%s] Placeholder - Configure step executed (but likely skipped).", a.Name())
	return nil
}

// CheckHealth performs a health check after configuration/availability check.
// For now, it simply re-uses the availability check.
func (a *AuthComponent) CheckHealth(ctx context.Context) error {
	// TODO: Implement a more specific health check if needed (e.g., attempt a basic auth operation).
	log.Printf("INFO: [%s] Performing health check (using availability check)...", a.Name())
	// For now, reuse availability check. This uses the same retry logic.
	return initTypes.WaitForCondition(ctx, a.Name(), "health", initTypes.MaxRetries, initTypes.RetryDelay, a.checkHTTPGet)
}

// Compile-time check to ensure AuthComponent implements the InitializableComponent interface
// from the imported 'types' package.
var _ initTypes.InitializableComponent = (*AuthComponent)(nil)
