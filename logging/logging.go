// Package logging routes xutils diagnostics to an application-provided logger.
// Configure the logger during startup, before creating background workers.
package logging

import "sync"

// Logger must support concurrent calls from xutils background workers.
type Logger interface {
	Debugf(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

// NopLogger discards all diagnostics.
type NopLogger struct{}

func (NopLogger) Debugf(string, ...any) {}
func (NopLogger) Warnf(string, ...any)  {}
func (NopLogger) Errorf(string, ...any) {}

// FuncLogger adapts package-level functions. Nil functions discard that level.
type FuncLogger struct {
	DebugFunc func(string, ...any)
	WarnFunc  func(string, ...any)
	ErrorFunc func(string, ...any)
}

func (l FuncLogger) Debugf(format string, args ...any) {
	if l.DebugFunc != nil {
		l.DebugFunc(format, args...)
	}
}

func (l FuncLogger) Warnf(format string, args ...any) {
	if l.WarnFunc != nil {
		l.WarnFunc(format, args...)
	}
}

func (l FuncLogger) Errorf(format string, args ...any) {
	if l.ErrorFunc != nil {
		l.ErrorFunc(format, args...)
	}
}

var (
	mu      sync.RWMutex
	current Logger = NopLogger{}
)

// SetLogger configures diagnostics during application startup. Nil restores
// silent logging. The application owns flushing and closing its logger.
// Pass nil itself, rather than an interface containing a typed nil pointer.
func SetLogger(logger Logger) {
	if logger == nil {
		logger = NopLogger{}
	}
	mu.Lock()
	current = logger
	mu.Unlock()
}

func getLogger() Logger {
	mu.RLock()
	logger := current
	mu.RUnlock()
	return logger
}

func Debugf(format string, args ...any) { getLogger().Debugf(format, args...) }
func Warnf(format string, args ...any)  { getLogger().Warnf(format, args...) }
func Errorf(format string, args ...any) { getLogger().Errorf(format, args...) }
