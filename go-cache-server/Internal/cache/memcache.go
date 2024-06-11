package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemCache struct {
	client *memcache.Client
	ttl    int32
}

// type CacheItem struct {
// 	Value     string `json:"value"`
// 	Timestamp int64  `json:"timestamp"`
// 	TTL       int32  `json:"ttl"`
// }

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

// func (m *MemCache) Get(key string) (string, error) {
// 	item, err := m.client.Get(key)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(item.Value), nil
// }

// GetWithTTL retrieves the value and its TTL for the specified key
func (m *MemCache) GetWithTTL(key string) (interface{}, time.Duration, error) {
	// First, get the actual value
	item, err := m.client.Get(key)
	if err != nil {
		return nil, 0, err
	}

	// Unmarshal the value part
	var data interface{}
	if err = json.Unmarshal(item.Value, &data); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal value: %v", err)
	}

	// Attempt to retrieve the TTL from a specially formed TTL key
	ttlKey := fmt.Sprintf("%s_ttl", key)
	ttlItem, err := m.client.Get(ttlKey)
	if err != nil {
		return data, 0, nil // If TTL key is not found, return the data without TTL
	}

	// Assuming TTL was stored as an integer number of seconds
	var ttlSeconds int
	if err = json.Unmarshal(ttlItem.Value, &ttlSeconds); err != nil {
		return data, 0, fmt.Errorf("failed to unmarshal TTL: %v", err)
	}

	// Convert seconds to time.Duration
	ttl := time.Duration(ttlSeconds) * time.Second
	return data, ttl, nil
}

// func (m *MemCache) GetWithTTL(key string) (interface{}, int, error) {
// 	item, err := m.client.Get(key)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	log.Printf("Raw value from cache: %s", string(item.Value)) // Add this line for debugging

// 	var cacheItem CacheItem
// 	if err := json.Unmarshal(item.Value, &cacheItem); err != nil {
// 		return "", 0, err
// 	}
// 	ttl := cacheItem.TTL - int32(time.Now().Unix()-cacheItem.Timestamp)
// 	if ttl < 0 {
// 		return "", 0, errors.New("item has expired")
// 	}
// 	return cacheItem.Value, int(ttl), nil
// }

// func (m *MemCache) Set(key string, value interface{}) error {
// 	val, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}
// 	return m.client.Set(&memcache.Item{Key: key, Value: val, Expiration: m.ttl})
// 	//return m.client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: m.ttl})
// }

func (m *MemCache) Set(key string, value interface{}, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(&memcache.Item{Key: key, Value: val, Expiration: int32(ttl.Seconds())})
}

// func (m *MemCache) SetWithTTL(key string, value string, expiration int) error {
// 	fmt.Printf("2- Setting key: %s with value: %s and TTL: %d seconds\n", key, value, expiration) // Add this line for logging
// 	return m.client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expiration)})
// }

func (m *MemCache) Delete(key string) error {
	return m.client.Delete(key)
}

func (m *MemCache) ClearAll() error {
	return m.client.FlushAll()
}
