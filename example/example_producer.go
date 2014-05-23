package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.SharedWriter) {

	for {
		if lower, upper := writer.Reserve(ItemsToPublish); upper != disruptor.Gating {
			for sequence := lower; sequence <= upper; sequence++ {
				ringBuffer[sequence&RingMask] = sequence
			}
			writer.Commit(lower, upper)
		}
	}
}
