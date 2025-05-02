// /Users/anil/workspace/opensecurity/jobs/app-init/types/interface.go
package types

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// Constants for retry logic, exported for use by components and callers.
const (
	MaxRetries     = 5                // Number of retries after the initial attempt fails
	RetryDelay     = 5 * time.Second  // Delay between retries
	RequestTimeout = 15 * time.Second // Timeout for individual check attempts (HTTP, DB ping, gRPC calls etc.)
	// Define longer timeouts specifically for potentially slow operations like Dex calls or DB DML
	DexTimeout = 30 * time.Second
	DBTimeout  = 30 * time.Second
)

// ErrAlreadyInitialized is a sentinel error.
// NOTE: In the final implementation, this specific error is NOT returned by CheckIfInitializationIsRequired.
// CheckIfInitializationIsRequired now returns nil ONLY if the DB exists, otherwise it returns a specific error.
// This constant remains defined but unused by the current AuthComponent logic.
var ErrAlreadyInitialized = errors.New("component prerequisite met, proceed with configuration check")

// InitializableComponent defines the standard lifecycle methods for a dependency
// that needs to be checked and potentially configured during app initialization.
type InitializableComponent interface {
	// Name returns a human-readable name for the component (used in logs).
	Name() string
	// CheckAvailability verifies if the component's external dependencies are reachable (e.g., network connectivity to DB/API).
	// Should use retries.
	CheckAvailability(ctx context.Context) error
	// CheckIfInitializationIsRequired checks if the core prerequisite for configuration is met (e.g., database exists).
	// Returns nil ONLY if the prerequisite is met, allowing the Configure step to run.
	// Returns a specific error if the prerequisite is NOT met, halting initialization for this component.
	CheckIfInitializationIsRequired(ctx context.Context) error
	// Configure performs the necessary setup actions IF CheckIfInitializationIsRequired returned nil.
	// (e.g., ensure Dex clients, create initial user).
	Configure(ctx context.Context) error
	// CheckHealth verifies the component is operational after potential configuration.
	// Should use retries.
	CheckHealth(ctx context.Context) error
}

// --- Exported Helper for Retries ---

// WaitForCondition implements generic retry logic for checks like availability or health.
// It calls the checkFunc up to maxRetries times after the first failed attempt.
// It applies the RequestTimeout to the context passed to each individual checkFunc call.
func WaitForCondition(ctx context.Context, componentName string, actionName string, maxRetries int, delay time.Duration, checkFunc func(ctx context.Context) error) error {
	log.Printf("INFO: [%s] Checking %s...", componentName, actionName)
	var lastErr error

	// Loop runs for the initial attempt (attempt 0) + maxRetries
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Check for overall context cancellation before each attempt
		select {
		case <-ctx.Done():
			log.Printf("WARN: [%s] Context cancelled while waiting for %s check.", componentName, actionName)
			// Wrap context error for clarity
			return fmt.Errorf("context cancelled waiting for %s %s: %w", componentName, actionName, ctx.Err())
		default:
			// Continue with the attempt
		}

		// Apply delay *before* retries (i.e., starting from the second attempt)
		if attempt > 0 {
			log.Printf("INFO: [%s] %s check failed (Attempt %d/%d). Retrying in %v...", componentName, actionName, attempt, maxRetries, delay)
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
		// Use the standard RequestTimeout for these checks
		checkCtx, checkCancel := context.WithTimeout(ctx, RequestTimeout)
		err := checkFunc(checkCtx)
		// It's crucial to cancel the context to release resources, especially in loops.
		checkCancel()

		if err == nil {
			log.Printf("INFO: [%s] %s check successful!", componentName, actionName)
			return nil // Success
		}

		// Store the last error encountered
		lastErr = err
		log.Printf("WARN: [%s] Attempt %d/%d failed for %s: %v", componentName, attempt+1, maxRetries+1, actionName, err)

	} // End retry loop

	// If the loop finishes, all attempts failed
	log.Printf("ERROR: [%s] %s failed after %d retries.", componentName, actionName, maxRetries)
	// Wrap the last error for better context
	return fmt.Errorf("%s %s failed after %d retries; last error: %w", componentName, actionName, maxRetries, lastErr)
}
