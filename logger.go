package logger

import (
	"go.uber.org/zap"
)

// New returns a new logger, ready to use in our services.
func New(service, version string, zaplog *zap.Logger) *zap.Logger {
	fields := zap.Fields(zap.String("service", service), zap.String("version", version))

	return zaplog.WithOptions(fields)
}
