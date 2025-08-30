package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once          sync.Once
	redisInstance *RedisCache
)

type RedisCache struct {
	client *redis.Client
}

func GetRedisClient() *RedisCache {

	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Sprintf("failed to connect to redis: %v", err))
		}
		fmt.Println("connected to redis ")
		redisInstance = &RedisCache{client: rdb}

	})

	return redisInstance
}

func (r *RedisCache) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}
