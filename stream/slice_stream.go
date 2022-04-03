package stream

import "sort"

// compile-time interface check
var _ Stream[int] = (*SliceStream[int])(nil)

type SliceStream[V any] struct {
	elements      []V
	closeHandlers []func()
	isParallel    bool
}

// Empty creates an empty stream.
//
// java: static <T> Stream<T> empty()
func Empty[V any]() Stream[V] {
	return Of([]V{})
}

// Of returns a stream containing the given elements.
//
// java: static <T> Stream<T> of(T... values)
func Of[V any](elements []V) Stream[V] {
	return &SliceStream[V]{elements: elements}
}

// OfSingle returns a stream containing a single element.
//
// java: static <T> Stream<T> of(T t)
func OfSingle[V any](element V) Stream[V] {
	return Of([]V{element})
}

// Count returns the number of elements of the stream.
//
// java: long count()
func (s *SliceStream[V]) Count() int {
	return len(s.elements)
}

// FindAny returns a reference to any of the elements of the stream or nil if the stream is empty.
//
// java: Optional<T> findAny()
func (s *SliceStream[V]) FindAny() *V {
	return s.FindFirst()
}

// FindAny returns a reference to the first element of the stream or nil if the stream is empty.
//
// java: Optional<T> findFirst()
func (s *SliceStream[V]) FindFirst() *V {
	if len(s.elements) == 0 {
		return nil
	}
	return &s.elements[0]
}

// ForEach performs the given action for each element of the stream.
//
// java: void forEach(Consumer<? super T> action)
func (s *SliceStream[V]) ForEach(f func(V)) {
	for _, v := range s.elements {
		f(v)
	}
}

// ForEach performs an action for each element of this stream,
// in the encounter order of the stream if the stream has a defined encounter order.
//
// java: void forEachOrdered(Consumer<? super T> action)
func (s *SliceStream[V]) ForEachOrdered(f func(V)) {
	s.ForEach(f)
}

// Peek returns a stream consisting of the elements of this stream, additionally performing the provided action on each element as elements are consumed from the resulting stream.
//
// java: Stream<T> peek(Consumer<? super T> action)
func (s *SliceStream[V]) Peek(f func(V)) Stream[V] {
	s.ForEach(f)
	return s
}

// AnyMatch return true if any of the elements of the stream match the given predicate.
//
// java: boolean anyMatch(Predicate<? super T> predicate)
func (s *SliceStream[V]) AnyMatch(predicate func(V) bool) bool {
	for _, v := range s.elements {
		if predicate(v) {
			return true
		}
	}
	return false
}

