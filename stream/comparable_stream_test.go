package stream_test

import (
	"testing"

	"github.com/asankov/go-streams/stream"
	"github.com/stretchr/testify/require"
)

func TestComparableStream(t *testing.T) {
	str := stream.ComparableStreamOf([]int{1, 1, 2, 2, 3})
	require.Equal(t, 5, str.Count())

	distinct := str.Distinct()
	require.Equal(t, 3, distinct.Count())
}
