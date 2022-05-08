package zeus

import (
	"context"

	"go.uber.org/zap"
)

// Logger returns the logger that is associated with this context. If no logger has been set on the
// context explicitly, a default logger is returned.
func Logger(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		if logger != nil {
			return logger
		}
	}
	return rootLogger
}

// WithName adds the specified name to the logger that is currently associated with the context.
func WithName(ctx context.Context, name string) context.Context {
	return WithLogger(ctx, Logger(ctx).Named(name))
}

// WithFields adds the specified fields to the logger associated with the context.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return WithLogger(ctx, Logger(ctx).With(fields...))
}

// WithNopLogger sets a logger on the context which does not write anything. All calls to the
// logger functions are simply noops.
func WithNopLogger(ctx context.Context) context.Context {
	return WithLogger(ctx, zap.NewNop())
}

// WithLogger sets the specified logger in the context. Typically, you don't need this function
// and should prefer functions modifying the logger directly (e.g. `WithName`, `WithFields`, ...).
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
