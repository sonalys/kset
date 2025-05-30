package kset

import (
	"iter"
	"maps"
)

type unsafeMapStore[Key comparable, Value any] struct {
	store map[Key]Value
}

// UnsafeHashMapKeyValue is a thread-unsafe hash table key-value set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(1)		O(logN^2)
//	Insert			O(1)		O(logN^2)
//	Delete			O(1)		O(n)
//
// Space complexity
//
//	Space			O(n)		O(n)
func UnsafeHashMapKeyValue[Key comparable, Value any](selector func(Value) Key, values ...Value) KeyValueSet[Key, Value] {
	data := make(map[Key]Value, len(values))

	for i := range values {
		data[selector(values[i])] = values[i]
	}

	return &keyValueSet[Key, Value, *unsafeMapStore[Key, Value]]{
		store: &unsafeMapStore[Key, Value]{
			store: data,
		},
		selector: selector,
	}
}

// UnsafeHashMapKey is a thread-unsafe hash table key set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(1)		O(logN^2)
//	Insert			O(1)		O(logN^2)
//	Delete			O(1)		O(n)
//
// Space complexity
//
//	Space			O(n)		O(n)
func UnsafeHashMapKey[Key comparable](values ...Key) KeySet[Key] {
	data := make(map[Key]empty, len(values))
	for _, value := range values {
		data[value] = empty{}
	}

	return &keySet[Key, *unsafeMapStore[Key, empty]]{
		store: &unsafeMapStore[Key, empty]{
			store: data,
		},
	}
}

func (m *unsafeMapStore[Key, Value]) Clear() {
	for k := range m.store {
		delete(m.store, k)
	}
}

func (m *unsafeMapStore[Key, Value]) Contains(key Key) bool {
	_, ok := m.store[key]
	return ok
}

func (m *unsafeMapStore[Key, Value]) Delete(keys ...Key) {
	for _, key := range keys {
		delete(m.store, key)
	}
}

func (m *unsafeMapStore[Key, Value]) Get(key Key) (Value, bool) {
	value, ok := m.store[key]
	return value, ok
}

func (m *unsafeMapStore[Key, Value]) Len() int {
	return len(m.store)
}

func (m *unsafeMapStore[Key, Value]) Upsert(key Key, value Value) {
	m.store[key] = value
}

func (m *unsafeMapStore[Key, Value]) Iter() iter.Seq2[Key, Value] {
	return func(yield func(Key, Value) bool) {
		for key, value := range m.store {
			if !yield(key, value) {
				return
			}
		}
	}
}

func (m *unsafeMapStore[Key, Value]) Clone() Storage[Key, Value] {
	return &unsafeMapStore[Key, Value]{
		store: maps.Clone(m.store),
	}
}

var _ Storage[string, string] = &unsafeMapStore[string, string]{}
