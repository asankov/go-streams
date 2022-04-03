package stream_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/asankov/go-streams/stream"
	"github.com/stretchr/testify/require"
)

func TestStream(t *testing.T) {
	str := stream.Of([]int{1, 2, 3})
	testStream(t, str)

	calledTimes := 0
	peaked := str.Peek(func(i int) {
		calledTimes++
		if i < 1 || i > 3 {
			t.Fatalf("ForEach: unexpected value %d", i)
		}
	})
	require.Equal(t, 3, calledTimes)

	testStream(t, peaked)
}

func TestAllMatch(t *testing.T) {
	builder := stream.NewBuilder[int]().
		Add(1).
		Add(1)
	builder.Accept(1)
	str := builder.Build()

	allMatch := str.AllMatch(func(i int) bool { return i == 1 })
	require.True(t, allMatch)
}

func TestEmpty(t *testing.T) {
	str := stream.Empty[int]()
	require.Nil(t, str.FindAny())
	require.Nil(t, str.FindFirst())
	require.Nil(t, str.Min(func(first, second int) int { return 0 }))
	require.Nil(t, str.Max(func(first, second int) int { return 0 }))
}

func TestOfSingle(t *testing.T) {
	str := stream.OfSingle(1)

	require.Equal(t, 1, str.Count())
	require.True(t, str.AllMatch(func(i int) bool { return i == 1 }))
	require.NotNil(t, str.FindAny())
	require.NotNil(t, str.FindFirst())
}

