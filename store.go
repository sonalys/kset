package kset

import (
	"iter"
)

type (
	Storage[Key, Value any] interface {
		Len() int
		Clear()
		Delete(...Key)
		Contains(Key) bool
		Get(Key) (Value, bool)
		Upsert(Key, Value)
		Iter() iter.Seq2[Key, Value]
		Clone() Storage[Key, Value]
	}

	empty = struct{}
)
