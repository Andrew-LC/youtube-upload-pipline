package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap.Logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new configured logger
func NewLogger(serviceName string, debug bool) (*Logger, error) {
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	config.InitialFields = map[string]interface{}{
		"service": serviceName,
	}

	// Output to stdout
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	zapLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{Logger: zapLogger}, nil
}

// GetZapLogger returns the underlying zap.Logger
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.Logger
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// DefaultLogger returns a basic logger for immediate use before configuration
func DefaultLogger() *Logger {
	config := zap.NewProductionConfig()
	logger, _ := config.Build()
	return &Logger{Logger: logger}
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1).
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.Logger.Fatal(msg, fields...)
	os.Exit(1)
}
