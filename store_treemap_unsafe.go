package kset

import (
	"iter"

	"github.com/igrmk/treemap/v2"
	"golang.org/x/exp/constraints"
)

type unsafeTreeMapStore[K constraints.Ordered, V any] struct {
	store *treemap.TreeMap[K, V]
}

func newStoreUnsafeTreeMapKV[K constraints.Ordered, V any](selector func(V) K, values ...V) *unsafeTreeMapStore[K, V] {
	store := &unsafeTreeMapStore[K, V]{
		store: treemap.New[K, V](),
	}

	for _, value := range values {
		store.Upsert(selector(value), value)
	}

	return store
}

func newStoreUnsafeTreeMapK[K constraints.Ordered](values ...K) *unsafeTreeMapStore[K, struct{}] {
	store := &unsafeTreeMapStore[K, struct{}]{
		store: treemap.New[K, struct{}](),
	}

	for _, value := range values {
		store.Upsert(value, struct{}{})
	}

	return store
}

func (t *unsafeTreeMapStore[K, V]) Clear() {

	t.store.Clear()
}

func (t *unsafeTreeMapStore[K, V]) Clone() store[K, V] {

	clone := treemap.New[K, V]()

	for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
		clone.Set(iter.Key(), iter.Value())
	}

	return &unsafeTreeMapStore[K, V]{
		store: clone,
	}
}

func (t *unsafeTreeMapStore[K, V]) Contains(key K) bool {
	return t.store.Contains(key)
}

func (t *unsafeTreeMapStore[K, V]) Delete(key K) {
	t.store.Del(key)
}

func (t *unsafeTreeMapStore[K, V]) Get(key K) (V, bool) {
	return t.store.Get(key)
}

func (t *unsafeTreeMapStore[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for iter := t.store.Iterator(); iter.Valid(); iter.Next() {
			if !yield(iter.Key(), iter.Value()) {
				return
			}
		}
	}
}

func (t *unsafeTreeMapStore[K, V]) Len() int {
	return t.store.Len()
}

func (t *unsafeTreeMapStore[K, V]) Upsert(key K, value V) {
	t.store.Set(key, value)
}

var _ store[string, string] = &unsafeTreeMapStore[string, string]{}
