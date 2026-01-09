// Package logger defines a common logging interface and provides implementations for various logging libraries.
package logger

import "context"

// Logger is an interface for logging.
//
//nolint:interfacebloat
type Logger interface {
	// With returns a new logger with the given fields.
	With(args ...any) Logger
	// Error logs a message at the error level.
	Error(args ...any)
	// Errorf logs a message at the error level.
	Errorf(format string, args ...any)
	// Fatalf logs a message at the fatal level and calls os.Exit(1).
	Fatalf(format string, args ...any)
	// Fatal logs a message at the fatal level and calls os.Exit(1).
	Fatal(args ...any)
	// Info logs a message at the info level.
	Info(args ...any)
	// Infof logs a message at the info level.
	Infof(format string, args ...any)
	// Warn logs a message at the warn level.
	Warn(args ...any)
	// Warnf logs a message at the warn level.
	Warnf(format string, args ...any)
	// Debug logs a message at the debug level.
	Debug(args ...any)
	// Debugf logs a message at the debug level.
	Debugf(format string, args ...any)
}

type ctxKey struct{}

// FromContext returns the Logger associated with the context.
// If no Logger is found, it returns a Nop logger.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}

	return NewNop()
}

// ToContext returns a new context with the Logger attached.
func ToContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}
