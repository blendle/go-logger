package logger_test

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	logger "github.com/blendle/go-logger"
	"github.com/blendle/go-logger/stackdriver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	logger.Warn("", stackdriver.LogUser("test"))

	require.Len(t, logs.All(), 1)
	assert.NotNil(t, logs.All()[0].ContextMap()["context"])
}

func TestLogger_Stackdriver_ServiceContext(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNew(t)
	logger.Warn("")

	require.Len(t, logs.All(), 1)
	assert.NotNil(t, logs.All()[0].ContextMap()["serviceContext"])
}

func TestLogger_Stackdriver_Labels(t *testing.T) {
	t.Parallel()

	want := map[string]interface{}{"foo": "bar", "baz": "qux"}

	logger, logs := logger.TestNew(t)
	logger.Warn("", stackdriver.LogLabels("foo", "bar", "baz", "qux"))

	require.Len(t, logs.All(), 1)
	assert.EqualValues(t, want, logs.All()[0].ContextMap()["labels"])
}

func TestLogger_Development_NonProduction(t *testing.T) {
	_ = os.Setenv("ENV", "non-production")
	defer func() { _ = os.Unsetenv("ENV") }()

	logger, logs := logger.TestNew(t)
	fn := func() { logger.DPanic("") }

	// because we've set the environment to anything other than "production", any
	// `DPanic` log call will cause a panic.
	assert.Panics(t, fn)
	assert.Len(t, logs.All(), 1)
}

func TestLogger_Development_Production(t *testing.T) {
	_ = os.Setenv("ENV", "production")
	defer func() { _ = os.Unsetenv("ENV") }()

	logger, logs := logger.TestNew(t)
	fn := func() { logger.DPanic("") }

	// because we've set the environment to "production", any `DPanic` log call
	// will not cause a panic.
	assert.NotPanics(t, fn)
	assert.Len(t, logs.All(), 1)
}

func TestLogger_Development_Unset(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNew(t)
	fn := func() { logger.DPanic("") }

	// because we haven't set the "ENV" environment variable, any `DPanic` log
	// call will not cause a panic.
	assert.NotPanics(t, fn)
	assert.Len(t, logs.All(), 1)
}

func TestLogger_Debug_Enabled(t *testing.T) {
	_ = os.Setenv("DEBUG", "1")
	defer func() { _ = os.Unsetenv("DEBUG") }()

	logger := logger.New("", "")

	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))
}

func TestLogger_Debug_Disabled(t *testing.T) {
	t.Parallel()

	logger := logger.New("", "")

	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))
	assert.True(t, logger.Core().Enabled(zapcore.InfoLevel))
}

func TestLogger_LevelToggler(t *testing.T) {
	t.Parallel()

	logger := logger.New("", "")
	time.Sleep(1 * time.Millisecond)

	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))

	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(1 * time.Millisecond)

	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))

	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(1 * time.Millisecond)

	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))
}

func TestMust(t *testing.T) {
	t.Parallel()

	assert.IsType(t, &zap.Logger{}, logger.Must(zap.NewProduction()))
}

func TestMust_Panic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() { logger.Must(zap.NewNop(), errors.New("panic")) })
}
