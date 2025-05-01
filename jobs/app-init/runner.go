// /Users/anil/workspace/opensecurity/jobs/app-init/runner.go
package app_init

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	// Import the types package from the sibling 'types' directory
	initTypes "github.com/opengovern/opensecurity/jobs/app-init/types" // <<< ADJUST PATH AS NEEDED
)

// Runner handles the ordered execution of initialization components.
type Runner struct {
	logger *log.Logger
}

// NewRunner creates a new Runner instance.
func NewRunner(logger *log.Logger) *Runner {
	// Use default logger if nil is passed
	if logger == nil {
		logger = log.New(os.Stdout, "APP-INIT-RUNNER: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
	return &Runner{logger: logger}
}

// Run executes the initialization lifecycle for the provided components in order.
// It stops and returns an error on the first failure.
func (r *Runner) Run(ctx context.Context, components []initTypes.InitializableComponent) error {
	r.logger.Println("INFO: Beginning component initialization sequence...")
	if len(components) == 0 {
		r.logger.Println("WARN: No components provided to initialize.")
		return nil
	}

	for _, comp := range components {
		componentName := comp.Name()
		r.logger.Printf("----- Initializing Component: %s -----", componentName)

		// 1. Check Availability
		r.logger.Printf("INFO: [%s] Checking availability...", componentName)
		if err := comp.CheckAvailability(ctx); err != nil {
			r.logger.Printf("ERROR: [%s] Availability check failed.", componentName)
			return fmt.Errorf("%s availability check failed: %w", componentName, err) // Stop on first failure
		}
		r.logger.Printf("INFO: [%s] Component is available.", componentName)

		// 2. Check if Configuration is Required
		r.logger.Printf("INFO: [%s] Checking if configuration is required...", componentName)
		initErr := comp.CheckIfInitializationIsRequired(ctx)

		if initErr == nil {
			// Configuration IS required
			r.logger.Printf("INFO: [%s] Configuration required. Proceeding...", componentName)
			if err := comp.Configure(ctx); err != nil {
				r.logger.Printf("ERROR: [%s] Configuration failed.", componentName)
				return fmt.Errorf("%s configuration failed: %w", componentName, err) // Stop on failure
			}
			r.logger.Printf("INFO: [%s] Component configured successfully.", componentName)
		} else if errors.Is(initErr, initTypes.ErrAlreadyInitialized) {
			// Configuration is NOT required, skip Configure step
			r.logger.Printf("INFO: [%s] Configuration step skipped (already initialized).", componentName)
		} else {
			// An actual error occurred during the check itself
			r.logger.Printf("ERROR: [%s] Failed to check if configuration is required.", componentName)
			return fmt.Errorf("%s check for required configuration failed: %w", componentName, initErr) // Stop on failure
		}

		// 3. Check Health Post-Configuration/Check
		r.logger.Printf("INFO: [%s] Performing post-configuration health check...", componentName)
		if err := comp.CheckHealth(ctx); err != nil {
			r.logger.Printf("ERROR: [%s] Health check failed.", componentName)
			return fmt.Errorf("%s health check failed: %w", componentName, err) // Stop on failure
		}
		r.logger.Printf("INFO: [%s] Component is healthy.", componentName)

		r.logger.Printf("----- Component Initialized Successfully: %s -----", componentName)
	} // End loop

	r.logger.Println("INFO: All components initialized successfully.")
	return nil // Overall Success
}
