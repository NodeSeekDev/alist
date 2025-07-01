package cache

import "time"

type Config[V any] struct {
	shards          int
	expiredCallback ExpiredCallback[V]
	hash            IHash
	clearInterval   time.Duration
}

func NewConfig[V any]() *Config[V] {
	return &Config[V]{shards: 1024, hash: newDefaultHash(), clearInterval: 10 * time.Minute}
}
