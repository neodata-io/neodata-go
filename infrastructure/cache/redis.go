// infrastructure/cache/redis_cache.go
package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{ /* config */ })
	return &RedisCache{client: client}
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

func (c *RedisCache) Set(key string, value string) error {
	return c.client.Set(context.Background(), key, value, 0).Err()
}
