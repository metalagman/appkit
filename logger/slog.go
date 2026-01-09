package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Slog is a logger that uses the standard library's log/slog package.
type Slog struct {
	l *slog.Logger
}

var _ Logger = (*Slog)(nil)

// NewSlog creates a new Slog.
func NewSlog() *Slog {
	return &Slog{
		l: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
}

// With returns a new logger with the given fields.
func (s *Slog) With(args ...any) Logger {
	return &Slog{l: s.l.With(args...)}
}

// Error logs a message at the error level.
func (s *Slog) Error(args ...any) {
	s.l.Error(fmt.Sprint(args...))
}

// Errorf logs a message at the error level.
func (s *Slog) Errorf(format string, args ...any) {
	s.l.Error(fmt.Sprintf(format, args...))
}

// Fatalf logs a message at the fatal level and calls os.Exit(1).
func (s *Slog) Fatalf(format string, args ...any) {
	s.l.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

// Fatal logs a message at the fatal level and calls os.Exit(1).
func (s *Slog) Fatal(args ...any) {
	s.l.Error(fmt.Sprint(args...))
	os.Exit(1)
}

// Info logs a message at the info level.
func (s *Slog) Info(args ...any) {
	s.l.Info(fmt.Sprint(args...))
}

// Infof logs a message at the info level.
func (s *Slog) Infof(format string, args ...any) {
	s.l.Info(fmt.Sprintf(format, args...))
}

// Warn logs a message at the warn level.
func (s *Slog) Warn(args ...any) {
	s.l.Warn(fmt.Sprint(args...))
}

// Warnf logs a message at the warn level.
func (s *Slog) Warnf(format string, args ...any) {
	s.l.Warn(fmt.Sprintf(format, args...))
}

// Debug logs a message at the debug level.
func (s *Slog) Debug(args ...any) {
	s.l.Debug(fmt.Sprint(args...))
}

// Debugf logs a message at the debug level.
func (s *Slog) Debugf(format string, args ...any) {
	s.l.Debug(fmt.Sprintf(format, args...))
}

// WithContext returns a logger that is bound to a specific context.
// Currently it just returns the same logger since we don't have context-specific fields yet.
func (s *Slog) WithContext(_ context.Context) *Slog {
	return s
}
