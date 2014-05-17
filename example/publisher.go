package main

import "bitbucket.org/jonathanoliver/go-disruptor"

func publish(sequencer *disruptor.SingleProducerSequencer) {
	for sequence := int64(0); sequence < MaxIterations; sequence++ {
		sequencer.Next(1)
		ringBuffer[sequence&RingMask] = sequence
		sequencer.Publish(sequence)
	}
}

const MaxIterations = disruptor.MaxSequenceValue
const Mod = 1000000 * 10 // 1 million * N
