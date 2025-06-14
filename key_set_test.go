package kset_test

import (
	"testing"

	"github.com/sonalys/kset"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func forEachStoreK[K constraints.Ordered](t *testing.T, f func(t *testing.T, constructor func(keys ...K) kset.KeySet[K])) {
	type tc struct {
		name string
		f    func(keys ...K) kset.KeySet[K]
	}

	stores := []tc{
		{name: "HashMapKey", f: kset.HashMapKey[K]},
		{name: "UnsafeHashMapKey", f: kset.UnsafeHashMapKey[K]},
		{name: "TreeMapKey", f: kset.TreeMapKey[K]},
		{name: "UnsafeTreeMapKey", f: kset.UnsafeTreeMapKey[K]},
	}

	for _, tc := range stores {
		t.Run(tc.name, func(t *testing.T) {
			f(t, tc.f)
		})
	}
}

func Test_KeySet_Append(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("new value", func(t *testing.T) {
			set := constructor(1, 2)
			count := set.Append(3)
			assert.Equal(t, 1, count)
			assert.True(t, set.ContainsKeys(3))
		})

		t.Run("duplicate", func(t *testing.T) {
			set := constructor(1, 2, 3)
			count := set.Append(3)
			assert.Equal(t, 0, count)
			assert.True(t, set.ContainsKeys(3))
		})
	})
}

func Test_KeySet_Clear(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		set := constructor(1, 2)
		set.Clear()
		assert.Equal(t, 0, set.Len())
	})
}

func Test_KeySet_Clone(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		set := constructor(1, 2)
		clone := set.Clone()

		assert.True(t, set.Equal(clone))

		clone.Clear()

		assert.False(t, set.IsEmpty())
	})
}

func Test_KeySet_ContainsKeys(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("contains", func(t *testing.T) {
			set := constructor(1, 2)
			assert.True(t, set.ContainsKeys(1))
		})

		t.Run("not contains", func(t *testing.T) {
			set := constructor(1, 2)
			assert.False(t, set.ContainsKeys(3))
		})
	})
}

func Test_KeySet_ContainsAnyKey(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("contains", func(t *testing.T) {
			set := constructor(1, 2)
			assert.True(t, set.ContainsAnyKey(3, 1))
		})

		t.Run("not contains", func(t *testing.T) {
			set := constructor(1, 2)
			assert.False(t, set.ContainsAnyKey(3, 4))
		})
	})
}

func Test_KeySet_Intersects(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(2, 3)

			assert.True(t, set1.Intersects(set2))
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(3, 4)

			assert.False(t, set1.Intersects(set2))
		})
	})
}

func Test_KeySet_Difference(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(2, 3)

			assert.ElementsMatch(t, []int{1}, set1.Difference(set2).Slice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(3, 4)

			assert.ElementsMatch(t, []int{1, 2}, set1.Difference(set2).Slice())
		})
	})
}

func Test_KeySet_Equal(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("equal", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(2, 1)

			assert.True(t, set1.Equal(set2))
		})

		t.Run("not equal", func(t *testing.T) {
			set1 := constructor(1, 2, 3)
			set2 := constructor(2, 1)

			assert.False(t, set1.Equal(set2))
		})
	})
}

func Test_KeySet_Intersect(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(2, 3)

			assert.ElementsMatch(t, []int{2}, set1.Intersect(set2).Slice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(3, 4)

			assert.Empty(t, set1.Intersect(set2).Slice())
		})
	})
}

func Test_KeySet_IsEmpty(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor()

			assert.True(t, set.IsEmpty())
		})

		t.Run("not empty", func(t *testing.T) {
			set := constructor(1)

			assert.False(t, set.IsEmpty())
		})
	})
}

func Test_KeySet_Subset(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		set1 := constructor(1, 2, 3, 4)
		subset1 := constructor(1, 4)

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

func Test_KeySet_Len(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor()

			assert.Zero(t, set.Len())
		})

		t.Run("not empty", func(t *testing.T) {
			set := constructor(1)

			assert.Equal(t, 1, set.Len())
		})
	})
}

func Test_KeySet_Pop(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor()

			value, ok := set.Pop()
			assert.False(t, ok)
			assert.Zero(t, value)
		})

		t.Run("not empty", func(t *testing.T) {
			set := constructor(1)

			value, ok := set.Pop()
			assert.True(t, ok)
			assert.Equal(t, 1, value)
		})
	})
}

func Test_KeySet_RemoveKeys(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("empty", func(t *testing.T) {
			set := constructor()
			set.RemoveKeys(1)
		})

		t.Run("not found", func(t *testing.T) {
			set := constructor(1)
			set.RemoveKeys(2)
			assert.Equal(t, 1, set.Len())
		})

		t.Run("found", func(t *testing.T) {
			set := constructor(1)
			set.RemoveKeys(1)
			assert.Equal(t, 0, set.Len())
		})
	})
}

func Test_KeySet_SymmetricDifference(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(2, 3)

			assert.ElementsMatch(t, []int{1, 3}, set1.SymmetricDifference(set2).Slice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(3, 4)

			assert.ElementsMatch(t, []int{1, 2, 3, 4}, set1.SymmetricDifference(set2).Slice())
		})
	})
}

func Test_KeySet_Slice(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("not empty", func(t *testing.T) {
			set := constructor(1)
			assert.Equal(t, []int{1}, set.Slice())
		})
	})
}

func Test_KeySet_Union(t *testing.T) {
	forEachStoreK(t, func(t *testing.T, constructor func(values ...int) kset.KeySet[int]) {
		t.Run("intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(2, 3)

			assert.ElementsMatch(t, []int{1, 2, 3}, set1.Union(set2).Slice())
		})

		t.Run("not intersects", func(t *testing.T) {
			set1 := constructor(1, 2)
			set2 := constructor(3, 4)

			assert.ElementsMatch(t, []int{1, 2, 3, 4}, set1.Union(set2).Slice())
		})
	})
}
