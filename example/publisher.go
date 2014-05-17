package main

import "github.com/smartystreets/go-disruptor"

func publish(sequencer *disruptor.SingleProducerSequencer) {
	for {
		sequence := sequencer.Next(1)
		ringBuffer[sequence&RingMask] = sequence
		sequencer.Publish(sequence)
	}
}
