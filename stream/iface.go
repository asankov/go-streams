package stream

type Stream[T any] interface {

	// AllMatch returns whether all elements of this stream match the provided predicate.
	//
	// 	java: boolean allMatch(Predicate<? super T> predicate)
	AllMatch(predicate func(T) bool) bool

	// AnyMatch returns whether any elements of this stream match the provided predicate.
	//
	// 	java: boolean anyMatch(Predicate<? super T> predicate)
	AnyMatch(predicate func(T) bool) bool

	// Count returns the count of elements in this stream.
	//
	// 	java: long count()
	Count() int64

	// Distinct returns a stream consisting of the distinct elements (according to the "==" operator) of this stream.
	//
	// NOTE: In Java all objects can be compared via Objects.equals.
	// In Go, that is not the case, and not everything can be compared via ==.
	// In order for this method to be able to be implemented the constraint of the generic type should be "comparable", not "any".
	//
	// TODO(asankov): create ComparableStream interface that is constrained by "comparable" and link it here.
	//
	// 	java: Stream<T> distinct()
	Distinct() Stream[T]

	// Filter returns a stream consisting of the elements of this stream that match the given predicate.
	//
	// 	java: Stream<T> filter(Predicate<? super T> predicate)
	Filter(predicate func(T) bool) Stream[T]

	// FindAny returns a pointer describing some element of the stream, or a nil pointer if the stream is empty.
	//
	// 	java: Optional<T> findAny()
	FindAny() *T

	// FindFirst returns a pointer describing the first element of this stream, or a nil pointer if the stream is empty.
	//
	// 	java: Optional<T> findFirst()
	FindFirst() *T

	// FlatMapToInt returns an Interface[int] consisting of the results of replacing each element
	// of this stream with the contents of a mapped stream produced by applying the provided mapping
	// function to each element.
	//
	// 	java: IntStream flatMapToInt(Function<? super T,? extends IntStream> mapper)
	// 	java: LongStream flatMapToLong(Function<? super T,? extends LongStream> mapper)
	FlatMapToInt(mapper func(T) Stream[int64]) Stream[int64]

	// FlatMapToDouble returns an Interface[float64] consisting of the results of replacing each element
	// of this stream with the contents of a mapped stream produced by applying the provided mapping
	// function to each element.
	//
	// 	java: DoubleStream flatMapToDouble(Function<? super T,? extends DoubleStream> mapper)
	FlatMapToDouble(mapper func(T) Stream[float64]) Stream[float64]

	// ForEach performs an action for each element of this stream.
	//
	// 	java: void forEach(Consumer<? super T> action)
	ForEach(consumer func(T))

	// ForEachOrdered performs an action for each element of this stream, in the encounter order of the stream if the stream has a defined encounter order.
	//
	// 	java: void forEachOrdered(Consumer<? super T> action)
	ForEachOrdered(consumer func(T))

	// Limit returns a stream consisting of the elements of this stream, truncated to be no longer than maxSize in length.
	//
	// 	java: Stream<T> limit(long maxSize)
	Limit(maxSize int64) Stream[T]

	// MapToInt returns an Interface[int64] consisting of the results of applying the given function to the elements of this stream.
	//
	// 	java: IntStream mapToInt(ToIntFunction<? super T> mapper)
	//	java: LongStream mapToLong(ToLongFunction<? super T> mapper)
	MapToInt(mapper func(T) int64) Stream[int64]

	// MapToDouble returns a DoubleStream consisting of the results of applying the given function to the elements of this stream.
	//
	// 	java: DoubleStream mapToDouble(ToDoubleFunction<? super T> mapper)
	MapToDouble(mapper func(T) float64) Stream[float64]

	// Max returns the maximum element of this stream according to the provided comparator.
	//
	// 	java: Optional<T> max(Comparator<? super T> comparator)
	Max(comparator func(T, T) int) *T

	// Min returns the minimum element of this stream according to the provided comparator.
	//
	// 	java: Optional<T> min(Comparator<? super T> comparator)
	Min(comparator func(T, T) int) *T

	// NoneMatch returns whether no elements of this stream match the provided predicate.
	//
	// 	java: boolean noneMatch(Predicate<? super T> predicate)
	NoneMatch(predicate func(T) bool) bool

	// Peek returns a stream consisting of the elements of this stream, additionally performing the provided action on each element as elements are consumed from the resulting stream.
	//
	// 	java: Stream<T> peek(Consumer<? super T> action)
	Peek(action func(T)) Stream[T]

	// Reduce performs a reduction on the elements of this stream, using an associative accumulation function, and returns an Optional describing the reduced value, if any.
	//
	// 	java: Optional<T> reduce(BinaryOperator<T> accumulator)
	Reduce(accumulator func(T, T) T) *T

	// ReduceWithIdentity performs a reduction on the elements of this stream, using the provided identity value and an associative accumulation function, and returns the reduced value.
	//
	// NOTE: In Java this method overloads the "reduce" method, but Go does not support method overloads, so we need to change the name.
	//
	// 	java: T reduce(T identity, BinaryOperator<T> accumulator)
	ReduceWithIdentity(identity T, accumulator func(T, T) T) T

	// Skip returns a stream consisting of the remaining elements of this stream after discarding the first n elements of the stream.
	//
	// java: Stream<T> skip(long n)
	Skip(n int64) Stream[T]

	// Sorted returns a stream consisting of the elements of this stream, sorted according to natural order.
	//
	// NOTE: This method will not be able to implemented in Go, because Go does not have an interface that support "<" and ">".
	// Even comparable supports only ==.
	// The ony way to implement this would be define a custom interface that support these checks and use it as a constraint.
	//
	// 	java: Stream<T> sorted()
	Sorted() Stream[T]

	// Sorted returns a stream consisting of the elements of this stream, sorted according to the provided Comparator.
	//
	// NOTE: In Java this method overloads the "sorted" method, but Go does not support method overloads, so we need to change the name.
	//
	// 	java: Stream<T> sorted(Comparator<? super T> comparator)
	SortedWithComparator(comparator func(T, T) int) Stream[T]

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
	ToArray() []T

	// Methods inherited from BaseStream:

	// Close closes this stream, causing all close handlers for this stream pipeline to be called.
	//
	// 	java: void close()
	Close()

	// IsParallel returns whether this stream, if a terminal operation were to be executed, would execute in parallel.
	//
	// 	java: boolean isParallel()
	IsParallel() bool

	// Iterator returns an iterator for the elements of this stream.
	//
	// 	java: Iterator<T> iterator()
	Iterator() Iterator[T]

	// OnClose returns an equivalent stream with an additional close handler.
	//
	// 	java: S onClose(Runnable closeHandler)
	OnClose(closeHandler func()) Stream[T]

	// Parallel returns an equivalent stream that is parallel.
	//
	// 	java: S parallel()
	Parallel() Stream[T]

	// Sequential returns an equivalent stream that is sequential.
	//
	// 	java: S sequential()
	Sequential() Stream[T]

	// Spliterator returns a spliterator for the elements of this stream.
	//
	// 	java: Spliterator<T> spliterator()
	Spliterator() Spliterator[T]

	// Unordered returns an equivalent stream that is unordered.
	//
	// 	java: S unordered()
	Unordered() Stream[T]
}
