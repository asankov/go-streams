package stream

import (
	"sort"
)

// compile-time interface check
var _ Stream[int] = (*SliceStream[int])(nil)

type SliceStream[T any] struct {
	elements      []T
	closeHandlers []func()
}

func newSliceStream[T any](elements ...T) *SliceStream[T] {
	return &SliceStream[T]{elements: elements}
}

// AllMatch returns whether all elements of this stream match the provided predicate.
//
//	java: boolean allMatch(Predicate<? super T> predicate)
func (s *SliceStream[T]) AllMatch(predicate func(T) bool) bool {
	for _, el := range s.elements {
		if !predicate(el) {
			return false
		}
	}
	return true
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
//
//	java: boolean anyMatch(Predicate<? super T> predicate)
func (s *SliceStream[T]) AnyMatch(predicate func(T) bool) bool {
	for _, el := range s.elements {
		if predicate(el) {
			return true
		}
	}
	return false
}

// NoneMatch returns whether no elements of this stream match the provided predicate.
//
//	java: boolean noneMatch(Predicate<? super T> predicate)
func (s *SliceStream[T]) NoneMatch(predicate func(T) bool) bool {
	for _, el := range s.elements {
		if predicate(el) {
			return false
		}
	}
	return true
}

// Count returns the count of elements in this stream.
//
//	java: long count()
func (s *SliceStream[T]) Count() int64 {
	return int64(len(s.elements))
}

// Distinct returns a stream consisting of the distinct elements (according to the "==" operator) of this stream.
//
// NOTE: In Java all objects can be compared via Objects.equals.
// In Go, that is not the case, and not everything can be compared via ==.
// In order for this method to be able to be implemented the constraint of the generic type should be "comparable", not "any".
//
// TODO(asankov): create ComparableStream interface that is constrained by "comparable" and link it here.
//
//	java: Stream<T> distinct()
func (s *SliceStream[T]) Distinct() Stream[T] {
	panic(`stream: Distinct cannot be called on a Stream containted by "any"`)
}

// Filter returns a stream consisting of the elements of this stream that match the given predicate.
//
//	java: Stream<T> filter(Predicate<? super T> predicate)
func (s *SliceStream[T]) Filter(predicate func(T) bool) Stream[T] {
	newEl := make([]T, 0, len(s.elements))
	for _, el := range s.elements {
		if predicate(el) {
			newEl = append(newEl, el)
		}
	}
	return newSliceStream(newEl...)
}

// FindAny returns a pointer describing some element of the stream, or a nil pointer if the stream is empty.
//
//	java: Optional<T> findAny()
func (s *SliceStream[T]) FindAny() *T {
	return s.FindFirst()
}

// FindFirst returns a pointer describing the first element of this stream, or a nil pointer if the stream is empty.
//
//	java: Optional<T> findFirst()
func (s *SliceStream[T]) FindFirst() *T {
	if len(s.elements) > 0 {
		return &s.elements[0]
	}
	return nil
}

// FlatMapToInt returns an Stream[int] consisting of the results of replacing each element
// of this stream with the contents of a mapped stream produced by applying the provided mapping
// function to each element.
//
//	java: IntStream flatMapToInt(Function<? super T,? extends IntStream> mapper)
//	java: LongStream flatMapToLong(Function<? super T,? extends LongStream> mapper)
func (s *SliceStream[T]) FlatMapToInt(mapper func(T) Stream[int64]) Stream[int64] {
	return FlatMap[T](s, mapper)
}

// FlatMapToDouble returns an Stream[float64] consisting of the results of replacing each element
// of this stream with the contents of a mapped stream produced by applying the provided mapping
// function to each element.
//
//	java: DoubleStream flatMapToDouble(Function<? super T,? extends DoubleStream> mapper)
func (s *SliceStream[T]) FlatMapToDouble(mapper func(T) Stream[float64]) Stream[float64] {
	return FlatMap[T](s, mapper)

}

// ForEach performs an action for each element of this stream.
//
//	java: void forEach(Consumer<? super T> action)
func (s *SliceStream[T]) ForEach(consumer func(T)) {
	for _, el := range s.elements {
		consumer(el)
	}
}

// ForEachOrdered performs an action for each element of this stream, in the encounter order of the stream if the stream has a defined encounter order.
//
//	java: void forEachOrdered(Consumer<? super T> action)
func (s *SliceStream[T]) ForEachOrdered(consumer func(T)) {
	s.ForEach(consumer)
}

// Limit returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
//
//	java: Stream<T> limit(long maxSize)
func (s *SliceStream[T]) Limit(maxSize int64) Stream[T] {
	if int(maxSize) > len(s.elements) {
		return s
	}
	return &SliceStream[T]{elements: s.elements[0:maxSize]}
}

// MapToInt returns an Stream[int64] consisting of the results of applying the given function to the elements of this stream.
//
//	java: IntStream mapToInt(ToIntFunction<? super T> mapper)
//	java: LongStream mapToLong(ToLongFunction<? super T> mapper)
func (s *SliceStream[T]) MapToInt(mapper func(T) int64) Stream[int64] {
	return Map[T](s, mapper)
}

// MapToDouble returns a DoubleStream consisting of the results of applying the given function to the elements of this stream.
//
//	java: DoubleStream mapToDouble(ToDoubleFunction<? super T> mapper)
func (s *SliceStream[T]) MapToDouble(mapper func(T) float64) Stream[float64] {
	return Map[T](s, mapper)
}

// Max returns the maximum element of this stream according to the provided comparator.
//
//	java: Optional<T> max(Comparator<? super T> comparator)
func (s *SliceStream[T]) Max(comparator func(T, T) int) *T {
	if len(s.elements) == 0 {
		return nil
	}
	max := s.elements[0]
	s.ForEach(func(t T) {
		if comparator(t, max) > 0 {
			max = t
		}
	})
	return &max
}

// Min returns the minimum element of this stream according to the provided comparator.
//
//	java: Optional<T> min(Comparator<? super T> comparator)
func (s *SliceStream[T]) Min(comparator func(T, T) int) *T {
	if s == nil || len(s.elements) == 0 {
		return nil
	}
	min := s.elements[0]
	s.ForEach(func(t T) {
		if comparator(t, min) < 0 {
			min = t
		}
	})
	return &min
}

// Peek returns a stream consisting of the elements of this stream, additionally performing the provided action on each element as elements are consumed from the resulting stream.
//
//	java: Stream<T> peek(Consumer<? super T> action)
func (s *SliceStream[T]) Peek(action func(T)) Stream[T] {
	for _, el := range s.elements {
		action(el)
	}
	return s
}

// Reduce performs a reduction on the elements of this stream, using an associative accumulation function, and returns an Optional describing the reduced value, if any.
//
//	java: Optional<T> reduce(BinaryOperator<T> accumulator)
func (s *SliceStream[T]) Reduce(accumulator func(T, T) T) *T {
	if len(s.elements) == 0 {
		return nil
	}
	res := s.elements[0]
	for i, el := range s.elements {
		if i == 0 {
			continue
		}
		res = accumulator(res, el)
	}
	return &res
}

// ReduceWithIdentity performs a reduction on the elements of this stream, using the provided identity value and an associative accumulation function, and returns the reduced value.
//
// NOTE: In Java this method overloads the "reduce" method, but Go does not support method overloads, so we need to change the name.
//
//	java: T reduce(T identity, BinaryOperator<T> accumulator)
func (s *SliceStream[T]) ReduceWithIdentity(identity T, accumulator func(T, T) T) T {
	result := identity
	s.ForEach(func(t T) {
		result = accumulator(result, t)
	})
	return result
}

// Skip returns a stream consisting of the remaining elements of this stream after discarding the first n elements of the stream.
//
// java: Stream<T> skip(long n)
func (s *SliceStream[T]) Skip(n int64) Stream[T] {
	if len(s.elements) < int(n) {
		return newSliceStream[T]()
	}
	return newSliceStream(s.elements[n:]...)
}

// Sorted returns a stream consisting of the elements of this stream, sorted according to natural order.
//
// NOTE: This method will not be able to implemented in Go, because Go does not have an interface that support "<" and ">".
// Even comparable supports only ==.
// The ony way to implement this would be define a custom interface that support these checks and use it as a constraint.
//
//	java: Stream<T> sorted()
func (s *SliceStream[T]) Sorted() Stream[T] {
	panic("stream: Sorted can only be implemented on types that have defined their ways of being sorted. Use SortedWithComparator.")
}

// Sorted returns a stream consisting of the elements of this stream, sorted according to the provided Comparator.
//
// NOTE: In Java this method overloads the "sorted" method, but Go does not support method overloads, so we need to change the name.
//
//	java: Stream<T> sorted(Comparator<? super T> comparator)
func (s *SliceStream[T]) SortedWithComparator(comparator func(T, T) int) Stream[T] {
	sortable := sortable[T]{data: s.ToArray(), comparator: comparator}
	sort.Sort(sortable)

	return newSliceStream(sortable.data...)
}

// ToArray returns an array containing the elements of this stream.
//
// NOTE: In Java there are 2 "toArray" methods -
// one that receives no arguments and returns an Object[] (hence erasing the original generic type),
// and one that receives an array generator and returns an array of the same type (hence saving the original generic type).
// That is because Java generics are just compile-type checks, and all generic information is erased at runtime.
// This is not the case in Go, and we do not need to receive an array generator in order to be able to preserve the original generic type.
//
//	java: Object[] toArray()
//	java: <A> A[] toArray(IntFunction<A[]> generator)
func (s *SliceStream[T]) ToArray() []T {
	return s.elements
}

// Methods inherited from BaseStream:

// Close closes this stream, causing all close handlers for this stream pipeline to be called.
//
//	java: void close()
func (s *SliceStream[T]) Close() {
	for _, closeHandler := range s.closeHandlers {
		closeHandler()
	}
}

// IsParallel returns whether this stream, if a terminal operation were to be executed, would execute in parallel.
//
// SliceStream is always sequential, so this function will always return false.
//
//	java: boolean isParallel()
func (s *SliceStream[T]) IsParallel() bool {
	return false
}

// Iterator returns an iterator for the elements of this stream.
//
//	java: Iterator<T> iterator()
func (s *SliceStream[T]) Iterator() Iterator[T] {
	panic("TODO: implement")
}

// OnClose returns an equivalent stream with an additional close handler.
//
//	java: S onClose(Runnable closeHandler)
func (s *SliceStream[T]) OnClose(closeHandler func()) Stream[T] {
	return &SliceStream[T]{
		elements:      s.elements,
		closeHandlers: append(s.closeHandlers, closeHandler),
	}
}

// Parallel returns an equivalent stream that is parallel.
//
// SliceStream is always sequential, so this function will always
// return the same stream without doing anything.
//
//	java: S parallel()
func (s *SliceStream[T]) Parallel() Stream[T] {
	return s
}

// Sequential returns an equivalent stream that is sequential.
//
// SliceStream is always sequential, so this function will always
// return the same stream without doing anything.
//
//	java: S sequential()
func (s *SliceStream[T]) Sequential() Stream[T] {
	return s
}

// Spliterator returns a spliterator for the elements of this stream.
//
//	java: Spliterator<T> spliterator()
func (s *SliceStream[T]) Spliterator() Spliterator[T] {
	panic("TODO: implement")
}

// Unordered returns an equivalent stream that is unordered.
//
// SliceStream is always sequential and ordered, so this function will always
// return the same stream without doing anything.
//
//	java: S unordered()
func (s *SliceStream[T]) Unordered() Stream[T] {
	return s
}
