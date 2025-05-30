package kset_test

import (
	"testing"

	"github.com/sonalys/kset"
)

func setupData(size int) []int {
	output := make([]int, size)
	for i := 0; i < size; i++ {
		output[i] = i
	}
	return output
}

func BenchmarkHashMapKey_1000(b *testing.B) {
	data := setupData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kset.HashMapKey(data...)
	}
}

func BenchmarkHashMapKey_Difference_1000(b *testing.B) {
	data := setupData(1000)

	set := kset.HashMapKey(data...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Difference(kset.HashMapKey(1))
	}
}
