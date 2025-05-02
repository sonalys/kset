# KSET - Generic Set Implementation for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/sonalys/kset.svg)](https://pkg.go.dev/github.com/sonalys/kset)
[![Tests](https://github.com/sonalys/kset/actions/workflows/test.yml/badge.svg)](https://github.com/sonalys/kset/actions/workflows/test.yml)

`kset` provides a flexible and type-safe implementation of a mathematical set data structure in Go, leveraging generics introduced in Go 1.18.

It allows you to work with sets of keys or key-values.  
It easily converts from slices  
The key set and key-value set are compatible with one another  
You can use sets of different data types as long as their keys are from the same data type.

## Features

*   **Type-Safe:** Compile-time type checking thanks to Go generics.
*   **Key-Based:** Uniqueness is determined by a user-provided selector function.
*   **Comprehensive API:** Implements standard set operations like Union, Intersection, Difference, Symmetric Difference, Subset/Superset checks, etc.
*   **Iterable:** Provides an `iter.Seq[V]` method compatible with Go 1.22+.
*   **Selectable Underlying Data Structure** Allows you to choose either a hash-table or a tree-map structure for the sets.

## Installation

```bash
go get github.com/sonalys/kset
```

## Usage

Using Custom Types
You can use kset with your own structs by providing an appropriate KeyFunc.

```go
package main

import (
	"fmt"
	"slices"

	"github.com/sonalys/kset"
)

func ExampleNewStoreMapKeyValue() {
	type User struct {
		ID   int
		Name string
	}

	userIDSelector := func(u User) int { return u.ID }

	userSet := kset.NewKeyValueSet(kset.TreeMap, userIDSelector,
		User{ID: 1, Name: "Alice"},
		User{ID: 2, Name: "Bob"},
	)

	// Adding a user with the same key will upsert the user.
	// Returns true if the entry is introducing a new key.
	addedCount := userSet.Append(User{ID: 1, Name: "Alice Smith"})
	fmt.Printf("Added: %v\n", addedCount)

	elements := userSet.ToSlice()
	slices.SortFunc(elements, func(a, b User) int {
		return a.ID - b.ID
	})

	fmt.Printf("Set Elements: %+v\n", elements)
	// Output:
	// Added: 0
	// Set Elements: [{ID:1 Name:Alice Smith} {ID:2 Name:Bob}]
}
```

Here's a basic example using a set of integers:

```go
func ExampleNewStoreMapKey() {
	setA := kset.NewKeySet(kset.TreeMap, 1, 2, 3, 1)

	sortSlice := func(slice []int) []int {
		slices.Sort(slice)
		return slice
	}

	fmt.Printf("Set: %v\n", sortSlice(setA.ToSlice()))
	fmt.Printf("Length: %d\n", setA.Len())
	fmt.Printf("Contains 2? %t\n", setA.ContainsKeys(2))
	fmt.Printf("Contains 4? %t\n", setA.ContainsKeys(4))

	setB := kset.NewKeySet(kset.TreeMap, 3, 4, 5)
	setB.Append(3, 4, 5)

	// Set operations
	union := setA.Union(setB)
	intersection := setA.Intersect(setB)
	difference := setA.Difference(setB) // Elements in intSet but not in otherSet
	symDifference := setA.SymmetricDifference(setB)

	fmt.Printf("Other Set: %v\n", sortSlice(setB.ToSlice()))
	fmt.Printf("Union: %v\n", sortSlice(union.ToSlice()))
	fmt.Printf("Intersection: %v\n", sortSlice(intersection.ToSlice()))
	fmt.Printf("Difference (setA - setB): %v\n", sortSlice(difference.ToSlice()))
	fmt.Printf("Symmetric Difference: %v\n", sortSlice(symDifference.ToSlice()))

	// Output:
	// Set: [1 2 3]
	// Length: 3
	// Contains 2? true
	// Contains 4? false
	// Other Set: [3 4 5]
	// Union: [1 2 3 4 5]
	// Intersection: [3]
	// Difference (setA - setB): [1 2]
	// Symmetric Difference: [1 2 4 5]
}
```