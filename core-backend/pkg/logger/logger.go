package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Options holds configuration for the JSON slog logger.
type Options struct {
	Level       string
	ServiceName string
	Environment string
}

// New creates a structured slog JSON logger writing to stdout.
func New(opts Options) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(opts.Level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handlerOpts := &slog.HandlerOptions{
		Level: lvl,
	}

	handler := slog.NewJSONHandler(os.Stdout, handlerOpts)
	logger := slog.New(handler)

	attrs := make([]any, 0, 4)
	if opts.ServiceName != "" {
		attrs = append(attrs, "service", opts.ServiceName)
	}
	if opts.Environment != "" {
		attrs = append(attrs, "env", opts.Environment)
	}

	if len(attrs) > 0 {
		logger = logger.With(attrs...)
	}

	return logger
}
