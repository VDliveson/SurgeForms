package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	Client *redis.Client
}

type CacheInterface interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}
