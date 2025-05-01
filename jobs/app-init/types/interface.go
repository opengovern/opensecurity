// /Users/anil/workspace/opensecurity/jobs/app-init/types/interface.go
package types

import (
	"context"
	"errors" // <<< Add/Ensure 'errors' package is imported
	"fmt"
	"log"
	"time"
)

// Constants for retry logic (Exported for use by components/callers)
const (
	MaxRetries     = 3
	RetryDelay     = 5 * time.Second
	RequestTimeout = 15 * time.Second
)

// *** ADD Sentinel Error Definition ***
// ErrAlreadyInitialized is returned by CheckIfInitializationIsRequired when configuration is not needed.
var ErrAlreadyInitialized = errors.New("component already initialized/configured")

// InitializableComponent defines the standard lifecycle methods for a dependency
// that needs to be checked and potentially configured during app initialization.
type InitializableComponent interface {
	Name() string
	CheckAvailability(ctx context.Context) error
	CheckIfInitializationIsRequired(ctx context.Context) error // Returns nil if init required, ErrAlreadyInitialized if not, other error on failure
	Configure(ctx context.Context) error
	CheckHealth(ctx context.Context) error
}

// AvailabilityChecker can be used if only the first step is needed initially.
// It's a subset of InitializableComponent.
type AvailabilityChecker interface {
	Name() string
	CheckAvailability(ctx context.Context) error
}

// --- Exported Helper for Retries ---

// WaitForCondition implements generic retry logic for checks like availability or health.
// It calls the checkFunc up to maxRetries+1 times with delays in between.
// It applies the RequestTimeout to the context passed to each individual checkFunc call.
func WaitForCondition(ctx context.Context, componentName string, actionName string, maxRetries int, delay time.Duration, checkFunc func(ctx context.Context) error) error {
	log.Printf("INFO: [%s] Checking %s...", componentName, actionName)
	var lastErr error

	// The loop runs for the initial attempt (attempt 0) plus maxRetries
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check for overall context cancellation before each attempt
		select {
		case <-ctx.Done():
			log.Printf("WARN: [%s] Context cancelled while waiting for %s check.", componentName, actionName)
			return fmt.Errorf("context cancelled waiting for %s %s: %w", componentName, actionName, ctx.Err())
		default:
			// Continue with the attempt
		}

		// Apply delay *before* retries (i.e., starting from the second attempt)
		if attempt > 0 {
			log.Printf("INFO: [%s] %s check failed. Retrying in %v (attempt %d/%d)", componentName, actionName, delay, attempt, maxRetries)
			// Wait for the delay, but also listen for context cancellation
			select {
			case <-time.After(delay):
				// Delay finished, continue to next attempt
			case <-ctx.Done():
				log.Printf("WARN: [%s] Context cancelled during retry delay for %s.", componentName, actionName)
				return fmt.Errorf("context cancelled during retry delay for %s %s: %w", componentName, actionName, ctx.Err())
			}
		}

		// Execute the specific check function with its own timeout context
		log.Printf("INFO: [%s] Attempt %d/%d: Performing %s check...", componentName, attempt+1, maxRetries+1, actionName)
		checkCtx, checkCancel := context.WithTimeout(ctx, RequestTimeout) // Apply attempt timeout
		err := checkFunc(checkCtx)
		checkCancel() // Release the timeout context promptly

		if err == nil {
			log.Printf("INFO: [%s] %s check successful!", componentName, actionName)
			return nil // Success
		}

		// Store the last error encountered
		lastErr = err
		log.Printf("WARN: [%s] Attempt %d failed for %s: %v", componentName, attempt+1, actionName, err)

	} // End retry loop

	// If the loop finishes, all attempts failed
	log.Printf("ERROR: [%s] %s failed after %d retries.", componentName, actionName, maxRetries)
	return fmt.Errorf("%s %s failed after %d retries; last error: %w", componentName, actionName, maxRetries, lastErr)
}
