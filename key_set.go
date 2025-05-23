package kset

import (
	"iter"
)

// KeySet defines a key-only set.
type KeySet[Key any] interface {
	Set[Key]

	// Append upserts multiple elements to the set.
	// It returns the number of elements that were actually added (i.e., were not already present).
	// Example:
	//  s := NewPrimitive(1)
	//  count := s.Append(1, 2, 3) // count is 2
	Append(values ...Key) int

	// Clone creates a shallow copy of the set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := s1.Clone() // s2 is {1, 2}, independent of s1
	//  s2.Add(3)
	//  // s1 is {1, 2}, s2 is {1, 2, 3}
	Clone() KeySet[Key]

	// Difference returns a new set containing elements that are in the current set but not in the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  diff := s1.Difference(s2) // diff is {1, 2}
	Difference(other Set[Key]) KeySet[Key]

	// Intersect returns a new set containing elements that are common to both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  intersection := s1.Intersect(s2) // intersection is {3}
	Intersect(other Set[Key]) KeySet[Key]

	// Each executes the given function fn for each element in the set.
	// Iteration stops if fn returns false.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  sum := 0
	//  s.Each(func(v int) bool {
	//      sum += v
	//      return true // Continue iteration
	//  }) // sum will be 6
	Each(fn func(Key) bool)

	// Iter returns an iterator (iter.Seq) over the elements of the set.
	// The order of iteration is not guaranteed.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  for v := range s.Iter() {
	//      fmt.Println(v) // Prints 1, 2, 3 in some order
	//  }
	Iter() iter.Seq[Key]

	// Remove removes the specified elements from the set.
	// Example:
	//  s := NewPrimitive(1, 2, 3, 4)
	//  s.Remove(2, 4) // s is {1, 3}
	Remove(v ...Key)

	// SymmetricDifference returns a new set containing elements that are in either the current set or the other set, but not both.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  symDiff := s1.SymmetricDifference(s2) // symDiff is {1, 2, 4, 5}
	SymmetricDifference(other KeySet[Key]) KeySet[Key]

	// Union returns a new set containing all elements from both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 3)
	//  union := s1.Union(s2) // union is {1, 2, 3}
	Union(other KeySet[Key]) KeySet[Key]

	// Pop removes and returns an arbitrary element from the set.
	// It returns the removed element and true if the set was not empty, otherwise it returns the zero value of V and false.
	// Example:
	//  s := NewPrimitive(1, 2)
	//  v, ok := s.Pop() // v could be 1 or 2, ok is true
	//  v, ok = s.Pop() // v is the remaining element, ok is true
	//  v, ok = s.Pop() // v is 0, ok is false
	Pop() (Key, bool)

	// ToSlice returns a slice containing all elements of the set.
	// The order of elements in the slice is not guaranteed.
	// Example:
	//  s := NewPrimitive(3, 1, 2)
	//  slice := s.ToSlice() // slice could be []int{1, 2, 3}, []int{3, 1, 2}, etc.
	ToSlice() []Key
}

// keySet is an implementation of KeySet.
// K is the key, must be ordered.
// S is just a generic type parameter for removing the store abstraction and accessing the implementation directly.
type keySet[Key any, Store Storage[Key, empty]] struct {
	store    Store
	newStore func(...Key) Store
}

// Append adds keys to the set. Returns the number of new keys added.
func (k keySet[Key, Store]) Append(keys ...Key) int {
	prevLen := k.store.Len()
	for _, key := range keys {
		k.store.Upsert(key, empty{})
	}
	return k.store.Len() - prevLen
}

// Len returns the number of keys in the set.
func (k keySet[Key, Store]) Len() int {
	return k.store.Len()
}

// Clear removes all keys from the set.
func (k keySet[Key, Store]) Clear() {
	k.store.Clear()
}

// Clone creates a copy of the set.
func (k keySet[Key, Store]) Clone() KeySet[Key] {
	return keySet[Key, Store]{
		store: k.store.Clone().(Store),
	}
}

