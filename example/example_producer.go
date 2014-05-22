package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.SharedWriter) {
	for {
		sequence := writer.Reserve(ItemsToPublish)

		if sequence != disruptor.Gating {
			for lower := sequence - ItemsToPublish; lower < sequence; {
				lower++
				ringBuffer[(lower)&RingMask] = lower
			}

			writer.Commit(sequence)
		}
	}
}
