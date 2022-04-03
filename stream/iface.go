package stream

// Stream is a sequence of elements supporting sequential and parallel aggregate operations.
type Stream[V any] interface {

	// Count returns the number of elements inside the stream.
	//
	// java: long count()
	Count() int

	// FindAny returns a reference to any of the elements of the stream or nil if the stream is empty.
	//
	// java: Optional<T> findAny()
	FindAny() *V

	// FindAny returns a reference to the first element of the stream or nil if the stream is empty.
	//
	// java: Optional<T> findFirst()
	FindFirst() *V

	// ForEach performs the given action for each element of the stream.
	//
	// java: void forEach(Consumer<? super T> action)
	ForEach(f func(V))

	// ForEach performs an action for each element of this stream,
	// in the encounter order of the stream if the stream has a defined encounter order.
	//
	// java: void forEachOrdered(Consumer<? super T> action)
	ForEachOrdered(f func(V))

	// Peek returns a stream consisting of the elements of this stream, additionally performing the provided action on each element as elements are consumed from the resulting stream.
	//
	// java: Stream<T> peek(Consumer<? super T> action)
	Peek(f func(V)) Stream[V]

	// AnyMatch return true if any of the elements of the stream match the given predicate.
	//
	// java: boolean anyMatch(Predicate<? super T> predicate)
	AnyMatch(predicate func(V) bool) bool

	// AllMatch returns true if all the elements of the stream match the given predicate.
	//
	// java: boolean allMatch(Predicate<? super T> predicate)
	AllMatch(predicate func(V) bool) bool

	// NoneMatch returns true if no elements of this stream match the provided predicate.
	//
	// java: boolean noneMatch(Predicate<? super T> predicate)
	NoneMatch(predicate func(V) bool) bool

	// Filter returns a stream consisting of the elements of this stream that match the given predicate.
	//
	// java: Stream<T> filter(Predicate<? super T> predicate)
	Filter(filter func(V) bool) Stream[V]

	// MapToInt returns a Stream[int] mapping the elements of the existing stream to int by applying the given function to each element.
	//
	// java: IntStream mapToInt(ToIntFunction<? super T> mapper)
	MapToInt(toIntFunc func(V) int) Stream[int]

	// MapToFloat returns a Stream[float64] mapping the elements of the existing stream to float64 by applying the given function to each element.
	//
	// java: LongStream mapToLong(ToLongFunction<? super T> mapper)
	// java: DoubleStream mapToDouble(ToDoubleFunction<? super T> mapper)
	MapToFloat(toFloatFunc func(V) float64) Stream[float64]

	// ToArray returns an array containing the elements of the stream.
	//
	// NOTE: Because the Java generics are not actual generics but just compile-time syntactic sugar
	// if you want to get an array of the actual type you need to pass an array generator, like `[]String::new` or `[]Integer::new`,
	// otherwise you will get an array of type Object.
	// Go does not have that limitation and we don't need to pass such parameter, hence we don't need the second method.
	//
	// java: Object[] toArray()
	// java: <A> A[] toArray(IntFunction<A[]> generator)
	ToArray() []V

	// Limit returns a new stream with the elements of the existing one limited to `limit` size.
	//
	// java: Stream<T> limit(long maxSize)
	Limit(limit int) Stream[V]

	// Skip returns a stream consisting of the remaining elements of this stream after discarding the first n elements of the stream.
	//
	// java: Stream<T> skip(long n)
	Skip(n int) Stream[V]

	// Min returns the minimum element of this stream according to the provided Comparator.
	//
	// java: min(Comparator<? super T> comparator)
	Min(comparator func(first, second V) int) *V

	// Max returns the maximum element of this stream according to the provided Comparator.
	//
	// java: Optional<T> max(Comparator<? super T> comparator)
	Max(comparator func(first, second V) int) *V

	// Sorted returns a stream consisting of the elements of this stream, sorted according to the given comparator.
	//
	// NOTE: Due to how Go types work it is not possible to have a version of this function without a comparator.
	// In order to have this we would have to limit V to comparable, which we currently don't do.
	//
	// java: Stream<T> sorted()
	// java: Stream<T> sorted(Comparator<? super T> comparator)
	Sorted(comparator func(first, second V) int) Stream[V]

	// Creates a concatenated stream whose elements are all the elements of the first stream followed by all the elements of the second stream.
	//
	// java: static <T> Stream<T> concat(Stream<? extends T> a, Stream<? extends T> b)
	// func Concat[V any](first, second Stream[V]) Stream[V] {
	// 	return Of(append(first.elements, second.elements...))
	// }

	// FlatMapToInt returns a Stream[int] consisting of the results of replacing each element of this stream
	// with the contents of a mapped stream produced by applying the provided mapping function to each element.
	//
	// java: IntStream flatMapToInt(Function<? super T,? extends IntStream> mapper)
	FlatMapToInt(mapper func(el V) Stream[int]) Stream[int]

	// FlatMapToLong returns a Stream[float64] consisting of the results of replacing each element of this stream
	// with the contents of a mapped stream produced by applying the provided mapping function to each element.
	//
	// java: LongStream flatMapToLong(Function<? super T,? extends LongStream> mapper)
	// java: DoubleStream flatMapToDouble(Function<? super T,? extends DoubleStream> mapper)
	FlatMapToLong(mapper func(el V) Stream[float64]) Stream[float64]

	// Reduce performs a reduction on the elements of this stream,
	// using an associative accumulation function,
	// and returns a reference to the reduced value
	//
	// java: Optional<T> reduce(BinaryOperator<T> accumulator)
	Reduce(operator func(first, second V) V) *V

	// ReduceWithIdentity performs a reduction on the elements of this stream,
	// using the provided identity value and an associative accumulation function,
	// and returns the reduced value.
	//
	// java: T reduce(T identity, BinaryOperator<T> accumulator)
	ReduceWithIdentity(identity V, operator func(first, second V) V) V

	// Reduce performs a reduction on the elements of this stream,
	// using the provided identity, accumulation and combining functions.
	//
	// java: <U> U reduce(U identity, BiFunction<U,? super T,U> accumulator, BinaryOperator<U> combiner)
	ReduceWithIdentityAndCombiner(identity V, operator func(first, second V) V, combiner func(first, second V) V) V

	// Close closes this stream.
	//
	// java: void close()
	Close()

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
	OnClose(closeHandler func()) Stream[V]

	// IsParallel returns whether this stream, if a terminal operation were to be executed, would execute in parallel.
	// Calling this method after invoking an terminal stream operation method may yield unpredictable results.
	//
	// java: boolean isParallel()
	IsParallel() bool

	// Parallel returns an equivalent stream that is parallel.
	// May return itself, either because the stream was already parallel,
	// or because the underlying stream state was modified to be parallel.
	//
	// java: S parallel()
	Parallel() Stream[V]

	// Sequential returns an equivalent stream that is sequential.
	// May return itself, either because the stream was already sequential,
	// or because the underlying stream state was modified to be sequential.
	//
	// java: S sequential()
	Sequential() Stream[V]

	// Unordered returns an equivalent stream that is unordered.
	//
	// java: S unordered()
	Unordered() Stream[V]
}
