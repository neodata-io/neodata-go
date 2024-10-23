// config/viper_config_manager.go
package config

type ViperConfigManager struct {
	config *AppConfig
}

// NewViperConfigManager loads the configuration and returns an instance of ViperConfigManager
func NewViperConfigManager(configPath string) (*ViperConfigManager, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return &ViperConfigManager{config: config}, nil
}

// GetAppConfig implements the ConfigManager interface
func (v *ViperConfigManager) GetAppConfig() *AppConfig {
	return v.config
}
