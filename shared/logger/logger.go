package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init initializes the global logger
func Init(dev bool) {
	var err error
	if dev {
		// Development: human-readable console logs
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		Log, err = cfg.Build()
	} else {
		// Production: structured JSON logs
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		Log, err = cfg.Build()
	}
	if err != nil {
		panic(err)
	}

	// Include file:line info and stacktraces for errors
	Log = Log.WithOptions(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

// Sync flushes any buffered log entries
func Sync() {
	_ = Log.Sync()
}

// Sugar provides the sugared logger for printf-style logging
func Sugar() *zap.SugaredLogger {
	return Log.Sugar()
}
