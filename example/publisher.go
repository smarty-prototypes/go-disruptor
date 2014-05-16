package main

import (
	"fmt"
	"time"

	"bitbucket.org/jonathanoliver/go-disruptor"
)

func publish(sequencer *disruptor.SingleProducerSequencer) {
	started := time.Now()
	for i := int64(0); i < MaxIterations; i++ {
		sequencer.Next(1)
		//ringBuffer[i&RingMask] = i
		sequencer.Publish(i)
		if i%Mod == 0 && i > 0 {
			finished := time.Now()
			elapsed := finished.Sub(started)
			fmt.Println(i, elapsed)
			started = time.Now()
		}
	}
}

const MaxIterations = disruptor.MaxSequenceValue
const Mod = 1000000 * 100 // 1 million * N
