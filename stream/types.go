package stream

// This file contains more types that have to be defined in order for the Stream methods to be able to be implemented.

// TODO(asankov): this is to be implemented
type Iterator[T any] interface{}

// TODO(asankov): this is to be implemented
type Spliterator[T any] interface{}

type Supplier[T any] func() T

type BiConsumer[T any, R any] func(T, R)

type Collector[T any, A any, R any] interface {
	Supplier() Supplier[R]
	Accumulator() BiConsumer[T, R]
	Combiner() BiConsumer[R, R]
}

type StreamBuilder[T any] struct{}

type UnaryOperator[T any] struct{}
