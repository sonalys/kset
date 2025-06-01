package kset_test

import (
	"testing"

	"github.com/sonalys/kset"
)

func BenchmarkHashMapKeyValue_Clone(b *testing.B) {
	// Difference scales with o(n) as it copies the whole set.
	runBenchmark := func(size int) func(b *testing.B) {
		data := setupData(size)
		set := kset.HashMapKeyValue(func(v int) int { return v }, data...)
		return func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				set.Clone()
			}
		}
	}

	b.Run("10", runBenchmark(10))
	b.Run("100", runBenchmark(100))
	b.Run("1000", runBenchmark(1000))
	b.Run("10000", runBenchmark(10000))
	// Output:
	// goos: linux
	// goarch: amd64
	// pkg: github.com/sonalys/kset
	// cpu: AMD Ryzen 9 5950X 16-Core Processor
	// BenchmarkHashMapKeyValue_Clone/10-32 	 1868337	       641.8 ns/op	     432 B/op	       7 allocs/op
	// BenchmarkHashMapKeyValue_Clone/100-32         	  312880	      3672 ns/op	    2448 B/op	       7 allocs/op
	// BenchmarkHashMapKeyValue_Clone/1000-32        	   31053	     38714 ns/op	   37048 B/op	       9 allocs/op
	// BenchmarkHashMapKeyValue_Clone/10000-32       	    2853	    421953 ns/op	  295661 B/op	      37 allocs/op
}

func BenchmarkHashMapKeyValue_1000(b *testing.B) {
	data := setupData(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kset.HashMapKeyValue(func(v int) int { return v }, data...)
	}
}
