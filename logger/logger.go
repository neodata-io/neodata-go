package logger

import (
	"github.com/neodata-io/neodata-go/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a logger based on the log level and environment.
func NewLogger(logLevel zapcore.Level, environment string) (*zap.Logger, error) {
	var config zap.Config

	// Choose between PRD and DEV configurations
	if environment == "PRD" {
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
		EncodeLevel:  zapcore.CapitalLevelEncoder, // Uppercase levels
		EncodeTime:   zapcore.ISO8601TimeEncoder,  // Standard time format
		EncodeCaller: zapcore.ShortCallerEncoder,  // Short file path for caller info
	}

	// Build and return the logger
	return config.Build()
}

// InitServiceLogger creates a base logger and attaches a service-specific field
func InitServiceLogger(cfg *config.AppConfig) (*zap.Logger, error) {

	// Create the logger based on environment and log level
	logger, err := NewLogger(cfg.Logger.LogLevel, cfg.App.Env)
	if err != nil {
		return nil, err
	}

	// Add the service name as a field for every log entry
	return logger.With(zap.String("service", cfg.App.Name)), nil
}