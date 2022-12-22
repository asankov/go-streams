# go-streams

This is a Go implementation of the Java [Stream API](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html).

It uses generics to enable the consumers to use the Stream with whatever data type they want.

## Usage

If you want to use this package in your code:

```console
go get github.com/asankov/go-stream
```

and then:

```go
import "github.com/asankov/go-streams/stream"

var s = stream.Of(1, 2, 3)

s.Map(s, func(i int) int { return i * 2}).
    ForEach(fmt.Println)
```

## Blog post

If you want to learn more details about my motivation to write this and follow my steps in doing so, check out [my blog post](https://asankov.dev/blog/2022/12/22/implementing-the-java-stream-api-with-go-generics-part-1/) on the topic.
