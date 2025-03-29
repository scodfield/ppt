package cache

import "time"

// CacheLoader 缓存加载
type CacheLoader[K comparable, V any] interface {
	Load(key K) (value V, err error)
}

type Cache[K comparable, V any] struct {
	ttl    time.Duration
	Loader CacheLoader[K, V]
}
