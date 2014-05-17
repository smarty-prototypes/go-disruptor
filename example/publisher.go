package main

import "github.com/smartystreets/go-disruptor"

func publish(writer *disruptor.Writer) {
	for {
		sequence := writer.Next(1)
		ringBuffer[sequence&RingMask] = sequence
		writer.Publish(sequence)
	}
}
