// neodata-go/config/config.go
package config

import (
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App       AppConfigDetails `mapstructure:"app"`
	Database  DatabaseConfig   `mapstructure:"database,omitempty"`
	Auth      AuthConfig       `mapstructure:"auth,omitempty"`
	Messaging MessagingConfig  `mapstructure:"messaging,omitempty"`
	Logger    LoggerConfig     `mapstructure:"logger,omitempty"`
	Redis     RedisConfig      `mapstructure:"redis,omitempty"`
}

type AuthConfig struct {
	JwtSecret   string               `mapstructure:"jwtSecret"`
	TokenExpiry time.Duration        `mapstructure:"tokenExpiry"`
	Policy      *PolicyManagerConfig `yaml:"policy_manager,omitempty"` // PolicyManager is optional
}

// AppConfigDetails holds specific app settings
type AppConfigDetails struct {
	Name           string        `mapstructure:"name"`
	Port           int           `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	Env            string        `mapstructure:"env"`
	RateLimit      int           `mapstructure:"rate_limit"`
	Secret         string        `mapstructure:"secret"`
	UserServiceURL string        `mapstructure:"user_service_url"`
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
	//viper.BindEnv("database.secret", "SECRET")

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
