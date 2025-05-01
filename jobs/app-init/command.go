// /Users/anil/workspace/opensecurity/jobs/app-init/command.go
package app_init

import (
	"fmt"
	"log"
	"os"
	"strings" // Added for Join and HasPrefix

	// Required by component implementations implicitly
	"github.com/spf13/cobra"

	// Import local sub-packages
	// <<< ADJUST IMPORT PATHS AS NEEDED >>>
	"github.com/opengovern/opensecurity/jobs/app-init/configurators"
	initTypes "github.com/opengovern/opensecurity/jobs/app-init/types"
)

// Environment variable names
const (
	// Auth Service URL Components
	envAuthServiceName = "AUTH_SERVICE_NAME" // e.g., "auth-service"
	envAuthNamespace   = "AUTH_NAMESPACE"    // e.g., "opensecurity"
	envAuthPort        = "AUTH_SERVICE_PORT" // e.g., "8251"
	envAuthHealthPath  = "AUTH_HEALTH_PATH"  // e.g., "/health" or "/healthz/ready" (Optional, defaults to /health)

	// Database
	//envDatabaseURL = "DATABASE_URL"
	// Add other env var names here
)

// buildComponentList reads config, constructs URLs, and creates the ordered list of components.
func buildComponentList() ([]initTypes.InitializableComponent, error) {
	log.Println("INFO: Building component list...")

	// --- Get Configuration ---
	authServiceName := os.Getenv(envAuthServiceName)
	authNamespace := os.Getenv(envAuthNamespace)
	authPort := os.Getenv(envAuthPort)
	authHealthPath := os.Getenv(envAuthHealthPath) // Optional

	//dbURL := os.Getenv(envDatabaseURL)
	// ... get other configs ...

	// --- Validate Configuration ---
	var missingVars []string
	if authServiceName == "" {
		missingVars = append(missingVars, envAuthServiceName)
	}
	if authNamespace == "" {
		missingVars = append(missingVars, envAuthNamespace)
	}
	if authPort == "" {
		missingVars = append(missingVars, envAuthPort)
	}

	// ... validate other configs ...

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("missing required environment variable(s): %s", strings.Join(missingVars, ", "))
	}

	// --- Construct Auth URL ---
	// Default path if not provided
	if authHealthPath == "" {
		authHealthPath = "/health" // Default health path
		log.Printf("INFO: Environment variable %s not set, using default path: %s", envAuthHealthPath, authHealthPath)
	}
	// Ensure path starts with a slash
	if !strings.HasPrefix(authHealthPath, "/") {
		authHealthPath = "/" + authHealthPath
	}
	// Format: http://<service>.<namespace>.svc.cluster.local:<port><path>
	// Note: Assumes http protocol. Change to https if needed.
	authURL := fmt.Sprintf("http://%s.%s.svc.cluster.local:%s%s",
		authServiceName,
		authNamespace,
		authPort,
		authHealthPath,
	)
	log.Printf("INFO: Constructed Auth Service Check URL: %s", authURL)
	log.Printf("INFO: Database URL configured (details omitted).")

	// --- Create Components ---
	authComp, err := configurators.NewAuthComponent(authURL) // Use the constructed URL
	if err != nil {
		return nil, fmt.Errorf("failed to create auth component: %w", err)
	}

	// Postgres Component creation (assuming it was removed as requested)
	// pgComp, err := configurators.NewPostgresComponent(dbURL)
	// if err != nil { return nil, fmt.Errorf("failed to create postgres component: %w", err) }

	// ... create other components ...

	// --- Define Initialization Order ---
	components := []initTypes.InitializableComponent{
		// pgComp, // Removed as requested
		authComp, // Check Auth service
		// Add other components here in the desired order
	}

	if len(components) == 0 {
		log.Println("WARN: No components configured for initialization.")
	} else {
		log.Printf("INFO: Prepared %d components for initialization.", len(components))
	}
	return components, nil
}

// Command creates the cobra command for the app-init job.
func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app-init",
		Short: "Checks availability and configures dependent services.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Setup logger for runner
			logger := log.New(os.Stdout, "APP-INIT: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

			logger.Println("INFO: Starting App Init Job RunE...")
			ctx := cmd.Context()

			// Build the list of components to process
			componentsToRun, err := buildComponentList() // This now reads env vars, constructs URLs, creates components
			if err != nil {
				logger.Printf("ERROR: Failed to build component list: %v", err)
				return err // Return config/build error
			}

			// Create and execute the runner
			runner := NewRunner(logger)            // Assumes NewRunner is defined in runner.go
			err = runner.Run(ctx, componentsToRun) // Call the Run method
			if err != nil {
				logger.Printf("ERROR: Component initialization sequence failed: %v", err)
				return err // Return execution error
			}

			logger.Println("INFO: App Init Job finished successfully.")
			return nil // Success
		},
	}
	cmd.SilenceUsage = true // Prevent cobra from printing usage when RunE returns an error
	return cmd
}

// Note: This assumes NewRunner(logger *log.Logger) *Runner and
// func (r *Runner) Run(ctx context.Context, components []initTypes.InitializableComponent) error
// are defined in runner.go within the same app_init package.
