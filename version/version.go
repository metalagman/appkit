// Package version provides version information for the application.
package version

import (
	"fmt"
	"runtime"
)

var (
	// version is the current version of the application.
	version = "dev"
	// metadata is extra build time data.
	metadata = ""
	// gitCommit is the git sha1.
	gitCommit = ""
	// buildDate is the date the application was built.
	buildDate = ""
)

// BuildInfo describes the compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"git_commit,omitempty"`
	// BuildDate is the date the application was built.
	BuildDate string `json:"build_date,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"go_version,omitempty"`
	// Platform is the platform the application was built for.
	Platform string `json:"platform,omitempty"`
}

// GetVersion returns the semver string of the version.
func GetVersion() string {
	if metadata == "" {
		return version
	}

	return version + "+" + metadata
}

// Get returns build info.
func Get() BuildInfo {
	return BuildInfo{
		Version:   GetVersion(),
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a string representation of the version information.
func String() string {
	return fmt.Sprintf("%s (commit: %s, buildDate: %s)", GetVersion(), gitCommit, buildDate)
}