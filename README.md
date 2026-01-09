# AppKit

AppKit is a collection of opinionated Go packages for building robust and scalable applications. It aims to provide standard patterns for common problems like lifecycle management, configuration, and more.

## Packages

### Lifecycle

The [`lifecycle`](./lifecycle) package provides primitives for managing the start/stop sequences of application components. It introduces `Runnable` and `Lifecycle` interfaces to standardize long-running processes and their coordination.

- **Runnable:** Represents a blocking process (e.g., HTTP server, worker loop).
- **Lifecycle:** Represents a component with non-blocking Start/Stop phases.
- **Adapters:** Helpers to convert `Lifecycle` components into `Runnable` processes that can be easily managed with `errgroup`.