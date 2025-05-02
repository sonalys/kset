package kset

import (
	"fmt"
	"iter"

	"golang.org/x/exp/constraints"
)

// KeyOnlySet defines a key-only set.
// K is the constraints.Ordered type used for the underlying map keys.
type KeyOnlySet[K constraints.Ordered] interface {
	Set[K]

	// Append upserts multiple elements to the set.
	// It returns the number of elements that were actually added (i.e., were not already present).
	// Example:
	//  s := NewPrimitive(1)
	//  count := s.Append(1, 2, 3) // count is 2
	Append(values ...K) int

	// Clone creates a shallow copy of the set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := s1.Clone() // s2 is {1, 2}, independent of s1
	//  s2.Add(3)
	//  // s1 is {1, 2}, s2 is {1, 2, 3}
	Clone() KeyOnlySet[K]

	// Difference returns a new set containing elements that are in the current set but not in the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  diff := s1.Difference(s2) // diff is {1, 2}
	Difference(other Set[K]) KeyOnlySet[K]

	// Intersect returns a new set containing elements that are common to both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  intersection := s1.Intersect(s2) // intersection is {3}
	Intersect(other Set[K]) KeyOnlySet[K]

	// Each executes the given function fn for each element in the set.
	// Iteration stops if fn returns false.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  sum := 0
	//  s.Each(func(v int) bool {
	//      sum += v
	//      return true // Continue iteration
	//  }) // sum will be 6
	Each(fn func(K) bool)

	// Iter returns an iterator (iter.Seq) over the elements of the set.
	// The order of iteration is not guaranteed.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  for v := range s.Iter() {
	//      fmt.Println(v) // Prints 1, 2, 3 in some order
	//  }
	Iter() iter.Seq[K]

	// Remove removes the specified elements from the set.
	// Example:
	//  s := NewPrimitive(1, 2, 3, 4)
	//  s.Remove(2, 4) // s is {1, 3}
	Remove(v ...K)

	// SymmetricDifference returns a new set containing elements that are in either the current set or the other set, but not both.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  symDiff := s1.SymmetricDifference(s2) // symDiff is {1, 2, 4, 5}
	SymmetricDifference(other KeyOnlySet[K]) KeyOnlySet[K]

	// Union returns a new set containing all elements from both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 3)
	//  union := s1.Union(s2) // union is {1, 2, 3}
	Union(other KeyOnlySet[K]) KeyOnlySet[K]

	// Pop removes and returns an arbitrary element from the set.
	// It returns the removed element and true if the set was not empty, otherwise it returns the zero value of V and false.
	// Example:
	//  s := NewPrimitive(1, 2)
	//  v, ok := s.Pop() // v could be 1 or 2, ok is true
	//  v, ok = s.Pop() // v is the remaining element, ok is true
	//  v, ok = s.Pop() // v is 0, ok is false
	Pop() (K, bool)

	// ToSlice returns a slice containing all elements of the set.
	// The order of elements in the slice is not guaranteed.
	// Example:
	//  s := NewPrimitive(3, 1, 2)
	//  slice := s.ToSlice() // slice could be []int{1, 2, 3}, []int{3, 1, 2}, etc.
	ToSlice() []K
}

// keySet is an implementation of KeySet.
// K is the key, must be ordered.
// S is just a generic type parameter for removing the store abstraction and accessing the implementation directly.
type keySet[K constraints.Ordered, S store[K, struct{}]] struct {
	store    S
	newStore func(len int) S
}

// NewKey creates a key-only set from any given slice.
// It requires a selector function that extracts a constraints.Ordered key K from a value V.
// Optionally, it can be initialized with one or more values.
// The returned set is not safe for concurrent use by multiple goroutines.
//
// Example:
//
//	// Create a set of User structs, using ID as the key.
//	set := kset.NewKey(kset.HashMap, func(u User) int { return u.ID }, user1, user2)
func NewKey[K constraints.Ordered](storeType StoreType, values ...K) KeyOnlySet[K] {
	switch storeType {
	case HashMap:
		return &keySet[K, *safeMapStore[K, struct{}]]{
			store: newStoreMapK(values...),
			newStore: func(len int) *safeMapStore[K, struct{}] {
				return newStoreMapK[K]()
			},
		}
	case HashMapUnsafe:
		return &keySet[K, *unsafeMapStore[K, struct{}]]{
			store: newStoreUnsafeMapK(values...),
			newStore: func(len int) *unsafeMapStore[K, struct{}] {
				return newStoreUnsafeMapK[K]()
			},
		}
	case TreeMap:
		return &keySet[K, *treeMapStore[K, struct{}]]{
			store: newStoreTreeMapK(values...),
			newStore: func(len int) *treeMapStore[K, struct{}] {
				return newStoreTreeMapK[K]()
			},
		}
	case TreeMapUnsafe:
		return &keySet[K, *unsafeTreeMapStore[K, struct{}]]{
			store: newStoreUnsafeTreeMapK(values...),
			newStore: func(len int) *unsafeTreeMapStore[K, struct{}] {
				return newStoreUnsafeTreeMapK[K]()
			},
		}
	default:
		panic(fmt.Sprintf("type not supported: %s", storeType))
	}
}

