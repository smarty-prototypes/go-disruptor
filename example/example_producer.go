package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.Writer) {
	for {
		_, upper := writer.Reserve(1)
		if upper == disruptor.Gating {
			continue
		}

		// ringBuffer[(upper-2)&RingMask] = upper - 2
		// ringBuffer[(upper-1)&RingMask] = upper - 1
		ringBuffer[upper&RingMask] = upper

		writer.Commit(upper)
	}
}
