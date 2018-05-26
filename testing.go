package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestNew calls New, but returns both the logger, and an observer that can be
// used to fetch and compare delivered logs.
func TestNew(tb testing.TB, options ...zap.Option) (*zap.Logger, *observer.ObservedLogs) {
	tb.Helper()

	return TestNewWithLevel(tb, zapcore.DebugLevel, options...)
}

// TestNewWithLevel is equal to TestNew, except that it takes an extra argument,
// dictating the minimum log level required to record an entry in the recorder.
func TestNewWithLevel(tb testing.TB, level zapcore.LevelEnabler, options ...zap.Option) (*zap.Logger, *observer.ObservedLogs) {
	tb.Helper()

	core, logs := observer.New(level)
	opt := zap.WrapCore(func(_ zapcore.Core) zapcore.Core { return core })

	zaplog := Must(New("test", "v0.0.1", append([]zap.Option{opt}, options...)...))

	return zaplog, logs
}
