package main

import (
	"context"
	"fmt"
	"github.com/opengovern/opensecurity/services/tasks"
	"os"
	"os/signal"
	"syscall"
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

	if err := tasks.Command().ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
