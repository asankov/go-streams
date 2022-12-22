package stream_test

import (
	"testing"

	"github.com/asankov/go-streams/stream"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	s := stream.Empty[int]()

	require.Equal(t, int64(0), s.Count())
}

func TestOfSingle(t *testing.T) {
	s := stream.OfSingle(1)

	require.Equal(t, int64(1), s.Count())
	require.Equal(t, 1, *s.FindFirst())
}
