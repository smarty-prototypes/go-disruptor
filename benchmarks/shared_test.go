package benchmarks

import (
	"log"

	"github.com/smartystreets-prototypes/go-disruptor"
)

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	var message int64
	for lower <= upper {
		message = ringBuffer[lower&RingBufferMask]
		if message != lower {
			log.Panicf("race condition: Sequence: %d, Message: %d", lower, message)
		}
		lower++
	}
}

func build(consumers ...disruptor.Consumer) (disruptor.Sequencer, disruptor.ListenCloser) {
	return disruptor.New(
		disruptor.WithCapacity(RingBufferSize),
		disruptor.WithConsumerGroup(consumers...))
}

const (
	RingBufferSize = 1024 * 64
	RingBufferMask = RingBufferSize - 1
	ReserveOne     = 1
	ReserveMany    = 16
)

var ringBuffer = [RingBufferSize]int64{}
