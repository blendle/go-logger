package logger_test

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	logger "github.com/blendle/go-logger"
	"github.com/blendle/zapdriver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	t.Parallel()

	logger, err := logger.New("", "")
	require.NoError(t, err)

	assert.IsType(t, &zap.Logger{}, logger)
}

func TestLogger_Stackdriver_Labels(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNew(t)
	logger.Warn("", zapdriver.Label("hello", "world"))

	require.Len(t, logs.All(), 1)
	assert.NotNil(t, logs.All()[0].ContextMap()["labels"])
}

func TestLogger_DevelopmentEnv(t *testing.T) {
	_ = os.Setenv("ENV", "development")
	defer func() { _ = os.Unsetenv("ENV") }()

	logger, logs := logger.TestNew(t)
	fn := func() { logger.DPanic("") }

	// because we've set the environment to "development", any `DPanic` log call
	// will cause a panic.
	assert.Panics(t, fn)
	assert.Len(t, logs.All(), 1)
}

func TestLogger_DevelopmentWrongEnv(t *testing.T) {
	_ = os.Setenv("ENV", "not-development")
	defer func() { _ = os.Unsetenv("ENV") }()

	logger, logs := logger.TestNew(t)
	fn := func() { logger.DPanic("") }

	// because we've set the environment to anything other than "development", any
	// `DPanic` log call will not cause a panic.
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

func TestLogger_Development_Option(t *testing.T) {
	t.Parallel()

	logger, logs := logger.TestNew(t, zap.Development())
	fn := func() { logger.DPanic("") }

	// because we've passed in `zap.Development()`, any `DPanic` log call will
	// cause a panic.
	assert.Panics(t, fn)
	assert.Len(t, logs.All(), 1)
}

func TestLogger_Debug_Enabled(t *testing.T) {
	_ = os.Setenv("DEBUG", "1")
	defer func() { _ = os.Unsetenv("DEBUG") }()

	logger := logger.Must(logger.New("", ""))

	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))
}

func TestLogger_Debug_Explicitly_Disabled(t *testing.T) {
	_ = os.Setenv("DEBUG", "0")
	defer func() { _ = os.Unsetenv("DEBUG") }()

	logger := logger.Must(logger.New("", ""))

	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))
	assert.True(t, logger.Core().Enabled(zapcore.InfoLevel))
}

func TestLogger_Debug_Disabled(t *testing.T) {
	t.Parallel()

	logger := logger.Must(logger.New("", ""))

	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))
	assert.True(t, logger.Core().Enabled(zapcore.InfoLevel))
}

func TestLogger_LevelToggler(t *testing.T) {
	t.Parallel()

	logger := logger.Must(logger.New("", ""))
	time.Sleep(10 * time.Millisecond)

	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(10 * time.Millisecond)

	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))

	_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	time.Sleep(10 * time.Millisecond)

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
