package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.Writer) {
	// runtime.LockOSThread()

	for {
		upper := writer.Reserve(3)
		if upper != disruptor.Gating {

			ringBuffer[(upper-2)&RingMask] = upper - 2
			ringBuffer[(upper-1)&RingMask] = upper - 1
			ringBuffer[upper&RingMask] = upper
			writer.Commit(upper)
		}
	}
}
