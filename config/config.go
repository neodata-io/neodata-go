// config/config.go
package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	AppName      string `mapstructure:"app_name"`
	AppPort      string `mapstructure:"app_port"`
	DatabaseDSN  string `mapstructure:"database_dsn"`
	RedisAddress string `mapstructure:"redis_address"`
}

// LoadConfig loads configuration from a config file or environment variables.
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // Automatically overrides config with environment variables

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Printf("Configuration loaded successfully: %s", config.AppName)
	return &config, nil
}
