// Internal/cache/redis_cache.go
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{client: client, ttl: ttl}
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *RedisCache) GetWithTTL(key string) (string, int, error) {
	ctx := context.Background()
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", 0, err
	}
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return "", 0, err
	}
	return value, int(ttl.Seconds()), nil
}

func (r *RedisCache) Set(key string, value string) error {
	return r.client.Set(context.Background(), key, value, r.ttl).Err()
}

func (r *RedisCache) SetWithTTL(key string, value string, expiration int) error {
	expireTime := time.Duration(expiration) * time.Second
	fmt.Printf("Setting key: %s with value: %s and TTL: %d seconds\n", key, value, expiration) // Add this line for logging
	return r.client.Set(context.Background(), key, value, expireTime).Err()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisCache) ClearAll() error {
	return r.client.FlushDB(context.Background()).Err()
}
