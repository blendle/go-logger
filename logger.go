package logger

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a new logger, ready to use in our services.
func New(service, version string, options ...zap.Option) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	config := &zap.Config{
		Level:            level,
		Encoding:         "json",
		EncoderConfig:    zapdriver.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if env, ok := os.LookupEnv("ENV"); ok && env != "production" {
		config.Development = true
	}

	if v := os.Getenv("DEBUG"); v != "" && v != "0" && v != "false" {
		level.SetLevel(zap.DebugLevel)
	}

	stackcore := zapdriver.WrapCore()

	fields := zap.Fields(
		zapdriver.Label("service", service),
		zapdriver.Label("version", version),
	)

	go levelToggler(level)

	return config.Build(append(options, stackcore, fields)...)
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

func levelToggler(level zap.AtomicLevel) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)

	for {
		<-ch

		if level.Level() == zap.DebugLevel {
			level.SetLevel(zap.InfoLevel)
		} else {
			level.SetLevel(zap.DebugLevel)
		}
	}
}
