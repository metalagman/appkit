package workerpool

import (
	"context"
	"fmt"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/metalagman/appkit/logger"
)

// AddPanicRecovery wraps a job with panic recovery. If a panic occurs, it converts the panic into an error.
func AddPanicRecovery(job Job) Job {
	return func(ctx context.Context) (retErr error) {
		defer func() {
			p := recover()

			if p != nil {
				retErr = fmt.Errorf("panic recovered: %v", p)
			}
		}()

		return job(ctx)
	}
}

// AddLogger replaces the logger in the job context with the specified one.
func AddLogger(job Job, l logger.Logger) Job {
	return func(ctx context.Context) error {
		ctx = logger.ToContext(ctx, l)

		return job(ctx)
	}
}

// AddRetry wraps a job with retry strategies.
//
// See https://github.com/Rican7/retry for details.
func AddRetry(job Job, strategies ...strategy.Strategy) Job {
	return func(ctx context.Context) error {
		return retry.Retry(
			func(attempt uint) error {
				l := logger.FromContext(ctx).With("attempt", attempt)

				return AddLogger(job, l)(ctx)
			},
			strategies...,
		)
	}
}

// AddTimeout wraps a job with a timeout.
// The timeout will be set on the job context.
func AddTimeout(job Job, timeout time.Duration) Job {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return job(ctx)
	}
}

// AddPostRun wraps a job with a hook that is executed after the job completes or fails.
// The err parameter in the hook will contain the error returned by the job.
func AddPostRun(job Job, hook func(err error)) Job {
	return func(ctx context.Context) error {
		err := job(ctx)

		hook(err)

		return err
	}
}
