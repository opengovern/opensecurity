// /Users/anil/workspace/opensecurity/jobs/app-init/runner.go
package app_init

import (
	"context"
	"fmt"
	"log"
	"os"

	// Use the types package from the sibling directory
	initTypes "github.com/opengovern/opensecurity/jobs/app-init/types"
)

// Runner handles the ordered execution of initialization components.
type Runner struct {
	logger *log.Logger
}

// NewRunner creates a new Runner instance.
func NewRunner(logger *log.Logger) *Runner {
	if logger == nil {
		// Fallback to default logger if none provided, ensuring logger is never nil.
		logger = log.New(os.Stdout, "APP-INIT-RUNNER: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		logger.Println("WARN: No logger provided to NewRunner, using default.")
	}
	return &Runner{logger: logger}
}

// Run executes the initialization lifecycle for the provided components sequentially.
// It stops and returns an error on the first failure in any step for any component.
func (r *Runner) Run(ctx context.Context, components []initTypes.InitializableComponent) error {
	r.logger.Printf("INFO: Beginning component initialization sequence (%d components)...", len(components))
	if len(components) == 0 {
		r.logger.Println("WARN: No components provided to initialize.")
		return nil // Nothing to do, not an error.
	}

	for i, comp := range components {
		componentName := comp.Name()
		r.logger.Printf("----- Initializing Component %d/%d: %s -----", i+1, len(components), componentName)

		// Check for overall context cancellation before starting component
		select {
		case <-ctx.Done():
			r.logger.Printf("WARN: [%s] Context cancelled before starting initialization.", componentName)
			return fmt.Errorf("context cancelled before initializing component %s: %w", componentName, ctx.Err())
		default:
			// Continue initialization
		}

		// --- Step 1: Check Availability ---
		r.logger.Printf("INFO: [%s] Step 1/4: Checking availability...", componentName)
		if err := comp.CheckAvailability(ctx); err != nil {
			r.logger.Printf("ERROR: [%s] Availability check failed.", componentName)
			// Wrap error for context
			return fmt.Errorf("component '%s' availability check failed: %w", componentName, err)
		}
		r.logger.Printf("INFO: [%s] Availability check successful.", componentName)

		// --- Step 2: Check if Initialization is Required ---
		r.logger.Printf("INFO: [%s] Step 2/4: Checking if initialization is required...", componentName)
		// This checks the core prerequisite (e.g., DB exists).
		// It returns nil ONLY if the prerequisite is met.
		initRequiredErr := comp.CheckIfInitializationIsRequired(ctx)

		// --- Step 3: Configure (if needed) ---
		if initRequiredErr == nil {
			// Prerequisite met (e.g., DB exists), proceed to configuration steps (Dex clients, user).
			r.logger.Printf("INFO: [%s] Step 3/4: Prerequisite met. Proceeding with configuration...", componentName)
			if err := comp.Configure(ctx); err != nil {
				r.logger.Printf("ERROR: [%s] Configuration failed.", componentName)
				return fmt.Errorf("component '%s' configuration failed: %w", componentName, err)
			}
			r.logger.Printf("INFO: [%s] Configuration successful.", componentName)
		} else {
			// Prerequisite NOT met (e.g., DB doesn't exist). This is a fatal error for this component.
			r.logger.Printf("ERROR: [%s] Step 3/4: Prerequisite check failed. Halting initialization for this component.", componentName)
			// Return the specific error from CheckIfInitializationIsRequired.
			return fmt.Errorf("component '%s' prerequisite check failed: %w", componentName, initRequiredErr)
		}

		// --- Step 4: Check Health Post-Configuration ---
		r.logger.Printf("INFO: [%s] Step 4/4: Performing post-configuration health check...", componentName)
		if err := comp.CheckHealth(ctx); err != nil {
			r.logger.Printf("ERROR: [%s] Health check failed.", componentName)
			return fmt.Errorf("component '%s' health check failed: %w", componentName, err)
		}
		r.logger.Printf("INFO: [%s] Health check successful.", componentName)

		r.logger.Printf("----- Component Initialized Successfully: %s -----", componentName)
	} // End component loop

	r.logger.Println("INFO: All components initialized successfully.")
	return nil // Overall Success
}
