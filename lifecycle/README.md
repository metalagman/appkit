# Lifecycle

The `lifecycle` package provides primitives for managing the lifecycle of application components, standardizing how long-running processes start, run, and stop.

## Interfaces

### Runnable

`Runnable` represents a long-running process that blocks until the context is canceled or a fatal error occurs. This is the core interface for the main application loop.

```go
type Runnable interface {
    Run(ctx context.Context) error
}
```

### Lifecycle

`Lifecycle` represents a component that has distinct start and stop phases. Unlike `Runnable`, the `Start` and `Stop` methods are expected to be non-blocking.

```go
type Lifecycle interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

## Usage

### Converting Lifecycle to Runnable

The `ToRunnable` function adapts a `Lifecycle` component into a `Runnable`. This handles the orchestration of starting the component, waiting for the context to be canceled, and then stopping the component.

It ensures that `Stop` is called even if the application context is canceled, allowing for graceful shutdowns.

```go
type MyService struct{}

func (s *MyService) Start(ctx context.Context) error {
    // Initialize resources...
    return nil
}

func (s *MyService) Stop(ctx context.Context) error {
    // Cleanup resources...
    return nil
}

func main() {
    svc := &MyService{}
    
    // Convert to Runnable with optional timeouts
    runnable := lifecycle.ToRunnable(svc, 
        lifecycle.WithStartTimeout(5*time.Second),
        lifecycle.WithStopTimeout(10*time.Second),
    )

    // Run blocks until context is done or error occurs
    if err := runnable.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

### Functional Runnable

You can also create a `Runnable` from a simple function using `RunnableFunc`:

```go
runnable := lifecycle.RunnableFunc(func(ctx context.Context) error {
    // Do work...
    <-ctx.Done()
    return nil
})
```

### Running Multiple Components

`Runnable` is designed to work seamlessly with `errgroup` for managing multiple concurrent processes. If any process returns an error, the context is canceled, triggering a graceful shutdown for all other processes.

```go
func main() {
    g, ctx := errgroup.WithContext(context.Background())

    // List of components to run
    runnables := []lifecycle.Runnable{
        lifecycle.ToRunnable(&HttpServer{}),
        lifecycle.ToRunnable(&WorkerPool{}),
        lifecycle.RunnableFunc(func(ctx context.Context) error {
            // Some other background task
            return nil
        }),
    }

    for _, r := range runnables {
        r := r // capture range variable
        g.Go(func() error {
            return r.Run(ctx)
        })
    }

    if err := g.Wait(); err != nil {
        log.Fatalf("application error: %v", err)
    }
}
```
