package stream

// sortable is a wrapper type that implements the sort.Interface interface.
type sortable[T any] struct {
	data       []T
	comparator func(T, T) int
}

func (s sortable[T]) Len() int           { return len(s.data) }
func (s sortable[T]) Less(i, j int) bool { return s.comparator(s.data[i], s.data[j]) < 0 }
func (s sortable[T]) Swap(i, j int)      { s.data[i], s.data[j] = s.data[j], s.data[i] }
