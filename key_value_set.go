package kset

import (
	"fmt"
	"iter"

	"golang.org/x/exp/constraints"
)

// KeyValueSet defines the interface for a generic set data structure.
// K is the comparable type used for the underlying map keys.
// V is the type of the elements stored in the set.
type KeyValueSet[K comparable, V any] interface {
	Set[K]

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
	Clone() KeyValueSet[K, V]

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
	Difference(other Set[K]) KeyValueSet[K, V]

	// Intersect returns a new set containing elements that are common to both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2, 3)
	//  s2 := NewPrimitive(3, 4, 5)
	//  intersection := s1.Intersect(s2) // intersection is {3}
	Intersect(other Set[K]) KeyValueSet[K, V]

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
	SymmetricDifference(other KeyValueSet[K, V]) KeyValueSet[K, V]

	// Union returns a new set containing all elements from both the current set and the other set.
	// Example:
	//  s1 := NewPrimitive(1, 2)
	//  s2 := NewPrimitive(2, 3)
	//  union := s1.Union(s2) // union is {1, 2, 3}
	Union(other KeyValueSet[K, V]) KeyValueSet[K, V]

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

type keyValueSet[K comparable, V any] struct {
	data     Store[K, V]
	selector func(V) K
	newStore func(len int) Store[K, V]
}

// NewKeyValueSet creates a new key-value set implementation.
// It requires a selector function that extracts a comparable key K from a value V.
// Optionally, it can be initialized with one or more values.
//
// Example:
//
//	// Create an unsafe set of User structs, using ID as the key.
//	userSet := kset.NewKeyValueUnsafe(func(u User) int { return u.ID }, user1, user2)
func NewKeyValueSet[K constraints.Ordered, V any](storeType StoreType, selector func(V) K, values ...V) KeyValueSet[K, V] {
	var m Store[K, V]

	switch storeType {
	case HashMap:
		m = NewStoreMapKeyValue(selector, values...)
	case HashMapUnsafe:
		m = NewUnsafeMapStore(selector, values...)
	case TreeMap:
		m = NewStoreTreeMapKeyValue(selector, values...)
	case TreeMapUnsafe:
		m = NewUnsafeStoreTreeMapKeyValue(selector, values...)
	default:
		panic(fmt.Sprintf("type not supported: %s", storeType))
	}

	return &keyValueSet[K, V]{
		data:     m,
		selector: m.Selector(),
		newStore: func(len int) Store[K, V] {
			return &safeMapStore[K, V]{
				store:    make(map[K]V, len),
				selector: m.Selector(),
			}
		},
	}
}

func (k *keyValueSet[K, V]) Selector(value V) K {
	return k.selector(value)
}

func (k *keyValueSet[K, V]) Append(values ...V) int {
	prevLen := k.data.Len()
	for _, val := range values {
		key := k.selector(val)
		k.data.Upsert(key, val)
	}
	return k.data.Len() - prevLen
}

func (k *keyValueSet[K, V]) Len() int {
	return k.data.Len()
}

func (k *keyValueSet[K, V]) Clear() {
	k.data.Clear()
}

func (k *keyValueSet[K, V]) Clone() KeyValueSet[K, V] {
	return &keyValueSet[K, V]{
		data:     k.data.Clone(),
		selector: k.selector,
	}
}

func (k *keyValueSet[K, V]) Contains(values ...V) bool {
	for _, val := range values {
		key := k.selector(val)
		if !k.data.Contains(key) {
			return false
		}
	}
	return true
}

func (k *keyValueSet[K, V]) ContainsKeys(keys ...K) bool {
	for _, key := range keys {
		if !k.data.Contains(key) {
			return false
		}
	}
	return true
}

func (k *keyValueSet[K, V]) ContainsAny(values ...V) bool {
	for _, val := range values {
		key := k.selector(val)
		if k.data.Contains(key) {
			return true
		}
	}
	return false
}

func (k *keyValueSet[K, V]) ContainsAnyKey(keys ...K) bool {
	for _, key := range keys {
		if k.data.Contains(key) {
			return true
		}
	}
	return false
}

func (k *keyValueSet[K, V]) Intersects(other Set[K]) bool {
	for key := range k.data.Iter() {
		if other.ContainsKeys(key) {
			return true
		}
	}

	return false
}

func (k *keyValueSet[K, V]) Difference(other Set[K]) KeyValueSet[K, V] {
	diff := &keyValueSet[K, V]{
		data: k.data.Clone(),
	}

	for key, value := range k.data.Iter() {
		if !other.ContainsKeys(key) {
			diff.Append(value)
		}
	}

	return diff
}

func (k *keyValueSet[K, V]) Each(f func(V) bool) {
	for _, elem := range k.data.Iter() {
		if !f(elem) {
			break
		}
	}
}

func (k *keyValueSet[K, V]) Equal(other Set[K]) bool {
	if k.Len() != other.Len() {
		return false
	}

	for key := range k.data.Iter() {
		if !other.ContainsKeys(key) {
			return false
		}
	}

	return true
}

func (k *keyValueSet[K, V]) Intersect(other Set[K]) KeyValueSet[K, V] {
	intersection := &keyValueSet[K, V]{
		data:     k.newStore(k.data.Len()),
		selector: k.selector,
		newStore: k.newStore,
	}

	for key, value := range k.data.Iter() {
		if other.ContainsKeys(key) {
			intersection.Append(value)
		}
	}
	return intersection
}

func (k *keyValueSet[K, V]) IsEmpty() bool {
	return k.Len() == 0
}

func (k *keyValueSet[K, V]) IsProperSubset(other Set[K]) bool {
	return k.Len() < other.Len() && k.IsSubset(other)
}

func (k *keyValueSet[K, V]) IsProperSuperset(other Set[K]) bool {
	return k.Len() > other.Len() && k.IsSuperset(other)
}

func (k *keyValueSet[K, V]) IsSubset(other Set[K]) bool {
	if k.Len() > other.Len() {
		return false
	}

	for key := range k.data.Iter() {
		if !other.ContainsKeys(key) {
			return false
		}
	}

	return true
}

func (k *keyValueSet[K, V]) IsSuperset(other Set[K]) bool {
	return other.IsSubset(k)
}

func (k *keyValueSet[K, V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, elem := range k.data.Iter() {
			if !yield(elem) {
				break
			}
		}
	}
}

func (k *keyValueSet[K, V]) IterKeys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for key := range k.data.Iter() {
			if !yield(key) {
				break
			}
		}
	}
}

func (k *keyValueSet[K, V]) Pop() (V, bool) {
	for key, value := range k.data.Iter() {
		k.data.Delete(key)
		return value, true
	}

	var zero V
	return zero, false
}

func (k *keyValueSet[K, V]) Remove(values ...V) {
	for _, val := range values {
		key := k.selector(val)
		k.data.Delete(key)
	}
}

func (k *keyValueSet[K, V]) SymmetricDifference(other KeyValueSet[K, V]) KeyValueSet[K, V] {
	sd := &keyValueSet[K, V]{
		data:     k.newStore(k.data.Len()),
		selector: k.selector,
		newStore: k.newStore,
	}

	for _, elem := range k.data.Iter() {
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

func (k *keyValueSet[K, V]) ToSlice() []V {
	result := make([]V, 0, k.Len())
	for _, elem := range k.data.Iter() {
		result = append(result, elem)
	}
	return result
}

func (k *keyValueSet[K, V]) Union(other KeyValueSet[K, V]) KeyValueSet[K, V] {
	union := k.Clone()

	for elem := range other.Iter() {
		union.Append(elem)
	}

	return union
}

var _ KeyValueSet[string, string] = &keyValueSet[string, string]{}
