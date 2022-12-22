package stream_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/asankov/go-streams/stream"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	s := stream.Of(1, 2, 3)

	mapped := stream.Map(s, strconv.Itoa)

	require.Equal(t, s.Count(), mapped.Count())
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "1" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "2" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "3" }))
}

func TestFlatMap(t *testing.T) {
	s := stream.Of(1, 2, 3)

	mapped := stream.FlatMap(s, func(i int) stream.Stream[string] {
		val := strconv.Itoa(i)
		return stream.Of(val, val+val, val+val+val)
	})

	require.Equal(t, int64(9), mapped.Count())
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "1" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "2" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "3" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "11" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "22" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "33" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "111" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "222" }))
	require.True(t, mapped.AnyMatch(func(s string) bool { return s == "333" }))
}

func TestReduce(t *testing.T) {
	s := stream.Of(1, 2, 3)

	reduced := stream.ReduceWithIdentityAndCombiner(s, "", func(s string, i int) string { return s + strconv.Itoa(i) }, nil)

	require.Equal(t, "123", reduced)
}

func TestCollect(t *testing.T) {
	s := stream.Of(1, 2, 3)

	t.Run("Collect - sum of ints", func(t *testing.T) {
		c := stream.Collect(s,
			func() *int { i := 0; return &i },
			func(i1 int, i2 *int) { *i2 = i1 + *i2 }, nil)

		require.Equal(t, 6, *c)
	})

	t.Run("CollectWithCollector - int to string and concat", func(t *testing.T) {
		sb := stream.CollectWithCollector[int, *strings.Builder, *strings.Builder](s, &simpleCollector[int, *strings.Builder, *strings.Builder]{
			supplier:    func() *strings.Builder { var s strings.Builder; return &s },
			accumulator: func(i1 int, s *strings.Builder) { _, _ = s.Write([]byte(strconv.Itoa(i1))) },
			combiner:    nil,
		})

		require.Equal(t, "123", sb.String())
	})
}

// compile-time interface check
var _ stream.Collector[int, int, int] = (*simpleCollector[int, int, int])(nil)

type simpleCollector[T any, A any, R any] struct {
	supplier    stream.Supplier[R]
	accumulator stream.BiConsumer[T, R]
	combiner    stream.BiConsumer[R, R]
}

func (s *simpleCollector[T, A, R]) Supplier() stream.Supplier[R]         { return s.supplier }
func (s *simpleCollector[T, A, R]) Accumulator() stream.BiConsumer[T, R] { return s.accumulator }
func (s *simpleCollector[T, A, R]) Combiner() stream.BiConsumer[R, R]    { return s.combiner }
