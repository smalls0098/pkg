package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// DefaultMessageKey default message key.
var DefaultMessageKey = "msg"

// Option is Helper option.
type Option func(*Logger)

// WithMessageKey with message key.
func WithMessageKey(k string) Option {
	return func(opts *Logger) {
		opts.msgKey = k
	}
}

type Logger struct {
	log    *zap.Logger
	msgKey string
}

func NewLogger(zLog *zap.Logger, opts ...Option) *Logger {
	options := &Logger{
		msgKey: DefaultMessageKey,
		log:    zLog,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

func (l *Logger) Log(level Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	msg := ""
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		k := fmt.Sprint(keyvals[i])
		v := fmt.Sprint(keyvals[i+1])
		if k == l.msgKey {
			msg = v
		} else {
			data = append(data, zap.Any(k, v))
		}
	}

	switch level {
	case LevelDebug:
		l.log.Debug(msg, data...)
	case LevelInfo:
		l.log.Info(msg, data...)
	case LevelWarn:
		l.log.Warn(msg, data...)
	case LevelError:
		l.log.Error(msg, data...)
	case LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

// Debug logs a message at debug level.
func (l *Logger) Debug(a ...interface{}) {
	_ = l.Log(LevelDebug, l.msgKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (l *Logger) Debugf(format string, a ...interface{}) {
	_ = l.Log(LevelDebug, l.msgKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (l *Logger) Debugw(keyvals ...interface{}) {
	_ = l.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (l *Logger) Info(a ...interface{}) {
	_ = l.Log(LevelInfo, l.msgKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (l *Logger) Infof(format string, a ...interface{}) {
	_ = l.Log(LevelInfo, l.msgKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (l *Logger) Infow(keyvals ...interface{}) {
	_ = l.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (l *Logger) Warn(a ...interface{}) {
	_ = l.Log(LevelWarn, l.msgKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (l *Logger) Warnf(format string, a ...interface{}) {
	_ = l.Log(LevelWarn, l.msgKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (l *Logger) Warnw(keyvals ...interface{}) {
	_ = l.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (l *Logger) Error(a ...interface{}) {
	_ = l.Log(LevelError, l.msgKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (l *Logger) Errorf(format string, a ...interface{}) {
	_ = l.Log(LevelError, l.msgKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (l *Logger) Errorw(keyvals ...interface{}) {
	_ = l.Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func (l *Logger) Fatal(a ...interface{}) {
	_ = l.Log(LevelFatal, l.msgKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (l *Logger) Fatalf(format string, a ...interface{}) {
	_ = l.Log(LevelFatal, l.msgKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func (l *Logger) Fatalw(keyvals ...interface{}) {
	_ = l.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