func TestLimit(t *testing.T) {
	str := stream.Of([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	limited := str.Limit(5)
	require.Equal(t, 5, limited.Count())

	limited = str.Limit(50)
	require.Equal(t, 10, limited.Count())
}

func TestConcat(t *testing.T) {
	first := stream.Of([]int{1, 2, 3, 4, 5})
	second := stream.Of([]int{6, 7, 8, 9, 10})

	concat := stream.Concat(first, second)
	require.Equal(t, 10, concat.Count())
	for i := 1; i <= 10; i++ {
		require.True(t, concat.AnyMatch(func(i int) bool { return i == i }))
	}
}

func testStream(t *testing.T, str stream.Stream[int]) {
	require.Equal(t, 3, str.Count())

	anyMatch := str.AnyMatch(func(i int) bool { return i == 1 })
	require.True(t, anyMatch)

	anyMatch = str.AnyMatch(func(i int) bool { return i == 5 })
	require.False(t, anyMatch)

	allMatch := str.AllMatch(func(i int) bool { return i == 1 })
	require.False(t, allMatch)

	noneMatch := str.NoneMatch(func(i int) bool { return i == 1 })
	require.False(t, noneMatch)

	noneMatch = str.NoneMatch(func(i int) bool { return i == 5 })
	require.True(t, noneMatch)

	toInt := str.MapToInt(func(i int) int { return i })
	require.Equal(t, 3, str.Count())
	require.True(t, toInt.AnyMatch(func(i int) bool { return i == 1 }))
	require.True(t, toInt.AnyMatch(func(i int) bool { return i == 2 }))
	require.True(t, toInt.AnyMatch(func(i int) bool { return i == 3 }))

	toFloat := str.MapToFloat(func(i int) float64 { return float64(i) })
	require.Equal(t, 3, str.Count())
	require.True(t, toFloat.AnyMatch(func(f float64) bool { return f == float64(1) }))
	require.True(t, toFloat.AnyMatch(func(f float64) bool { return f == float64(2) }))
	require.True(t, toFloat.AnyMatch(func(f float64) bool { return f == float64(3) }))

	toString := stream.Map(str, func(v int) string { return strconv.Itoa(v) })
	require.Equal(t, 3, str.Count())
	require.True(t, toString.AnyMatch(func(s string) bool { return s == "1" }))
	require.True(t, toString.AnyMatch(func(s string) bool { return s == "2" }))
	require.True(t, toString.AnyMatch(func(s string) bool { return s == "3" }))

	filtered := str.Filter(func(i int) bool { return i == 1 || i == 3 })
	require.Equal(t, 2, filtered.Count())
	require.True(t, filtered.AnyMatch(func(i int) bool { return i == 1 }))
	require.True(t, filtered.AnyMatch(func(i int) bool { return i == 3 }))
	require.False(t, filtered.AnyMatch(func(i int) bool { return i == 2 }))

	require.NotNil(t, str.FindAny())
	require.NotNil(t, str.FindFirst())

	flatMappedInt := str.FlatMapToInt(func(el int) stream.Stream[int] {
		return stream.Of([]int{el * 5, el * 10})
	})
	require.True(t, flatMappedInt.AnyMatch(func(i int) bool { return i == 5 }))
	require.True(t, flatMappedInt.AnyMatch(func(i int) bool { return i == 10 }))
	require.True(t, flatMappedInt.AnyMatch(func(i int) bool { return i == 15 }))
	require.True(t, flatMappedInt.AnyMatch(func(i int) bool { return i == 10 }))
	require.True(t, flatMappedInt.AnyMatch(func(i int) bool { return i == 20 }))
	require.True(t, flatMappedInt.AnyMatch(func(i int) bool { return i == 30 }))

	flatMappedDouble := str.FlatMapToLong(func(el int) stream.Stream[float64] {
		return stream.Of([]float64{float64(el) + 0.5, float64(el) + 0.75})
	})
	require.True(t, flatMappedDouble.AnyMatch(func(i float64) bool { return i == 1.5 }))
	require.True(t, flatMappedDouble.AnyMatch(func(i float64) bool { return i == 1.75 }))
	require.True(t, flatMappedDouble.AnyMatch(func(i float64) bool { return i == 2.5 }))
	require.True(t, flatMappedDouble.AnyMatch(func(i float64) bool { return i == 2.75 }))
	require.True(t, flatMappedDouble.AnyMatch(func(i float64) bool { return i == 3.5 }))
	require.True(t, flatMappedDouble.AnyMatch(func(i float64) bool { return i == 3.75 }))

	flatMappedString := stream.FlatMap(str, func(el int) stream.Stream[string] {
		s := strconv.Itoa(el)
		return stream.Of([]string{s, s + "."})
	})
	require.True(t, flatMappedString.AnyMatch(func(s string) bool { return s == "1" }))
	require.True(t, flatMappedString.AnyMatch(func(s string) bool { return s == "1." }))
	require.True(t, flatMappedString.AnyMatch(func(s string) bool { return s == "2" }))
	require.True(t, flatMappedString.AnyMatch(func(s string) bool { return s == "2." }))
	require.True(t, flatMappedString.AnyMatch(func(s string) bool { return s == "3" }))
	require.True(t, flatMappedString.AnyMatch(func(s string) bool { return s == "3." }))

	reducedSum := str.ReduceWithIdentity(0, func(value, el int) int { return el + value })
	require.Equal(t, 6, reducedSum)

	reducedSum = str.ReduceWithIdentityAndCombiner(0, func(value, el int) int { return el + value }, func(first, second int) int { return 0 })
	require.Equal(t, 6, reducedSum)

	t.Run("ToArray", func(t *testing.T) {
		arr := str.ToArray()

		require.Equal(t, 3, len(arr))
		require.Contains(t, arr, 1)
		require.Contains(t, arr, 2)
		require.Contains(t, arr, 3)
	})

	t.Run("Reduce", func(t *testing.T) {
		res := stream.Empty[int]().Reduce(func(first, second int) int {
			if first > second {
				return first
			}
			return second
		})
		require.Nil(t, res)

		res = stream.Of([]int{1}).Reduce(func(first, second int) int {
			if first > second {
				return first
			}
			return second
		})
		require.NotNil(t, res)
		require.Equal(t, 1, *res)

		res = stream.Of([]int{1, 2, 3}).Reduce(func(first, second int) int {
			if first > second {
				return first
			}
			return second
		})
		require.NotNil(t, res)
		require.Equal(t, 3, *res)
	})

	t.Run("ForEach", func(t *testing.T) {
		calledTimes := 0
		str.ForEach(func(i int) {
			calledTimes++
			if i < 1 || i > 3 {
				t.Fatalf("ForEach: unexpected value %d", i)
			}
		})
		require.Equal(t, 3, calledTimes)
	})

	t.Run("ForEachOrdered", func(t *testing.T) {
		calledTimes := 0
		str.ForEachOrdered(func(i int) {
			calledTimes++
			if i < 1 || i > 3 {
				t.Fatalf("ForEach: unexpected value %d", i)
			}
		})
		require.Equal(t, 3, calledTimes)
	})

	min := str.Min(func(first, second int) int {
		return first - second
	})
	require.NotNil(t, min)
	require.Equal(t, 1, *min)

	newStr := stream.Of([]int{3, 2, 1})
	min = newStr.Min(func(first, second int) int {
		return first - second
	})
	require.NotNil(t, min)
	require.Equal(t, 1, *min)

	max := str.Max(func(first, second int) int {
		if first > second {
			return 1
		}
		return -1
	})
	require.NotNil(t, max)
	require.Equal(t, 3, *max)

	skipped := str.Skip(2)
	require.Equal(t, 1, skipped.Count())
	require.True(t, skipped.AnyMatch(func(i int) bool { return i == 3 }))

	skipped = str.Skip(5)
	require.Equal(t, 0, skipped.Count())

	sorted := str.Sorted(func(first, second int) int { return first - second })
	j := 0
	sorted.ForEach(func(i int) {
		j++
		require.Equal(t, j, i)
	})
	require.Equal(t, 3, j)

	t.Run("Close - no handlers", func(t *testing.T) {
		str := stream.Empty[int]()
		require.NotPanics(t, func() {
			str.Close()
		})
	})

	t.Run("Close - handlers, no panic", func(t *testing.T) {
		str := stream.Empty[int]()

		called := 0
		str = str.OnClose(func() { called++ })
		require.NotPanics(t, func() {
			str.Close()
		})
		require.Equal(t, 1, called)
	})

	t.Run("Close - handlers, panic with error", func(t *testing.T) {
		str := stream.Empty[int]()
		str.OnClose(func() {
			panic(fmt.Errorf("panic on close"))
		})
		require.Panics(t, func() {
			str.Close()
		})
	})

	t.Run("Close - handlers, panic with non-error", func(t *testing.T) {
		str := stream.Empty[int]()
		str.OnClose(func() {
			panic("panic on close")
		})
		require.Panics(t, func() {
			str.Close()
		})
	})

	t.Run("Parallel/IsParallel/Sequential", func(t *testing.T) {
		str := stream.Empty[int]()
		require.False(t, str.IsParallel())

		str = str.Sequential()
		require.False(t, str.IsParallel())

		str = str.Parallel()
		require.True(t, str.IsParallel())

		str = str.Sequential()
		require.False(t, str.IsParallel())
	})
}
