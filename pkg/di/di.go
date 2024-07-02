package di

import (
	"sync"
)

// NewInstance returns an object of type *di.instance typed by a definition of it's content
func NewInstance[T any]() *instance[T] {
	return &instance[T]{}
}

type instance[T any] struct {
	instance T
	once     sync.Once
}

// GetOrSet receives a function that returns a generic value of type T which is stored within the receiver.
// The second argument defines if the value should be overwritten in case that it already exists
func (r *instance[T]) GetOrSet(c func() T, overwrite bool) T {
	if overwrite {
		return c()
	}
	r.once.Do(func() {
		r.instance = c()
	})

	return r.instance
}
