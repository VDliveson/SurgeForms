package cache

import (
	"context"
	"errors"

	"github.com/VDliveson/SurgeForms/go-producer/utils"
	"github.com/go-redis/redis"
)

func ConnectCache(ctx context.Context, cacheType string) (CacheInterface, error) {
	switch cacheType {
	case "redis":
		client := redis.NewClient(&redis.Options{
			Addr:     utils.GetEnv("REDIS_HOST", "localhost:6379"),
			Password: "",
			DB:       0,
		})

		// Test connection
		if _, err := client.Ping().Result(); err != nil {
			return nil, err
		}

		return &RedisCache{Client: client}, nil

	default:
		return nil, errors.New("unsupported cache type")
	}
}
