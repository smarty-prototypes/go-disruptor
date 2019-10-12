package main

import (
	"fmt"
	"io"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func main() {
	writer, reader := disruptor.New(
		disruptor.WithCapacity(BufferSize),
		disruptor.WithConsumerGroup(MyConsumer{}))

	go publish(writer, reader)

	reader.Read()
}

func publish(writer disruptor.Writer, closer io.Closer) {
	for sequence := int64(0); sequence <= Iterations; {
		sequence = writer.Reserve(Reservations)

		for lower := sequence - Reservations + 1; lower <= sequence; lower++ {
			ringBuffer[lower&BufferMask] = lower
		}

		writer.Commit(sequence-Reservations+1, sequence)
	}

	_ = closer.Close()
}

// ////////////////////

type MyConsumer struct{}

func (this MyConsumer) Consume(lower, upper int64) {
	for ; lower <= upper; lower++ {
		message := ringBuffer[lower&BufferMask]
		if message != lower {
			panic(fmt.Errorf("race condition: %d %d", message, lower))
		}
	}
}

const (
	BufferSize   = 1024 * 64
	BufferMask   = BufferSize - 1
	Iterations   = 128 * 1024 * 32
	Reservations = 1
)

var ringBuffer = [BufferSize]int64{}
