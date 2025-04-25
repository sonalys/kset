package kset

import (
	"iter"
	"sync"
)

type safeSet[K comparable, V any] struct {
	lock   sync.RWMutex
	unsafe *unsafeSet[K, V]
}

func New[K comparable, V any](selector func(V) K, values ...V) Set[K, V] {
	return &safeSet[K, V]{
		unsafe: newUnsafe(selector, values...),
	}
}

func NewPrimitive[K comparable](values ...K) Set[K, K] {
	return &safeSet[K, K]{
		unsafe: newUnsafePrimitive(values...),
	}
}

func (s *safeSet[K, V]) Add(v V) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.unsafe.Add(v)
}

func (s *safeSet[K, V]) Append(v ...V) int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.unsafe.Append(v...)
}

func (s *safeSet[K, V]) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.unsafe.Clear()
}

func (s *safeSet[K, V]) Clone() Set[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Clone()
}

func (s *safeSet[K, V]) Contains(v ...V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Contains(v...)
}

func (s *safeSet[K, V]) ContainsKeys(keys ...K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ContainsKeys(keys...)
}

func (s *safeSet[K, V]) ContainsAnyKey(keys ...K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ContainsAnyKey(keys...)
}

func (s *safeSet[K, V]) ContainsAny(v ...V) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ContainsAny(v...)
}

func (s *safeSet[K, V]) Intersects(other Set[K, V]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Intersects(other)
}

func (s *safeSet[K, V]) Difference(other Set[K, V]) Set[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Difference(other)
}

func (s *safeSet[K, V]) Each(f func(V) bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	s.unsafe.Each(f)
}

func (s *safeSet[K, V]) Equal(other Set[K, V]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Equal(other)
}

func (s *safeSet[K, V]) Intersect(other Set[K, V]) Set[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Intersect(other)
}

func (s *safeSet[K, V]) IsEmpty() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsEmpty()
}

func (s *safeSet[K, V]) IsProperSubset(other Set[K, V]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsProperSubset(other)
}

func (s *safeSet[K, V]) IsProperSuperset(other Set[K, V]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsProperSuperset(other)
}

func (s *safeSet[K, V]) IsSubset(other Set[K, V]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsSubset(other)
}

func (s *safeSet[K, V]) IsSuperset(other Set[K, V]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.IsSuperset(other)
}

func (s *safeSet[K, V]) Iter() iter.Seq[V] {
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

func (s *safeSet[K, V]) Selector(v V) K {
	return s.unsafe.Selector(v)
}

func (s *safeSet[K, V]) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Len()
}

func (s *safeSet[K, V]) Pop() (V, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.unsafe.Pop()
}

func (s *safeSet[K, V]) Remove(v ...V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.unsafe.Remove(v...)
}

func (s *safeSet[K, V]) SymmetricDifference(other Set[K, V]) Set[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.SymmetricDifference(other)
}

func (s *safeSet[K, V]) ToSlice() []V {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.ToSlice()
}

func (s *safeSet[K, V]) Union(other Set[K, V]) Set[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.unsafe.Union(other)
}

var _ Set[string, string] = &safeSet[string, string]{}
