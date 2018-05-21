package logger_test

import (
	"testing"

	logger "github.com/blendle/go-logger"
	"github.com/blendle/go-logger/stackdriver"
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

func TestMust(t *testing.T) {
	t.Parallel()

	assert.IsType(t, &zap.Logger{}, logger.Must(zap.NewProduction()))
}
