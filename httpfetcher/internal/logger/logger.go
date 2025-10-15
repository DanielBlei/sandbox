package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	LoggerKey contextKey = "logger"
)

// WithLogger - adds a logger to the context
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// FromContext - extracts logger from context with fallback
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		return logger
	}
	// Fallback to a no-op logger if none found
	return zap.NewNop()
}

// Init creates a new logger
// Not applygin production for this project, but keeping it for reference or future use
func Init(debug bool) *zap.Logger {
	if debug {
		return zap.Must(zap.NewDevelopment())
	}
	return zap.Must(zap.NewProduction())
}
