package kset

import (
	"iter"
	"slices"
)

// KeyValueSet is a key-value set.
// It performs mathematical set operations on a batch of values.
// It uses from a selector to extract keys from any given value.
// The underlying data structure used for the set is dependable on the used constructor.
type KeyValueSet[Key, Value any] interface {
	Set[Key]

	// Append upserts multiple elements to the set.
	// It returns the number of elements that were actually added (i.e., were not already present).
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 1)
	//  count := s.Append(1, 2, 3) // count is 2
	Append(values ...Value) int

	// Clone creates a shallow copy of the set.
	// Example:
	//  s1 := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2)
	//  s2 := s1.Clone() // s2 is {1, 2}, independent of s1
	//  s2.Add(3)
	//  // s1 is {1, 2}, s2 is {1, 2, 3}
	Clone() KeyValueSet[Key, Value]

	// Contains checks if all specified elements are present in the set.
	// It returns true if all elements v are in the set, false otherwise.
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3)
	//  hasAll := s.Contains(1, 2) // hasAll is true
	//  hasAll = s.Contains(1, 4) // hasAll is false
	Contains(values ...Value) bool

	// ContainsAny checks if any of the specified elements are present in the set.
	// It returns true if at least one element v is in the set, false otherwise.
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2)
	//  hasAny := s.ContainsAny(2, 4) // hasAny is true
	//  hasAny = s.ContainsAny(4, 5) // hasAny is false
	ContainsAny(values ...Value) bool

	// Difference returns a new set containing elements that are in the current set but not in the other set.
	// Example:
	//  s1 := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3)
	//  s2 := kset.HashMapKeyValue(func(v int) int { return v }, 3, 4, 5)
	//  diff := s1.Difference(s2) // diff is {1, 2}
	Difference(other Set[Key]) KeyValueSet[Key, Value]

	// DifferenceKeys returns a copy of the current set, excluding the given keys.
	// Example:
	//  s1 := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3)
	//  diff := s1.DifferenceKeys(3) // diff is {1, 2}
	DifferenceKeys(keys ...Key) KeyValueSet[Key, Value]

	// Intersect returns a new set containing elements that are common to both the current set and the other set.
	// Example:
	//  s1 := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3)
	//  s2 := kset.HashMapKeyValue(func(v int) int { return v }, 3, 4, 5)
	//  intersection := s1.Intersect(s2) // intersection is {3}
	Intersect(other Set[Key]) KeyValueSet[Key, Value]

	// KeyValues returns an iterator (iter.Seq) over the elements of the set.
	// The order of iteration is not guaranteed.
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3)
	//  for v := range s.KeyValues() {
	//      fmt.Println(v) // Prints 1, 2, 3 in some order
	//  }
	KeyValues() iter.Seq2[Key, Value]

	// Remove removes the specified elements from the set.
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3, 4)
	//  s.Remove(2, 4) // s is {1, 3}
	Remove(v ...Value)

	// RemoveKeys removes the specified keys from the set.
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3, 4)
	//  s.RemoveKeys(2, 4) // s is {1, 3}
	RemoveKeys(v ...Key)

	// SymmetricDifference returns a new set containing elements that are in either the current set or the other set, but not both.
	// Example:
	//  s1 := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2, 3)
	//  s2 := kset.HashMapKeyValue(func(v int) int { return v }, 3, 4, 5)
	//  symDiff := s1.SymmetricDifference(s2) // symDiff is {1, 2, 4, 5}
	SymmetricDifference(other KeyValueSet[Key, Value]) KeyValueSet[Key, Value]

	// Union returns a new set containing all elements from both the current set and the other set.
	// Example:
	//  s1 := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2)
	//  s2 := kset.HashMapKeyValue(func(v int) int { return v }, 2, 3)
	//  union := s1.Union(s2) // union is {1, 2, 3}
	Union(other KeyValueSet[Key, Value]) KeyValueSet[Key, Value]

	// Pop removes and returns an arbitrary element from the set.
	// It returns the removed element and true if the set was not empty, otherwise it returns the zero value of V and false.
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 1, 2)
	//  v, ok := s.Pop() // v could be 1 or 2, ok is true
	//  v, ok = s.Pop() // v is the remaining element, ok is true
	//  v, ok = s.Pop() // v is 0, ok is false
	Pop() (Value, bool)

	// Slice returns a slice containing all elements of the set.
	// The order of elements in the slice is not guaranteed.
	// Example:
	//  s := kset.HashMapKeyValue(func(v int) int { return v }, 3, 1, 2)
	//  slice := s.Slice() // slice could be []int{1, 2, 3}, []int{3, 1, 2}, etc.
	Slice() []Value
}

type keyValueSet[Key, Value any, Store Storage[Key, Value]] struct {
	store    Store
	selector func(Value) Key
}

func (k *keyValueSet[Key, Value, Store]) Append(values ...Value) int {
	prevLen := k.store.Len()
	for _, val := range values {
		key := k.selector(val)
		k.store.Upsert(key, val)
	}
	return k.store.Len() - prevLen
}

