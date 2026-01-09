package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunnableFunc(t *testing.T) {
	called := false
	errExpected := errors.New("expected error")

	f := RunnableFunc(func(ctx context.Context) error {
		called = true
		return errExpected
	})

	err := f.Run(context.Background())
	if !called {
		t.Error("RunnableFunc did not call the underlying function")
	}
	if err != errExpected {
		t.Errorf("expected error %v, got %v", errExpected, err)
	}
}

// mockLifecycle is a helper for testing ToRunnable
type mockLifecycle struct {
	startFunc func(ctx context.Context) error
	stopFunc  func(ctx context.Context) error
	started   bool
	stopped   bool
}

func (m *mockLifecycle) Start(ctx context.Context) error {
	m.started = true
	if m.startFunc != nil {
		return m.startFunc(ctx)
	}
	return nil
}

func (m *mockLifecycle) Stop(ctx context.Context) error {
	m.stopped = true
	if m.stopFunc != nil {
		return m.stopFunc(ctx)
	}
	return nil
}

func TestToRunnable_Success(t *testing.T) {
	ml := &mockLifecycle{}
	r := ToRunnable(ml)

	ctx, cancel := context.WithCancel(context.Background())
	
	// Run in a separate goroutine so we can cancel it
	errCh := make(chan error)
	go func() {
		errCh <- r.Run(ctx)
	}()

	// Wait a bit to ensure Start has been called and it's blocking
	time.Sleep(50 * time.Millisecond)

	if !ml.started {
		t.Error("Start was not called")
	}
	if ml.stopped {
		t.Error("Stop was called too early")
	}

	// Cancel context to trigger shutdown
	cancel()

	err := <-errCh
	if err != nil {
		t.Errorf("Run returned unexpected error: %v", err)
	}

	if !ml.stopped {
		t.Error("Stop was not called")
	}
}

func TestToRunnable_StartError(t *testing.T) {
	expectedErr := errors.New("start failure")
	ml := &mockLifecycle{
		startFunc: func(ctx context.Context) error {
			return expectedErr
		},
	}
	r := ToRunnable(ml)

	err := r.Run(context.Background())
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if !ml.started {
		t.Error("Start was not called")
	}
	// Stop should NOT be called if Start failed
	if ml.stopped {
		t.Error("Stop was called after start failure")
	}
}

func TestToRunnable_StopError(t *testing.T) {
	expectedErr := errors.New("stop failure")
	ml := &mockLifecycle{
		stopFunc: func(ctx context.Context) error {
			return expectedErr
		},
	}
	r := ToRunnable(ml)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := r.Run(ctx)
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if !ml.started {
		t.Error("Start was not called")
	}
	if !ml.stopped {
		t.Error("Stop was not called")
	}
}

func TestToRunnable_StartTimeout(t *testing.T) {
	ml := &mockLifecycle{
		startFunc: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		},
	}
	// Set a short timeout
	r := ToRunnable(ml, WithStartTimeout(10*time.Millisecond))

	err := r.Run(context.Background())
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestToRunnable_StopTimeout(t *testing.T) {
	ml := &mockLifecycle{
		stopFunc: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		},
	}
	// Set a short timeout
	r := ToRunnable(ml, WithStopTimeout(10*time.Millisecond))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := r.Run(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
