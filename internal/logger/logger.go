// Package logger ...
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is no-operation logger by default
// Init in Init() function
var Log *zap.Logger = zap.NewNop()

// Init and configure logger
func Init(level string) error {
	// set log level
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl
	return nil
}
