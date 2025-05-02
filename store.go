package kset

import (
	"fmt"
	"iter"
)

type (
	Store[K comparable, V any] interface {
		Len() int
		Clear()
		Delete(K)
		Contains(K) bool
		Get(K) (V, bool)
		Upsert(K, V)
		Iter() iter.Seq2[K, V]
		Clone() Store[K, V]
	}

	StoreType int
)

const (
	StoreTypeUndefined StoreType = iota
	// HashMap is a thread-safe hash table store implementation.
	//	Operation		Average		WorstCase
	//	Search			O(1)		O(logN^2)
	//	Insert			O(1)		O(logN^2)
	//	Delete			O(1)		O(n)
	// Space complexity
	//	Space			O(n)		O(n)
	HashMap
	// HashMapUnsafe is a thread-unsafe hash table store implementation.
	//	Operation		Average		WorstCase
	//	Search			O(1)		O(logN^2)
	//	Insert			O(1)		O(logN^2)
	//	Delete			O(1)		O(n)
	// Space complexity
	//	Space			O(n)		O(n)
	HashMapUnsafe
	// TreeMap is a thread-safe red-black tree store implementation.
	//	Operation		Average		WorstCase
	//	Search			O(logN)		O(logN)
	//	Insert			O(logN)		O(logN)
	//	Delete			O(logN)		O(logN)
	// Space complexity
	//	Space			O(n)		O(n)
	TreeMap
	// TreeMapUnsafe is a thread-unsafe red-black tree store implementation.
	//	Operation		Average		WorstCase
	//	Search			O(logN)		O(logN)
	//	Insert			O(logN)		O(logN)
	//	Delete			O(logN)		O(logN)
	// Space complexity
	//	Space			O(n)		O(n)
	TreeMapUnsafe
)

var storeTypeString = map[StoreType]string{
	HashMap:       "hashMap",
	HashMapUnsafe: "unsafeHashMap",
	TreeMap:       "treeMap",
	TreeMapUnsafe: "unsafeTreeMap",
}

func (t StoreType) String() string {
	val, ok := storeTypeString[t]
	if ok {
		return val
	}

	return fmt.Sprintf("invalid(%v)", int(t))
}
