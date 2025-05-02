package kset

import (
	"fmt"
	"iter"

	"golang.org/x/exp/constraints"
)

// KeyOnlySet defines a key-only set.
// K is the comparable type used for the underlying map keys.
type KeyOnlySet[K comparable] interface {
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

// Set defines the interface of the behavior expected from only comparing keys, and not values.
// This interface is useful for comparing sets that shares the same key type, but not the same value.
type Set[K comparable] interface {
	// Len returns the number of elements in the set.
	// Example:
	//  s := NewPrimitive(1, 2)
	//  length := s.Len() // length is 2
	Len() int

	// Clear removes all elements from the set.
	// Example:
	//  s := NewPrimitive(1, 2)
	//  s.Clear() // s is {}
	//  length := s.Len() // length is 0
	Clear()

	// Contains checks if all specified elements are present in the set.
	// It returns true if all elements v are in the set, false otherwise.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  hasAll := s.ContainsKeys(1, 2) // hasAll is true
	//  hasAll = s.ContainsKeys(1, 4) // hasAll is false
	ContainsKeys(keys ...K) bool

	// ContainsAny checks if any of the specified elements are present in the set.
	// It returns true if at least one element v is in the set, false otherwise.
	// Example:
	//  s := NewPrimitive(1, 2)
	//  hasAny := s.ContainsAnyKey(2, 4) // hasAny is true
	//  hasAny = s.ContainsAnyKey(4, 5) // hasAny is false
	ContainsAnyKey(keys ...K) bool

	// IsEmpty checks if the set contains no elements.
	// Example:
	//  s := NewPrimitive[int]()
	//  s.IsEmpty() // empty is true
	//  s.Add(1)
	//  empty = s.IsEmpty() // empty is false
	IsEmpty() bool

	// IsProperSubset checks if the set is a proper subset of another set.
	// A proper subset is a subset that is not equal to the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(1, 2, 3)
	//  s3 := NewPrimitive(1, 2)
	//  isProper := s1.IsProperSubset(s2) // isProper is true
	//  isProper = s1.IsProperSubset(s3) // isProper is false
	IsProperSubset(other Set[K]) bool

	// IsProperSuperset checks if the set is a proper superset of another set.
	// A proper superset is a superset that is not equal to the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(1, 2)
	//  s3 := NewPrimitive(1, 2, 3)
	//  isProper := s1.IsProperSuperset(s2) // isProper is true
	//  isProper = s1.IsProperSuperset(s3) // isProper is false
	IsProperSuperset(other Set[K]) bool

	// IsSubset checks if the set is a subset of another set (i.e., all elements of the current set are also in the other set).
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(1, 2, 3)
	//  s3 := NewPrimitive(1, 3)
	//  isSub := s1.IsSubset(s2) // isSub is true
	//  isSub = s1.IsSubset(s3) // isSub is false
	IsSubset(other Set[K]) bool

	// IsSuperset checks if the set is a superset of another set (i.e., all elements of the other set are also in the current set).
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(1, 2)
	//  s3 := NewPrimitive(1, 4)
	//  isSuper := s1.IsSuperset(s2) // isSuper is true
	//  isSuper = s1.IsSuperset(s3) // isSuper is false
	IsSuperset(other Set[K]) bool

	// Intersects checks if the set has at least one element in common with another set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 3)
	//  s3 := NewPrimitive(4, 5)
	//  intersects := s1.Intersects(s2) // intersects is true
	//  intersects = s1.Intersects(s3) // intersects is false
	Intersects(other Set[K]) bool

	// Equal checks if the set is equal to another set (i.e., contains the same elements).
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 1)
	//  s3 := NewPrimitive(1, 3)
	//  isEqual := s1.Equal(s2) // isEqual is true
	//  isEqual = s1.Equal(s3) // isEqual is false
	Equal(other Set[K]) bool

	IterKeys() iter.Seq[K]
}

