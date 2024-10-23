// config/viper_config_manager.go
package config

// ConfigManager defines an interface for fetching configuration values
type ConfigManager interface {
	GetAppConfig() *AppConfig
}

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

// GetAppConfig implements the ConfigManager interface
func (c *ViperConfigManager) GetAppConfig() *AppConfig {
	return c.config
}
