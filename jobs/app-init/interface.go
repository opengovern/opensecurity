// Proposed file: jobs/app-init/interface.go

package app_init // Use the package name you chose

import (
	"context"
	"fmt"
	"log"
	"time"
)

// InitializableComponent defines the standard lifecycle methods for a dependency
// that needs to be checked and potentially configured during app initialization.
type InitializableComponent interface {
	// Name returns a human-readable name for the component (used for logging).
	Name() string

	// CheckAvailability verifies if the component is reachable (e.g., network connectivity, basic endpoint check).
	// This method should implement necessary retry logic.
	CheckAvailability(ctx context.Context) error

	// Configure performs any necessary setup or configuration actions for the component.
	// This might be a no-op for components that don't require active configuration by this job.
	Configure(ctx context.Context) error

	// CheckHealth verifies if the component is fully operational and healthy after configuration.
	// This might involve deeper checks than CheckAvailability. It could be a no-op or reuse CheckAvailability.
	CheckHealth(ctx context.Context) error
}

// --- Helper for Retries (Optional but Recommended) ---

// waitForCondition implements generic retry logic for checks.
func waitForCondition(ctx context.Context, componentName string, actionName string, maxRetries int, delay time.Duration, checkFunc func(ctx context.Context) error) error {
	log.Printf("INFO: Checking %s for %s...", actionName, componentName)
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check context cancellation before attempt
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled waiting for %s %s: %w", componentName, actionName, ctx.Err())
		default: // Continue
		}

		if attempt > 0 {
			log.Printf("INFO: %s %s check failed. Retrying in %v (attempt %d/%d)", componentName, actionName, delay, attempt, maxRetries)
			// Wait respecting context cancellation
			select {
			case <-time.After(delay): // Wait for the delay
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry delay for %s %s: %w", componentName, actionName, ctx.Err())
			}
		}

		// Execute the specific check function passed in
		log.Printf("INFO: Attempt %d/%d: Performing %s check for %s...", attempt+1, maxRetries+1, actionName, componentName)
		err := checkFunc(ctx)
		if err == nil {
			log.Printf("INFO: %s %s check successful!", componentName, actionName)
			return nil // Success
		}

		// Store last error and log warning before retrying
		lastErr = err
		log.Printf("WARN: Attempt %d failed for %s %s: %v", attempt+1, componentName, actionName, err)

	} // End retry loop

	// If loop finishes, all retries failed
	log.Printf("ERROR: %s %s failed after %d retries.", componentName, actionName, maxRetries)
	return fmt.Errorf("%s %s failed after %d retries; last error: %w", componentName, actionName, maxRetries, lastErr)
}