// Append adds keys to the set. Returns the number of new keys added.
func (k keySet[K, S]) Append(keys ...K) int {
	prevLen := k.store.Len()
	for _, key := range keys {
		k.store.Upsert(key, struct{}{})
	}
	return k.store.Len() - prevLen
}

// Len returns the number of keys in the set.
func (k keySet[K, S]) Len() int {
	return k.store.Len()
}

// Clear removes all keys from the set.
func (k keySet[K, S]) Clear() {
	k.store.Clear()
}

// Clone creates a copy of the set.
func (k keySet[K, S]) Clone() KeyOnlySet[K] {
	return keySet[K, S]{
		store: k.store.Clone().(S),
	}
}

// ContainsKeys checks if all specified keys are present in the set.
func (k keySet[K, S]) ContainsKeys(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.store.Get(key); !ok {
			return false
		}
	}
	return true
}

// ContainsAnyKey checks if any of the specified keys are present in the set.
func (k keySet[K, S]) ContainsAnyKey(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.store.Get(key); ok {
			return true
		}
	}
	return false
}

// Intersects checks if this set shares any keys with the other set.
func (k keySet[K, S]) Intersects(other Set[K]) bool {
	for key := range k.store.Iter() {
		if other.ContainsKeys(key) {
			return true
		}
	}
	return false
}

// Difference returns a new set with keys in this set but not in the other.
func (k keySet[K, S]) Difference(other Set[K]) KeyOnlySet[K] {
	diff := &keySet[K, S]{
		store:    k.newStore(k.Len()),
		newStore: k.newStore,
	}
	for key := range k.store.Iter() {
		if !other.ContainsKeys(key) {
			diff.Append(key)
		}
	}
	return diff
}

// Each executes a function for each key in the set until the function returns false.
func (k keySet[K, S]) Each(f func(K) bool) {
	for key := range k.store.Iter() {
		if !f(key) {
			break
		}
	}
}

// Equal checks if this set is equal to another set (contains the same keys).
func (k keySet[K, S]) Equal(other Set[K]) bool {
	if k.Len() != other.Len() {
		return false
	}
	for key := range k.store.Iter() {
		if !other.ContainsKeys(key) {
			return false
		}
	}
	return true
}

// Intersect returns a new set with keys common to both this set and the other.
func (k keySet[K, S]) Intersect(other Set[K]) KeyOnlySet[K] {
	intersection := &keySet[K, S]{
		store:    k.newStore(k.Len()),
		newStore: k.newStore,
	}

	for key := range k.store.Iter() {
		if other.ContainsKeys(key) {
			intersection.Append(key)
		}
	}

	return intersection
}

// IsEmpty checks if the set is empty.
func (k keySet[K, S]) IsEmpty() bool {
	return k.Len() == 0
}

// IsProperSubset checks if this set is a proper subset of the other set.
func (k keySet[K, S]) IsProperSubset(other Set[K]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

// IsProperSuperset checks if this set is a proper superset of the other set.
func (k keySet[K, S]) IsProperSuperset(other Set[K]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

// IsSubset checks if this set is a subset of the other set.
func (k keySet[K, S]) IsSubset(other Set[K]) bool {
	if k.Len() > other.Len() {
		return false
	}
	for key := range k.store.Iter() {
		if !other.ContainsKeys(key) {
			return false
		}
	}
	return true
}

// IsSuperset checks if this set is a superset of the other set.
func (k keySet[K, S]) IsSuperset(other Set[K]) bool {
	if k.Len() > other.Len() {
		return false
	}

	for key := range k.store.Iter() {
		if !other.ContainsKeys(key) {
			return false
		}
	}

	return true
}

// Iter returns an iterator for the keys in the set.
func (k keySet[K, S]) Iter() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range k.store.Iter() {
			if !yield(key) {
				break
			}
		}
	}
}

// Iter returns an iterator for the keys in the set.
func (k keySet[K, S]) IterKeys() iter.Seq[K] {
	return k.Iter()
}

// Pop removes and returns an arbitrary key from the set.
// The second return value indicates if a key was removed (true) or if the set was empty (false).
func (k keySet[K, S]) Pop() (K, bool) {
	for key := range k.store.Iter() {
		k.store.Delete(key)
		return key, true
	}
	var zero K
	return zero, false
}

// Remove removes the specified keys from the set.
func (k keySet[K, S]) Remove(keys ...K) {
	for _, key := range keys {
		k.store.Delete(key)
	}
}

// SymmetricDifference returns a new set with keys in either this set or the other, but not both.
func (k keySet[K, S]) SymmetricDifference(other KeyOnlySet[K]) KeyOnlySet[K] {
	sd := &keySet[K, S]{
		store:    k.newStore(k.Len()),
		newStore: k.newStore,
	}
	for key := range k.store.Iter() {
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
func (k keySet[K, S]) ToSlice() []K {
	result := make([]K, 0, k.Len())
	for key := range k.store.Iter() {
		result = append(result, key)
	}
	return result
}

// Union returns a new set with all keys from both this set and the other.
func (k keySet[K, S]) Union(other KeyOnlySet[K]) KeyOnlySet[K] {
	union := k.Clone()
	for key := range other.Iter() {
		union.Append(key)
	}
	return union
}

// Ensure unsafeKeySet implements KeySet at compile time.
var _ KeyOnlySet[string] = keySet[string, *treeMapStore[string, struct{}]]{}
