package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.Writer) {
	for {
		sequence := writer.Reserve(4)
		if sequence != disruptor.Gating {
			ringBuffer[(sequence-3)&RingMask] = sequence - 3
			ringBuffer[(sequence-2)&RingMask] = sequence - 2
			ringBuffer[(sequence-1)&RingMask] = sequence - 1
			ringBuffer[sequence&RingMask] = sequence
			writer.Commit(sequence)
		}
	}
}
