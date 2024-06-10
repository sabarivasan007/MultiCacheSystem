package cache

type CacheLibrary interface {
	Get(key string) (string, error)
	//GetWithTTL(key string) (string, error)
	Set(key string, value string) error
	SetWithTTL(key string, value string, ttl int) error
	Delete(key string) error
	ClearAll() error
}
