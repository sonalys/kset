package kset

import "iter"

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
	store := &unsafeMapStore[Key, Value]{
		store: make(map[Key]Value, len(values)),
	}

	for i := range values {
		store.store[selector(values[i])] = values[i]
	}

	return &keyValueSet[Key, Value, *unsafeMapStore[Key, Value]]{
		store:    store,
		selector: selector,
		newStore: func(i int) *unsafeMapStore[Key, Value] {
			return &unsafeMapStore[Key, Value]{
				store: make(map[Key]Value, len(values)),
			}
		},
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
	store := &unsafeMapStore[Key, empty]{
		store: make(map[Key]empty, len(values)),
	}

	for _, value := range values {
		store.Upsert(value, empty{})
	}

	return &keySet[Key, *unsafeMapStore[Key, empty]]{
		store: store,
		newStore: func(k ...Key) *unsafeMapStore[Key, empty] {
			return &unsafeMapStore[Key, empty]{
				store: make(map[Key]empty, len(values)),
			}
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

func (m *unsafeMapStore[Key, Value]) Delete(key Key) {
	delete(m.store, key)
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
	store := &unsafeMapStore[Key, Value]{
		store: make(map[Key]Value, m.Len()),
	}

	for key, value := range m.store {
		store.Upsert(key, value)
	}

	return store
}

var _ Storage[string, string] = &unsafeMapStore[string, string]{}
