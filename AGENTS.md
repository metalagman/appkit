# Agent Integration Guide for AppKit

This document provides a comprehensive guide for AI agents, CI/CD systems, and automated tools on how to interact with, extend, and utilize the AppKit repository.

## Table of Contents

- [Versioning Strategy (version)](#versioning-strategy-version)
- [Application Lifecycle (lifecycle)](#application-lifecycle-lifecycle)
- [Logging Standards (logger)](#logging-standards-logger)
- [Concurrent Job Execution (workerpool)](#concurrent-job-execution-workerpool)

---

## Versioning Strategy (version)

The `version` package is modeled after the Helm project, utilizing unexported variables injected at build time via `-ldflags`.

### ldflags Targets
Agents should target the following internal variables within the `github.com/metalagman/appkit/version` package:

| Variable | Description | Recommended Agent Action |
| :--- | :--- | :--- |
| `version` | The semver tag | `git describe --tags --always` |
| `gitCommit` | The commit hash | `git rev-parse HEAD` |
| `metadata` | Prerelease info | Append build/branch info |
| `buildDate` | Build timestamp | `date -u +%Y-%m-%dT%H:%M:%SZ` |

### Usage Example
```bash
go build -ldflags "-X github.com/metalagman/appkit/version.version=v1.0.0" -o my-app
```

---

## Application Lifecycle (lifecycle)

The `lifecycle` package standardizes how components start and stop, ensuring graceful shutdowns.

### Core Interfaces
- **`Runnable`**: For blocking processes (e.g., servers). Implement `Run(ctx) error`.
- **`Lifecycle`**: For components with non-blocking phases. Implement `Start(ctx)` and `Stop(ctx)`.

### Agent Integration Patterns
1. **Component Adaptation**: Use `lifecycle.ToRunnable(myLifecycle)` to convert non-blocking components into managed processes.
2. **Orchestration**: Use `errgroup.Group` to manage multiple `Runnable` instances. When one fails, the context is canceled, triggering a cascade of graceful shutdowns.
3. **Timeouts**: Always use `WithStartTimeout` and `WithStopTimeout` when creating runnables to prevent hung processes during CI/CD health checks.

---

## Logging Standards (logger)

AppKit provides a pluggable logging interface to ensure consistent output across all components.

### Implementation Choices
- `logger.NewZerolog(zlog)`: Recommended for production (structured JSON).
- `logger.NewSlog(slog)`: Recommended for standard Go 1.21+ compatibility.
- `logger.NewNop()`: Use in unit tests to reduce noise.

### Context Integration
Agents should favor passing loggers through `context.Context` using `logger.ToContext(ctx, l)` and retrieving them via `logger.FromContext(ctx)`. This allows for request-scoped logging and tracing.

---

## Concurrent Job Execution (workerpool)

The `workerpool` package manages background task execution and is fully integrated with the `lifecycle` pattern.

### Key Features for Agents
- **Lifecycle Integration**: The `Pool` implements `Start` and `Stop` methods. It can be managed directly by the `lifecycle` package.
- **Context Merging**: Jobs submitted to the pool are automatically canceled if either the job-specific context expires or the pool itself is stopped.
- **Middleware**: Always wrap jobs with provided middleware:
    - `AddPanicRecovery`: Prevents a single job failure from crashing the entire worker.
    - `AddLogger`: Automatically attaches job IDs and durations to logs.

### Submission Pattern
```go
pool.Submit(ctx, func(ctx context.Context) error {
    // Perform task
    return nil
})
```