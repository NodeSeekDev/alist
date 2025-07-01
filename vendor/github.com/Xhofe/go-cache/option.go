package cache

import "time"

// SetIOption The option used to cache set
type SetIOption[V any] func(ICache[V], string, IItem) bool

//WithEx Set the specified expire time, in time.Duration.
func WithEx[V any](d time.Duration) SetIOption[V] {
	return func(c ICache[V], k string, v IItem) bool {
		v.SetExpireAt(time.Now().Add(d))
		return true
	}
}

//WithExAt Set the specified expire deadline, in time.Time.
func WithExAt[V any](t time.Time) SetIOption[V] {
	return func(c ICache[V], k string, v IItem) bool {
		v.SetExpireAt(t)
		return true
	}
}

// ICacheOption The option used to create the cache object
type ICacheOption[V any] func(conf *Config[V])

//WithShards set custom size of sharding. Default is 1024
//The larger the size, the smaller the lock force, the higher the concurrency performance,
//and the higher the memory footprint, so try to choose a size that fits your business scenario
func WithShards[V any](shards int) ICacheOption[V] {
	if shards <= 0 {
		panic("Invalid shards")
	}
	return func(conf *Config[V]) {
		conf.shards = shards
	}
}

//WithExpiredCallback set custom expired callback function
//This callback function is called when the key-value pair expires
func WithExpiredCallback[V any](ec ExpiredCallback[V]) ICacheOption[V] {
	return func(conf *Config[V]) {
		conf.expiredCallback = ec
	}
}

//WithHash set custom hash key function
func WithHash[V any](hash IHash) ICacheOption[V] {
	return func(conf *Config[V]) {
		conf.hash = hash
	}
}

//WithClearInterval set custom clear interval.
//Interval for clearing expired key-value pairs. The default value is 1 second
//If the d is 0, the periodic clearing function is disabled
func WithClearInterval[V any](d time.Duration) ICacheOption[V] {
	return func(conf *Config[V]) {
		conf.clearInterval = d
	}
}
