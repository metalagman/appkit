package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Zerolog is a logger that uses the zerolog package.
type Zerolog struct {
	l zerolog.Logger
}

var _ Logger = (*Zerolog)(nil)

// NewZerolog creates a new Zerolog.
func NewZerolog(l zerolog.Logger) *Zerolog {
	return &Zerolog{l: l}
}

// With returns a new logger with the given fields.
func (z *Zerolog) With(args ...any) Logger {
	if len(args) == 0 {
		return z
	}

	ctx := z.l.With()

	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			ctx = ctx.Interface(fmt.Sprint(args[i]), args[i+1])
		} else {
			ctx = ctx.Interface(fmt.Sprint(args[i]), nil)
		}
	}

	return &Zerolog{l: ctx.Logger()}
}

// Error logs a message at the error level.
func (z *Zerolog) Error(args ...any) {
	z.l.Error().Msg(fmt.Sprint(args...))
}

// Errorf logs a message at the error level.
func (z *Zerolog) Errorf(format string, args ...any) {
	z.l.Error().Msgf(format, args...)
}

// Fatalf logs a message at the fatal level and calls os.Exit(1).
func (z *Zerolog) Fatalf(format string, args ...any) {
	z.l.Fatal().Msgf(format, args...)
}

// Fatal logs a message at the fatal level and calls os.Exit(1).
func (z *Zerolog) Fatal(args ...any) {
	z.l.Fatal().Msg(fmt.Sprint(args...))
}

// Info logs a message at the info level.
func (z *Zerolog) Info(args ...any) {
	z.l.Info().Msg(fmt.Sprint(args...))
}

// Infof logs a message at the info level.
func (z *Zerolog) Infof(format string, args ...any) {
	z.l.Info().Msgf(format, args...)
}

// Warn logs a message at the warn level.
func (z *Zerolog) Warn(args ...any) {
	z.l.Warn().Msg(fmt.Sprint(args...))
}

// Warnf logs a message at the warn level.
func (z *Zerolog) Warnf(format string, args ...any) {
	z.l.Warn().Msgf(format, args...)
}

// Debug logs a message at the debug level.
func (z *Zerolog) Debug(args ...any) {
	z.l.Debug().Msg(fmt.Sprint(args...))
}

// Debugf logs a message at the debug level.
func (z *Zerolog) Debugf(format string, args ...any) {
	z.l.Debug().Msgf(format, args...)
}
