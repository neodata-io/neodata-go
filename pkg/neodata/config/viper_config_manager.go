// config/viper_config_manager.go
package config

import "github.com/spf13/viper"

type ViperConfigManager struct {
	config *AppConfig
}

// NewConfigManager loads the configuration and returns an instance of ConfigManager
func NewConfigManager(configPath string) (*ViperConfigManager, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return &ViperConfigManager{config: config}, nil
}

// Validate checks that required configuration values are set
func (c *ViperConfigManager) IsEnabled(configKey string) bool {
	return viper.IsSet(configKey)
}

// AppName returns the application's name
func (c *ViperConfigManager) AppName() string {
	return c.config.App.Name
}

// AppPort returns the application's port
func (c *ViperConfigManager) AppPort() int {
	return c.config.App.Port
}

// GetAppConfig implements the ConfigManager interface
func (c *ViperConfigManager) GetAppConfig() *AppConfig {
	return c.config
}

// GetDatabaseConfig returns the database configuration struct
func (c *ViperConfigManager) GetDatabaseConfig() *DatabaseConfig {
	return &c.config.Database
}

// GetAuthConfig returns the entire AuthConfig struct
func (c *ViperConfigManager) GetAuthConfig() *AuthConfig {
	return &c.config.Auth
}

// GetMessagingConfig returns the entire MessagingConfig struct
func (c *ViperConfigManager) GetMessagingConfig() *MessagingConfig {
	return &c.config.Messaging
}

// GetLoggerConfig returns the entire LoggerConfig struct
func (c *ViperConfigManager) GetLoggerConfig() *LoggerConfig {
	return &c.config.Logger
}

// GetRedisConfig returns the entire RedisConfig struct
func (c *ViperConfigManager) GetRedisConfig() *RedisConfig {
	return &c.config.Redis
}
