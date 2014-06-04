package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/smartystreets/go-disruptor"
)

const (
	BufferSize   = 1024 * 64
	BufferMask   = BufferSize - 1
	Iterations   = 1000000 * 100
	Reservations = 1
)

var ring = [BufferSize]int64{}

func main() {
	runtime.GOMAXPROCS(2)

	controller := disruptor.
		Configure(BufferSize).
		WithConsumerGroup(SampleConsumer{}).
		Build()

	controller.Start()

	started := time.Now()
	publish(controller.Writer())
	finished := time.Now()

	controller.Stop()
	fmt.Println(Iterations, finished.Sub(started))
}

func publish(writer *disruptor.Writer) {
	sequence := disruptor.InitialSequenceValue
	for sequence <= Iterations {
		sequence = writer.Reserve(Reservations)
		ring[sequence&BufferMask] = sequence
		writer.Commit(sequence, sequence)
	}
}

// func publish(writer *disruptor.Writer) {
// 	sequence := disruptor.InitialSequenceValue
// 	for sequence <= Iterations {
// 		sequence += Reservations
// 		writer.Await(sequence)
// 		ring[sequence&BufferMask] = sequence
// 		writer.Commit(sequence, sequence)
// 	}
// }

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
