package cache

import (
	"ppt/cache/base"
	"time"
)

// CacheLoader 缓存加载
type CacheLoader[K comparable, V any] interface {
	Load(key K) (value V, err error)
}

type Cache[K comparable, V any] struct {
	ttl    time.Duration
	cache  *base.Cache[K, V]
	Loader CacheLoader[K, V]
}

func (c *Cache[K, V]) Get(key K) (any, error) {
	value, exists := c.cache.Get(key)
	if !exists {
		load, err := c.Loader.Load(key)
		if err != nil {
			return load, err
		}
		c.cache.Set(key, load, c.ttl)
		return load, nil
	}
	return value, nil
}

func (c *Cache[K, V]) Set(key K, value any, ttl time.Duration) {
	c.cache.Set(key, value, ttl)
}

func (c *Cache[K, V]) Range(f func(key K, value V) bool) {
	c.cache.Range(f)
}

func (c *Cache[K, V]) SetEvicted(f func(key K, value V)) {
	c.cache.OnEvicted(f)
}

func (c *Cache[K, V]) Delete(key K) {
	c.cache.Delete(key)
}

func (c *Cache[K, V]) StopCache() {
	c.cache.StopJanitor()
}

func (c *Cache[K, V]) Load(key K) (V, error) {
	return c.Loader.Load(key)
}

func NewCache[K comparable, V any](defaultExpiration, cleanupInterval time.Duration, loader CacheLoader[K, V], isRefresh bool) *Cache[K, V] {
	if defaultExpiration <= 0 {
		defaultExpiration = 1 * time.Minute
	}
	if cleanupInterval <= 0 {
		cleanupInterval = 1 * time.Minute
	}
	return &Cache[K, V]{
		ttl:    defaultExpiration,
		cache:  base.New[K, V](defaultExpiration, cleanupInterval, isRefresh),
		Loader: loader,
	}
}
