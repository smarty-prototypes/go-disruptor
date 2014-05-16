package main

import (
	"fmt"
	"time"

	"bitbucket.org/jonathanoliver/go-disruptor"
)

func publish(sequencer *disruptor.SingleProducerSequencer) {
	started := time.Now()
	for sequence := int64(0); sequence < MaxIterations; sequence++ {
		sequencer.Next(1)
		ringBuffer[sequence&RingMask] = sequence
		sequencer.Publish(sequence)
		if sequence%Mod == 0 && sequence > 0 {
			finished := time.Now()
			elapsed := finished.Sub(started)
			fmt.Println(sequence, elapsed)
			started = time.Now()
		}
	}
}

const MaxIterations = disruptor.MaxSequenceValue
const Mod = 1000000 * 100 // 1 million * N
