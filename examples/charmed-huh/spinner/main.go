// examples/charmed-huh/spinner/main.go
// A spinner for work that happens after the form. Action runs the work;
// ActionWithErr propagates its error and honors the context.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"charm.land/huh/v2/spinner"
)

func simulateWork(ctx context.Context) error {
	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := spinner.New().
		Type(spinner.Dots).
		Title("Simulating local work...").
		Context(ctx).
		ActionWithErr(simulateWork).
		Run()
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Fprintln(os.Stderr, "timed out")
		os.Exit(1)
	case err != nil:
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("done")
}
