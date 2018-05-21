package logger_test

import (
	"testing"

	logger "github.com/blendle/go-logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNew(t *testing.T) {
	t.Parallel()

	assert.IsType(t, &zap.Logger{}, logger.New("", ""))
}

func TestNew_Stackdriver(t *testing.T) {
	t.Parallel()

	core, logs := observer.New(zapcore.WarnLevel)
	opt := zap.WrapCore(func(_ zapcore.Core) zapcore.Core { return core })

	logger := logger.New("", "", opt)
	logger.Warn("")

	require.Len(t, logs.All(), 1)
	assert.NotNil(t, logs.All()[0].ContextMap()["context"])
}

func TestMust(t *testing.T) {
	t.Parallel()

	assert.IsType(t, &zap.Logger{}, logger.Must(zap.NewProduction()))
}
