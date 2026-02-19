package logger

import (
	"go.uber.org/zap"
)

// Log is the global application logger
var Log *zap.Logger = zap.NewNop()

// Initialize configures the global logger with the given log level
func Initialize(level string) error {
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	config := zap.NewProductionConfig()
	config.Level = logLevel
	zapLogger, err := config.Build()
	if err != nil {
		return err
	}

	Log = zapLogger
	return nil
}
