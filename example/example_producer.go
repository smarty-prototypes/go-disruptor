package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.Writer) {

	for {
		if lower, upper := writer.Reserve(ItemsToPublish); upper != disruptor.Gating {
			// for sequence := lower; sequence <= upper; sequence++ {
			// 	ringBuffer[sequence&RingMask] = sequence
			// }
			// ringBuffer[lower&RingMask] = lower
			writer.Commit(lower, upper)
		}
	}
}
