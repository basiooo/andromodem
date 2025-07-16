package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache implements the ICache interface using patrickmn/go-cache
type Cache struct {
	cache *cache.Cache
}

var (
	instance *Cache
	once     sync.Once
)

func NewCache(defaultExpiration, cleanupInterval time.Duration) ICache {
	once.Do(func() {
		instance = &Cache{
			cache: cache.New(defaultExpiration, cleanupInterval),
		}
	})
	return instance
}

func GetInstance() ICache {
	if instance == nil {
		return NewCache(5*time.Minute, 10*time.Minute)
	}
	return instance
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.cache.Set(key, value, expiration)
}

func (c *Cache) Delete(key string) {
	c.cache.Delete(key)
}

func (c *Cache) Flush() {
	c.cache.Flush()
}
