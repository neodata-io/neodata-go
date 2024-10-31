// config/redis_config.go
package config

// RedisConfig holds Redis-specific configuration fields
type RedisConfig struct {
	Address string `mapstructure:"address"`
}
