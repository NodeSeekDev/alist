package cache

import (
	"sync"
	"time"
)

// ExpiredCallback Callback the function when the key-value pair expires
// Note that it is executed after expiration
type ExpiredCallback[V any] func(k string, v V) error

type memCacheShard[V any] struct {
	hashmap         map[string]Item[V]
	lock            sync.RWMutex
	expiredCallback ExpiredCallback[V]
}

func newMemCacheShard[V any](conf *Config[V]) *memCacheShard[V] {
	return &memCacheShard[V]{
		expiredCallback: conf.expiredCallback,
		hashmap:         map[string]Item[V]{},
	}
}

func (c *memCacheShard[V]) set(k string, item *Item[V]) {
	c.lock.Lock()
	c.hashmap[k] = *item
	c.lock.Unlock()
	return
}

func getZero[T any]() T {
	var result T
	return result
}

func (c *memCacheShard[V]) get(k string) (V, bool) {
	c.lock.RLock()
	item, exist := c.hashmap[k]
	c.lock.RUnlock()
	if !exist {
		return getZero[V](), false
	}
	if !item.Expired() {
		return item.v, true
	}
	if c.delExpired(k) {
		return getZero[V](), false
	}
	return c.get(k)
}

func (c *memCacheShard[V]) getSet(k string, item *Item[V]) (V, bool) {
	defer c.set(k, item)
	return c.get(k)
}

func (c *memCacheShard[V]) getDel(k string) (V, bool) {
	defer c.del(k)
	return c.get(k)
}

func (c *memCacheShard[V]) del(k string) int {
	var count int
	c.lock.Lock()
	v, found := c.hashmap[k]
	if found {
		delete(c.hashmap, k)
		if !v.Expired() {
			count++
		}
	}
	c.lock.Unlock()
	return count
}

// clear all
func (c *memCacheShard[V]) clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.hashmap = map[string]Item[V]{}
}

//delExpired Only delete when key expires
func (c *memCacheShard[V]) delExpired(k string) bool {
	c.lock.Lock()
	item, found := c.hashmap[k]
	if !found || !item.Expired() {
		c.lock.Unlock()
		return false
	}
	delete(c.hashmap, k)
	c.lock.Unlock()
	if c.expiredCallback != nil {
		_ = c.expiredCallback(k, item.v)
	}
	return true
}

func (c *memCacheShard[V]) ttl(k string) (time.Duration, bool) {
	c.lock.RLock()
	v, found := c.hashmap[k]
	c.lock.RUnlock()
	if !found || !v.CanExpire() || v.Expired() {
		return 0, false
	}
	return v.expire.Sub(time.Now()), true
}

func (c *memCacheShard[V]) checkExpire() {
	var expiredKeys []string
	c.lock.RLock()
	for k, item := range c.hashmap {
		if item.Expired() {
			expiredKeys = append(expiredKeys, k)
		}
	}
	c.lock.RUnlock()
	for _, k := range expiredKeys {
		c.delExpired(k)
	}
}
