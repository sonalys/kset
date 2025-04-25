package kset

import (
	"iter"
	"sync"
)

type keySet[K comparable] struct {
	sync.RWMutex
	data *unsafeKeySet[K]
}

// New creates a new thread-safe key-only set from any slice of comparable values.
// Optionally, it can be initialized with one or more values.
// The returned set is safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create a set of User structs, using ID as the key.
//	set := kset.New("id1", "id2")
func New[K comparable](values ...K) KeyOnlySet[K] {
	return &keySet[K]{
		data: newKeySetUnsafe(values...),
	}
}

// NewFrom creates a new thread-safe key-only set from any given slice.
// It requires a selector function that extracts a comparable key K from a value V.
// Optionally, it can be initialized with one or more values.
// The returned set is safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create a set of User structs, using ID as the key.
//	set := kset.NewFrom(func(u User) int { return u.ID }, user1, user2)
func NewFrom[K comparable, V any](selector func(V) K, values ...V) KeyOnlySet[K] {
	keys := make([]K, 0, len(values))

	for i := range values {
		keys = append(keys, selector(values[i]))
	}

	return New(keys...)
}

func (k *keySet[K]) Append(values ...K) int {
	k.Lock()
	defer k.Unlock()

	return k.data.Append(values...)
}

func (k *keySet[K]) Clear() {
	k.Lock()
	defer k.Unlock()

	k.data.Clear()
}

func (k *keySet[K]) Clone() KeyOnlySet[K] {
	k.RLock()
	defer k.RUnlock()

	return k.data.Clone()
}

func (k *keySet[K]) ContainsAnyKey(keys ...K) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.ContainsAnyKey(keys...)
}

func (k *keySet[K]) ContainsKeys(keys ...K) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.ContainsKeys(keys...)
}

func (k *keySet[K]) Difference(other KeySet[K]) KeyOnlySet[K] {
	k.RLock()
	defer k.RUnlock()

	return k.data.Difference(other)
}

func (k *keySet[K]) Each(fn func(K) bool) {
	k.RLock()
	defer k.RUnlock()

	k.data.Each(fn)
}

func (k *keySet[K]) Equal(other KeySet[K]) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.Equal(other)
}

func (k *keySet[K]) Intersect(other KeySet[K]) KeyOnlySet[K] {
	k.RLock()
	defer k.RUnlock()

	return k.data.Intersect(other)
}

func (k *keySet[K]) Intersects(other KeySet[K]) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.Intersects(other)
}

func (k *keySet[K]) IsEmpty() bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.IsEmpty()
}

func (k *keySet[K]) IsProperSubset(other KeySet[K]) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.IsProperSubset(other)
}

func (k *keySet[K]) IsProperSuperset(other KeySet[K]) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.IsProperSuperset(other)
}

func (k *keySet[K]) IsSubset(other KeySet[K]) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.IsSubset(other)
}

func (k *keySet[K]) IsSuperset(other KeySet[K]) bool {
	k.RLock()
	defer k.RUnlock()

	return k.data.IsSuperset(other)
}

func (k *keySet[K]) Iter() iter.Seq[K] {
	k.RLock()
	defer k.RUnlock()

	return k.data.Iter()
}

func (k *keySet[K]) IterKeys() iter.Seq[K] {
	k.RLock()
	defer k.RUnlock()

	return k.data.IterKeys()
}

func (k *keySet[K]) Len() int {
	k.RLock()
	defer k.RUnlock()

	return k.data.Len()
}

func (k *keySet[K]) Pop() (K, bool) {
	k.Lock()
	defer k.Unlock()

	return k.data.Pop()
}

func (k *keySet[K]) Remove(v ...K) {
	k.Lock()
	defer k.Unlock()

	k.data.Remove(v...)
}

func (k *keySet[K]) SymmetricDifference(other KeyOnlySet[K]) KeyOnlySet[K] {
	k.RLock()
	defer k.RUnlock()

	return k.data.SymmetricDifference(other)
}

func (k *keySet[K]) ToSlice() []K {
	k.RLock()
	defer k.RUnlock()

	return k.data.ToSlice()
}

func (k *keySet[K]) Union(other KeyOnlySet[K]) KeyOnlySet[K] {
	k.RLock()
	defer k.RUnlock()

	return k.data.Union(other)
}