// AllMatch returns true if all the elements of the stream match the given predicate.
//
// java: boolean allMatch(Predicate<? super T> predicate)
func (s *SliceStream[V]) AllMatch(predicate func(V) bool) bool {
	for _, v := range s.elements {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// NoneMatch returns true if no elements of this stream match the provided predicate.
//
// java: boolean noneMatch(Predicate<? super T> predicate)
func (s *SliceStream[V]) NoneMatch(predicate func(V) bool) bool {
	return s.AllMatch(func(v V) bool {
		return !predicate(v)
	})
}

// Filter returns a stream consisting of the elements of this stream that match the given predicate.
//
// java: Stream<T> filter(Predicate<? super T> predicate)
func (s *SliceStream[V]) Filter(filter func(V) bool) Stream[V] {
	elements := make([]V, 0, len(s.elements))
	for _, v := range s.elements {
		if filter(v) {
			elements = append(elements, v)
		}
	}
	return &SliceStream[V]{elements: elements}
}

// MapToInt returns a *SliceStream[int] mapping the elements of the existing stream to int by applying the given function to each element.
//
// java: IntStream mapToInt(ToIntFunction<? super T> mapper)
func (s *SliceStream[V]) MapToInt(toIntFunc func(V) int) Stream[int] {
	return Map[V](s, toIntFunc)
}

// MapToFloat returns a *SliceStream[float64] mapping the elements of the existing stream to float64 by applying the given function to each element.
//
// java: LongStream mapToLong(ToLongFunction<? super T> mapper)
// java: DoubleStream mapToDouble(ToDoubleFunction<? super T> mapper)
func (s *SliceStream[V]) MapToFloat(toFloatFunc func(V) float64) Stream[float64] {
	return Map[V](s, toFloatFunc)
}

// Map returns a stream consisting of the results of applying the given function to the elements of the given stream.
//
// NOTE: Due to a limitation in the generics spec, this function cannot have *SliceStream[V] as a receiver, because then we cannot use generics for T.
// e.g.
// 	func (s *SliceStream[V]) Map[T any](mapFunc func(V) T) *SliceStream[T]
// is not a valid Go code
//
// java: <R> Stream<R> map(Function<? super T,? extends R> mapper)
func Map[V, T any](s Stream[V], mapFunc func(V) T) Stream[T] {
	elements := make([]T, 0, s.Count())
	s.ForEach(func(v V) {
		elements = append(elements, mapFunc(v))
	})
	return Of(elements)
}

// ToArray returns an array containing the elements of the stream.
//
// NOTE: Because the Java generics are not actual generics but just compile-time syntactic sugar
// if you want to get an array of the actual type you need to pass an array generator, like `[]String::new` or `[]Integer::new`,
// otherwise you will get an array of type Object.
// Go does not have that limitation and we don't need to pass such parameter, hence we don't need the second method.
//
// java: Object[] toArray()
// java: <A> A[] toArray(IntFunction<A[]> generator)
func (s *SliceStream[V]) ToArray() []V {
	return s.elements
}

// Limit returns a new stream with the elements of the existing one limited to `limit` size.
//
// java: Stream<T> limit(long maxSize)
func (s *SliceStream[V]) Limit(limit int) Stream[V] {
	if limit > len(s.elements) {
		return s
	}
	return Of(s.elements[0:limit])
}

// Skip returns a stream consisting of the remaining elements of this stream after discarding the first n elements of the stream.
//
// java: Stream<T> skip(long n)
func (s *SliceStream[V]) Skip(n int) Stream[V] {
	if n > len(s.elements) {
		return Empty[V]()
	}
	return Of(s.elements[n:])
}

// Min returns the minimum element of this stream according to the provided Comparator.
//
// java: min(Comparator<? super T> comparator)
func (s *SliceStream[V]) Min(comparator func(first, second V) int) *V {
	if len(s.elements) == 0 {
		return nil
	}
	sorted := s.Sorted(comparator)
	return sorted.FindFirst()
}

// Max returns the maximum element of this stream according to the provided Comparator.
//
// java: Optional<T> max(Comparator<? super T> comparator)
func (s *SliceStream[V]) Max(comparator func(first, second V) int) *V {
	if len(s.elements) == 0 {
		return nil
	}
	max := s.elements[0]
	for _, v := range s.elements {
		if comparator(max, v) < 0 {
			max = v
		}
	}
	return &max
}

// Sorted returns a stream consisting of the elements of this stream, sorted according to the given comparator.
//
// NOTE: Due to how Go types work it is not possible to have a version of this function without a comparator.
// In order to have this we would have to limit V to comparable, which we currently don't do.
//
// java: Stream<T> sorted()
// java: Stream<T> sorted(Comparator<? super T> comparator)
func (s *SliceStream[V]) Sorted(comparator func(first, second V) int) Stream[V] {
	sort.Slice(s.elements, func(i, j int) bool {
		return comparator(s.elements[i], s.elements[j]) < 1
	})
	return Of(s.elements)
}

// Creates a concatenated stream whose elements are all the elements of the first stream followed by all the elements of the second stream.
//
// java: static <T> Stream<T> concat(Stream<? extends T> a, Stream<? extends T> b)
func Concat[V any](first, second Stream[V]) Stream[V] {
	newElements := make([]V, 0, first.Count()+second.Count())
	first.ForEach(func(el V) { newElements = append(newElements, el) })
	second.ForEach(func(el V) { newElements = append(newElements, el) })
	return Of(newElements)
}

// FlatMapToInt returns a *SliceStream[int] consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to each element.
//
// java: IntStream flatMapToInt(Function<? super T,? extends IntStream> mapper)
func (s *SliceStream[V]) FlatMapToInt(mapper func(el V) Stream[int]) Stream[int] {
	return FlatMap[V](s, mapper)
}

// FlatMapToLong returns a *SliceStream[float64] consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to each element.
//
// java: LongStream flatMapToLong(Function<? super T,? extends LongStream> mapper)
// java: DoubleStream flatMapToDouble(Function<? super T,? extends DoubleStream> mapper)
func (s *SliceStream[V]) FlatMapToLong(mapper func(el V) Stream[float64]) Stream[float64] {
	return FlatMap[V](s, mapper)
}

// FlatMap returns a stream consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to each element.
//
// NOTE: Due to a limitation in the generics spec, this function cannot have *SliceStream[V] as a receiver, because then we cannot use generics for T.
// e.g.
// 	func (s *SliceStream[V]) FlatMap[T any](mapFunc func(V) T) *SliceStream[T]
// is not a valid Go code
//
// java: <R> Stream<R> flatMap(Function<? super T,? extends Stream<? extends R>> mapper)
func FlatMap[V, T any](stream Stream[V], mapper func(el V) Stream[T]) Stream[T] {
	streams := make([]Stream[T], 0, stream.Count())
	stream.ForEach(func(el V) {
		streams = append(streams, mapper(el))
	})

	elements := make([]T, 0, stream.Count())
	for _, stream := range streams {
		stream.ForEach(func(el T) {
			elements = append(elements, el)
		})
	}

	return Of(elements)
}

// Reduce performs a reduction on the elements of this stream,
// using an associative accumulation function,
// and returns a reference to the reduced value
//
// java: Optional<T> reduce(BinaryOperator<T> accumulator)
func (s *SliceStream[V]) Reduce(operator func(first, second V) V) *V {
	if len(s.elements) == 0 {
		return nil
	}

	res := s.ReduceWithIdentity(s.elements[0], operator)
	return &res
}

// ReduceWithIdentity performs a reduction on the elements of this stream,
// using the provided identity value and an associative accumulation function,
// and returns the reduced value.
//
// java: T reduce(T identity, BinaryOperator<T> accumulator)
func (s *SliceStream[V]) ReduceWithIdentity(identity V, operator func(first, second V) V) V {
	res := identity
	for _, el := range s.elements {
		res = operator(res, el)
	}
	return res
}

// Reduce performs a reduction on the elements of this stream,
// using the provided identity, accumulation and combining functions.
//
// java: <U> U reduce(U identity, BiFunction<U,? super T,U> accumulator, BinaryOperator<U> combiner)
func (s *SliceStream[V]) ReduceWithIdentityAndCombiner(identity V, operator func(first, second V) V, combiner func(first, second V) V) V {
	// the combiner only makes sense in a paraller execution
	// we still don't support that so just ignore
	return s.ReduceWithIdentity(identity, operator)
}

// Close closes this stream.
//
// java: void close()
func (s *SliceStream[V]) Close() {
	errors := []error{}
	for _, closeHandler := range s.closeHandlers {
		err := runWithDefer(closeHandler)
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) != 0 {
		panic(&MultiError{errors: errors})
	}
}

func runWithDefer(f func()) (err error) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		if e, isErr := recovered.(error); isErr {
			err = e
		} else {
			// not an error, propagate the panic
			panic(recovered)
		}
	}()

	f()
	return
}

