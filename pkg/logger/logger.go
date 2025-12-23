package logger

import (
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerSettings struct {
	Level     zap.AtomicLevel
	Logger    logr.Logger
	ZapLogger *zap.Logger
}

var logSettings = LoggerSettings{
	Level: zap.NewAtomicLevelAt(zapcore.DebugLevel),
}

func NewStructuredLogger() logr.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logSettings.ZapLogger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		logSettings.Level,
	))

	logSettings.Logger = zapr.NewLogger(logSettings.ZapLogger)
	return logSettings.Logger
}

func NewLogger() logr.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logSettings.ZapLogger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		logSettings.Level,
	))

	logSettings.Logger = zapr.NewLogger(logSettings.ZapLogger)
	return logSettings.Logger
}

func SetLevel(level int8) {
	logSettings.Level.SetLevel(zapcore.Level(-1 * level))
}

func Logger() logr.Logger {
	return logSettings.Logger
}

func LoggerWithName(name string) logr.Logger {
	return logSettings.Logger.WithName(name)
}

func LoggerWithContext(name string) logr.Logger {
	return logSettings.Logger.WithName(name)
}
