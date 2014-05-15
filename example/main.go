package main

import (
	"runtime"

	"bitbucket.org/jonathanoliver/go-disruptor"
)

func main() {
	runtime.GOMAXPROCS(2)

	producerSequence := disruptor.NewSequence()
	producerBarrier := disruptor.NewBarrier(producerSequence)

	consumerSequences := []*disruptor.Sequence{}
	for i := 0; i < 1; i++ {
		sequence := disruptor.NewSequence()
		consumerSequences = append(consumerSequences, sequence)
		go consume(producerBarrier, producerSequence, sequence)
	}

	consumerBarrier := disruptor.NewBarrier(consumerSequences...)
	sequencer := disruptor.NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)

	publish(sequencer)
}
