// Package contexts provides utilities for working with Go contexts.
package contexts

import (
	"context"
)

// MergeCancel returns a context that contains the values of parent and is canceled
// when either parent or other is canceled.
//
// If parent is canceled, the returned context is canceled with the parent's cause.
// If other is canceled, the returned context is canceled with the other's cause.
func MergeCancel(parent, other context.Context) (context.Context, context.CancelCauseFunc) {
	ctx, cancel := context.WithCancelCause(parent)
	stop := context.AfterFunc(other, func() {
		cancel(context.Cause(other))
	})

	return ctx, func(err error) {
		stop()

		if err == nil {
			err = context.Canceled
		}

		cancel(err)
	}
}
