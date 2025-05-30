package kset

import (
	"iter"

	"github.com/igrmk/treemap/v2"
	"golang.org/x/exp/constraints"
)

type unsafeTreeMapStore[Key constraints.Ordered, Value any] struct {
	store *treemap.TreeMap[Key, Value]
}

// UnsafeTreeMapKeyValue is a thread-unsafe red-black tree key-value set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(logN)		O(logN)
//	Insert			O(logN)		O(logN)
//	Delete			O(logN)		O(logN)
//
// Space complexity
//
//	Space			O(n)		O(n)
func UnsafeTreeMapKeyValue[Key constraints.Ordered, Value any](selector func(Value) Key, values ...Value) KeyValueSet[Key, Value] {
	store := &unsafeTreeMapStore[Key, Value]{
		store: treemap.New[Key, Value](),
	}

	for i := range values {
		store.store.Set(selector(values[i]), values[i])
	}

	return &keyValueSet[Key, Value, *unsafeTreeMapStore[Key, Value]]{
		store:    store,
		selector: selector,
	}
}

// UnsafeTreeMapKey is a thread-unsafe red-black tree key set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(logN)		O(logN)
//	Insert			O(logN)		O(logN)
//	Delete			O(logN)		O(logN)
//
// Space complexity
//
//	Space			O(n)		O(n)
func UnsafeTreeMapKey[Key constraints.Ordered](keys ...Key) KeySet[Key] {
	store := &unsafeTreeMapStore[Key, empty]{
		store: treemap.New[Key, empty](),
	}

	for _, value := range keys {
		store.Upsert(value, empty{})
	}

	return &keySet[Key, *unsafeTreeMapStore[Key, empty]]{
		store: store,
	}
}

func (t *unsafeTreeMapStore[Key, Value]) Clear() {
	t.store.Clear()
}

func (t *unsafeTreeMapStore[Key, Value]) Clone() Storage[Key, Value] {
	clone := treemap.New[Key, Value]()

	for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
		clone.Set(iter.Key(), iter.Value())
	}

	return &unsafeTreeMapStore[Key, Value]{
		store: clone,
	}
}

func (t *unsafeTreeMapStore[Key, Value]) Contains(key Key) bool {
	return t.store.Contains(key)
}

func (t *unsafeTreeMapStore[Key, Value]) Delete(keys ...Key) {
	for _, key := range keys {
		t.store.Del(key)
	}
}

func (t *unsafeTreeMapStore[Key, Value]) Get(key Key) (Value, bool) {
	return t.store.Get(key)
}

func (t *unsafeTreeMapStore[Key, Value]) Iter() iter.Seq2[Key, Value] {
	return func(yield func(Key, Value) bool) {
		for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
			if !yield(iter.Key(), iter.Value()) {
				return
			}
		}
	}
}

func (t *unsafeTreeMapStore[Key, Value]) Len() int {
	return t.store.Len()
}

func (t *unsafeTreeMapStore[Key, Value]) Upsert(key Key, value Value) {
	t.store.Set(key, value)
}

var _ Storage[string, string] = &unsafeTreeMapStore[string, string]{}
