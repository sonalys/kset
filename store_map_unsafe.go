package kset

import "iter"

type unsafeMapStore[K comparable, V any] struct {
	store    map[K]V
	selector func(V) K
}

func NewUnsafeMapStore[K comparable, V any](selector func(V) K, values ...V) *unsafeMapStore[K, V] {
	store := &unsafeMapStore[K, V]{
		store:    make(map[K]V, len(values)),
		selector: selector,
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

// NewKeyValueSet creates a new non-thread-safe set.
// It requires a selector function that extracts a comparable key K from a value V.
// Optionally, it can be initialized with one or more values.
// The returned set is *not* safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create an unsafe set of User structs, using ID as the key.
//	userSet := store.NewKeyValueSet()
func (m *unsafeMapStore[K, V]) NewKeyValueSet() KeyValueSet[K, V] {
	return &keyValueSet[K, V]{
		data:     m,
		selector: m.selector,
		newStore: func(len int) Store[K, V] {
			return &unsafeMapStore[K, V]{
				store:    make(map[K]V, len),
				selector: m.selector,
			}
		},
	}
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
	selector := m.selector

	store := &unsafeMapStore[K, V]{
		store:    make(map[K]V, m.Len()),
		selector: selector,
	}

	for _, value := range m.store {
		store.Upsert(selector(value), value)
	}

	return store
}

func (m *unsafeMapStore[K, V]) Selector() func(V) K {
	return m.selector
}

var _ Store[string, string] = &unsafeMapStore[string, string]{}
