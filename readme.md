This is a port of the LMAX Disruptor into the Go programming language.

It retains the essence and spirit of the Disruptor and utilizes a lot of the same abstractions and concepts, but does not maintain the same API.

On my MacBook Pro (early 2013) using Go 1.2.1, I was able to push 225 million messages per second from one goroutine to another goroutine. The message being transferred between the two CPU cores was simply the incrementing sequence number, but literally could be anything. Note that your mileage may vary and that different operating systems can introduce significant “jitter” into the application by taking control of the CPU. Linux and Windows have the ability to assign a given process to specific CPU cores which reduces jitter significantly.

Once initialized and running, one of the preeminent design considerations of the Disruptor is to produce no garbage thus avoiding the need for GC altogether and to avoid locks at all costs.

This code is pre-Alpha stage and is not supported or recommended for production environments.  That being said, it has been run non-stop for days without exposing any race conditions. Also, it does not yet contain any unit tests and is meant to be spike code to serve as a proof of concept that the Disruptor is, in fact possible, on the Go runtime despite some of the limits imposed by the Go memory model.

We are very interested to receive feedback on this project and how performance can be improved using subtle techniques such as cache line padding, utilizing a pointer vs a struct in a given location, replacing less optimal techniques with more optimal ones, especially in the performance critical paths of `Next` in the `Sequencer` and `Process` in the `Worker`

One last caveat worth noting.  In the Java-based Disruptor implementation, a ring buffer is created and preallocated and pre-populated with instances of the classes that serve as the message type to be transferred between threads.  Because Go lacks generics, we have opted to touching the ring buffer at all within the library code. This has the benefit of avoiding a require casting during the receipt of a given message from type `interface{}` to a concrete type.  This also means that it is the responsibility of the library user to create and populate their particular ring buffer during application wireup.

The reference to the ring buffer can easily be scoped as a package-level variable. The reason for this is any given application should only be running very few Disruptor instances. The instances are designed to be created at startup and killed at shutdown. They are not typically meant to be created “on the fly” and passed around. In any case, it is the responsibility of the application developer to manage references to the ring buffer instance such that the producer can push messages in and the consumers can receive messages out.
