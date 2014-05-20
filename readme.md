Notes
=====

This is a port of the LMAX Disruptor into the Go programming language. It retains the essence and spirit of the Disruptor and utilizes a lot of the same abstractions and concepts, but does not maintain the same API.

On my MacBook Pro (Intel Core i7 3740QM @ 2.7 Ghz) using Go 1.2.1, I was able to push over **700 million messages per second** (yes, you read that right) from one goroutine to another goroutine. The message being transferred between the two CPU cores was simply the incrementing sequence number, but literally could be anything. Note that your mileage may vary and that different operating systems can introduce significant “jitter” into the application by taking control of the CPU. Linux and Windows have the ability to assign a given process to specific CPU cores which reduces jitter significantly.  Parenthetically, when the Disruptor code is compiled and run on a Nexus 5, it can push about 15-20 million messages per second.

Once initialized and running, one of the preeminent design considerations of the Disruptor is to produce no garbage thus avoiding the need for GC altogether and to avoid locks at all costs. The current channel implementation maintains a big, fat lock around enqueue/dequeue operations and maxes out on the aforementioned hardware at about 25M messages per second for uncontended access-—more than an order of magnitude slower when compared to the Disruptor.  The same channel, when contended between OS threads (GOMAXPROCS=2 or more) only pushes about 7 million messages per second.

When In Doubt, Use Channels
----------------------------
Despite Go's channels being significantly slower than the Disruptor, channels should still be considered the best and most desirable choice for the vast majority of all use cases. The Disruptor's target use case us ultra-low latency environments where application response times are measured in nanoseconds and where stable, consistent latency is paramount.

Pre-Alpha
---------
This code is pre-Alpha stage and is not supported or recommended for production environments. That being said, it has been run non-stop for days without exposing any race conditions. Also, it does not yet contain any unit tests and is meant to be spike code to serve as a proof of concept that the Disruptor is, in fact possible, on the Go runtime despite some of the limits imposed by the Go memory model. The goal is to have an alpha release by mid June 2014 and a series of beta releases each month thereafter until we are satisfied. Following this a release will be created and supported moving forward.

We are very interested to receive feedback on this project and how performance can be improved using subtle techniques such as additional cache line padding, utilizing a pointer vs a struct in a given location, replacing less optimal techniques with more optimal ones, especially in the performance critical paths of `Next` in the `Sequencer` and `Process` in the `Reader`

Caveats
-------
One last caveat worth noting.  In the Java-based Disruptor implementation, a ring buffer is created,  preallocated, and prepopulated with instances of the class which serve as the message type to be transferred between threads.  Because Go lacks generics, we have opted to not interact with ring buffers at all within the library code. This has the benefit of avoiding a required cast during the receipt of a given message from type `interface{}` to a concrete type.  It also means that it is the responsibility of the application developer to create and populate their particular ring buffer during application wireup.

The reference to the ring buffer can easily be scoped as a package-level variable. The reason for this is any given application should only be running very few Disruptor instances. The instances are designed to be created at startup and stopped during shutdown. They are not typically meant to be created “on the fly” and passed around. In any case, it is the responsibility of the application developer to manage references to the ring buffer instances such that the producer can push messages in and the consumers can receive messages out.