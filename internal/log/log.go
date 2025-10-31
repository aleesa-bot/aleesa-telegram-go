// Package log wraps log/slog for more easier and convinient use.
package log

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

var (
	// Handler exposes slog handler struct.
	Handler slog.Handler

	// Ctx log context. For now it is just sits here and do nothing.
	Ctx = context.Background()

	// Writer go-writer eported for capturing other types of logs, even for std log go facility.
	Writer io.Writer

	// reader is required for Writer to work properly. It is connected via io.Pipe() to Writer.
	reader io.Reader
)

// Init setup logger stuff.
// level can be error, warn, info, debug if something other supplied info level selected.
// fileDescriptor should be opened before supplying it to Init().
func Init(level string, fileDescriptor *os.File) {
	var loglevel slog.Level

	switch level {
	case "debug":
		loglevel = slog.LevelDebug

	case "info":
		loglevel = slog.LevelInfo

	case "warn":
		loglevel = slog.LevelWarn

	case "error":
		loglevel = slog.LevelError

	default:
		loglevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{ //nolint: exhaustruct
		// Use the ReplaceAttr function on the handler options to be able to replace any single attribute in the log
		// output.
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr { //nolint: revive
			// Check that we are handling the time key.
			if a.Key != slog.TimeKey {
				return a
			}

			t := a.Value.Time()

			// Change the value from a time.Time to the String. The string has correct time format.
			a.Value = slog.StringValue(t.Format(time.DateTime))

			return a
		},

		Level: loglevel,
	}

	Handler = slog.NewTextHandler(fileDescriptor, opts)

	slog.SetDefault(
		slog.New(Handler),
	)

	reader, Writer = io.Pipe()

	go readIoPipe()
}

// readIoPipe reads pipe and writes it to slog.
func readIoPipe() {
	scanner := bufio.NewScanner(reader)

	for {
		for scanner.Scan() {
			line := scanner.Text() // Get the current line as a string.
			slog.Debug(line)
		}

		if err := scanner.Err(); err != nil {
			// io.EOF is expected when done reading, e.g. on exit.
			if errors.Is(err, io.EOF) {
				return
			}

			// something weird happen, we should it log to guess later what it was.
			Errorf("Error during scanning log writer pipe: %v\n", err)
		}
	}
}

// GetLevel returns current loglevel.
func GetLevel() string {
	switch {
	case Handler.Enabled(Ctx, slog.LevelDebug):
		return "debug"

	case Handler.Enabled(Ctx, slog.LevelInfo):
		return "info"

	case Handler.Enabled(Ctx, slog.LevelWarn):
		return "warn"

	case Handler.Enabled(Ctx, slog.LevelError):
		return "error"

	default:
		return "unknown"
	}
}

// Fatal passes message directly to slog.Error() and exits with error code 1.
func Fatal(message string) {
	slog.Error(message)
	os.Exit(1)
}

// Fatalf logs message on error log level and exits with error code 1. It allows to supply arguments in printf() style .
func Fatalf(format string, a ...any) {
	slog.Error(fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Error passes message directly to slog.Error().
func Error(message string) {
	slog.Error(message)
}

// Errorf logs message on error log level and allows to supply arguments in printf() style.
func Errorf(format string, a ...any) {
	slog.Error(fmt.Sprintf(format, a...))
}

// Warn passes message directly to slog.Warn().
func Warn(message string) {
	slog.Warn(message)
}

// Warnf logs message on warn log level and allows to supply arguments in printf() style.
func Warnf(format string, a ...any) {
	slog.Warn(fmt.Sprintf(format, a...))
}

// Info passes message directly to slog.Info().
func Info(message string) {
	slog.Info(message)
}

// Infof logs message on info log level and allows to supply arguments in printf() style.
func Infof(format string, a ...any) {
	slog.Info(fmt.Sprintf(format, a...))
}

// Debug passes message directly to slog.Debug().
func Debug(message string) {
	slog.Debug(message)
}

// Debugf logs message on debug log level and allows to supply arguments in printf() style.
func Debugf(format string, a ...any) {
	slog.Debug(fmt.Sprintf(format, a...))
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
