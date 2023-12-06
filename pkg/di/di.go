package di

import (
	"sync"
)

func NewSingleton[T any]() *singleton[T] {
	return &singleton[T]{}
}

type singleton[T any] struct {
	mx       sync.RWMutex
	isset    bool
	instance T
}

func (r *singleton[T]) GetOrSet(c func() T, overwrite bool) T {
	r.mx.Lock()
	defer r.mx.Unlock()

	if overwrite || !r.isset {
		r.instance = c()
		r.isset = true
	}

	return r.instance
}
