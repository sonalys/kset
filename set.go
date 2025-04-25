package kset

import "iter"

// KeySet defines the interface of the behavior expected from only comparing keys, and not values.
// This interface is useful for comparing sets that shares the same key type, but not the same value.
type KeySet[K comparable] interface {
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
	IsProperSubset(other KeySet[K]) bool

	// IsProperSuperset checks if the set is a proper superset of another set.
	// A proper superset is a superset that is not equal to the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(1, 2)
	//  s3 := NewPrimitive(1, 2, 3)
	//  isProper := s1.IsProperSuperset(s2) // isProper is true
	//  isProper = s1.IsProperSuperset(s3) // isProper is false
	IsProperSuperset(other KeySet[K]) bool

	// IsSubset checks if the set is a subset of another set (i.e., all elements of the current set are also in the other set).
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(1, 2, 3)
	//  s3 := NewPrimitive(1, 3)
	//  isSub := s1.IsSubset(s2) // isSub is true
	//  isSub = s1.IsSubset(s3) // isSub is false
	IsSubset(other KeySet[K]) bool

	// IsSuperset checks if the set is a superset of another set (i.e., all elements of the other set are also in the current set).
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(1, 2)
	//  s3 := NewPrimitive(1, 4)
	//  isSuper := s1.IsSuperset(s2) // isSuper is true
	//  isSuper = s1.IsSuperset(s3) // isSuper is false
	IsSuperset(other KeySet[K]) bool

	// Intersects checks if the set has at least one element in common with another set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 3)
	//  s3 := NewPrimitive(4, 5)
	//  intersects := s1.Intersects(s2) // intersects is true
	//  intersects = s1.Intersects(s3) // intersects is false
	Intersects(other KeySet[K]) bool

	// Equal checks if the set is equal to another set (i.e., contains the same elements).
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 1)
	//  s3 := NewPrimitive(1, 3)
	//  isEqual := s1.Equal(s2) // isEqual is true
	//  isEqual = s1.Equal(s3) // isEqual is false
	Equal(other KeySet[K]) bool
}

// Set defines the interface for a generic set data structure.
// K is the comparable type used for the underlying map keys.
// V is the type of the elements stored in the set.
type Set[K comparable, V any] interface {
	KeySet[K]

	// Append upserts multiple elements to the set.
	// It returns the number of elements that were actually added (i.e., were not already present).
	// Example:
	//  s := NewPrimitive(1)
	//  count := s.Append(1, 2, 3) // count is 2
	Append(values ...V) int

	// Clone creates a shallow copy of the set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := s1.Clone() // s2 is {1, 2}, independent of s1
	//  s2.Add(3)
	//  // s1 is {1, 2}, s2 is {1, 2, 3}
	Clone() Set[K, V]

	// Contains checks if all specified elements are present in the set.
	// It returns true if all elements v are in the set, false otherwise.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  hasAll := s.Contains(1, 2) // hasAll is true
	//  hasAll = s.Contains(1, 4) // hasAll is false
	Contains(values ...V) bool

	// ContainsAny checks if any of the specified elements are present in the set.
	// It returns true if at least one element v is in the set, false otherwise.
	// Example:
	//  s := NewPrimitive(1, 2)
	//  hasAny := s.ContainsAny(2, 4) // hasAny is true
	//  hasAny = s.ContainsAny(4, 5) // hasAny is false
	ContainsAny(values ...V) bool

	// Difference returns a new set containing elements that are in the current set but not in the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  diff := s1.Difference(s2) // diff is {1, 2}
	Difference(other KeySet[K]) Set[K, V]

	// Intersect returns a new set containing elements that are common to both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  intersection := s1.Intersect(s2) // intersection is {3}
	Intersect(other KeySet[K]) Set[K, V]

	// Each executes the given function fn for each element in the set.
	// Iteration stops if fn returns false.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  sum := 0
	//  s.Each(func(v int) bool {
	//      sum += v
	//      return true // Continue iteration
	//  }) // sum will be 6
	Each(fn func(V) bool)

	// Iter returns an iterator (iter.Seq) over the elements of the set.
	// The order of iteration is not guaranteed.
	// Example:
	//  s := NewPrimitive(1, 2, 3)
	//  for v := range s.Iter() {
	//      fmt.Println(v) // Prints 1, 2, 3 in some order
	//  }
	Iter() iter.Seq[V]

	// Remove removes the specified elements from the set.
	// Example:
	//  s := NewPrimitive(1, 2, 3, 4)
	//  s.Remove(2, 4) // s is {1, 3}
	Remove(v ...V)

	// SymmetricDifference returns a new set containing elements that are in either the current set or the other set, but not both.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  symDiff := s1.SymmetricDifference(s2) // symDiff is {1, 2, 4, 5}
	SymmetricDifference(other Set[K, V]) Set[K, V]

	// Union returns a new set containing all elements from both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 3)
	//  union := s1.Union(s2) // union is {1, 2, 3}
	Union(other Set[K, V]) Set[K, V]

	// Pop removes and returns an arbitrary element from the set.
	// It returns the removed element and true if the set was not empty, otherwise it returns the zero value of V and false.
	// Example:
	//  s := NewPrimitive(1, 2)
	//  v, ok := s.Pop() // v could be 1 or 2, ok is true
	//  v, ok = s.Pop() // v is the remaining element, ok is true
	//  v, ok = s.Pop() // v is 0, ok is false
	Pop() (V, bool)

	// ToSlice returns a slice containing all elements of the set.
	// The order of elements in the slice is not guaranteed.
	// Example:
	//  s := NewPrimitive(3, 1, 2)
	//  slice := s.ToSlice() // slice could be []int{1, 2, 3}, []int{3, 1, 2}, etc.
	ToSlice() []V

	// Selector returns the function used to derive the key K from an element V.
	// This function is typically provided when the set is created.
	// Example:
	//  keyFn := func(i int) int { return i+1 }
	//  s := New(keyFn)
	//  key := s.Selector(5) // key is 6
	Selector(V) K
}