func (k *keyValueSet[Key, Value, Store]) Len() int {
	return k.store.Len()
}

func (k *keyValueSet[Key, Value, Store]) Clear() {
	k.store.Clear()
}

func (k *keyValueSet[Key, Value, Store]) Clone() KeyValueSet[Key, Value] {
	return &keyValueSet[Key, Value, Store]{
		store:    k.store.Clone().(Store),
		selector: k.selector,
	}
}

func (k *keyValueSet[Key, Value, Store]) Contains(values ...Value) bool {
	for _, val := range values {
		key := k.selector(val)
		if !k.store.Contains(key) {
			return false
		}
	}
	return true
}

func (k *keyValueSet[Key, Value, Store]) ContainsKeys(keys ...Key) bool {
	for _, key := range keys {
		if !k.store.Contains(key) {
			return false
		}
	}
	return true
}

func (k *keyValueSet[Key, Value, Store]) ContainsAny(values ...Value) bool {
	for _, val := range values {
		key := k.selector(val)
		if k.store.Contains(key) {
			return true
		}
	}
	return false
}

func (k *keyValueSet[Key, Value, Store]) ContainsAnyKey(keys ...Key) bool {
	return slices.ContainsFunc(keys, k.store.Contains)
}

func (k *keyValueSet[Key, Value, Store]) Intersects(other Set[Key]) bool {
	for key := range k.store.Iter() {
		if other.ContainsKeys(key) {
			return true
		}
	}

	return false
}

func (k *keyValueSet[Key, Value, Store]) Difference(other Set[Key]) KeyValueSet[Key, Value] {
	diff := k.Clone()
	diff.RemoveKeys(bufferedCollect(other.Keys(), other.Len())...)
	return diff
}

func (k *keyValueSet[Key, Value, Store]) DifferenceKeys(keys ...Key) KeyValueSet[Key, Value] {
	diff := k.Clone()
	diff.RemoveKeys(keys...)
	return diff
}

func (k *keyValueSet[Key, Value, Store]) Equal(other Set[Key]) bool {
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

func (k *keyValueSet[Key, Value, Store]) Intersect(other Set[Key]) KeyValueSet[Key, Value] {
	intersection := k.Clone()

	outerKeys := make([]Key, 0, other.Len())

	for key := range k.store.Iter() {
		if !other.ContainsKeys(key) {
			outerKeys = append(outerKeys, key)
		}
	}

	intersection.RemoveKeys(outerKeys...)

	return intersection
}

func (k *keyValueSet[Key, Value, Store]) IsEmpty() bool {
	return k.Len() == 0
}

func (k *keyValueSet[Key, Value, Store]) IsProperSubset(other Set[Key]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

func (k *keyValueSet[Key, Value, Store]) IsProperSuperset(other Set[Key]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

func (k *keyValueSet[Key, Value, Store]) IsSubset(other Set[Key]) bool {
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

func (k *keyValueSet[Key, Value, Store]) IsSuperset(other Set[Key]) bool {
	return other.IsSubset(k)
}

func (k *keyValueSet[Key, Value, Store]) KeyValues() iter.Seq2[Key, Value] {
	return k.store.Iter()
}

func (k *keyValueSet[Key, Value, Store]) Keys() iter.Seq[Key] {
	return func(yield func(Key) bool) {
		for key := range k.store.Iter() {
			if !yield(key) {
				break
			}
		}
	}
}

func (k *keyValueSet[Key, Value, Store]) Pop() (Value, bool) {
	for key, value := range k.store.Iter() {
		defer k.store.Delete(key)
		return value, true
	}

	var zero Value
	return zero, false
}

func (k *keyValueSet[Key, Value, Store]) Remove(values ...Value) {
	keys := make([]Key, 0, len(values))
	for _, val := range values {
		keys = append(keys, k.selector(val))
	}
	k.store.Delete(keys...)
}

func (k *keyValueSet[Key, Value, Store]) RemoveKeys(keys ...Key) {
	k.store.Delete(keys...)
}

func (k *keyValueSet[Key, Value, Store]) SymmetricDifference(other KeyValueSet[Key, Value]) KeyValueSet[Key, Value] {
	sd := k.Clone()

	innerKeys := make([]Key, 0, other.Len())
	outerValues := make([]Value, 0, other.Len())

	for key, value := range other.KeyValues() {
		if !k.ContainsKeys(key) {
			outerValues = append(outerValues, value)
			continue
		}
		innerKeys = append(innerKeys, key)
	}

	sd.RemoveKeys(innerKeys...)
	sd.Append(outerValues...)

	return sd
}

func (k *keyValueSet[Key, Value, Store]) Slice() []Value {
	result := make([]Value, 0, k.Len())
	for _, elem := range k.store.Iter() {
		result = append(result, elem)
	}
	return result
}

func (k *keyValueSet[Key, Value, Store]) Union(other KeyValueSet[Key, Value]) KeyValueSet[Key, Value] {
	union := k.Clone()
	union.Append(other.Slice()...)
	return union
}

var _ KeyValueSet[string, string] = &keyValueSet[string, string, *safeMapStore[string, string]]{}
