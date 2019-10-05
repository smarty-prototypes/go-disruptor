package benchmarks

import (
	"log"

	"github.com/smartystreets-prototypes/go-disruptor"
)

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := ringBuffer[lower&RingBufferMask]
		if message != lower {
			log.Panicf("race condition: Sequence: %d, Message: %d", lower, message)
		}
		lower++
	}
}

func build(consumers ...disruptor.Consumer) (disruptor.Sequencer, disruptor.ListenCloser) {
	wireup, err := disruptor.New(
		disruptor.WithCapacity(RingBufferSize),
		disruptor.WithConsumerGroup(consumers...),
	)
	if err != nil {
		panic(err)
	}

	return wireup.Build()
}

const (
	RingBufferSize   = 1024 * 64
	RingBufferMask   = RingBufferSize - 1
	ReserveOne       = 1
	ReserveMany      = 16
	ReserveManyDelta = ReserveMany - 1
)

var ringBuffer = [RingBufferSize]int64{}
