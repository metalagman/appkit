package logger

// Nop is a logger that does nothing.
type Nop struct{}

var _ Logger = Nop{}

// NewNop creates a new Nop.
func NewNop() Nop {
	return Nop{}
}

// With returns a new logger with the given fields.
func (Nop) With(_ ...any) Logger { return Nop{} }

// Error logs a message at the error level.
func (Nop) Error(_ ...any) {}

// Errorf logs a message at the error level.
func (Nop) Errorf(_ string, _ ...any) {}

// Fatalf logs a message at the fatal level and calls os.Exit(1).
func (Nop) Fatalf(_ string, _ ...any) {}

// Fatal logs a message at the fatal level and calls os.Exit(1).
func (Nop) Fatal(_ ...any) {}

// Info logs a message at the info level.
func (Nop) Info(_ ...any) {}

// Infof logs a message at the info level.
func (Nop) Infof(_ string, _ ...any) {}

// Warn logs a message at the warn level.
func (Nop) Warn(_ ...any) {}

// Warnf logs a message at the warn level.
func (Nop) Warnf(_ string, _ ...any) {}

// Debug logs a message at the debug level.
func (Nop) Debug(_ ...any) {}

// Debugf logs a message at the debug level.
func (Nop) Debugf(_ string, _ ...any) {}
