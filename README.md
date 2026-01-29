# AppKit

AppKit is a collection of opinionated Go packages for building robust and scalable applications. It aims to provide standard patterns for common problems like lifecycle management, configuration, and more.

## Packages

### Lifecycle

The [`lifecycle`](./lifecycle) package provides primitives for managing the start/stop sequences of application components. It introduces `Runnable` and `Lifecycle` interfaces to standardize long-running processes and their coordination.

- **Runnable:** Represents a blocking process (e.g., HTTP server, worker loop).
- **Lifecycle:** Represents a component with non-blocking Start/Stop phases.
- **Adapters:** Helpers to convert `Lifecycle` components into `Runnable` processes that can be easily managed with `errgroup`.

### Logger

The [`logger`](./logger) package defines a common logging interface and provides implementations for various logging libraries.

- **Interface:** A standard `Logger` interface for consistent logging across the application.
- **Implementations:** Includes `zerolog`, `slog`, and `Nop` logger implementations.

### Version

The [`version`](./version) package provides a way to manage and retrieve application version information, which can be configured at build time using `ldflags`.

- **Build Info:** Access structured build information including version, git commit, and build date.
- **ldflags:** Support for setting version metadata during the `go build` process.

### Worker Pool

The [`workerpool`](./workerpool) package provides a service for running small parts of code (called jobs) in the background. It is compatible with the `Lifecycle` interface.

- **Jobs:** Functions with context support, timeouts, and rich retry strategies.
- **Pool:** Manages a set of workers to execute jobs concurrently.
- **Middleware:** Support for panic recovery, logging, retries, and timeouts.