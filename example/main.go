package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const MaxConsumers = 1

func main() {
	runtime.GOMAXPROCS(MaxConsumers + 1)

	producerSequence := disruptor.NewSequence()
	producerBarrier := disruptor.NewBarrier(producerSequence)

	consumers := startConsumers(producerBarrier, producerSequence)
	consumerBarrier := disruptor.NewBarrier(consumers...)

	sequencer := disruptor.NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)
	publish(sequencer)
}
func startConsumers(barrier disruptor.Barrier, sequence *disruptor.Sequence) (consumers []*disruptor.Sequence) {
	for i := 0; i < MaxConsumers; i++ {
		sequence := disruptor.NewSequence()
		consumers = append(consumers, sequence)
		go consume(barrier, sequence, sequence)
	}

	return consumers
}
