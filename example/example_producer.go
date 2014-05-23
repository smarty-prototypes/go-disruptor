package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.Writer) {

	for {
		if lower, upper := writer.Reserve(ItemsToPublish); upper != disruptor.Gating {

			ringBuffer[lower&RingMask] = lower
			// for sequence := lower; sequence <= upper; sequence++ {
			// 	ringBuffer[sequence&RingMask] = sequence
			// }
			writer.Commit(lower, upper)
		}
	}
}
