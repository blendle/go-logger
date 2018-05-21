package logger_test

import (
	"testing"

	logger "github.com/blendle/go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Parallel()

	assert.IsType(t, &zap.Logger{}, logger.New("", ""))
}

func TestTestNew(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNew(t)
	logger.Debug("")

	assert.Len(t, logs.All(), 1)
}

func TestLogger_Stackdriver_Context(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNew(t)
	logger.Warn("")

	require.Len(t, logs.All(), 1)
	assert.NotNil(t, logs.All()[0].ContextMap()["context"])
}

func TestMust(t *testing.T) {
	t.Parallel()

	assert.IsType(t, &zap.Logger{}, logger.Must(zap.NewProduction()))
}
