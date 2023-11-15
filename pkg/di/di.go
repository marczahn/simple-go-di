package di

import (
	"sync"
)

var mainMx sync.Mutex
var mxs = map[string]*sync.Mutex{}
var reg = map[string]any{}

func getMx(k string) *sync.Mutex {
	mainMx.Lock()
	defer mainMx.Unlock()
	if mx, ok := mxs[k]; ok {
		return mx
	}
	mxs[k] = &sync.Mutex{}

	return mxs[k]
}

func set[T any](k string, c func() T) T {
	v := c()
	reg[k] = v

	return v
}

func get[T any](k string) (T, bool) {
	item, ok := reg[k]
	if !ok {
		var v T

		return v, false
	}

	if out, casted := item.(T); casted {
		return out, true
	}

	panic("could not cast object with name " + k)
}

func ResetReg() {
	mainMx.Lock()
	defer mainMx.Unlock()
	mxs = map[string]*sync.Mutex{}
	reg = map[string]any{}
}

func GetOrSet[T any](k string, c func() T, overwrite bool) T {
	mx := getMx(k)
	mx.Lock()
	defer mx.Unlock()

	if overwrite {
		return set[T](k, c)
	}
	if item, ok := get[T](k); ok {
		return item
	}

	return set[T](k, c)
}
