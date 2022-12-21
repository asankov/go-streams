package stream

// This file contains more types that have to be defined in order for the Stream methods to be able to be implemented.

// TODO(asankov): this is to be implemented
type Iterator[T any] interface{}

// TODO(asankov): this is to be implemented
type Spliterator[T any] interface{}

// Supplier represents a supplier of results.
//
// There is no requirement that a new or distinct result be returned each time the supplier is invoked.
//
// NOTE: In Java this is a functional interface with only one method:
//
//	T get()
//
// I think that there is no point to make it like this in Go,
// that is why this type is not an alias for a function that returns a result of the generic type T.
//
//	java: interface Supplier<T>
type Supplier[T any] func() T

// BiConsumer represents an operation that accepts two input arguments and returns no result.
// This is the two-arity specialization of Consumer.
// Unlike most other functional interfaces, BiConsumer is expected to operate via side-effects.
//
// NOTE: In Java this is a functional interface with two methods:
//
//	void accept(T t, U u)
//	BiConsumer<T, U> andThen(BiConsumer<? super T,? super U> after)
//
// For the Stream use-case we are only using the first one, that is why I made this type
// to be a type-alias for a function that accepts two-arguments of generics types T and R,
// so it it currently lacking the 'andThen' method.
// This could be added in the future.
//
// NOTE: Since this method does not return a value it is expected to operator only via side-effects.
// That is why it's important to be used with mutable types (like pointers, or types like strings.Builder).
//
//	java: interface BiConsumer<T,U>
type BiConsumer[T any, R any] func(T, R)

// Collector represents a a mutable reduction operation that accumulates input elements into a mutable result container,
// optionally transforming the accumulated result into a final representation after all input elements have been processed.
// Reduction operations can be performed either sequentially or in parallel.
//
// NOTE: In Java this interface contains a few more methods that are used for the parallel use-cases of the Streams.
// For the sake of simplicity I have decidec to not add them for now.
//
//	java: interface Collector<T,A,R>
type Collector[T any, A any, R any] interface {
	Supplier() Supplier[R]
	Accumulator() BiConsumer[T, R]
	Combiner() BiConsumer[R, R]
}

// TODO(asankov): this is to be implemented
type StreamBuilder[T any] struct{}

// TODO(asankov): this is to be implemented
type UnaryOperator[T any] struct{}
