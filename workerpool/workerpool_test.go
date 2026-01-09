package workerpool_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Rican7/retry/strategy"
	"github.com/metalagman/appkit/lifecycle"
	"github.com/metalagman/appkit/logger"
	"github.com/metalagman/appkit/workerpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLifecycleCompatibility(t *testing.T) {
	pool, _ := workerpool.New()

	var _ lifecycle.Lifecycle = pool

	ctx, cancel := context.WithCancel(context.Background())
	runnable := lifecycle.ToRunnable(pool)

	done := make(chan error)

	go func() {
		done <- runnable.Run(ctx)
	}()

	// Give it some time to start
	time.Sleep(10 * time.Millisecond)

	jobDone := make(chan struct{})
	ok := pool.Submit(context.Background(), func(ctx context.Context) error {
		close(jobDone)

		return nil
	})
	assert.True(t, ok)

	select {
	case <-jobDone:
		// success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("job was not executed")
	}

	cancel()

	err := <-done
	assert.NoError(t, err)
}

func Example() {
	pool, _ := workerpool.New()

	pool.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	job := func(ctx context.Context) error {
		fmt.Println("hello")
		wg.Done()

		return nil
	}

	pool.Submit(context.Background(), job)
	wg.Wait()
	pool.Stop(context.Background())
	// Output: hello
}

func Example_advancedUsage() {
	pool, _ := workerpool.New()

	pool.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	job := func(ctx context.Context) error {
		fmt.Println("job executed")
		wg.Done()

		return nil
	}
	// add 3 seconds timeout for a job execution
	job = workerpool.AddTimeout(job, time.Second*3)
	// retry job execution within 5 attempts
	job = workerpool.AddRetry(job, strategy.Limit(5))

	pool.Submit(context.Background(), job)
	wg.Wait()
	pool.Stop(context.Background())
	// Output: job executed
}

func TestAddPanicRecovery(t *testing.T) {
	job := func(ctx context.Context) error {
		panic("oops")
	}
	job = workerpool.AddPanicRecovery(job)
	err := job(context.Background())
	assert.EqualError(t, err, "panic recovered: oops")
}

func ExampleAddPanicRecovery() {
	pool, _ := workerpool.New(workerpool.WithLogger(logger.NewNop())) // Nop logger to avoid output noise
	pool.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	job := func(ctx context.Context) error {
		defer wg.Done()

		panic("oops")
	}

	pool.Submit(context.Background(), job)
	wg.Wait()
	// Panic is recovered by middleware and logged by the pool worker.

	pool.Stop(context.Background())
	// Output:
}

func TestSubmitToStoppedPool(t *testing.T) {
	pool, _ := workerpool.New()

	pool.Start(context.Background())
	pool.Stop(context.Background())

	ok := pool.Submit(context.Background(), func(ctx context.Context) error { return nil })
	assert.False(t, ok, "Submitting to a stopped pool should return false")
}

func ExampleAddRetry() {
	pool, _ := workerpool.New(workerpool.WithLogger(logger.NewNop()))

	pool.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	retryCount := 0

	job := func(ctx context.Context) error {
		retryCount++

		if retryCount < 5 {
			return errors.New("temporary error")
		}

		wg.Done()

		return nil
	}
	job = workerpool.AddRetry(job, strategy.Limit(5))

	pool.Submit(context.Background(), job)
	wg.Wait()
	pool.Stop(context.Background())

	fmt.Println(retryCount)
	// Output: 5
}

func ExampleAddPostRun() {
	pool, _ := workerpool.New(workerpool.WithLogger(logger.NewNop()))

	pool.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	job := func(ctx context.Context) error {
		return errors.New("unrecoverable error")
	}
	job = workerpool.AddPostRun(job, func(err error) {
		if err != nil {
			fmt.Println(err.Error())

			wg.Done()
		}
	})

	pool.Submit(context.Background(), job)
	wg.Wait()
	pool.Stop(context.Background())

	// Output: unrecoverable error
}

func ExampleAddTimeout() {
	pool, _ := workerpool.New(workerpool.WithLogger(logger.NewNop()))

	pool.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	job := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return nil
		}
	}
	// adding timeout to job
	job = workerpool.AddTimeout(job, time.Millisecond*10)
	// adding post run hook to job to output timeout error
	job = workerpool.AddPostRun(job, func(err error) {
		if err != nil {
			fmt.Println(err.Error())

			wg.Done()
		}
	})

	pool.Submit(context.Background(), job)
	wg.Wait()
	pool.Stop(context.Background())

	// Output: context deadline exceeded
}

func ExamplePool_Submit() {
	p, _ := workerpool.New(workerpool.WithLogger(logger.NewNop()))

	p.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)

	ctxKey := struct{}{}

	job := func(ctx context.Context) error {
		if v, ok := ctx.Value(ctxKey).(string); ok {
			fmt.Println(v)
		}

		wg.Done()

		return nil
	}

	// submitting with custom context directly
	ctx := context.WithValue(context.Background(), ctxKey, "from context")
	p.Submit(ctx, job)

	wg.Wait()
	p.Stop(context.Background())

	// Output: from context
}

func ExampleAddLogger() {
	// mock time for test purposes
	zerolog.TimestampFunc = func() time.Time {
		t, _ := time.Parse("2006-01-02", "2021-01-01")

		return t
	}

	zl := zerolog.New(os.Stdout).Level(zerolog.WarnLevel).With().Timestamp().Logger()
	l := logger.NewZerolog(zl)
	l1 := l.With("logger", "logger1")
	l2 := l.With("logger", "logger2")

	pool, _ := workerpool.New(workerpool.WithLogger(l))

	pool.Start(context.Background())

	var wg sync.WaitGroup

	wg.Add(2)

	job1 := func(ctx context.Context) error {
		logger.FromContext(ctx).Warnf("hello from job1")

		wg.Done()

		return nil
	}
	job2 := func(ctx context.Context) error {
		logger.FromContext(ctx).Warnf("hello from job2")

		wg.Done()

		return nil
	}
	// adding custom jobs
	job1 = workerpool.AddLogger(job1, l1)
	job2 = workerpool.AddLogger(job2, l2)

	pool.Submit(context.Background(), job1)
	pool.Submit(context.Background(), job2)

	wg.Wait()
	pool.Stop(context.Background())

	// Unordered output:
	// {"level":"warn","logger":"logger1","time":"2021-01-01T00:00:00Z","message":"hello from job1"}
	// {"level":"warn","logger":"logger2","time":"2021-01-01T00:00:00Z","message":"hello from job2"}
}

func TestNewValidation(t *testing.T) {
	pool, err := workerpool.New(workerpool.WithNumWorkers(1))
	assert.Error(t, err)
	assert.Nil(t, pool)
	assert.Contains(t, err.Error(), "numWorkers")
}
