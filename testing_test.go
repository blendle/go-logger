package logger_test

import (
	"testing"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"

	logger "github.com/blendle/go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestNew(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNew(t)
	logger.Debug("")

	assert.Len(t, logs.All(), 1)
}

func TestTestNewWithLevel(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNewWithLevel(t, zapcore.ErrorLevel)
	logger.Debug("")
	logger.Warn("")
	logger.Error("error")
	logger.DPanic("dpanic")

	require.Len(t, logs.All(), 2)
	assert.Equal(t, zapcore.ErrorLevel, logs.All()[0].Level)
	assert.Equal(t, zapcore.DPanicLevel, logs.All()[1].Level)
}

func TestTestNewWithLevel_WithOptions(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNewWithLevel(t, zapcore.WarnLevel, zap.Fields(zap.Int("1", 1)))
	logger.Error("error")

	require.Len(t, logs.All(), 1)
	assert.Equal(t, int64(1), logs.All()[0].ContextMap()["1"])
}
