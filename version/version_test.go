package version_test

import (
	"testing"

	"github.com/metalagman/appkit/version"
)

func TestGet(t *testing.T) {
	info := version.Get()
	if info.Version != "dev" {
		t.Errorf("expected version to be 'dev', got %s", info.Version)
	}
}

func TestGetVersion(t *testing.T) {
	v := version.GetVersion()
	if v != "dev" {
		t.Errorf("expected version to be 'dev', got %s", v)
	}
}

func TestString(t *testing.T) {
	s := version.String()
	expected := "dev (commit: , buildDate: )"

	if s != expected {
		t.Errorf("expected %s, got %s", expected, s)
	}
}