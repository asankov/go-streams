package stream

type ComparableStream[V comparable] struct {
	*Stream[V]
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
	set := map[V]interface{}{}
	for _, el := range s.elements {
		set[el] = struct{}{}
	}
	uniqueElements := make([]V, 0, len(set))
	for uniqueElement := range set {
		uniqueElements = append(uniqueElements, uniqueElement)
	}
	return ComparableStreamOf(uniqueElements)
}
