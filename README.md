# go-stream

I implemented (most of) the [Java Stream API](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html) with Go generics... because why not.

## Yet to be implemented

```java
/* Performs a mutable reduction operation on the elements of this stream using a Collector. */
<R,A> R collect(Collector<? super T,A,R> collector)

/* Performs a mutable reduction operation on the elements of this stream. */
<R> R collect(Supplier<R> supplier, BiConsumer<R,? super T> accumulator, BiConsumer<R,R> combiner)

/* Returns an infinite sequential unordered stream where each element is generated by the provided Supplier. */
static <T> Stream<T> generate(Supplier<T> s)

/* Returns an infinite sequential ordered Stream produced by iterative application of a function f to an initial element seed, producing a Stream consisting of seed, f(seed), f(f(seed)), etc. */
static <T> Stream<T> iterate(T seed, UnaryOperator<T> f)

/* Performs a reduction on the elements of this stream, using the provided identity, accumulation and combining functions. */
<U> U reduce(U identity, BiFunction<U,? super T,U> accumulator, BinaryOperator<U> combiner)
```

## Caveats

### Methods that cannot be implemented by the base type

```java
/* Returns a stream consisting of the distinct elements (according to Object.equals(Object)) of this stream. */
Stream<T> distinct()
```

The `distinct` function cannot be implemented if the allowed types if the type `V` starts from `any`.
That is because unlike in Java, some types in Go cannot be compared.

In Java each object has the `equals` method which can be used for comparison.
In Go, there is no such thing and comparing values by using `==`.
Types that can be compared via `==` implement the `comparable` interface.

Hence, in order to implement this we need to limit `V` to `comparable`.

I did not want to do this for the base `Stream` type, where I wanted to allow all types, which no limitation.

Instead, this is done on a separate type that has this limitation - [`ComparableStream`](stream/comparable_stream.go).
This type embeds `Stream` and has all the methods that `Stream` has.