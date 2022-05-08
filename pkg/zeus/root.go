package zeus

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var rootLogger *zap.Logger = func() *zap.Logger {
	logger, err := initLogger(
		strings.ToLower(os.Getenv("GO_LOG")),
		strings.ToLower(os.Getenv("GO_LOG_FORMAT")),
	)
	if err != nil {
		panic(err)
	}
	return logger
}()

func initLogger(level string, format string) (*zap.Logger, error) {
	options := []zap.Option{
		zap.AddStacktrace(zap.FatalLevel),
	}

	// Get level
	zapLevel, warn := func() (zapcore.Level, bool) {
		switch level {
		case "debug":
			return zap.DebugLevel, false
		case "warn":
			return zap.WarnLevel, false
		case "error":
			return zap.ErrorLevel, false
		case "fatal":
			return zap.FatalLevel, false
		default:
			return zap.InfoLevel, level != "" && level != "info"
		}
	}()

	// Create logger
	logger, err := func() (*zap.Logger, error) {
		switch format {
		case "json":
			config := zap.NewProductionConfig()
			config.Level.SetLevel(zapLevel)
			return zap.NewProduction(options...)
		default:
			config := zap.NewDevelopmentConfig()
			config.EncoderConfig.CallerKey = ""
			config.Level.SetLevel(zapLevel)
			return config.Build(options...)
		}
	}()
	if err != nil {
		return nil, err
	}

	// Log if unknown configuration values were encountered
	if warn {
		logger.Warn("unknown log level, falling back to 'info'", zap.String("level", level))
	}
	if format != "" && format != "human" && format != "json" {
		logger.Warn("unknown log format, falling back to 'human'", zap.String("format", format))
	}
	return logger, nil
}

// Sync should be deferred at the beginning of a program's main function. It ensures that all logs
// are printed prior to program shutdown.
func Sync() {
	rootLogger.Sync() // nolint:errcheck
}
