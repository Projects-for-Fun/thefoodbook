package logging

import (
	"context"

	"github.com/rs/zerolog"
)

// LoggerKey The key used in context for logger
var LoggerKey = "LoggerKey"

func AttachLogger(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

func GetLogger(ctx context.Context) zerolog.Logger {
	logger, ok := ctx.Value(LoggerKey).(zerolog.Logger)

	if !ok {
		return zerolog.Logger{}
	}

	return logger
}