// Returns an equivalent stream with an additional close handler.
// Close handlers are run when the Close() method is called on the stream,
// and are executed in the order they were added.
// All close handlers are run, even if earlier close handlers returned errors.
// The errors returned by the close handlers will be bundled and if any errors ocurred the Close method will panic
// with the bundled errors. TODO(asankov): reevaluate whether this is the best approach
// This is needed to satisfy the Java interface, but still have a way to relay the errors to the caller.
// May return itself.
//
// java: S onClose(Runnable closeHandler)
func (s *SliceStream[V]) OnClose(closeHandler func()) Stream[V] {
	s.closeHandlers = append(s.closeHandlers, closeHandler)
	return s
}

// IsParallel returns whether this stream, if a terminal operation were to be executed, would execute in parallel.
// Calling this method after invoking an terminal stream operation method may yield unpredictable results.
//
// java: boolean isParallel()
func (s *SliceStream[V]) IsParallel() bool {
	return s.isParallel
}

// Parallel returns an equivalent stream that is parallel.
// May return itself, either because the stream was already parallel,
// or because the underlying stream state was modified to be parallel.
//
// java: S parallel()
func (s *SliceStream[V]) Parallel() Stream[V] {
	s.isParallel = true
	return s
}

// Sequential returns an equivalent stream that is sequential.
// May return itself, either because the stream was already sequential,
// or because the underlying stream state was modified to be sequential.
//
// java: S sequential()
func (s *SliceStream[V]) Sequential() Stream[V] {
	s.isParallel = false
	return s
}

// Unordered returns an equivalent stream that is unordered.
//
// java: S unordered()
func (s *SliceStream[V]) Unordered() Stream[V] {
	return s
}
