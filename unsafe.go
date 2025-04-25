package kset

import (
	"iter"
)

type unsafeSet[K comparable, V any] struct {
	data     map[K]V
	selector func(V) K
}

// NewUnsafe creates a new non-thread-safe set.
// It requires a selector function that extracts a comparable key K from a value V.
// Optionally, it can be initialized with one or more values.
// The returned set is *not* safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create an unsafe set of User structs, using ID as the key.
//	userSet := kset.NewUnsafe(func(u User) int { return u.ID }, user1, user2)
func NewUnsafe[K comparable, V any](selector func(V) K, values ...V) Set[K, V] {
	return newUnsafe(selector, values...)
}

// NewUnsafePrimitive creates a new non-thread-safe set for primitive types
// (or types where the value itself is the key).
// It uses an identity function (func(k K) K { return k }) as the selector.
// Optionally, it can be initialized with one or more values.
// The returned set is *not* safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create an unsafe set of integers.
//	intSet := kset.NewUnsafePrimitive(1, 2, 3, 2) // Resulting set: {1, 2, 3}
func NewUnsafePrimitive[K comparable](values ...K) Set[K, K] {
	return newUnsafePrimitive(values...)
}

func newUnsafePrimitive[K comparable](values ...K) *unsafeSet[K, K] {
	return newUnsafe(func(k K) K { return k }, values...)
}

func newUnsafe[K comparable, V any](selector func(V) K, values ...V) *unsafeSet[K, V] {
	data := make(map[K]V, len(values))

	for i := range values {
		value := values[i]
		key := selector(value)
		data[key] = value
	}

	return &unsafeSet[K, V]{
		selector: selector,
		data:     data,
	}
}

func (k *unsafeSet[K, V]) Selector(value V) K {
	return k.selector(value)
}

func (k *unsafeSet[K, V]) Append(values ...V) int {
	prevLen := len(k.data)
	for _, val := range values {
		key := k.selector(val)
		k.data[key] = val
	}
	return len(k.data) - prevLen
}

func (k *unsafeSet[K, V]) Len() int {
	return len(k.data)
}

func (k *unsafeSet[K, V]) Clear() {
	for key := range k.data {
		delete(k.data, key)
	}
}

func (k *unsafeSet[K, V]) Clone() Set[K, V] {
	clonedSet := NewUnsafe(k.selector)

	for value := range k.Iter() {
		clonedSet.Append(value)
	}

	return clonedSet
}

func (k *unsafeSet[K, V]) Contains(values ...V) bool {
	for _, val := range values {
		key := k.selector(val)
		if _, ok := k.data[key]; !ok {
			return false
		}
	}
	return true
}

func (k *unsafeSet[K, V]) ContainsKeys(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.data[key]; !ok {
			return false
		}
	}
	return true
}

func (k *unsafeSet[K, V]) ContainsAny(values ...V) bool {
	for _, val := range values {
		key := k.selector(val)
		if _, ok := k.data[key]; ok {
			return true
		}
	}
	return false
}

func (k *unsafeSet[K, V]) ContainsAnyKey(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.data[key]; ok {
			return true
		}
	}
	return false
}

func (k *unsafeSet[K, V]) Intersects(other Set[K, V]) bool {
	if k.Len() < other.Len() {
		for _, elem := range k.data {
			if !other.Contains(elem) {
				return true
			}
		}
		return false
	}

	for elem := range other.Iter() {
		if k.Contains(elem) {
			return true
		}
	}

	return false
}

func (k *unsafeSet[K, V]) Difference(other Set[K, V]) Set[K, V] {
	diff := NewUnsafe(k.selector)

	for _, elem := range k.data {
		if !other.Contains(elem) {
			diff.Append(elem)
		}
	}

	return diff
}

func (k *unsafeSet[K, V]) Each(f func(V) bool) {
	for _, elem := range k.data {
		if !f(elem) {
			break
		}
	}
}

func (k *unsafeSet[K, V]) Equal(other Set[K, V]) bool {
	if k.Len() != other.Len() {
		return false
	}

	for _, elem := range k.data {
		if !other.Contains(elem) {
			return false
		}
	}

	return true
}

func (k *unsafeSet[K, V]) Intersect(other Set[K, V]) Set[K, V] {
	intersection := NewUnsafe(k.selector)

	if k.Len() < other.Len() {
		for _, elem := range k.data {
			if other.Contains(elem) {
				intersection.Append(elem)
			}
		}
		return intersection
	}

	for elem := range other.Iter() {
		if k.Contains(elem) {
			intersection.Append(elem)
		}
	}
	return intersection
}

func (k *unsafeSet[K, V]) IsEmpty() bool {
	return k.Len() == 0
}

func (k *unsafeSet[K, V]) IsProperSubset(other Set[K, V]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

func (k *unsafeSet[K, V]) IsProperSuperset(other Set[K, V]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

func (k *unsafeSet[K, V]) IsSubset(other Set[K, V]) bool {
	if k.Len() > other.Len() {
		return false
	}

	for _, elem := range k.data {
		if !other.Contains(elem) {
			return false
		}
	}

	return true
}

func (k *unsafeSet[K, V]) IsSuperset(other Set[K, V]) bool {
	return other.IsSubset(k)
}

func (k *unsafeSet[K, V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, elem := range k.data {
			if !yield(elem) {
				break
			}
		}
	}
}

func (k *unsafeSet[K, V]) Pop() (V, bool) {
	for key, value := range k.data {
		delete(k.data, key)
		return value, true
	}

	var zero V
	return zero, false
}

func (k *unsafeSet[K, V]) Remove(values ...V) {
	for _, val := range values {
		key := k.selector(val)
		delete(k.data, key)
	}
}

func (k *unsafeSet[K, V]) SymmetricDifference(other Set[K, V]) Set[K, V] {
	sd := NewUnsafe(k.selector)

	for _, elem := range k.data {
		if !other.Contains(elem) {
			sd.Append(elem)
		}
	}

	for elem := range other.Iter() {
		if !k.Contains(elem) {
			sd.Append(elem)
		}
	}

	return sd
}

func (k *unsafeSet[K, V]) ToSlice() []V {
	result := make([]V, 0, k.Len())
	for _, elem := range k.data {
		result = append(result, elem)
	}
	return result
}

func (k *unsafeSet[K, V]) Union(other Set[K, V]) Set[K, V] {
	union := k.Clone()

	for elem := range other.Iter() {
		union.Append(elem)
	}

	return union
}

var _ Set[string, string] = &unsafeSet[string, string]{}
