package kset

import "iter"

type unsafeMapStore[K comparable, V any] struct {
	store map[K]V
}

func NewUnsafeMapStore[K comparable, V any](selector func(V) K, values ...V) *unsafeMapStore[K, V] {
	store := &unsafeMapStore[K, V]{
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
func NewUnsafeStoreMapKey[K comparable](values ...K) *unsafeMapStore[K, struct{}] {
	store := &unsafeMapStore[K, struct{}]{
		store: make(map[K]struct{}, len(values)),
	}

	for _, value := range values {
		store.Upsert(value, struct{}{})
	}

	return store
}

func (m *unsafeMapStore[K, V]) Clear() {
	for k := range m.store {
		delete(m.store, k)
	}
}

func (m *unsafeMapStore[K, V]) Contains(key K) bool {
	_, ok := m.store[key]
	return ok
}

func (m *unsafeMapStore[K, V]) Delete(key K) {
	delete(m.store, key)
}

func (m *unsafeMapStore[K, V]) Get(key K) (V, bool) {
	value, ok := m.store[key]
	return value, ok
}

func (m *unsafeMapStore[K, V]) Len() int {
	return len(m.store)
}

func (m *unsafeMapStore[K, V]) Upsert(key K, value V) {
	m.store[key] = value
}

func (m *unsafeMapStore[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, value := range m.store {
			if !yield(key, value) {
				return
			}
		}
	}
}

func (m *unsafeMapStore[K, V]) Clone() Store[K, V] {
	store := &unsafeMapStore[K, V]{
		store: make(map[K]V, m.Len()),
	}

	for key, value := range m.store {
		store.Upsert(key, value)
	}

	return store
}

var _ Store[string, string] = &unsafeMapStore[string, string]{}
