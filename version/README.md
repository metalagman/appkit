# Version Package

This package provides version information for the application. The variables can be configured at build time using `ldflags`.

## Usage

Import the package in your code:

```go
import "github.com/metalagman/appkit/version"

fmt.Println(version.String())
```

## Build Configuration

Use `-ldflags` to set the version information during build:

```bash
go build -ldflags "-X github.com/metalagman/appkit/version.version=v1.0.0 \
                   -X github.com/metalagman/appkit/version.metadata=beta.1 \
                   -X github.com/metalagman/appkit/version.gitCommit=$(git rev-parse HEAD) \
                   -X github.com/metalagman/appkit/version.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
```

## Variables

- `version`: The version of the application (e.g., `v1.2.3`).
- `metadata`: Extra build time data (e.g., `beta.1`).
- `gitCommit`: The source control revision (e.g., git commit hash).
- `buildDate`: The timestamp when the binary was built.

