package kset

import (
	"iter"
	"sync"

	"github.com/igrmk/treemap/v2"
	"golang.org/x/exp/constraints"
)

type treeMapStore[K constraints.Ordered, V any] struct {
	mutex sync.RWMutex
	store *treemap.TreeMap[K, V]
}

func newStoreTreeMapKV[K constraints.Ordered, V any](selector func(V) K, values ...V) *treeMapStore[K, V] {
	store := &treeMapStore[K, V]{
		store: treemap.New[K, V](),
	}

	for _, value := range values {
		store.Upsert(selector(value), value)
	}

	return store
}

func newStoreTreeMapK[K constraints.Ordered](values ...K) *treeMapStore[K, struct{}] {
	store := &treeMapStore[K, struct{}]{
		store: treemap.New[K, struct{}](),
	}

	for _, value := range values {
		store.Upsert(value, struct{}{})
	}

	return store
}

func (t *treeMapStore[K, V]) Clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.store.Clear()
}

func (t *treeMapStore[K, V]) Clone() store[K, V] {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	clone := treemap.New[K, V]()

	for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
		clone.Set(iter.Key(), iter.Value())
	}

	return &treeMapStore[K, V]{
		store: clone,
	}
}

func (t *treeMapStore[K, V]) Contains(key K) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.store.Contains(key)
}

func (t *treeMapStore[K, V]) Delete(key K) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.store.Del(key)
}

func (t *treeMapStore[K, V]) Get(key K) (V, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.store.Get(key)
}

func (t *treeMapStore[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		t.mutex.RLock()
		defer t.mutex.RUnlock()
		for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
			if !yield(iter.Key(), iter.Value()) {
				return
			}
		}
	}
}

func (t *treeMapStore[K, V]) Len() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.store.Len()
}

func (t *treeMapStore[K, V]) Upsert(key K, value V) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.store.Set(key, value)
}

var _ store[string, string] = &treeMapStore[string, string]{}
