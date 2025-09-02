package redis

import (
	"context"
	"sync"
	"time"

	"github.com/Arpitmovers/reviewservice/internal/config"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	once          sync.Once
	redisInstance *RedisCache
	redisInitErr  error
)

type RedisCache struct {
	client *redis.Client
}

func GetRedisClient(cfg *config.Config) (*RedisCache, error) {

	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
			Password: "", // set if needed
			DB:       0,
		})

		if pingErr := rdb.Ping(context.Background()).Err(); pingErr != nil {
			logger.Logger.Error("error connecting to redis", zap.Error(pingErr))
			redisInitErr = pingErr
			return
		}

		logger.Logger.Info("connected to redis")
		redisInstance = &RedisCache{client: rdb}
	})

	if redisInitErr != nil {
		return nil, redisInitErr
	}
	return redisInstance, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// SetNX sets key only if it does not exist (atomic lock)
func (r *RedisCache) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, ttl).Result()
}
