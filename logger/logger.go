package logger

import (
	"fmt"
	"strings"

	"github.com/neodata-io/neodata-go/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	// other logging methods...
}

// NewLogger creates a logger based on the log level and environment.
func NewLogger(logLevel zapcore.Level, environment string) (*zap.Logger, error) {
	var config zap.Config

	// Choose between prd and dev configurations
	if environment == "prd" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Set the log level passed by the microservice
	config.Level = zap.NewAtomicLevelAt(logLevel)

	// Set encoder keys for structured logs (log level, message, caller)
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		EncodeTime:   zapcore.RFC3339TimeEncoder,
		EncodeLevel:  zapcore.CapitalColorLevelEncoder, // Enables color for log levels
		EncodeCaller: zapcore.ShortCallerEncoder,       // Short file path for caller info

	}

	// Build and return the logger
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer l.Sync()

	return l, nil
}

// InitServiceLogger creates a base logger and attaches a service-specific field
func InitServiceLogger(cfg *config.AppConfig) (*zap.Logger, error) {
	logLevel, err := mapLogLevel(cfg.Logger.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to set log level: %w", err)
	}
	// Create the logger based on environment and log level
	logger, err := NewLogger(logLevel, cfg.App.Env)
	if err != nil {
		return nil, err
	}

	// Add the service name as a field for every log entry
	return logger.With(zap.String("service", cfg.App.Name)), nil
}

// mapLogLevel maps a string log level to zapcore.Level
func mapLogLevel(logLevel string) (zapcore.Level, error) {
	switch strings.ToLower(logLevel) {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "dpanic":
		return zapcore.DPanicLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s", logLevel)
	}
}
