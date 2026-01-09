package contexts_test

import (
	"context"
	"errors"
	"testing"

	"github.com/metalagman/appkit/internal/contexts"
	"github.com/stretchr/testify/assert"
)

func TestMergeCancel(t *testing.T) {
	t.Run("CancelParent", func(t *testing.T) {
		cause := errors.New("parent error")
		ctx1, cancel1 := context.WithCancelCause(context.Background())
		ctx2, cancel2 := context.WithCancelCause(context.Background())

		defer cancel2(nil)

		mergedCtx, mergedCancel := contexts.MergeCancel(ctx1, ctx2)
		defer mergedCancel(nil)

		cancel1(cause)
		<-mergedCtx.Done()
		assert.Equal(t, cause, context.Cause(mergedCtx))
		assert.Equal(t, context.Canceled, mergedCtx.Err())
	})

	t.Run("CancelOther", func(t *testing.T) {
		cause := errors.New("other error")
		ctx1, cancel1 := context.WithCancelCause(context.Background())

		defer cancel1(nil)

		ctx2, cancel2 := context.WithCancelCause(context.Background())

		mergedCtx, mergedCancel := contexts.MergeCancel(ctx1, ctx2)
		defer mergedCancel(nil)

		cancel2(cause)
		<-mergedCtx.Done()
		assert.Equal(t, cause, context.Cause(mergedCtx))
	})

	t.Run("ManualCancel", func(t *testing.T) {
		ctx1 := context.Background()
		ctx2 := context.Background()

		mergedCtx, mergedCancel := contexts.MergeCancel(ctx1, ctx2)
		mergedCancel(nil)

		<-mergedCtx.Done()
		assert.Equal(t, context.Canceled, mergedCtx.Err())
	})

	t.Run("ValuePreservation", func(t *testing.T) {
		type key string

		ctx1 := context.WithValue(context.Background(), key("foo"), "bar")
		ctx2 := context.Background()

		mergedCtx, mergedCancel := contexts.MergeCancel(ctx1, ctx2)
		defer mergedCancel(nil)

		assert.Equal(t, "bar", mergedCtx.Value(key("foo")))
	})
}
