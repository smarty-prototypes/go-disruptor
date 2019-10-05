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
	Iterations   = 100000 * 100
	Reservations = 16
)

var ring = [BufferSize]int64{}

func main() {
	runtime.GOMAXPROCS(2)

	controller := disruptor.
		Configure(BufferSize).
		WithConsumerGroup(SampleConsumer{}).
		Build()

	started := time.Now()

	go func() {
		publish(controller.Sequencer())
		_ = controller.Close()
	}()

	controller.Listen() // blocks until complete
	finished := time.Now()
	fmt.Println(finished.Sub(started))
}

func publish(sequencer disruptor.Sequencer) {
	sequence := disruptor.InitialCursorSequenceValue
	for sequence <= Iterations {
		sequence = sequencer.Reserve(Reservations)
		for lower := sequence - Reservations + 1; lower <= sequence; lower++ {
			ring[lower&BufferMask] = lower
		}

		sequencer.Commit(sequence-Reservations+1, sequence)
	}
}

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := ring[lower&BufferMask]
		if message != lower {
			panic(fmt.Errorf("race condition: %d %d", message, lower))
		}
		lower++
	}
}
