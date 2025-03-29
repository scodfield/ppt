package base

import (
	"fmt"
	"sync"
	"time"
)

const (
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Item[V any] struct {
	Object     V
	Expiration int64
}

func (item Item[V]) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}

type janitor[K comparable, V any] struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor[K, V]) Run(c *cache[K, V]) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

type cache[K comparable, V any] struct {
	defaultExpiration time.Duration
	items             map[K]*Item[V]
	mu                sync.RWMutex
	onEvicted         func(k K, v V)
	janitor           *janitor[K, V]
	isRefreshTTL      bool
	nilV              V
}

func (c *cache[K, V]) get(k K) (V, bool) {
	item, found := c.items[k]
	if !found {
		return c.nilV, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return c.nilV, false
		}

		if c.isRefreshTTL {
			item.Expiration = time.Now().Add(c.defaultExpiration).UnixNano()
		}
	}

	return item.Object, true
}

func (c *cache[K, V]) set(k K, v V, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = &Item[V]{
		Object:     v,
		Expiration: e,
	}
}

func (c *cache[K, V]) delete(k K) (V, bool) {
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return c.nilV, false
}

func (c *cache[K, V]) Set(k K, v V, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(k, v, d)
}

func (c *cache[K, V]) SetDefault(k K, v V) {
	c.Set(k, v, DefaultExpiration)
}

func (c *cache[K, V]) Add(k K, v V, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.get(k)
	if found {
		return fmt.Errorf("item %s already exists", k)
	}
	c.set(k, v, d)
	return nil
}

type keyAndValue[K comparable, V any] struct {
	key K
	val V
}

// DeleteExpired delete all expired items from cache
func (c *cache[K, V]) DeleteExpired() {
	var evictedItems []keyAndValue[K, V]
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, item := range c.items {
		if item.Expiration > 0 && now > item.Expiration {
			evictedItems = append(evictedItems, keyAndValue[K, V]{key: key, val: item.Object})
		}
	}
	for _, v := range evictedItems {
		_, evicted := c.delete(v.key)
		if evicted {
			c.onEvicted(v.key, v.val)
		}
	}
}

func (c *cache[K, V]) StopJanitor() {
	c.janitor.stop <- true
}

type Cache[K comparable, V any] struct {
	*cache[K, V]
}
