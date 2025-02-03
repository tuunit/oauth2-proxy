package structured

import (
	"io"
	"os"

	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

type LogType string

const (
	ConsoleLogger LogType = "console"
	JSONLogger    LogType = "json"
)

// Global base logger
var baseLogger zerologr.Logger

type levelWriter struct {
	stdout io.Writer
	stderr io.Writer
}

func (lw levelWriter) Write(p []byte) (n int, err error) {
	return lw.stdout.Write(p) // Default all logs to stdout
}

func (lw levelWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < zerolog.ErrorLevel {
		return lw.stdout.Write(p) // Logs below ERROR go to stdout
	}
	return lw.stderr.Write(p) // ERROR and above go to stderr
}

// func InitializeBaseLogger(format LogType) {
func init() {
	// TODO: make configurable
	format := JSONLogger
	// Unix by default because it is the fastest
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var stdoutWriter, stderrWriter io.Writer

	if format == ConsoleLogger {
		stdoutWriter = zerolog.ConsoleWriter{Out: os.Stdout}
		stderrWriter = zerolog.ConsoleWriter{Out: os.Stderr}
	} else if format == JSONLogger {
		stdoutWriter = os.Stdout
		stderrWriter = os.Stderr
	}

	// TODO: allow for
	// 1. writing to file
	// 2. writing to stdout or stderr only
	// 3. writing to file and stdout/err
	multiWriter := levelWriter{
		stdout: stdoutWriter,
		stderr: stderrWriter,
	}
	base := zerolog.New(multiWriter).With().Timestamp().Caller().Logger()
	baseLogger = zerologr.New(&base)
}

// GetLogger returns a child logger with the given name
// Usually the package name or struct name
func GetLogger(name string) zerologr.Logger {
	return baseLogger.WithName(name)
}
