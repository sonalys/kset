package kset

import (
	"iter"
)

type unsafeKeyValueSet[K comparable, V any] struct {
	data     map[K]V
	selector func(V) K
}

// NewKeyValueUnsafe creates a new non-thread-safe set.
// It requires a selector function that extracts a comparable key K from a value V.
// Optionally, it can be initialized with one or more values.
// The returned set is *not* safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create an unsafe set of User structs, using ID as the key.
//	userSet := kset.NewKeyValueUnsafe(func(u User) int { return u.ID }, user1, user2)
func NewKeyValueUnsafe[K comparable, V any](selector func(V) K, values ...V) KeyValueSet[K, V] {
	return newKeyValueUnsafe(selector, values...)
}

func newKeyValueUnsafe[K comparable, V any](selector func(V) K, values ...V) *unsafeKeyValueSet[K, V] {
	data := make(map[K]V, len(values))

	for i := range values {
		value := values[i]
		key := selector(value)
		data[key] = value
	}

	return &unsafeKeyValueSet[K, V]{
		selector: selector,
		data:     data,
	}
}

func (k *unsafeKeyValueSet[K, V]) Selector(value V) K {
	return k.selector(value)
}

func (k *unsafeKeyValueSet[K, V]) Append(values ...V) int {
	prevLen := len(k.data)
	for _, val := range values {
		key := k.selector(val)
		k.data[key] = val
	}
	return len(k.data) - prevLen
}

func (k *unsafeKeyValueSet[K, V]) Len() int {
	return len(k.data)
}

func (k *unsafeKeyValueSet[K, V]) Clear() {
	for key := range k.data {
		delete(k.data, key)
	}
}

func (k *unsafeKeyValueSet[K, V]) Clone() KeyValueSet[K, V] {
	clonedSet := NewKeyValueUnsafe(k.selector)

	for value := range k.Iter() {
		clonedSet.Append(value)
	}

	return clonedSet
}

func (k *unsafeKeyValueSet[K, V]) Contains(values ...V) bool {
	for _, val := range values {
		key := k.selector(val)
		if _, ok := k.data[key]; !ok {
			return false
		}
	}
	return true
}

func (k *unsafeKeyValueSet[K, V]) ContainsKeys(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.data[key]; !ok {
			return false
		}
	}
	return true
}

func (k *unsafeKeyValueSet[K, V]) ContainsAny(values ...V) bool {
	for _, val := range values {
		key := k.selector(val)
		if _, ok := k.data[key]; ok {
			return true
		}
	}
	return false
}

func (k *unsafeKeyValueSet[K, V]) ContainsAnyKey(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.data[key]; ok {
			return true
		}
	}
	return false
}

func (k *unsafeKeyValueSet[K, V]) Intersects(other KeySet[K]) bool {
	for key := range k.data {
		if other.ContainsKeys(key) {
			return true
		}
	}

	return false
}

func (k *unsafeKeyValueSet[K, V]) Difference(other KeySet[K]) KeyValueSet[K, V] {
	diff := NewKeyValueUnsafe(k.selector)

	for key, value := range k.data {
		if !other.ContainsKeys(key) {
			diff.Append(value)
		}
	}

	return diff
}

func (k *unsafeKeyValueSet[K, V]) Each(f func(V) bool) {
	for _, elem := range k.data {
		if !f(elem) {
			break
		}
	}
}

func (k *unsafeKeyValueSet[K, V]) Equal(other KeySet[K]) bool {
	if k.Len() != other.Len() {
		return false
	}

	for key := range k.data {
		if !other.ContainsKeys(key) {
			return false
		}
	}

	return true
}

func (k *unsafeKeyValueSet[K, V]) Intersect(other KeySet[K]) KeyValueSet[K, V] {
	intersection := NewKeyValueUnsafe(k.selector)

	for key, value := range k.data {
		if other.ContainsKeys(key) {
			intersection.Append(value)
		}
	}
	return intersection
}

func (k *unsafeKeyValueSet[K, V]) IsEmpty() bool {
	return k.Len() == 0
}

func (k *unsafeKeyValueSet[K, V]) IsProperSubset(other KeySet[K]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

func (k *unsafeKeyValueSet[K, V]) IsProperSuperset(other KeySet[K]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

func (k *unsafeKeyValueSet[K, V]) IsSubset(other KeySet[K]) bool {
	if k.Len() > other.Len() {
		return false
	}

	for key := range k.data {
		if !other.ContainsKeys(key) {
			return false
		}
	}

	return true
}

func (k *unsafeKeyValueSet[K, V]) IsSuperset(other KeySet[K]) bool {
	return other.IsSubset(k)
}

func (k *unsafeKeyValueSet[K, V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, elem := range k.data {
			if !yield(elem) {
				break
			}
		}
	}
}

func (k *unsafeKeyValueSet[K, V]) IterKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range k.data {
			if !yield(key) {
				break
			}
		}
	}
}

func (k *unsafeKeyValueSet[K, V]) Pop() (V, bool) {
	for key, value := range k.data {
		delete(k.data, key)
		return value, true
	}

	var zero V
	return zero, false
}

func (k *unsafeKeyValueSet[K, V]) Remove(values ...V) {
	for _, val := range values {
		key := k.selector(val)
		delete(k.data, key)
	}
}

func (k *unsafeKeyValueSet[K, V]) SymmetricDifference(other KeyValueSet[K, V]) KeyValueSet[K, V] {
	sd := NewKeyValueUnsafe(k.selector)

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

func (k *unsafeKeyValueSet[K, V]) ToSlice() []V {
	result := make([]V, 0, k.Len())
	for _, elem := range k.data {
		result = append(result, elem)
	}
	return result
}

func (k *unsafeKeyValueSet[K, V]) Union(other KeyValueSet[K, V]) KeyValueSet[K, V] {
	union := k.Clone()

	for elem := range other.Iter() {
		union.Append(elem)
	}

	return union
}

var _ KeyValueSet[string, string] = &unsafeKeyValueSet[string, string]{}
