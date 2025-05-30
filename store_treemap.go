package kset

import (
	"iter"
	"sync"

	"github.com/igrmk/treemap/v2"
	"golang.org/x/exp/constraints"
)

type treeMapStore[Key constraints.Ordered, Value any] struct {
	mutex sync.RWMutex
	store *treemap.TreeMap[Key, Value]
}

// TreeMapKeyValue is a thread-safe red-black tree key-value set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(logN)		O(logN)
//	Insert			O(logN)		O(logN)
//	Delete			O(logN)		O(logN)
//
// Space complexity
//
//	Space			O(n)		O(n)
func TreeMapKeyValue[Key constraints.Ordered, Value any](selector func(Value) Key, values ...Value) KeyValueSet[Key, Value] {
	store := &treeMapStore[Key, Value]{
		store: treemap.New[Key, Value](),
	}

	for i := range values {
		store.store.Set(selector(values[i]), values[i])
	}

	return &keyValueSet[Key, Value, *treeMapStore[Key, Value]]{
		store:    store,
		selector: selector,
	}
}

// TreeMapKey is a thread-safe red-black tree key set implementation.
//
//	Operation		Average		WorstCase
//	Search			O(logN)		O(logN)
//	Insert			O(logN)		O(logN)
//	Delete			O(logN)		O(logN)
//
// Space complexity
//
//	Space			O(n)		O(n)
func TreeMapKey[Key constraints.Ordered](keys ...Key) KeySet[Key] {
	store := &treeMapStore[Key, empty]{
		store: treemap.New[Key, empty](),
	}

	for _, value := range keys {
		store.Upsert(value, empty{})
	}

	return &keySet[Key, *treeMapStore[Key, empty]]{
		store: store,
	}
}

func (t *treeMapStore[Key, Value]) Clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.store.Clear()
}

func (t *treeMapStore[Key, Value]) Clone() Storage[Key, Value] {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	clone := treemap.New[Key, Value]()

	for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
		clone.Set(iter.Key(), iter.Value())
	}

	return &treeMapStore[Key, Value]{
		store: clone,
	}
}

func (t *treeMapStore[Key, Value]) Contains(key Key) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.store.Contains(key)
}

func (t *treeMapStore[Key, Value]) Delete(keys ...Key) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, key := range keys {
		t.store.Del(key)
	}
}

func (t *treeMapStore[Key, Value]) Get(key Key) (Value, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.store.Get(key)
}

func (t *treeMapStore[Key, Value]) Iter() iter.Seq2[Key, Value] {
	return func(yield func(Key, Value) bool) {
		t.mutex.RLock()
		defer t.mutex.RUnlock()
		for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
			if !yield(iter.Key(), iter.Value()) {
				return
			}
		}
	}
}

func (t *treeMapStore[Key, Value]) Len() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.store.Len()
}

func (t *treeMapStore[Key, Value]) Upsert(key Key, value Value) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.store.Set(key, value)
}

var _ Storage[string, string] = &treeMapStore[string, string]{}
