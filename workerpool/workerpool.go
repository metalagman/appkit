/*
Package workerpool provides a service for running small parts of code (called jobs) in the background.

Jobs can have contexts, timeouts, and retry strategies.
*/
package workerpool

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/metalagman/appkit/internal/contexts"
	"github.com/metalagman/appkit/logger"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

// Job is a function that receives a context and is run asynchronously by the worker pool.
type Job func(ctx context.Context) error

type task struct {
	ctx context.Context
	fn  Job
}

// Pool is a worker pool that manages a set of worker goroutines to execute jobs.
type Pool struct {
	wg   *sync.WaitGroup
	opts Options

	jobs       chan task
	numWorkers int

	mu     sync.RWMutex
	cancel context.CancelFunc
}

// Options contain configuration for the worker pool.
//go:generate go tool options-gen -from-struct Options -constructor public -out-setter-name Option -out-filename options.gen.go
type Options struct {
	// logger sets the logger used by the pool and its workers.
	logger logger.Logger
	// numWorkers sets the number of worker goroutines to start.
	numWorkers int `option:"default:0" validate:"min=2"`
}

const (
	minWorkers               = 2
	defaultWorkersMultiplier = 2
)

// New creates a new Pool with the provided options.
func New(options ...Option) (*Pool, error) {
	opts := NewOptions(options...)

	if opts.logger == nil {
		opts.logger = logger.NewZerolog(log.Logger)
	}

	// If not set, use default
	if opts.numWorkers == 0 {
		opts.numWorkers = runtime.GOMAXPROCS(0) * defaultWorkersMultiplier
		if opts.numWorkers < minWorkers {
			opts.numWorkers = minWorkers
		}
	}

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &Pool{
		wg:         &sync.WaitGroup{},
		opts:       opts,
		numWorkers: opts.numWorkers,
	}, nil
}

// Start initializes the worker goroutines.
func (s *Pool) Start(_ context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		return nil // Already started
	}

	s.opts.logger.With("worker_num", s.numWorkers).Debug("Starting workers")

	poolCtx, poolCancel := context.WithCancel(context.Background())
	s.cancel = poolCancel
	s.jobs = make(chan task)

	s.wg.Add(s.numWorkers)

	for i := 0; i < s.numWorkers; i++ {
		go s.worker(i, poolCtx)
	}

	s.opts.logger.Debug("Done starting workers")

	return nil
}

// Stop shuts down all workers and waits for them to complete their current jobs.
// It respects the provided context's deadline for the waiting phase.
func (s *Pool) Stop(ctx context.Context) error {
	s.mu.Lock()

	cancel := s.cancel
	if cancel == nil {
		s.mu.Unlock()

		return nil
	}

	s.cancel = nil // Mark as stopped to prevent new submissions
	cancel()       // Signal all workers to stop
	s.mu.Unlock()

	s.opts.logger.With("worker_num", s.numWorkers).Debug("Shutting down workers")

	// Wait for workers to finish in a goroutine so we can select with ctx.Done()
	waitDone := make(chan struct{})

	go func() {
		s.wg.Wait()
		close(waitDone)
	}()

	var err error

	select {
	case <-waitDone:
		s.opts.logger.Debug("All workers stopped gracefully")
	case <-ctx.Done():
		err = ctx.Err()
		s.opts.logger.With("error", err).Error("Workers shutdown timed out or canceled")
	}

	close(s.jobs)

	return err
}

// Submit sends a job to the pool for execution with the given context.
// It returns false if the pool is not running.
func (s *Pool) Submit(ctx context.Context, job Job) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cancel == nil {
		return false
	}

	s.jobs <- task{ctx: ctx, fn: job}

	return true
}

func (s *Pool) worker(workerID int, pCtx context.Context) {
	defer s.wg.Done()

	l := s.opts.logger.With("worker_id", workerID)

	for {
		select {
		case <-pCtx.Done():
			return
		case t, ok := <-s.jobs:
			if !ok {
				return
			}

			s.runJob(pCtx, t, l)
		}
	}
}

func (s *Pool) runJob(pCtx context.Context, t task, l logger.Logger) {
	id := xid.New()
	now := time.Now()
	ll := l.With("job_id", id.String())

	ll.Debug("Running job")

	// Senior Pattern: Merge the incoming job context with the pool's lifecycle context.
	// This ensures the job is canceled if EITHER its own context expires OR the pool stops.
	jobCtx, jobCancel := contexts.MergeCancel(t.ctx, pCtx)

	err := AddLogger(AddPanicRecovery(t.fn), ll)(jobCtx)
	if err != nil {
		// Don't log "context canceled" as an error during shutdown or if job context was canceled
		if !(errors.Is(err, context.Canceled) && (pCtx.Err() != nil || t.ctx.Err() != nil)) {
			ll.With("error", err).Error("Error running job")
		}
	}

	jobCancel(err)

	ll.With("job_duration", time.Since(now)).Info("Done running job")
}
