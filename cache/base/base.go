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

func (c *cache[K, V]) Get(k K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.get(k)
}

func (c *cache[K, V]) GetWithExpiration(k K) (V, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[k]
	if !found {
		return c.nilV, time.Time{}, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return c.nilV, time.Time{}, false
		}
		if c.isRefreshTTL {
			item.Expiration = time.Now().Add(c.defaultExpiration).UnixNano()
		}
		return item.Object, time.Unix(0, item.Expiration), true
	}
	return item.Object, time.Time{}, false
}

func (c *cache[K, V]) Set(k K, v V, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.set(k, v, d)
}

func (c *cache[K, V]) SetDefault(k K, v V) {
	c.Set(k, v, DefaultExpiration)
}

func (c *cache[K, V]) Delete(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, evicted := c.delete(k)
	if evicted {
		c.onEvicted(k, v)
	}
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

func (c *cache[K, V]) Replace(k K, v V, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, found := c.get(k)
	if !found {
		return fmt.Errorf("item %s does not exist", k)
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

// OnEvicted set optional function that is called when an item is evicted from the cache
func (c *cache[K, V]) OnEvicted(f func(k K, v V)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvicted = f
}

func (c *cache[K, V]) Range(f func(k K, v V) bool) {
	if f == nil {
		return
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	isContinue := true
	for k, v := range c.items {
		isContinue = f(k, v.Object)
		if !isContinue {
			return
		}
	}
}

func (c *cache[K, V]) Items() map[K]*Item[V] {
	c.mu.RLock()
	defer c.mu.RUnlock()
	now := time.Now().UnixNano()
	items := make(map[K]*Item[V])
	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			continue
		}
		items[k] = &Item[V]{}
	}
	return items
}

// ItemCount return the number of items in cache include expired but not cleaned up
func (c *cache[K, V]) ItemCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *cache[K, V]) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[K]*Item[V])
}

func runJanitor[K comparable, V any](c *cache[K, V], interval time.Duration) {
	c.janitor = &janitor[K, V]{
		Interval: interval,
		stop:     make(chan bool),
	}
	go c.janitor.Run(c)
}

func newCache[K comparable, V any](de time.Duration, items map[K]*Item[V], isRefresh bool) *cache[K, V] {
	if de == 0 {
		de = NoExpiration
	}
	return &cache[K, V]{
		defaultExpiration: de,
		items:             items,
		isRefreshTTL:      isRefresh,
	}
}

func newCacheWithJanitor[K comparable, V any](de time.Duration, ci time.Duration, items map[K]*Item[V], isRefresh bool) *Cache[K, V] {
	c := newCache(de, items, isRefresh)
	C := &Cache[K, V]{
		c,
	}
	if ci > 0 {
		runJanitor(c, ci)
	}
	return C
}

func (c *Cache[K, V]) StopJanitor() {
	c.janitor.stop <- true
}

type Cache[K comparable, V any] struct {
	*cache[K, V]
}

func New[K comparable, V any](defaultExpiration, cleanupInterval time.Duration, isRefresh bool) *Cache[K, V] {
	items := make(map[K]*Item[V])
	return newCacheWithJanitor[K, V](defaultExpiration, cleanupInterval, items, isRefresh)
}
