package kset_test

import (
	"testing"

	"github.com/sonalys/kset"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func testKeyer(v int) int {
	return v
}

func forEachStore[K constraints.Ordered, V any](t *testing.T, f func(t *testing.T, constructor func(selector func(V) K, values ...V) kset.KeyValueSet[K, V])) {
	stores := []kset.StoreType{
		kset.HashMap,
		kset.HashMapUnsafe,
		kset.TreeMap,
		kset.TreeMapUnsafe,
	}

	for _, storeType := range stores {
		t.Run(storeType.String(), func(t *testing.T) {
			f(t, func(selector func(V) K, values ...V) kset.KeyValueSet[K, V] {
				return kset.NewKeyValueSet(storeType, selector, values...)
			})
		})
	}
}

func Test_Append(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("new value", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2)
			count := set.Append(3)
			assert.Equal(t, 1, count)
			assert.True(t, set.Contains(3))
		})

		t.Run("duplicate", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2, 3)
			count := set.Append(3)
			assert.Equal(t, 0, count)
			assert.True(t, set.Contains(3))
		})
	})
}

func Test_Clear(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		set := constructor(testKeyer, 1, 2)
		set.Clear()
		assert.Equal(t, 0, set.Len())
	})
}

func Test_Clone(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		set := constructor(testKeyer, 1, 2)
		clone := set.Clone()

		assert.True(t, set.Equal(clone))

		clone.Clear()

		assert.False(t, set.IsEmpty())
	})
}

func Test_Contains(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("contains", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2)
			assert.True(t, set.Contains(1))
		})

		t.Run("not contains", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2)
			assert.False(t, set.Contains(3))
		})
	})
}

func Test_ContainsAny(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("contains", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2)
			assert.True(t, set.ContainsAny(3, 1))
		})

		t.Run("not contains", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2)
			assert.False(t, set.ContainsAny(3, 4))
		})
	})
}

func Test_Intersects(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 2, 3)

			assert.True(t, set1.Intersects(set2))
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 3, 4)

			assert.False(t, set1.Intersects(set2))
		})
	})
}

func Test_Difference(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 2, 3)

			assert.ElementsMatch(t, []int{1}, set1.Difference(set2).ToSlice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 3, 4)

			assert.ElementsMatch(t, []int{1, 2}, set1.Difference(set2).ToSlice())
		})
	})
}

func Test_Each(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("all elements", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2)

			resp := make([]int, 0, 2)
			set.Each(func(i int) bool {
				resp = append(resp, i)
				return true
			})

			assert.ElementsMatch(t, []int{1, 2}, resp)
		})

		t.Run("early return", func(t *testing.T) {
			set := constructor(testKeyer, 1, 2)

			resp := make([]int, 0, 2)
			set.Each(func(i int) bool {
				resp = append(resp, i)
				return false
			})

			assert.Len(t, resp, 1)
		})
	})
}

func Test_Equal(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("equal", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 2, 1)

			assert.True(t, set1.Equal(set2))
		})

		t.Run("not equal", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2, 3)
			set2 := constructor(testKeyer, 2, 1)

			assert.False(t, set1.Equal(set2))
		})
	})
}

func Test_Intersect(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 2, 3)

			assert.ElementsMatch(t, []int{2}, set1.Intersect(set2).ToSlice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 3, 4)

			assert.Empty(t, set1.Intersect(set2).ToSlice())
		})
	})
}

func Test_IsEmpty(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor(testKeyer)

			assert.True(t, set.IsEmpty())
		})

		t.Run("not empty", func(t *testing.T) {
			set := constructor(testKeyer, 1)

			assert.False(t, set.IsEmpty())
		})
	})
}

func Test_Subset(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		set1 := constructor(testKeyer, 1, 2, 3, 4)
		subset1 := constructor(testKeyer, 1, 4)

		t.Run("proper subset", func(t *testing.T) {
			assert.True(t, subset1.IsProperSubset(set1))
			assert.True(t, set1.IsProperSuperset(subset1))

			assert.False(t, subset1.IsProperSubset(subset1))
			assert.False(t, subset1.IsProperSuperset(subset1))
		})

		t.Run("subset", func(t *testing.T) {
			assert.True(t, subset1.IsSubset(set1))
			assert.True(t, set1.IsSuperset(subset1))

			assert.True(t, subset1.IsSubset(subset1))
			assert.True(t, subset1.IsSuperset(subset1))
		})
	})
}

func Test_Iter(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		set := constructor(testKeyer, 1, 2)
		t.Run("all values", func(t *testing.T) {
			sum := 0
			for value := range set.Iter() {
				sum += value
			}
			assert.Equal(t, sum, 3)
		})

		t.Run("first value", func(t *testing.T) {
			sum := 0
			for value := range set.Iter() {
				sum += value
				break
			}
			assert.NotEqual(t, sum, 3)
		})
	})
}

func Test_Len(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor(testKeyer)

			assert.Zero(t, set.Len())
		})

		t.Run("not empty", func(t *testing.T) {
			set := constructor(testKeyer, 1)

			assert.Equal(t, 1, set.Len())
		})
	})
}

func Test_Pop(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor(testKeyer)

			value, ok := set.Pop()
			assert.False(t, ok)
			assert.Zero(t, value)
		})

		t.Run("not empty", func(t *testing.T) {
			set := constructor(testKeyer, 1)

			value, ok := set.Pop()
			assert.True(t, ok)
			assert.Equal(t, 1, value)
		})
	})
}

func Test_Remove(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor(testKeyer)
			set.Remove(1)
		})

		t.Run("not found", func(t *testing.T) {
			set := constructor(testKeyer, 1)
			set.Remove(2)
			assert.Equal(t, 1, set.Len())
		})

		t.Run("found", func(t *testing.T) {
			set := constructor(testKeyer, 1)
			set.Remove(1)
			assert.Equal(t, 0, set.Len())
		})
	})
}

func Test_SymmetricDifference(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 2, 3)

			assert.ElementsMatch(t, []int{1, 3}, set1.SymmetricDifference(set2).ToSlice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 3, 4)

			assert.ElementsMatch(t, []int{1, 2, 3, 4}, set1.SymmetricDifference(set2).ToSlice())
		})
	})
}

func Test_ToSlice(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("not empty", func(t *testing.T) {
			set := constructor(testKeyer, 1)
			assert.Equal(t, []int{1}, set.ToSlice())
		})
	})
}

func Test_Union(t *testing.T) {
	forEachStore(t, func(t *testing.T, constructor func(selector func(int) int, values ...int) kset.KeyValueSet[int, int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 2, 3)

			assert.ElementsMatch(t, []int{1, 2, 3}, set1.Union(set2).ToSlice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(testKeyer, 1, 2)
			set2 := constructor(testKeyer, 3, 4)

			assert.ElementsMatch(t, []int{1, 2, 3, 4}, set1.Union(set2).ToSlice())
		})
	})
}
