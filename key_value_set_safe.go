package kset

import (
	"iter"
	"sync"
)

type keyValueSet[K comparable, V any] struct {
	lock   sync.RWMutex
	unsafe *unsafeKeyValueSet[K, V]
}

// NewKeyValue creates a new thread-safe set.
// It requires a selector function that extracts a comparable key K from a value V.
// Optionally, it can be initialized with one or more values.
// The returned set is safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create a set of User structs, using ID as the key.
//	userSet := kset.NewKeyValue(func(u User) int { return u.ID }, user1, user2)
func NewKeyValue[K comparable, V any](selector func(V) K, values ...V) KeyValueSet[K, V] {
	return &keyValueSet[K, V]{
		unsafe: newKeyValueUnsafe(selector, values...),
	}
}

func (s *keyValueSet[K, V]) Append(v ...V) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.unsafe.Append(v...)
}

func (s *keyValueSet[K, V]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.unsafe.Clear()
}

func (s *keyValueSet[K, V]) Clone() KeyValueSet[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Clone()
}

func (s *keyValueSet[K, V]) Contains(v ...V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Contains(v...)
}

func (s *keyValueSet[K, V]) ContainsKeys(keys ...K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ContainsKeys(keys...)
}

func (s *keyValueSet[K, V]) ContainsAnyKey(keys ...K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ContainsAnyKey(keys...)
}

func (s *keyValueSet[K, V]) ContainsAny(v ...V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ContainsAny(v...)
}

func (s *keyValueSet[K, V]) Intersects(other KeySet[K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Intersects(other)
}

func (s *keyValueSet[K, V]) Difference(other KeySet[K]) KeyValueSet[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Difference(other)
}

func (s *keyValueSet[K, V]) Each(f func(V) bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	s.unsafe.Each(f)
}

func (s *keyValueSet[K, V]) Equal(other KeySet[K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Equal(other)
}

func (s *keyValueSet[K, V]) Intersect(other KeySet[K]) KeyValueSet[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Intersect(other)
}

func (s *keyValueSet[K, V]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsEmpty()
}

func (s *keyValueSet[K, V]) IsProperSubset(other KeySet[K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsProperSubset(other)
}

func (s *keyValueSet[K, V]) IsProperSuperset(other KeySet[K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsProperSuperset(other)
}

func (s *keyValueSet[K, V]) IsSubset(other KeySet[K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsSubset(other)
}

func (s *keyValueSet[K, V]) IsSuperset(other KeySet[K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsSuperset(other)
}

func (s *keyValueSet[K, V]) Iter() iter.Seq[V] {
	s.lock.RLock()

	return func(yield func(V) bool) {
		defer s.lock.RUnlock()

		for value := range s.unsafe.Iter() {
			if !yield(value) {
				break
			}
		}
	}
}

func (s *keyValueSet[K, V]) IterKeys() iter.Seq[K] {
	s.lock.RLock()

	return func(yield func(K) bool) {
		defer s.lock.RUnlock()

		for key := range s.unsafe.IterKeys() {
			if !yield(key) {
				break
			}
		}
	}
}

func (s *keyValueSet[K, V]) Selector(v V) K {
	return s.unsafe.Selector(v)
}

func (s *keyValueSet[K, V]) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Len()
}

func (s *keyValueSet[K, V]) Pop() (V, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.unsafe.Pop()
}

func (s *keyValueSet[K, V]) Remove(v ...V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.unsafe.Remove(v...)
}

func (s *keyValueSet[K, V]) SymmetricDifference(other KeyValueSet[K, V]) KeyValueSet[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.SymmetricDifference(other)
}

func (s *keyValueSet[K, V]) ToSlice() []V {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ToSlice()
}

func (s *keyValueSet[K, V]) Union(other KeyValueSet[K, V]) KeyValueSet[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Union(other)
}

var _ KeyValueSet[string, string] = &keyValueSet[string, string]{}
