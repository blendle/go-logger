package logger

import (
	"github.com/blendle/go-logger/stackdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	stackcore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &stackdriver.Core{Core: core}
	})

	fields := zap.Fields(zap.String("service", service), zap.String("version", version))

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
