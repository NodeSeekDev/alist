package cache

import "time"

type IItem interface {
	Expired() bool
	CanExpire() bool
	SetExpireAt(t time.Time)
}

type Item[V any] struct {
	v      V
	expire time.Time
}

func (i *Item[V]) Expired() bool {
	if !i.CanExpire() {
		return false
	}
	return time.Now().Equal(i.expire) || time.Now().After(i.expire)
}

func (i *Item[V]) CanExpire() bool {
	return !i.expire.IsZero()
}

func (i *Item[V]) SetExpireAt(t time.Time) {
	i.expire = t
}
