package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/smartystreets-prototypes/go-disruptor"
)

const (
	BufferSize   = 1024 * 64
	BufferMask   = BufferSize - 1
	Iterations   = 1000000 * 100
	Reservations = 16
)

var ring = [BufferSize]int64{}

func main() {
	runtime.GOMAXPROCS(2)

	controller := disruptor.
		Configure(BufferSize).
		WithConsumerGroup(SampleConsumer{}).
		Build()

	go controller.Listen()

	started := time.Now()
	publish(controller.Writer())
	finished := time.Now()

	_ = controller.Close()
	fmt.Println(Iterations, finished.Sub(started))
}

func publish(writer disruptor.Writer) {
	sequence := disruptor.InitialSequenceValue
	for sequence <= Iterations {
		sequence = writer.Reserve(Reservations)
		for lower := sequence - Reservations + 1; lower <= sequence; lower++ {
			ring[lower&BufferMask] = lower
		}

		writer.Commit(sequence-Reservations+1, sequence)
	}
}

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := ring[lower&BufferMask]
		if message != lower {
			fmt.Println("Race condition", message, lower)
			panic("Race condition")
		}
		lower++
	}
}
