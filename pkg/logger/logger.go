package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

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
