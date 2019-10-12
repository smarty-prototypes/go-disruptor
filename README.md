Disruptor Overview
----------------------------

This is a port of the [LMAX Disruptor](https://github.com/LMAX-Exchange/disruptor) into the Go programming language. It retains the essence and spirit of the Disruptor and utilizes a lot of the same abstractions and concepts, but does not maintain the same API.

On a MacBook Pro (Intel Core i9-8950HK CPU @ 2.90GHz) using Go 1.13.1, it can pass many hundreds of millions of messages per second (yes, you read that right) from one goroutine to another goroutine.

Once initialized and running, one of the preeminent design considerations of the Disruptor is to process messages at a constant rate. It does this using two primary techniques. First, it avoids using locks at all costs which cause contention between CPU cores and prevents scaling the number of cores. Secondly, it produces no garbage by allowing the application to preallocate sequential space on a ring buffer. By avoiding garbage, the need for a garbage collection and the stop-the-world application pauses introduced can be almost entirely avoided.

In Go, the current implementation of a channel (`chan`) maintains a lock around send, receive, and `len` operations and it maxes out around 25 million messages per second for uncontended access&mdash;more than an orders of magnitude slower when compared to the Disruptor.  The same channel, when contended between OS threads only pushes about 7 million messages per second.

Example Usage
-------------

```
func main() {
    writer, reader := disruptor.New(
        disruptor.WithCapacity(BufferSize),
        disruptor.WithConsumerGroup(MyConsumer{}))

    // producer
    go func() {
        reservation := writer.Reserve(1)
        ringBuffer[sequence&RingBufferMask] = 42 // example of incoming value from a network operation such as HTTP, TCP, UDP, etc.
        writer.Commit(reservation, reservation)

        _ = reader.Close() // close the Reader once we're done producing messages
    }()

    reader.Read() // blocks until fully closed
}

type MyConsumer struct{}

func (m MyConsumer) Consume(lowerSequence, upperSequence int64) {
    for sequence := lowerSequence; sequence <= upperSequence; sequence++ {
        index := sequence&RingBufferMask
        message := ringBuffer[index]
        fmt.Println(message)
    }
}

var ringBuffer = [BufferSize]int

const (
    BufferSize = 1024 * 64
    BufferMask = BufferSize - 1
)
```

The above example is more complex than a typical channel implementation. When removing all the comments, wireup, and extra fluff to explain what's happening, the code is very concise.  In fact, a given "Publish" only takes three lines&mdash;`Reserve` a slot, update the ring buffer at that slot, and `Commit` the reserved sequence range.  On the consumer side, there's a `for`-loop to handle all incoming items into your application.  Again, this not quite as concise as a channel (nor as flexible), but it's much, much faster.

Benchmarks
----------------------------
Each of the following benchmark tests sends an incrementing sequence message from one goroutine to another. The receiving goroutine asserts that the message is received is the expected incrementing sequence value. Any failures cause a panic.

* CPU: `Intel Core i9-8950HK CPU @ 2.90GHz`
* Operation System: `OS X 10.14.6`
* Go Runtime: `Go 1.13.1`

Scenario | Per Operation Time
-------- | ------------------
Channels: Buffered, Blocking | 53 ns
Channels: Buffered, Blocking | 54 ns
Channels: Buffered, Blocking, Contended Write | 121 ns
Channels: Buffered, Non-blocking | 60 ns
Channels: Buffered, Non-blocking | 68 ns
Channels: Buffered, Non-blocking Contended Write | 205 ns
Disruptor: Sequencer, Reserve One | 9.73 ns
Disruptor: Sequencer, Reserve Many | 1.64 ns
Disruptor: Sequencer, Reserve One, Multiple Readers | 10.4 ns
Disruptor: Sequencer, Reserve Many, Multiple Readers | 1.59 ns

When In Doubt, Use Channels
----------------------------
Despite Go channels being significantly slower than the Disruptor, channels should still be considered the easiest, best, and most desirable choice for the vast majority of all use cases. The Disruptor's target use case is ultra-low latency environments where application response times are measured in nanoseconds and where stable, consistent latency is paramount and latency spikes cannot be tolerated.

Pre-Alpha
---------
This code is currently experimental and is not recommended for production environments. It does not have any unit tests and is only meant serve as a proof of concept that the Disruptor is possible on the Go runtime despite some of the limits imposed by the [Go memory model](http://golang.org/ref/mem).

We are very interested to receive feedback on this project and how performance can be improved using subtle techniques such as additional cache line padding, memory alignment, utilizing a pointer vs a struct in a given location, replacing less optimal techniques with more optimal ones, especially in the performance-critical operations `Reserve`/`Commit` in the `DefaultWriter` struct as well as the `Read` operation of the `DefaultReader`.

Caveats
-------
In the Java-based Disruptor implementation, a ring buffer is created, preallocated, and prepopulated with instances of the class which serve as the message type to be transferred between threads.  Because Go lacks generics, we have opted to not interact with ring buffers at all within the library code. This has the benefit of avoiding an unnecessary boxing/unboxing type conversions during the receipt of a given message.  It also means that it is the responsibility of the application developer to create and populate their particular ring buffer during application wireup. Pre-populating the ring buffer at startup should ensure contiguous memory allocation for all items in the various ring buffer slots, whereas on-the-fly creation may introduce gaps in the memory allocation and subsequent CPU cache misses which introduce latency spikes.

The reference to the ring buffer can (but need not be) be scoped as a package-level variable. The reason for this is that any given application should have very few Disruptor instances. The instances are designed to be created at startup and stopped during shutdown. They are not typically meant to be created ad-hoc and passed around like channel instances. It is the responsibility of the application developer to manage references to the ring buffer instances such that the producer can push messages to the buffer and consumers can receive messages from the buffer.
