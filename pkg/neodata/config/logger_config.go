// config/logger_config.go
package config

// LoggerConfig holds logger-specific configuration fields
type LoggerConfig struct {
	LogLevel string `mapstructure:"log_level"`
}
