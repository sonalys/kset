package kset

import (
	"iter"
)

// unsafeKeySet is a non-thread-safe implementation of KeySet using a map.
type unsafeKeySet[K comparable] map[K]struct{}

// NewKeySetUnsafe creates a new non-thread-safe key-only set.
// Optionally, it can be initialized with one or more keys.
// The returned set is *not* safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create an unsafe set of integers.
//	intSet := kset.NewKeySetUnsafe(1, 2, 3, 2) // Resulting set: {1, 2, 3}
func NewKeySetUnsafe[K comparable](keys ...K) KeyOnlySet[K] {
	return newKeySetUnsafe(keys...)
}

// NewNewFromUnsaferom creates a new non-thread-safe key-only set from any given slice.
// It requires a selector function that extracts a comparable key K from a value V.
// Optionally, it can be initialized with one or more values.
// The returned set is not safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create a set of User structs, using ID as the key.
//	set := kset.NewFromUnsafe(func(u User) int { return u.ID }, user1, user2)
func NewFromUnsafe[K comparable, V any](selector func(V) K, values ...V) KeyOnlySet[K] {
	keys := make([]K, 0, len(values))

	for i := range values {
		keys = append(keys, selector(values[i]))
	}

	return newKeySetUnsafe(keys...)
}

func newKeySetUnsafe[K comparable](keys ...K) *unsafeKeySet[K] {
	set := make(unsafeKeySet[K], len(keys))
	set.Append(keys...)
	return &set
}

// Append adds keys to the set. Returns the number of new keys added.
func (k unsafeKeySet[K]) Append(keys ...K) int {
	prevLen := len(k)
	for _, key := range keys {
		k[key] = struct{}{}
	}
	return len(k) - prevLen
}

// Len returns the number of keys in the set.
func (k unsafeKeySet[K]) Len() int {
	return len(k)
}

// Clear removes all keys from the set.
func (k unsafeKeySet[K]) Clear() {
	for key := range k {
		delete(k, key)
	}
}

// Clone creates a copy of the set.
func (k unsafeKeySet[K]) Clone() KeyOnlySet[K] {
	clonedSet := make(unsafeKeySet[K], len(k))
	for key := range k {
		clonedSet[key] = struct{}{}
	}
	return clonedSet
}

// ContainsKeys checks if all specified keys are present in the set.
func (k unsafeKeySet[K]) ContainsKeys(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k[key]; !ok {
			return false
		}
	}
	return true
}

// ContainsAnyKey checks if any of the specified keys are present in the set.
func (k unsafeKeySet[K]) ContainsAnyKey(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k[key]; ok {
			return true
		}
	}
	return false
}

// Intersects checks if this set shares any keys with the other set.
func (k unsafeKeySet[K]) Intersects(other KeySet[K]) bool {
	for key := range k {
		if other.ContainsKeys(key) {
			return true
		}
	}
	return false
}

// Difference returns a new set with keys in this set but not in the other.
func (k unsafeKeySet[K]) Difference(other KeySet[K]) KeyOnlySet[K] {
	diff := NewKeySetUnsafe[K]()
	for key := range k {
		if !other.ContainsKeys(key) {
			diff.Append(key)
		}
	}
	return diff
}

// Each executes a function for each key in the set until the function returns false.
func (k unsafeKeySet[K]) Each(f func(K) bool) {
	for key := range k {
		if !f(key) {
			break
		}
	}
}

// Equal checks if this set is equal to another set (contains the same keys).
func (k unsafeKeySet[K]) Equal(other KeySet[K]) bool {
	if k.Len() != other.Len() {
		return false
	}
	for key := range k {
		if !other.ContainsKeys(key) {
			return false
		}
	}
	return true
}

// Intersect returns a new set with keys common to both this set and the other.
func (k unsafeKeySet[K]) Intersect(other KeySet[K]) KeyOnlySet[K] {
	intersection := NewKeySetUnsafe[K]()

	for key := range k {
		if other.ContainsKeys(key) {
			intersection.Append(key)
		}
	}

	return intersection
}

// IsEmpty checks if the set is empty.
func (k unsafeKeySet[K]) IsEmpty() bool {
	return k.Len() == 0
}

// IsProperSubset checks if this set is a proper subset of the other set.
func (k unsafeKeySet[K]) IsProperSubset(other KeySet[K]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

// IsProperSuperset checks if this set is a proper superset of the other set.
func (k unsafeKeySet[K]) IsProperSuperset(other KeySet[K]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

// IsSubset checks if this set is a subset of the other set.
func (k unsafeKeySet[K]) IsSubset(other KeySet[K]) bool {
	if k.Len() > other.Len() {
		return false
	}
	for key := range k {
		if !other.ContainsKeys(key) {
			return false
		}
	}
	return true
}

// IsSuperset checks if this set is a superset of the other set.
func (k unsafeKeySet[K]) IsSuperset(other KeySet[K]) bool {
	if k.Len() > other.Len() {
		return false
	}

	for key := range k {
		if !other.ContainsKeys(key) {
			return false
		}
	}

	return true
}

// Iter returns an iterator for the keys in the set.
func (k unsafeKeySet[K]) Iter() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range k {
			if !yield(key) {
				break
			}
		}
	}
}

// Iter returns an iterator for the keys in the set.
func (k unsafeKeySet[K]) IterKeys() iter.Seq[K] {
	return k.Iter()
}

// Pop removes and returns an arbitrary key from the set.
// The second return value indicates if a key was removed (true) or if the set was empty (false).
func (k unsafeKeySet[K]) Pop() (K, bool) {
	for key := range k {
		delete(k, key)
		return key, true
	}
	var zero K
	return zero, false
}

// Remove removes the specified keys from the set.
func (k unsafeKeySet[K]) Remove(keys ...K) {
	for _, key := range keys {
		delete(k, key)
	}
}

// SymmetricDifference returns a new set with keys in either this set or the other, but not both.
func (k unsafeKeySet[K]) SymmetricDifference(other KeyOnlySet[K]) KeyOnlySet[K] {
	sd := NewKeySetUnsafe[K]()
	for key := range k {
		if !other.ContainsKeys(key) {
			sd.Append(key)
		}
	}
	for key := range other.Iter() {
		if !k.ContainsKeys(key) {
			sd.Append(key)
		}
	}
	return sd
}

// ToSlice returns a slice containing all the keys in the set. The order is not guaranteed.
func (k unsafeKeySet[K]) ToSlice() []K {
	result := make([]K, 0, k.Len())
	for key := range k {
		result = append(result, key)
	}
	return result
}

// Union returns a new set with all keys from both this set and the other.
func (k unsafeKeySet[K]) Union(other KeyOnlySet[K]) KeyOnlySet[K] {
	union := k.Clone()
	for key := range other.Iter() {
		union.Append(key)
	}
	return union
}

// Ensure unsafeKeySet implements KeySet at compile time.
var _ KeyOnlySet[string] = unsafeKeySet[string]{}
