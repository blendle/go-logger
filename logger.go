package logger

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config of the current logger instance
type Config struct {
	App         string
	Tier        string
	Production  bool
	Version     string
	Environment string
}

// L is the global logger instance.
var L = zap.NewNop()

// LogLevel is the current log level of the logger.
var LogLevel = zap.NewAtomicLevel()

// Init initializes the logger
func Init(config *Config, options ...func(zap.Config)) {
	var err error

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	zapconfig := zap.Config{
		Level:       LogLevel,
		Development: !config.Production,
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
	}

	for _, option := range options {
		option(zapconfig)
	}

	if os.Getenv("DEBUG") == "true" {
		zapconfig.Level.SetLevel(zap.DebugLevel)
	}

	go logLevelToggler()

	app := zap.String("app", config.App)
	tier := zap.String("tier", config.Tier)
	production := zap.Bool("production", config.Production)
	version := zap.String("version", config.Version)
	environment := zap.String("environment", config.Environment)

	L, err = zapconfig.Build(zap.Fields(app, tier, production, version, environment))
	if err != nil {
		log.Printf(`{"severity": "fatal", "msg": "%v"}`, err)
		os.Exit(1)
	}
}

func logLevelToggler() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGUSR1)

	for {
		<-ch

		if LogLevel.Level() == zap.DebugLevel {
			LogLevel.SetLevel(zap.InfoLevel)
		} else {
			LogLevel.SetLevel(zap.DebugLevel)
		}
	}
}
