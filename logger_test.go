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

	zaplog, err := zap.NewProduction()
	require.NoError(t, err)

	assert.IsType(t, &zap.Logger{}, logger.New("", "", zaplog))
}

func TestMust(t *testing.T) {
	t.Parallel()

	assert.IsType(t, &zap.Logger{}, logger.Must(zap.NewProduction()))
}
