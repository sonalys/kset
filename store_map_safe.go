package kset

import (
	"iter"
	"sync"
)

type safeMapStore[Key comparable, Value any] struct {
	mutex sync.RWMutex
	store map[Key]Value
}

// HashMapKeyValue is a thread-safe hash table key-value set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(1)		O(logN^2)
//	Insert			O(1)		O(logN^2)
//	Delete			O(1)		O(n)
//
// Space complexity
//
//	Space			O(n)		O(n)
func HashMapKeyValue[Key comparable, Value any](selector func(Value) Key, values ...Value) KeyValueSet[Key, Value] {
	store := &safeMapStore[Key, Value]{
		store: make(map[Key]Value, len(values)),
	}

	for i := range values {
		store.store[selector(values[i])] = values[i]
	}

	return &keyValueSet[Key, Value, *safeMapStore[Key, Value]]{
		store:    store,
		selector: selector,
	}
}

// HashMapKey is a thread-safe hash table key set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(1)		O(logN^2)
//	Insert			O(1)		O(logN^2)
//	Delete			O(1)		O(n)
//
// Space complexity
//
//	Space			O(n)		O(n)
func HashMapKey[Key comparable](values ...Key) KeySet[Key] {
	store := &safeMapStore[Key, empty]{
		store: make(map[Key]empty, len(values)),
	}

	for _, value := range values {
		store.Upsert(value, empty{})
	}

	return &keySet[Key, *safeMapStore[Key, empty]]{
		store: store,
	}
}

func (m *safeMapStore[Key, Value]) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for k := range m.store {
		delete(m.store, k)
	}
}

func (m *safeMapStore[Key, Value]) Contains(key Key) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.store[key]
	return ok
}

func (m *safeMapStore[Key, Value]) Delete(keys ...Key) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, key := range keys {
		delete(m.store, key)
	}
}

func (m *safeMapStore[Key, Value]) Get(key Key) (Value, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	value, ok := m.store[key]
	return value, ok
}

func (m *safeMapStore[Key, Value]) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.store)
}

func (m *safeMapStore[Key, Value]) Upsert(key Key, value Value) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.store[key] = value
}

func (m *safeMapStore[Key, Value]) Iter() iter.Seq2[Key, Value] {
	return func(yield func(Key, Value) bool) {
		m.mutex.RLock()
		defer m.mutex.RUnlock()
		for key, value := range m.store {
			if !yield(key, value) {
				return
			}
		}
	}
}

func (m *safeMapStore[Key, Value]) Clone() Storage[Key, Value] {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	store := &safeMapStore[Key, Value]{
		store: make(map[Key]Value, m.Len()),
	}

	for key, value := range m.store {
		store.Upsert(key, value)
	}

	return store
}

var _ Storage[string, string] = &safeMapStore[string, string]{}
