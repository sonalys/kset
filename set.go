package kset

import (
	"iter"
)

// Set defines the interface of the behavior expected from only comparing keys, and not values.
// This interface is useful for comparing sets that shares the same key type, but not the same value.
type Set[Key any] interface {
	// Len returns the number of elements in the set.
	// Example:
	//  s := kset.HashMapKey(1, 2)
	//  length := s.Len() // length is 2
	Len() int

	// Clear removes all elements from the set.
	// Example:
	//  s := kset.HashMapKey(1, 2)
	//  s.Clear() // s is {}
	//  length := s.Len() // length is 0
	Clear()

	// Contains checks if all specified elements are present in the set.
	// It returns true if all elements v are in the set, false otherwise.
	// Example:
	//  s := kset.HashMapKey(1, 2, 3)
	//  hasAll := s.ContainsKeys(1, 2) // hasAll is true
	//  hasAll = s.ContainsKeys(1, 4) // hasAll is false
	ContainsKeys(keys ...Key) bool

	// ContainsAny checks if any of the specified elements are present in the set.
	// It returns true if at least one element v is in the set, false otherwise.
	// Example:
	//  s := kset.HashMapKey(1, 2)
	//  hasAny := s.ContainsAnyKey(2, 4) // hasAny is true
	//  hasAny = s.ContainsAnyKey(4, 5) // hasAny is false
	ContainsAnyKey(keys ...Key) bool

	// IsEmpty checks if the set contains no elements.
	// Example:
	//  s := kset.HashMapKey[int]()
	//  s.IsEmpty() // empty is true
	//  s.Add(1)
	//  empty = s.IsEmpty() // empty is false
	IsEmpty() bool

	// IsProperSubset checks if the set is a proper subset of another set.
	// A proper subset is a subset that is not equal to the other set.
	// Example:
	//  s1 := kset.HashMapKey(1, 2)
	//  s2 := kset.HashMapKey(1, 2, 3)
	//  s3 := kset.HashMapKey(1, 2)
	//  isProper := s1.IsProperSubset(s2) // isProper is true
	//  isProper = s1.IsProperSubset(s3) // isProper is false
	IsProperSubset(other Set[Key]) bool

	// IsProperSuperset checks if the set is a proper superset of another set.
	// A proper superset is a superset that is not equal to the other set.
	// Example:
	//  s1 := kset.HashMapKey(1, 2, 3)
	//  s2 := kset.HashMapKey(1, 2)
	//  s3 := kset.HashMapKey(1, 2, 3)
	//  isProper := s1.IsProperSuperset(s2) // isProper is true
	//  isProper = s1.IsProperSuperset(s3) // isProper is false
	IsProperSuperset(other Set[Key]) bool

	// IsSubset checks if the set is a subset of another set (i.e., all elements of the current set are also in the other set).
	// Example:
	//  s1 := kset.HashMapKey(1, 2)
	//  s2 := kset.HashMapKey(1, 2, 3)
	//  s3 := kset.HashMapKey(1, 3)
	//  isSub := s1.IsSubset(s2) // isSub is true
	//  isSub = s1.IsSubset(s3) // isSub is false
	IsSubset(other Set[Key]) bool

	// IsSuperset checks if the set is a superset of another set (i.e., all elements of the other set are also in the current set).
	// Example:
	//  s1 := kset.HashMapKey(1, 2, 3)
	//  s2 := kset.HashMapKey(1, 2)
	//  s3 := kset.HashMapKey(1, 4)
	//  isSuper := s1.IsSuperset(s2) // isSuper is true
	//  isSuper = s1.IsSuperset(s3) // isSuper is false
	IsSuperset(other Set[Key]) bool

	// Intersects checks if the set has at least one element in common with another set.
	// Example:
	//  s1 := kset.HashMapKey(1, 2)
	//  s2 := kset.HashMapKey(2, 3)
	//  s3 := kset.HashMapKey(4, 5)
	//  intersects := s1.Intersects(s2) // intersects is true
	//  intersects = s1.Intersects(s3) // intersects is false
	Intersects(other Set[Key]) bool

	// Equal checks if the set is equal to another set (i.e., contains the same elements).
	// Example:
	//  s1 := kset.HashMapKey(1, 2)
	//  s2 := kset.HashMapKey(2, 1)
	//  s3 := kset.HashMapKey(1, 3)
	//  isEqual := s1.Equal(s2) // isEqual is true
	//  isEqual = s1.Equal(s3) // isEqual is false
	Equal(other Set[Key]) bool

	// Keys iterates through all keys stored in the set.
	// Example:
	//	set := kset.HashMapKey(1, 2, 3)
	//	set.Keys() // returns iter[1, 2, 3]
	Keys() iter.Seq[Key]
}
