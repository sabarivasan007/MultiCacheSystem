// Internal/cache/redis_cache.go
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	log.Print("Redis initialized")
	// Print the defaults
	fmt.Println("Default PoolSize:", client.Options().PoolSize)
	fmt.Println("Default MinIdleConns:", client.Options().MinIdleConns)
	fmt.Println("Default MaxRetries:", client.Options().MaxRetries)
	fmt.Println("Default DialTimeout:", client.Options().DialTimeout)
	fmt.Println("Default ReadTimeout:", client.Options().ReadTimeout)
	fmt.Println("Default WriteTimeout:", client.Options().WriteTimeout)
	return &RedisCache{client: client, ttl: ttl}
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("key does not exist")
		}
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
	//return r.client.Get(context.Background(), key).Result()
}

func (r *RedisCache) GetWithTTL(key string) (interface{}, time.Duration, error) {
	ctx := context.Background()

	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, 0, err
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, 0, err
	}

	var data interface{}
	err = json.Unmarshal([]byte(value), &data)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal value: %v", err)
	}

	return data, ttl, nil
}

// func (r *RedisCache) Set(key string, value interface{}) error {
// 	val, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}
// 	return r.client.Set(context.Background(), key, val, r.ttl).Err()
// }

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {

	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// Use the default TTL if the provided ttl is not provided
	// actualTTL := ttl
	// if ttl <= 0 {
	// 	actualTTL = r.ttl
	// }

	log.Printf("Setting KEY: %s with VALUE: %s and TTL: %v seconds\n", key, value, ttl) // Add this line for logging
	// return r.client.Set(context.Background(), key, val, actualTTL).Err()
	return r.client.Set(context.Background(), key, val, ttl).Err()

}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisCache) ClearAll() error {
	return r.client.FlushDB(context.Background()).Err()
}