// keySet is an implementation of KeySet.
type keySet[K comparable] struct {
	store    Store[K, struct{}]
	newStore func(len int) Store[K, struct{}]
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
func NewKeySet[K constraints.Ordered](storeType StoreType, values ...K) KeyOnlySet[K] {
	var m Store[K, struct{}]

	switch storeType {
	case HashMap:
		m = NewStoreMapKey(values...)
	case HashMapUnsafe:
		m = NewUnsafeStoreMapKey(values...)
	case TreeMap:
		m = NewStoreTreeMapKey(values...)
	case TreeMapUnsafe:
		m = NewUnsafeStoreTreeMapKey(values...)
	default:
		panic(fmt.Sprintf("type not supported: %s", storeType))
	}

	return &keySet[K]{
		store: m,
		newStore: func(len int) Store[K, struct{}] {
			return &safeMapStore[K, struct{}]{
				store: make(map[K]struct{}, len),
			}
		},
	}
}

// Append adds keys to the set. Returns the number of new keys added.
func (k keySet[K]) Append(keys ...K) int {
	prevLen := k.store.Len()
	for _, key := range keys {
		k.store.Upsert(key, struct{}{})
	}
	return k.store.Len() - prevLen
}

// Len returns the number of keys in the set.
func (k keySet[K]) Len() int {
	return k.store.Len()
}

// Clear removes all keys from the set.
func (k keySet[K]) Clear() {
	k.store.Clear()
}

// Clone creates a copy of the set.
func (k keySet[K]) Clone() KeyOnlySet[K] {
	return keySet[K]{
		store: k.store.Clone(),
	}
}

// ContainsKeys checks if all specified keys are present in the set.
func (k keySet[K]) ContainsKeys(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.store.Get(key); !ok {
			return false
		}
	}
	return true
}

// ContainsAnyKey checks if any of the specified keys are present in the set.
func (k keySet[K]) ContainsAnyKey(keys ...K) bool {
	for _, key := range keys {
		if _, ok := k.store.Get(key); ok {
			return true
		}
	}
	return false
}

// Intersects checks if this set shares any keys with the other set.
func (k keySet[K]) Intersects(other Set[K]) bool {
	for key := range k.store.Iter() {
		if other.ContainsKeys(key) {
			return true
		}
	}
	return false
}

// Difference returns a new set with keys in this set but not in the other.
func (k keySet[K]) Difference(other Set[K]) KeyOnlySet[K] {
	diff := &keySet[K]{
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
func (k keySet[K]) Each(f func(K) bool) {
	for key := range k.store.Iter() {
		if !f(key) {
			break
		}
	}
}

// Equal checks if this set is equal to another set (contains the same keys).
func (k keySet[K]) Equal(other Set[K]) bool {
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
func (k keySet[K]) Intersect(other Set[K]) KeyOnlySet[K] {
	intersection := &keySet[K]{
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
func (k keySet[K]) IsEmpty() bool {
	return k.Len() == 0
}

// IsProperSubset checks if this set is a proper subset of the other set.
func (k keySet[K]) IsProperSubset(other Set[K]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

// IsProperSuperset checks if this set is a proper superset of the other set.
func (k keySet[K]) IsProperSuperset(other Set[K]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

// IsSubset checks if this set is a subset of the other set.
func (k keySet[K]) IsSubset(other Set[K]) bool {
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
func (k keySet[K]) IsSuperset(other Set[K]) bool {
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
func (k keySet[K]) Iter() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range k.store.Iter() {
			if !yield(key) {
				break
			}
		}
	}
}

// Iter returns an iterator for the keys in the set.
func (k keySet[K]) IterKeys() iter.Seq[K] {
	return k.Iter()
}

// Pop removes and returns an arbitrary key from the set.
// The second return value indicates if a key was removed (true) or if the set was empty (false).
func (k keySet[K]) Pop() (K, bool) {
	for key := range k.store.Iter() {
		k.store.Delete(key)
		return key, true
	}
	var zero K
	return zero, false
}

// Remove removes the specified keys from the set.
func (k keySet[K]) Remove(keys ...K) {
	for _, key := range keys {
		k.store.Delete(key)
	}
}

// SymmetricDifference returns a new set with keys in either this set or the other, but not both.
func (k keySet[K]) SymmetricDifference(other KeyOnlySet[K]) KeyOnlySet[K] {
	sd := &keySet[K]{
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
func (k keySet[K]) ToSlice() []K {
	result := make([]K, 0, k.Len())
	for key := range k.store.Iter() {
		result = append(result, key)
	}
	return result
}

// Union returns a new set with all keys from both this set and the other.
func (k keySet[K]) Union(other KeyOnlySet[K]) KeyOnlySet[K] {
	union := k.Clone()
	for key := range other.Iter() {
		union.Append(key)
	}
	return union
}

// Ensure unsafeKeySet implements KeySet at compile time.
var _ KeyOnlySet[string] = keySet[string]{}
