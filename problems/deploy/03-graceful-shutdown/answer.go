//go:build ignore

package main

import (
	"context"
	"time"
)

type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

func GracefulShutdown(shutdowner Shutdowner, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return shutdowner.Shutdown(ctx)
}

func main() {}
