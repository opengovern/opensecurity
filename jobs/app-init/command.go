package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// <<< ADJUST IMPORT PATH to your actual job package location >>>
)

func main() {
	// --- Context Setup ---
	ctx := context.Background()
	// Create a cancellable context which finishes on interrupt/terminate signals
	ctx, cancel := context.WithCancel(ctx)

	// Setup signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Ensure cleanup happens
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	// Goroutine listens for signals or context cancellation
	go func() {
		select {
		case sig := <-signalChan:
			fmt.Fprintf(os.Stderr, "Received signal: %s. Shutting down...\n", sig)
			cancel()
		case <-ctx.Done():
			// Context already cancelled (e.g., command finished/errored)
		}
	}()

	// --- Command Execution ---
	fmt.Println("INFO: Executing service health check command...")
	// Get the command definition from the job package
	cmd := healthchecker.Command()

	// Execute the command with the cancellable context
	if err := cmd.ExecuteContext(ctx); err != nil {
		// Cobra usually prints the error, but we print it again for clarity if needed
		// fmt.Fprintf(os.Stderr, "ERROR: Command execution failed: %v\n", err)
		os.Exit(1) // Exit with error code if command returns error
	}

	// If ExecuteContext returns nil, the command's RunE returned nil (success)
	fmt.Println("INFO: Service health check command completed successfully.")
	os.Exit(0) // Exit successfully
}
