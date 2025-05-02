package kset

import (
	"iter"
	"sync"
)

type safeMapStore[K comparable, V any] struct {
	mutex sync.RWMutex
	store map[K]V
}

// NewStoreMapKeyValue creates a new map store for the given values.
// The map key is created by the selector function.
// The underlying data structure is a hash-map.
// This store is thread-safe.
func NewStoreMapKeyValue[K comparable, V any](selector func(V) K, values ...V) *safeMapStore[K, V] {
	store := &safeMapStore[K, V]{
		store: make(map[K]V, len(values)),
	}

	for _, value := range values {
		store.Upsert(selector(value), value)
	}

	return store
}

// NewStoreMapKey creates a new map store for the given keys.
// The underlying data structure is a hash-map.
// This store is thread-safe.
func NewStoreMapKey[K comparable](values ...K) *safeMapStore[K, struct{}] {
	store := &safeMapStore[K, struct{}]{
		store: make(map[K]struct{}, len(values)),
	}

	for _, value := range values {
		store.Upsert(value, struct{}{})
	}

	return store
}

func (m *safeMapStore[K, V]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k := range m.store {
		delete(m.store, k)
	}
}

func (m *safeMapStore[K, V]) Contains(key K) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.store[key]
	return ok
}

func (m *safeMapStore[K, V]) Delete(key K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.store, key)
}

func (m *safeMapStore[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.store[key]
	return value, ok
}

func (m *safeMapStore[K, V]) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.store)
}

func (m *safeMapStore[K, V]) Upsert(key K, value V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.store[key] = value
}

func (m *safeMapStore[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
		for key, value := range m.store {
			if !yield(key, value) {
				return
			}
		}
	}
}

func (m *safeMapStore[K, V]) Clone() Store[K, V] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	store := &safeMapStore[K, V]{
		store: make(map[K]V, m.Len()),
	}

	for key, value := range m.store {
		store.Upsert(key, value)
	}

	return store
}

var _ Store[string, string] = &safeMapStore[string, string]{}
