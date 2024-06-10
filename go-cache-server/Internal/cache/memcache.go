package cache

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type MemCache struct {
	client *memcache.Client
	ttl    int32
}

func NewMemCache(server string, ttl int32) *MemCache {
	client := memcache.New(server)
	return &MemCache{client: client, ttl: ttl}
}

func (m *MemCache) Get(key string) (string, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

func (m *MemCache) Set(key string, value string) error {
	return m.client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: m.ttl})
}

func (m *MemCache) SetWithTTL(key string, value string, expiration int) error {
	return m.client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expiration)})
}

func (m *MemCache) Delete(key string) error {
	return m.client.Delete(key)
}

func (m *MemCache) ClearAll() error {
	return m.client.FlushAll()
}
