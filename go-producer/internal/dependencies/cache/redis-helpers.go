package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// Set stores a key-value pair in Redis with an expiration time
func (c *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := c.Client.Set(key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %v", key, err)
	}
	return nil
}

// Get retrieves a value from Redis by key
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.Client.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key %s not found", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key %s: %v", key, err)
	}
	return value, nil
}

// Delete removes a key from Redis
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.Client.Del(key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %v", key, err)
	}
	return nil
}

// Exists checks if a key exists in Redis
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.Client.Exists(key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %v", key, err)
	}
	return exists > 0, nil
}

// Close shuts down the Redis connection
func (c *RedisCache) Close() error {
	return c.Client.Close()
}
