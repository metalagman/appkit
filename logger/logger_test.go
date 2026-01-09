package logger_test

import (
	"bytes"
	"testing"

	"github.com/metalagman/appkit/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestZerolog(t *testing.T) {
	var buf bytes.Buffer

	zl := zerolog.New(&buf)
	l := logger.NewZerolog(zl)

	l.Infof("test %s", "info")
	assert.Contains(t, buf.String(), "test info")

	buf.Reset()
	l.Debug("test debug")
	assert.Contains(t, buf.String(), "test debug")
}

func TestNop(t *testing.T) {
	l := logger.NewNop()

	l.Infof("test %s", "info")
	l.Debug("test debug")
	// Should not panic and do nothing
}