// ContainsKeys checks if all specified keys are present in the set.
func (k keySet[Key, Store]) ContainsKeys(keys ...Key) bool {
	for _, key := range keys {
		if _, ok := k.store.Get(key); !ok {
			return false
		}
	}
	return true
}

// ContainsAnyKey checks if any of the specified keys are present in the set.
func (k keySet[Key, Store]) ContainsAnyKey(keys ...Key) bool {
	for _, key := range keys {
		if _, ok := k.store.Get(key); ok {
			return true
		}
	}
	return false
}

// Intersects checks if this set shares any keys with the other set.
func (k keySet[Key, Store]) Intersects(other Set[Key]) bool {
	for key := range k.store.Iter() {
		if other.ContainsKeys(key) {
			return true
		}
	}
	return false
}

// Difference returns a new set with keys in this set but not in the other.
func (k keySet[Key, Store]) Difference(other Set[Key]) KeySet[Key] {
	diff := &keySet[Key, Store]{
		store:    k.newStore(),
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
func (k keySet[Key, Store]) Each(f func(Key) bool) {
	for key := range k.store.Iter() {
		if !f(key) {
			break
		}
	}
}

// Equal checks if this set is equal to another set (contains the same keys).
func (k keySet[Key, Store]) Equal(other Set[Key]) bool {
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
func (k keySet[Key, Store]) Intersect(other Set[Key]) KeySet[Key] {
	intersection := &keySet[Key, Store]{
		store:    k.newStore(),
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
func (k keySet[Key, Store]) IsEmpty() bool {
	return k.Len() == 0
}

// IsProperSubset checks if this set is a proper subset of the other set.
func (k keySet[Key, Store]) IsProperSubset(other Set[Key]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

// IsProperSuperset checks if this set is a proper superset of the other set.
func (k keySet[Key, Store]) IsProperSuperset(other Set[Key]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

// IsSubset checks if this set is a subset of the other set.
func (k keySet[Key, Store]) IsSubset(other Set[Key]) bool {
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
func (k keySet[Key, Store]) IsSuperset(other Set[Key]) bool {
	return other.IsSubset(k)
}

// Iter returns an iterator for the keys in the set.
func (k keySet[Key, Store]) Iter() iter.Seq[Key] {
	return func(yield func(Key) bool) {
		for key := range k.store.Iter() {
			if !yield(key) {
				break
			}
		}
	}
}

// Iter returns an iterator for the keys in the set.
func (k keySet[Key, Store]) IterKeys() iter.Seq[Key] {
	return k.Iter()
}

// Pop removes and returns an arbitrary key from the set.
// The second return Store indicates if a key was removed (true) or if the set was empty (false).
func (k keySet[Key, Store]) Pop() (Key, bool) {
	for key := range k.store.Iter() {
		defer k.store.Delete(key)
		return key, true
	}
	var zero Key
	return zero, false
}

// Remove removes the specified keys from the set.
func (k keySet[Key, Store]) Remove(keys ...Key) {
	for _, key := range keys {
		k.store.Delete(key)
	}
}

// SymmetricDifference returns a new set with keys in either this set or the other, but not both.
func (k keySet[Key, Store]) SymmetricDifference(other KeySet[Key]) KeySet[Key] {
	sd := &keySet[Key, Store]{
		store:    k.newStore(),
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
func (k keySet[Key, Store]) ToSlice() []Key {
	result := make([]Key, 0, k.Len())
	for key := range k.store.Iter() {
		result = append(result, key)
	}
	return result
}

// Union returns a new set with all keys from both this set and the other.
func (k keySet[Key, Store]) Union(other KeySet[Key]) KeySet[Key] {
	union := k.Clone()
	for key := range other.Iter() {
		union.Append(key)
	}
	return union
}

// Ensure unsafeKeySet implements KeySet at compile time.
var _ KeySet[string] = keySet[string, *treeMapStore[string, empty]]{}
