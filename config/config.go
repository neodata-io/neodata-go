// neodata-go/config/config.go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Name         string        `mapstructure:"name"`
		Port         int           `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		Env          string        `mapstructure:"env"`
		RateLimit    int           `mapstructure:"rate_limit"`
	} `mapstructure:"app"`

	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		SSLmode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	Auth struct {
		JwtSecret   string `mapstructure:"jwtSecret"`
		TokenExpiry int    `mapstructure:"tokenExpiry"`
	}

	Logger struct {
		LogLevel string `mapstructure:"log_level"`
	}

	Redis struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"redis"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	viper.SetConfigType("yaml")
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("config") // Fallback path
	}

	viper.AutomaticEnv() // Automatically overrides config with environment variables

	// Bind environment variables explicitly for sensitive fields
	//viper.BindEnv("database.password", "DB_PASSWORD")
	//viper.BindEnv("database.user", "DB_USER")
	//viper.BindEnv("database.jwtSecret", "JWT_SECRET")

	// Set default values
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.read_timeout", 10)
	viper.SetDefault("app.write_timeout", 10)
	viper.SetDefault("redis.address", "localhost:6379")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
