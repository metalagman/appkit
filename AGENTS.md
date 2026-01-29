# Versioning Strategy for Agents

This document describes how the `version` package is structured and how automated agents or CI/CD systems should interact with it to provide consistent versioning across the AppKit ecosystem.

## Design Philosophy

The `version` package follows the pattern established by the [Helm](https://github.com/helm/helm) project. It uses unexported variables that are injected at build time using Go's `-ldflags`. This approach ensures that the binary itself contains all necessary metadata without requiring external configuration files or environment variables at runtime.

## Core Components

### 1. BuildInfo Struct
The `BuildInfo` struct is the primary data container for versioning information. It is designed to be easily serialized to JSON for API or CLI consumption.

```go
type BuildInfo struct {
	Version   string // The semver version
	GitCommit string // The SHA1 of the commit
	BuildDate string // RFC3339 formatted build timestamp
	GoVersion string // Version of the Go compiler
	Platform  string // OS/Arch information
}
```

### 2. ldflags Targets
Agents should target the following internal variables within the `github.com/metalagman/appkit/version` package:

| Variable | Description | Recommended Agent Action |
| :--- | :--- | :--- |
| `version` | The semver tag | Use `git describe --tags --always` |
| `gitCommit` | The commit hash | Use `git rev-parse HEAD` |
| `metadata` | Prerelease info | Append build numbers or branch names |
| `buildDate` | Build timestamp | Use `date -u +%Y-%m-%dT%H:%M:%SZ` |

## Automated Build Example

When an agent is performing a release build, it should execute a command similar to the following:

```bash
VERSION_PKG="github.com/metalagman/appkit/version"
GIT_COMMIT=$(git rev-parse HEAD)
VERSION=$(git describe --tags --always)
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -ldflags "
  -X ${VERSION_PKG}.version=${VERSION} \
  -X ${VERSION_PKG}.gitCommit=${GIT_COMMIT} \
  -X ${VERSION_PKG}.buildDate=${BUILD_DATE}"
-o my-app ./main.go
```

## Integration Patterns

- **CLI Tools:** Use `version.String()` to implement a `--version` flag.
- **Web Services:** Use `version.Get()` to populate a `/version` or `/health` endpoint with JSON metadata.
- **Loggers:** At application startup, log `version.String()` to provide context in log aggregators.
