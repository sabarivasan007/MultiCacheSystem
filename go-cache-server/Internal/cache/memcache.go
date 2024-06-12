package cache

import (
	"encoding/json"
	"time"

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

func (m *MemCache) Get(key string) (interface{}, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}
	var data interface{}
	err = json.Unmarshal(item.Value, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetWithTTL retrieves the value and its TTL for the specified key
// func (m *MemCache) GetWithTTL(key string) (interface{}, time.Duration, error) {
// 	item, err := m.client.Get(key)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	var data interface{}
// 	if err = json.Unmarshal(item.Value, &data); err != nil {
// 		return nil, 0, fmt.Errorf("failed to unmarshal value: %v", err)
// 	}

// 	ttlKey := fmt.Sprintf("%s_ttl", key)
// 	ttlItem, err := m.client.Get(ttlKey)
// 	if err != nil {
// 		return data, 0, nil
// 	}

// 	var ttlSeconds int
// 	if err = json.Unmarshal(ttlItem.Value, &ttlSeconds); err != nil {
// 		return data, 0, fmt.Errorf("failed to unmarshal TTL: %v", err)
// 	}

// 	ttl := time.Duration(ttlSeconds) * time.Second
// 	return data, ttl, nil
// }

func (m *MemCache) Set(key string, value interface{}, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(&memcache.Item{Key: key, Value: val, Expiration: int32(ttl.Seconds())})
}

func (m *MemCache) Delete(key string) error {
	return m.client.Delete(key)
}

func (m *MemCache) ClearAll() error {
	return m.client.FlushAll()
}
