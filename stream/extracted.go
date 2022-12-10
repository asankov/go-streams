package stream

// This file contains method that in Java are part of the Stream interface,
// but that is not possible in Go, because Go does not allow generic parameters for interface methods.

// Map returns a stream consisting of the results of applying the given function to the elements of this stream.
//
// NOTE: In Java this method is part of the Stream interface.
// However, Go does not support generic parameters for interface methods.
// That is why we have extracted this method as a package method that accepts the stream as a first parameter,
// instead of a method on the interface.
//
//	java: <R> Stream<R> map(Function<? super T,? extends R> mapper)
func Map[T any, R any](stream Stream[T], mapper func(T) R) Stream[R] {
	panic("stream: there is no available Stream implementation yet.")
}

// FlatMap returns a stream consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to each element.
//
// NOTE: In Java this method is part of the Stream interface.
// However, Go does not support generic parameters for interface methods.
// That is why we have extracted this method as a package method that accepts the stream as a first parameter,
// instead of a method on the interface.
//
//	java: <R> Stream<R> flatMap(Function<? super T,? extends Stream<? extends R>> mapper)
func FlatMap[T any, R any](stream Stream[T], mapper func(T) Stream[R]) Stream[R] {
	panic("stream: there is no available Stream implementation yet.")
}

// Collect performs a mutable reduction operation on the elements of this stream.
//
// NOTE: In Java this method is part of the Stream interface.
// However, Go does not support generic parameters for interface methods.
// That is why we have extracted this method as a package method that accepts the stream as a first parameter,
// instead of a method on the interface.
//
//	java: <R> R collect(Supplier<R> supplier, BiConsumer<R,? super T> accumulator, BiConsumer<R,R> combiner)
func Collect[T any, R any](stream Stream[T], supplier Supplier[R], accumulator BiConsumer[T, R], combiner BiConsumer[R, R]) R {
	panic("stream: there is no available Stream implementation yet.")
}

// CollectWithCollector performs a mutable reduction operation on the elements of this stream using a Collector.
//
// NOTE: In Java this method is part of the Stream interface.
// However, Go does not support generic parameters for interface methods.
// That is why we have extracted this method as a package method that accepts the stream as a first parameter,
// instead of a method on the interface.
//
// NOTE: In Java this method overloads the "collect" method, but Go does not support method overloads, so we need to change the name.
//
//	java: <R,A> R collect(Collector<? super T,A,R> collector)
func CollectWithCollector[T any, A any, R any](stream Stream[T], collector Collector[T, A, R]) {
	panic("stream: there is no available Stream implementation yet.")
}

// ReduceWithIdentityAndCombiner performs a reduction on the elements of this stream, using the provided identity, accumulation and combining functions.
//
// NOTE: In Java this method is part of the Stream interface.
// However, Go does not support generic parameters for interface methods.
// That is why we have extracted this method as a package method that accepts the stream as a first parameter,
// instead of a method on the interface.
//
// NOTE: In Java this method overloads the "reduce" method, but Go does not support method overloads, so we need to change the name.
//
//	java: <U> U reduce(U identity, BiFunction<U,? super T,U> accumulator, BinaryOperator<U> combiner)
func ReduceWithIdentityAndCombiner[T any, U any](stream Stream[T], identity U, accumulator func(U, T) U, combiner func(U, U) U) U {
	panic("stream: there is no available Stream implementation yet.")
}
