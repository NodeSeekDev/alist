package gsync

import (
	"container/list"
	"errors"
	"sync"
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
	ErrQueueLess  = errors.New("queue len less than n")
)

type QueueOf[T any] struct {
	rw   sync.RWMutex
	list list.List
}

func NewQueue[T any]() *QueueOf[T] {
	return &QueueOf[T]{}
}

func (q *QueueOf[T]) Push(v T) {
	q.rw.Lock()
	defer q.rw.Unlock()
	q.list.PushBack(v)
}

func (q *QueueOf[T]) Pop() (T, error) {
	q.rw.Lock()
	defer q.rw.Unlock()
	if q.list.Len() == 0 {
		return GetZero[T](), ErrQueueEmpty
	}
	e := q.list.Front()
	q.list.Remove(e)
	return e.Value.(T), nil
}

func (q *QueueOf[T]) MustPop() T {
	v, err := q.Pop()
	if err != nil {
		panic(err)
	}
	return v
}

func (q *QueueOf[T]) Len() int {
	return q.list.Len()
}

func (q *QueueOf[T]) IsEmpty() bool {
	return q.Len() == 0
}

func (q *QueueOf[T]) Clear() {
	q.rw.Lock()
	defer q.rw.Unlock()
	q.list.Init()
}

func (q *QueueOf[T]) Peek() (T, error) {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if q.list.Len() == 0 {
		return GetZero[T](), ErrQueueEmpty
	}
	e := q.list.Front()
	return e.Value.(T), nil
}

func (q *QueueOf[T]) MustPeek() T {
	v, err := q.Peek()
	if err != nil {
		panic(err)
	}
	return v
}

func (q *QueueOf[T]) PeekN(n int) ([]T, error) {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if q.list.Len() < n {
		return nil, ErrQueueLess
	}
	var v []T
	var front = q.list.Front()
	for i := 0; i < n; i++ {
		v = append(v, front.Value.(T))
		front = front.Next()
	}
	return v, nil
}

func (q *QueueOf[T]) MustPeekN(n int) []T {
	v, err := q.PeekN(n)
	if err != nil {
		panic(err)
	}
	return v
}

func (q *QueueOf[T]) PopN(n int) ([]T, error) {
	q.rw.Lock()
	defer q.rw.Unlock()
	if q.list.Len() < n {
		return nil, ErrQueueLess
	}
	var v []T
	for i := 0; i < n; i++ {
		front := q.list.Front()
		v = append(v, front.Value.(T))
		q.list.Remove(front)
	}
	return v, nil
}

func (q *QueueOf[T]) MustPopN(n int) []T {
	v, err := q.PopN(n)
	if err != nil {
		panic(err)
	}
	return v
}

func (q *QueueOf[T]) PopAll() []T {
	q.rw.Lock()
	defer q.rw.Unlock()
	var v []T
	for q.list.Len() > 0 {
		front := q.list.Front()
		v = append(v, front.Value.(T))
		q.list.Remove(front)
	}
	return v
}

func (q *QueueOf[T]) PopWhile(f func(T) bool) []T {
	q.rw.Lock()
	defer q.rw.Unlock()
	var v []T
	front := q.list.Front()
	for front != nil && f(front.Value.(T)) {
		v = append(v, front.Value.(T))
		q.list.Remove(front)
	}
	return v
}

func (q *QueueOf[T]) PopUntil(f func(T) bool) []T {
	q.rw.Lock()
	defer q.rw.Unlock()
	var v []T
	front := q.list.Front()
	for front != nil && !f(front.Value.(T)) {
		v = append(v, front.Value.(T))
		q.list.Remove(front)
	}
	return v
}
