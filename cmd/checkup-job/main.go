package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/opengovern/opensecurity/jobs/checkup-job"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := checkup.WorkerCommand().ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
