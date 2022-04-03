package stream

type Builder[V any] struct {
	elements []V
}

// NewBuilder creates and returns a builder for a Stream.
//
// java: static <T> Stream.Builder<T> builder()
func NewBuilder[V any]() *Builder[V] {
	return &Builder[V]{}
}

// Accept adds an element to the stream being built.
//
// java: void accept(T t)
func (b *Builder[V]) Accept(value V) {
	b.elements = append(b.elements, value)
}

// Add adds an element to the stream being built and returns the builder.
//
// java: default Stream.Builder<T> add(T t)
func (b *Builder[V]) Add(value V) *Builder[V] {
	b.Accept(value)
	return b
}

// Build builds and returns the stream.
//
// java: Stream<T> build()
func (b *Builder[V]) Build() Stream[V] {
	return Of(b.elements)
}
