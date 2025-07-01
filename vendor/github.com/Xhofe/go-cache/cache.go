package cache

import (
	"runtime"
	"time"
)

type ICache[T any] interface {
	//Set key to hold the string value. If key already holds a value, it is overwritten, regardless of its type.
	//Any previous time to live associated with the key is discarded on successful SET operation.
	//Example:
	//c.Set("demo", 1)
	//c.Set("demo", 1, WithEx[int](10*time.Second))
	//c.Set("demo", 1, WithEx[int](10*time.Second), WithNx[int]())
	Set(k string, v T, opts ...SetIOption[T]) bool
	//Get the value of key.
	//If the key does not exist the special value nil,false is returned.
	//Example:
	//c.Get("demo") //nil, false
	//c.Set("demo", "value")
	//c.Get("demo") //"value", true
	Get(k string) (T, bool)
	//GetSet Atomically sets key to value and returns the old value stored at key.
	//Returns nil,false when key not exists.
	//Example:
	//c.GetSet("demo", 1) //nil,false
	//c.GetSet("demo", 2) //1,true
	GetSet(k string, v T, opts ...SetIOption[T]) (T, bool)
	//GetDel Get the value of key and delete the key.
	//This command is similar to GET, except for the fact that it also deletes the key on success.
	//Example:
	//c.Set("demo", "value")
	//c.GetDel("demo") //"value", true
	//c.GetDel("demo") //nil, false
	GetDel(k string) (T, bool)
	//Del Removes the specified keys. A key is ignored if it does not exist.
	//Return the number of keys that were removed.
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "1")
	//c.Del("demo1", "demo2", "demo3") //2
	Del(keys ...string) int
	//DelExpired Only delete when key expires
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "1", WithEx[string](1*time.Second))
	//time.Sleep(1*time.Second)
	//c.DelExpired("demo1", "demo2") //true
	DelExpired(k string) bool
	//Exists Returns if key exists.
	//Return the number of exists keys.
	//Example:
	//c.Set("demo1", "1")
	//c.Set("demo2", "1")
	//c.Exists("demo1", "demo2", "demo3") //2
	Exists(keys ...string) bool
	//Expire Set a timeout on key.
	//After the timeout has expired, the key will automatically be deleted.
	//Return false if the key not exist.
	//Example:
	//c.Expire("demo", 1*time.Second) // false
	//c.Set("demo", "1")
	//c.Expire("demo", 1*time.Second) // true
	Expire(k string, d time.Duration) bool
	//ExpireAt has the same effect and semantic as Expire, but instead of specifying the number of seconds representing the TTL (time to live),
	//it takes an absolute Unix Time (seconds since January 1, 1970). A Time in the past will delete the key immediately.
	//Return false if the key not exist.
	//Example:
	//c.ExpireAt("demo", time.Now().Add(10*time.Second)) // false
	//c.Set("demo", "1")
	//c.ExpireAt("demo", time.Now().Add(10*time.Second)) // true
	ExpireAt(k string, t time.Time) bool
	//Persist Remove the existing timeout on key.
	//Return false if the key not exist.
	//Example:
	//c.Persist("demo") // false
	//c.Set("demo", "1")
	//c.Persist("demo") // true
	Persist(k string) bool
	//Ttl Returns the remaining time to live of a key that has a timeout.
	//Returns 0,false if the key does not exist or if the key exist but has no associated expire.
	//Example:
	//c.Set("demo", "1")
	//c.Ttl("demo") // 0,false
	//c.Set("demo", "1", WithEx[string](10*time.Second))
	//c.Ttl("demo") // 10*time.Second,true
	Ttl(k string) (time.Duration, bool)
	// Clear all the keys
	Clear()
}

func NewMemCache[V any](opts ...ICacheOption[V]) ICache[V] {
	conf := NewConfig[V]()
	for _, opt := range opts {
		opt(conf)
	}

	c := &memCache[V]{
		shards:    make([]*memCacheShard[V], conf.shards),
		closed:    make(chan struct{}),
		shardMask: uint64(conf.shards - 1),
		config:    conf,
		hash:      conf.hash,
	}
	for i := 0; i < conf.shards; i++ {
		c.shards[i] = newMemCacheShard(conf)
	}
	if conf.clearInterval > 0 {
		go func() {
			ticker := time.NewTicker(conf.clearInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					for _, shard := range c.shards {
						shard.checkExpire()
					}
				case <-c.closed:
					return
				}
			}
		}()
	}
	cache := &MemCache[V]{c}
	// Associated finalizer function with obj.
	// When the obj is unreachable, close the obj.
	runtime.SetFinalizer(cache, func(cache *MemCache[V]) { close(cache.closed) })
	return cache
}

type MemCache[V any] struct {
	*memCache[V]
}

type memCache[V any] struct {
	shards    []*memCacheShard[V]
	hash      IHash
	shardMask uint64
	config    *Config[V]
	closed    chan struct{}
}

func (c *memCache[V]) Set(k string, v V, opts ...SetIOption[V]) bool {
	item := Item[V]{v: v}
	for _, opt := range opts {
		if pass := opt(c, k, &item); !pass {
			return false
		}
	}
	hashedKey := c.hash.Sum64(k)
	shard := c.getShard(hashedKey)
	shard.set(k, &item)
	return true
}

func (c *memCache[V]) Get(k string) (V, bool) {
	hashedKey := c.hash.Sum64(k)
	shard := c.getShard(hashedKey)
	return shard.get(k)
}

func (c *memCache[V]) GetSet(k string, v V, opts ...SetIOption[V]) (V, bool) {
	defer c.Set(k, v, opts...)
	return c.Get(k)
}

func (c *memCache[V]) GetDel(k string) (V, bool) {
	defer c.Del(k)
	return c.Get(k)
}

func (c *memCache[V]) Del(ks ...string) int {
	var count int
	for _, k := range ks {
		hashedKey := c.hash.Sum64(k)
		shard := c.getShard(hashedKey)
		count += shard.del(k)
	}
	return count
}

func (c *memCache[V]) DelExpired(k string) bool {
	hashedKey := c.hash.Sum64(k)
	shard := c.getShard(hashedKey)
	return shard.delExpired(k)
}

func (c *memCache[V]) Exists(ks ...string) bool {
	for _, k := range ks {
		if _, found := c.Get(k); !found {
			return false
		}
	}
	return true
}

func (c *memCache[V]) Expire(k string, d time.Duration) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v, WithEx[V](d))
}

func (c *memCache[V]) ExpireAt(k string, t time.Time) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v, WithExAt[V](t))
}

func (c *memCache[V]) Persist(k string) bool {
	v, found := c.Get(k)
	if !found {
		return false
	}
	return c.Set(k, v)
}

func (c *memCache[V]) Ttl(k string) (time.Duration, bool) {
	hashedKey := c.hash.Sum64(k)
	shard := c.getShard(hashedKey)
	return shard.ttl(k)
}

func (c *memCache[V]) getShard(hashedKey uint64) (shard *memCacheShard[V]) {
	return c.shards[hashedKey&c.shardMask]
}

func (c *memCache[V]) Clear() {
	for _, s := range c.shards {
		s.clear()
	}
}
