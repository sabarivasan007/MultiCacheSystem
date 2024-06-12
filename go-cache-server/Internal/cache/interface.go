package cache

import "time"

type CacheSystem interface {
	Get(key string) (interface{}, error)
	GetWithTTL(key string) (interface{}, time.Duration, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	ClearAll() error
}
