package logger

import (
	"go.uber.org/zap"
)

// New returns a new logger, ready to use in our services.
func New(service, version string, options ...func(zap.Config)) *zap.Logger {
	config := zap.NewProductionConfig()
	fields := zap.Fields(zap.String("service", service), zap.String("version", version))

	logger, err := config.Build(fields)
	if err != nil {
		panic(err)
	}

	return logger
}
