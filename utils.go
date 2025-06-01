package kset

import "iter"

func bufferedCollect[T any](seq iter.Seq[T], size int) []T {
	buffer := make([]T, 0, size)
	for value := range seq {
		buffer = append(buffer, value)
	}
	return buffer
}
