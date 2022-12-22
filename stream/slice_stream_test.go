package stream

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	compareIntFunc = func(i1, i2 int) int {
		if i1 < i2 {
			return -1
		}
		if i1 > i2 {
			return 1
		}
		return 0
	}
)

func TestSliceStream(t *testing.T) {
	s := SliceStream[int]{elements: []int{1, 2, 3}}

	t.Run("TestAllMatch", func(t *testing.T) {
		allMatch := s.AllMatch(func(i int) bool { return i == 1 })
		require.False(t, allMatch)

		allMatch = newSliceStream(1, 1, 1).AllMatch(func(i int) bool { return i == 1 })
		require.True(t, allMatch)
	})

	t.Run("TestAnyMatch", func(t *testing.T) {
		anyMatch := s.AnyMatch(func(i int) bool { return i == 1 })
		require.True(t, anyMatch)
		anyMatch = s.AnyMatch(func(i int) bool { return i == 2 })
		require.True(t, anyMatch)
		anyMatch = s.AnyMatch(func(i int) bool { return i == 3 })
		require.True(t, anyMatch)
		anyMatch = s.AnyMatch(func(i int) bool { return i == 4 })
		require.False(t, anyMatch)
	})

	t.Run("TestNoneMatch", func(t *testing.T) {
		noneMatch := s.NoneMatch(func(i int) bool { return i == 1 })
		require.False(t, noneMatch)
		noneMatch = s.NoneMatch(func(i int) bool { return i == 4 })
		require.True(t, noneMatch)
	})

	t.Run("TestCount", func(t *testing.T) {
		count := s.Count()
		require.Equal(t, int64(3), count)
	})

	t.Run("TestDistinct", func(t *testing.T) {
		require.Panics(t, func() {
			_ = s.Distinct()
		})
	})

	t.Run("TestFilter", func(t *testing.T) {
		filtered := s.Filter(func(i int) bool { return i%2 != 0 })

		require.Equal(t, int64(2), filtered.Count())
		require.True(t, filtered.AnyMatch(func(i int) bool { return i == 1 }))
		require.True(t, filtered.AnyMatch(func(i int) bool { return i == 3 }))
	})

	t.Run("TestFindAny", func(t *testing.T) {
		res := s.FindAny()

		require.NotNil(t, res)

		res = newSliceStream[int]().FindFirst()
		require.Nil(t, res)
	})

	t.Run("TestFindFirst", func(t *testing.T) {
		res := s.FindFirst()

		require.NotNil(t, res)
		require.Equal(t, 1, *res)

		res = newSliceStream[int]().FindFirst()
		require.Nil(t, res)
	})

	t.Run("TestFlatMapToInt", func(t *testing.T) {
		mapped := s.FlatMapToInt(func(i int) Stream[int64] {
			return newSliceStream(int64(i), int64(i*10), int64(i*100))
		})
		require.Equal(t, int64(9), mapped.Count())
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 1 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 10 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 100 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 2 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 20 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 200 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 3 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 30 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 300 }))
	})
	t.Run("TestFlatMapToDouble", func(t *testing.T) {
		mapped := s.FlatMapToDouble(func(i int) Stream[float64] {
			return newSliceStream(float64(i)+0.1, float64(i)+0.2, float64(i)+0.3)
		})
		require.Equal(t, int64(9), mapped.Count())
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 1.1 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 1.2 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 1.3 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 2.1 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 2.2 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 2.3 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 3.1 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 3.2 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 3.3 }))
	})

	t.Run("TestForEach", func(t *testing.T) {
		res := make([]int, 0, s.Count())
		s.ForEach(func(i int) {
			res = append(res, i)
		})

		require.Equal(t, 3, len(res))
		require.Equal(t, s.Count(), int64(len(res)))
		require.Contains(t, res, 1)
		require.Contains(t, res, 2)
		require.Contains(t, res, 3)
	})

	t.Run("TestForEachOrdered", func(t *testing.T) {
		res := make([]int, 0, s.Count())
		s.ForEachOrdered(func(i int) {
			res = append(res, i)
		})

		require.Equal(t, 3, len(res))
		require.Equal(t, s.Count(), int64(len(res)))
		require.Contains(t, res, 1)
		require.Contains(t, res, 2)
		require.Contains(t, res, 3)
	})

	t.Run("TestLimit", func(t *testing.T) {
		limited := s.Limit(2)

		require.Equal(t, int64(2), limited.Count())
		require.True(t, limited.AnyMatch(func(i int) bool { return i == 1 }))
		require.True(t, limited.AnyMatch(func(i int) bool { return i == 2 }))

		limited = s.Limit(4)
		require.Equal(t, s.Count(), limited.Count())
		s.ForEach(func(i int) {
			require.True(t, limited.AnyMatch(func(ii int) bool { return i == ii }))
		})
	})

	t.Run("TestMapToInt", func(t *testing.T) {
		mapped := s.MapToInt(func(i int) int64 {
			return int64(i * 2)
		})

		require.Equal(t, int64(3), mapped.Count())
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 2 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 4 }))
		require.True(t, mapped.AnyMatch(func(i int64) bool { return i == 6 }))
	})

	t.Run("TestMapToDouble", func(t *testing.T) {
		mapped := s.MapToDouble(func(i int) float64 {
			return float64(i) + 0.5
		})

		require.Equal(t, int64(3), mapped.Count())
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 1.5 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 2.5 }))
		require.True(t, mapped.AnyMatch(func(f float64) bool { return f == 3.5 }))
	})

	t.Run("TestMax", func(t *testing.T) {
		max := s.Max(compareIntFunc)
		require.NotNil(t, max)
		require.Equal(t, 3, *max)

		max = newSliceStream(3, 2, 1).Max(compareIntFunc)
		require.NotNil(t, max)
		require.Equal(t, 3, *max)

		max = newSliceStream[int]().Max(compareIntFunc)
		require.Nil(t, max)
	})

	t.Run("TestMin", func(t *testing.T) {
		min := s.Min(compareIntFunc)
		require.NotNil(t, min)
		require.Equal(t, 1, *min)

		min = newSliceStream(2, 3, 1).Min(compareIntFunc)
		require.NotNil(t, min)
		require.Equal(t, 1, *min)

		min = newSliceStream[int]().Min(compareIntFunc)
		require.Nil(t, min)
	})

	t.Run("TestPeek", func(t *testing.T) {
		res := make([]int, 0, s.Count())
		_ = s.Peek(func(i int) {
			res = append(res, i)
		})

		require.Equal(t, 3, len(res))
		require.Equal(t, s.Count(), int64(len(res)))
		require.Contains(t, res, 1)
		require.Contains(t, res, 2)
		require.Contains(t, res, 3)
	})

	sum := func(i1, i2 int) int { return i1 + i2 }

	t.Run("TestReduce", func(t *testing.T) {
		res := s.Reduce(sum)
		require.NotNil(t, res)
		require.Equal(t, 6, *res)

		res = newSliceStream[int]().Reduce(sum)
		require.Nil(t, res)

		res = newSliceStream(1).Reduce(sum)
		require.NotNil(t, res)
		require.Equal(t, 1, *res)

		res = newSliceStream(1, 2).Reduce(sum)
		require.NotNil(t, res)
		require.Equal(t, 3, *res)
	})

	t.Run("TestReduceWithIdentity", func(t *testing.T) {
		res := s.ReduceWithIdentity(0, sum)
		require.Equal(t, 6, res)

		res = newSliceStream[int]().ReduceWithIdentity(0, sum)
		require.Equal(t, 0, res)
	})

	t.Run("TestSkip", func(t *testing.T) {
		skipped := s.Skip(2)
		require.Equal(t, int64(1), skipped.Count())
		require.True(t, skipped.AnyMatch(func(i int) bool { return i == 3 }))

		skipped = s.Skip(5)
		require.Equal(t, int64(0), skipped.Count())
	})
	t.Run("TestSorted", func(t *testing.T) {
		require.Panics(t, func() {
			_ = s.Sorted()
		})
	})
	t.Run("TestSortedWithComparator", func(t *testing.T) {
		sorted := newSliceStream(2, 3, 1).SortedWithComparator(compareIntFunc)

		require.Equal(t, int64(3), sorted.Count())

		i := 0
		sorted.ForEach(func(n int) {
			defer func() {
				i++
			}()
			if i == 0 {
				require.Equal(t, 1, n)
				return
			}
			if i == 1 {
				require.Equal(t, 2, n)
				return
			}
			if i == 2 {
				require.Equal(t, 3, n)
				return
			}
		})
	})

	t.Run("TestToArray", func(t *testing.T) {
		arr := s.ToArray()

		require.Equal(t, s.Count(), int64(len(arr)))
		s.ForEach(func(i int) {
			require.Contains(t, arr, i)
		})
	})

	t.Run("TestClose", func(t *testing.T) {
		var called bool
		withHandlers := s.OnClose(func() { called = true })
		withHandlers.Close()

		require.True(t, called)
	})

	t.Run("TestIsParallel", func(t *testing.T) {
		isParallel := s.IsParallel()
		require.False(t, isParallel)
	})

	t.Run("TestJustForCoverage", func(t *testing.T) {
		_ = s.Unordered()
		_ = s.Sequential()
		_ = s.Parallel()
	})
}
