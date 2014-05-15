package main

import (
	"fmt"
	"time"

	"bitbucket.org/jonathanoliver/go-disruptor"
)

func publish(sequencer *disruptor.SingleProducerSequencer) {
	started := time.Now()
	for i := int64(0); i < MaxIterations; i++ {
		ticket := sequencer.Next(1)
		ringBuffer[ticket&RingMask] = ticket
		sequencer.Publish(ticket)
		if ticket%Mod == 0 && ticket > 0 {
			finished := time.Now()
			elapsed := finished.Sub(started)
			fmt.Println(ticket, elapsed)
			started = time.Now()
		}
	}
}

const MaxIterations = disruptor.MaxSequenceValue
const Mod = 1000000 * 10 // 1 million * N
