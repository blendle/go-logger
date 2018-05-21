package logger

import (
	"os"
	"testing"

	"github.com/blendle/go-logger/stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// New returns a new logger, ready to use in our services.
func New(service, version string, options ...zap.Option) *zap.Logger {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	config := &zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    stackdriver.EncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if env, ok := os.LookupEnv("ENV"); ok && env != "production" {
		config.Development = true
	}

	if _, ok := os.LookupEnv("DEBUG"); ok {
		level.SetLevel(zap.DebugLevel)
	}

	stackcore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &stackdriver.Core{Core: core}
	})

	fields := zap.Fields(stackdriver.LogServiceContext(&stackdriver.ServiceContext{
		Service: service,
		Version: version,
	}))

	return Must(config.Build(append(options, stackcore, fields)...))
}

// Must is a convenience function that takes a zaplog and error as input, panics
// if the error is not nil, and returns the passed in logger.
//
// This can be used for example with `Must(zap.NewProduction())`
func Must(zaplog *zap.Logger, err error) *zap.Logger {
	if err != nil {
		panic(err)
	}

	return zaplog
}

// TestNew calls New, but returns both the logger, and an observer that can be
// used to fetch and compare delivered logs.
func TestNew(tb testing.TB, options ...zap.Option) (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.DebugLevel)
	opt := zap.WrapCore(func(_ zapcore.Core) zapcore.Core { return core })

	zaplog := New("test", "v0.0.1", append(options, opt)...)

	return zaplog, logs
}
