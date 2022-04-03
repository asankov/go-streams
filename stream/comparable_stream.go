package stream

// compile-time interface check
var _ Stream[int] = (*ComparableStream[int])(nil)

type ComparableStream[V comparable] struct {
	Stream[V]
}

func ComparableStreamOf[V comparable](elements []V) *ComparableStream[V] {
	return &ComparableStream[V]{
		Stream: Of(elements),
	}
}

// Distinct returns a stream consisting of the distinct elements (according to ==) of this stream.
//
// java: Stream<T> distinct()
func (s *ComparableStream[V]) Distinct() *ComparableStream[V] {
	set := map[V]any{}
	s.ForEach(func(el V) {
		set[el] = struct{}{}
	})
	uniqueElements := make([]V, 0, len(set))
	for uniqueElement := range set {
		uniqueElements = append(uniqueElements, uniqueElement)
	}
	return ComparableStreamOf(uniqueElements)
}
