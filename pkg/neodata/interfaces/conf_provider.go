package interfaces

import "github.com/neodata-io/neodata-go/pkg/neodata/config"

type ConfigProvider interface {
	AppConfig() *config.AppConfig
	DatabaseConfig() *config.DatabaseConfig
	AuthConfig() *config.AuthConfig
	MessagingConfig() *config.MessagingConfig
	LoggerConfig() *config.LoggerConfig
	RedisConfig() *config.RedisConfig

	AppName() string
	AppPort() int
	IsEnabled(key string) bool
	// Other frequently used configuration methods
}
