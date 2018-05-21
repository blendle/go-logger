package logger

import (
	"github.com/blendle/go-logger/stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a new logger, ready to use in our services.
func New(service, version string, zaplog *zap.Logger) *zap.Logger {
	stackcore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &stackdriver.Core{Core: core}
	})

	fields := zap.Fields(zap.String("service", service), zap.String("version", version))

	return zaplog.WithOptions(stackcore, fields)
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
