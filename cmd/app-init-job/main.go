// /Users/anil/workspace/opensecurity/cmd/app-init-job/main.go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// --- ADJUST IMPORT PATH ---
	// Replace with the actual import path for your app_init job package
	appinit "github.com/opengovern/opensecurity/jobs/app-init"
)

func main() {
	// --- Root Context Setup ---
	// Create a base context
	ctx := context.Background()
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(ctx)

	// --- Signal Handling Setup ---
	// Create a channel to listen for OS signals
	signalChan := make(chan os.Signal, 1)
	// Notify the channel for Interrupt (Ctrl+C) and Terminate signals
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Ensure cleanup happens when main exits
	defer func() {
		signal.Stop(signalChan) // Stop listening for signals
		cancel()                // Cancel the context
	}()

	// --- Goroutine to Cancel Context on Signal ---
	// Start a background goroutine that waits for either a signal or context cancellation
	go func() {
		select {
		case sig := <-signalChan:
			// Received an OS signal
			fmt.Fprintf(os.Stderr, "\nReceived signal: %s. Initiating shutdown...\n", sig)
			cancel() // Trigger context cancellation
		case <-ctx.Done():
			// Context was cancelled elsewhere (e.g., command finished/errored)
			// This case prevents the goroutine from leaking if the command finishes normally.
		}
	}()

	// --- Execute the Cobra Command ---
	fmt.Println("INFO: Executing App Init command...")
	// Get the command definition from the job package (appinit)
	cmd := appinit.Command()

	// Execute the command, passing the cancellable context
	if err := cmd.ExecuteContext(ctx); err != nil {
		// Cobra typically prints the error already when RunE returns an error.
		// Printing it again might be redundant, but ensures it's visible.
		fmt.Fprintf(os.Stderr, "ERROR: Command execution failed: %v\n", err)
		os.Exit(1) // Exit with a non-zero code indicating failure
	}

	// If ExecuteContext returns nil, it means RunE returned nil (success)
	fmt.Println("INFO: App Init command completed successfully.")
	os.Exit(0) // Exit with success code
}
