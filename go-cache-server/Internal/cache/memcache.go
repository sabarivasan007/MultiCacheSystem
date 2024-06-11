package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemCache struct {
	client *memcache.Client
	ttl    int32
}

type CacheItem struct {
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
	TTL       int32  `json:"ttl"`
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

func (m *MemCache) GetWithTTL(key string) (string, int, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return "", 0, err
	}
	log.Printf("Raw value from cache: %s", string(item.Value)) // Add this line for debugging

	var cacheItem CacheItem
	if err := json.Unmarshal(item.Value, &cacheItem); err != nil {
		return "", 0, err
	}
	ttl := cacheItem.TTL - int32(time.Now().Unix()-cacheItem.Timestamp)
	if ttl < 0 {
		return "", 0, errors.New("item has expired")
	}
	return cacheItem.Value, int(ttl), nil
}

func (m *MemCache) Set(key string, value string) error {
	return m.client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: m.ttl})
}

func (m *MemCache) SetWithTTL(key string, value string, expiration int) error {
	fmt.Printf("2- Setting key: %s with value: %s and TTL: %d seconds\n", key, value, expiration) // Add this line for logging
	return m.client.Set(&memcache.Item{Key: key, Value: []byte(value), Expiration: int32(expiration)})
}

func (m *MemCache) Delete(key string) error {
	return m.client.Delete(key)
}

func (m *MemCache) ClearAll() error {
	return m.client.FlushAll()
}
