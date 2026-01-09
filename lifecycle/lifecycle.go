package lifecycle

import (
	"context"
	"time"
)

// Runnable represents a long-running process that blocks until the context is canceled or an error occurs.
type Runnable interface {
	// Run starts the process.
	// The context passed to Run is the application lifecycle context.
	// Run must block until the context is canceled or a fatal error occurs.
	Run(ctx context.Context) error
}

// RunnableFunc is a function adapter for the Runnable interface.
type RunnableFunc func(ctx context.Context) error

// Run calls f(ctx).
func (f RunnableFunc) Run(ctx context.Context) error {
	return f(ctx)
}

// Lifecycle represents a component that has a distinct start and stop phase.
type Lifecycle interface {
	// Start starts the component.
	// This method must be non-blocking.
	// The context is used for the start operation itself (e.g., timeout), not the application lifecycle.
	Start(ctx context.Context) error

	// Stop stops the component.
	// This method must be non-blocking.
	// The context is used for the stop operation itself (e.g., timeout), not the application lifecycle.
	Stop(ctx context.Context) error
}

// Option configures the adapter.
type Option func(*adapter)

// WithStartTimeout sets the timeout for the Start operation.
func WithStartTimeout(timeout time.Duration) Option {
	return func(a *adapter) {
		a.startTimeout = timeout
	}
}

// WithStopTimeout sets the timeout for the Stop operation.
func WithStopTimeout(timeout time.Duration) Option {
	return func(a *adapter) {
		a.stopTimeout = timeout
	}
}

type adapter struct {
	lifecycle    Lifecycle
	startTimeout time.Duration
	stopTimeout  time.Duration
}

// ToRunnable converts a Lifecycle to a Runnable.
// It handles the start and stop sequence, respecting the application lifecycle context.
func ToRunnable(l Lifecycle, opts ...Option) Runnable {
	a := &adapter{
		lifecycle:    l,
		startTimeout: 30 * time.Second,
		stopTimeout:  30 * time.Second,
	}

	for _, opt := range opts {
		opt(a)
	}

	return RunnableFunc(func(ctx context.Context) error {
		startCtx, startCancel := context.WithTimeout(ctx, a.startTimeout)
		defer startCancel()

		if err := a.lifecycle.Start(startCtx); err != nil {
			return err
		}

		<-ctx.Done()

		stopCtx, stopCancel := context.WithTimeout(context.Background(), a.stopTimeout)
		defer stopCancel()

		if err := a.lifecycle.Stop(stopCtx); err != nil {
			return err
		}

		return nil
	})
}
