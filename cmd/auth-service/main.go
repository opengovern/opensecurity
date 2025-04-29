package main

import (
	"context"
	"fmt"
	"os"
	"os/signal" // Added for checking env var
	"syscall"

	"github.com/opengovern/opensecurity/services/auth"
)

func main() {
	// Create base context
	ctx := context.Background()
	// Create cancellable context
	ctx, cancel := context.WithCancel(ctx)

	// Set up signal channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Ensure cleanup happens on exit
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	// Goroutine to handle signals
	go func() {
		select {
		case sig := <-c:
			fmt.Printf("\nReceived signal: %s. Shutting down...\n", sig) // Keep user feedback on signal
			cancel()
		case <-ctx.Done():
		}
	}()

	// Execute the root command defined locally
	if err := auth.Command().ExecuteContext(ctx); err != nil {
		// Print error consistent with compliance example
		fmt.Println(err)
		os.Exit(1) // Exit with non-zero code on error
	}
}
