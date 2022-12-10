package stream

// This file contains more types that have to be defined in order for the Stream methods to be able to be implemented.

type Iterator[T any] interface{}

type Spliterator[T any] interface{}

type Supplier[T any] interface{}

type BiConsumer[T any, R any] interface{}

type Collector[T any, R any, A any] interface{}

type StreamBuilder[T any] struct{}

type UnaryOperator[T any] struct{}
