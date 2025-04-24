# KSET - Generic Set Implementation for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/sonalys/kset.svg)](https://pkg.go.dev/github.com/sonalys/kset)

`kset` provides a flexible and type-safe implementation of a mathematical set data structure in Go, leveraging generics introduced in Go 1.18.

It allows you to create sets of any data type (`V`) by providing a `KeyFunc` that extracts a comparable key (`K`) from each element. This enables set operations on complex data types based on specific attributes.

## Features

*   **Generic:** Works with any data type for elements (`V`) and keys (`K`).
*   **Key-Based:** Uniqueness is determined by a user-provided key function (`KeyFunc`).
*   **Comprehensive API:** Implements standard set operations like Union, Intersection, Difference, Symmetric Difference, Subset/Superset checks, etc.
*   **Type-Safe:** Compile-time type checking thanks to Go generics.
*   **Iterable:** Provides an `Iter()` method compatible with Go 1.22+ `iter.Seq`.

## Installation

```bash
go get github.com/sonalys/kset
```

## Usage

Here's a basic example using a set of integers:

```go
package main

import (
    "fmt"

    "github.com/sonalys/kset"
)

func main() {
    selector := func(i int) int { return i }
    // Create a new set of integers.
    // The key function simply returns the integer itself.
    intSet := kset.New(selector)

    // Add elements
    intSet.Add(1)
    intSet.Append(2, 3, 2) // Append adds multiple, ignores duplicates

    fmt.Printf("Set: %v\n", intSet.ToSlice()) // Order not guaranteed
    fmt.Printf("Length: %d\n", intSet.Len())
    fmt.Printf("Contains 2? %t\n", intSet.Contains(2))
    fmt.Printf("Contains 4? %t\n", intSet.Contains(4))

    // Create another set
    otherSet := kset.New(selector)
    otherSet.Append(3, 4, 5)

    // Set operations
    union := intSet.Union(otherSet)
    intersection := intSet.Intersect(otherSet)
    difference := intSet.Difference(otherSet) // Elements in intSet but not in otherSet
    symDifference := intSet.SymmetricDifference(otherSet)

    fmt.Printf("Other Set: %v\n", otherSet.ToSlice())
    fmt.Printf("Union: %v\n", union.ToSlice())
    fmt.Printf("Intersection: %v\n", intersection.ToSlice())
    fmt.Printf("Difference (intSet - otherSet): %v\n", difference.ToSlice())
    fmt.Printf("Symmetric Difference: %v\n", symDifference.ToSlice())

    // Iterate over the set
    fmt.Println("Iterating:")
    for v := range intSet.Iter() {
        fmt.Printf("- %d\n", v)
    }
}

// Possible Output:
// Set: [1 2 3]
// Length: 3
// Contains 2? true
// Contains 4? false
// Other Set: [3 4 5]
// Union: [1 2 3 4 5]
// Intersection: [3]
// Difference (intSet - otherSet): [1 2]
// Symmetric Difference: [1 2 4 5]
// Iterating:
// - 1
// - 2
// - 3
```

Using Custom Types
You can use kset with your own structs by providing an appropriate KeyFunc.

```go
package main

import (
    "fmt"

    "github.com/sonalys/kset"
)

type User struct {
    ID   int
    Name string
}

// Key function extracts the User ID as the key
func userIDSelector(u User) int {
    return u.ID
}

func main() {
    userSet := kset.New(userIDSelector)

    userSet.Add(User{ID: 1, Name: "Alice"})
    userSet.Add(User{ID: 2, Name: "Bob"})
    // Adding a user with the same ID will not change the set
    added := userSet.Add(User{ID: 1, Name: "Alice Smith"}) // added will be false

    fmt.Printf("User Set Length: %d\n", userSet.Len())
    fmt.Printf("Added duplicate ID? %t\n", added)

    if user, ok := userSet.Pop(); ok {
        fmt.Printf("Popped user: %+v\n", user)
    }

    fmt.Println("Remaining users:")
    for user := range userSet.Iter() {
        fmt.Printf("- %+v\n", user)
    }
}

// Possible Output:
// User Set Length: 2
// Added duplicate ID? false
// Popped user: {ID:1 Name:Alice}
// Remaining users:
// - {ID:2 Name:Bob}
```