package main

import (
	"fmt"
	"runtime"
	"time"

	"bitbucket.org/jonathanoliver/go-disruptor"
)

func main() {
	runtime.GOMAXPROCS(2)

	producerSequence := disruptor.NewSequence()
	producerBarrier := disruptor.NewBarrier(producerSequence)

	consumerSequence1 := disruptor.NewSequence()
	// consumerSequence2 := disruptor.NewSequence()

	consumerBarrier := disruptor.NewBarrier(consumerSequence1) //, consumerSequence2)
	sequencer := disruptor.NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)

	go consume(producerBarrier, producerSequence, consumerSequence1)
	// go consume(producerBarrier, producerSequence, consumerSequence2)

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

func publish() {

}

const MaxIterations = disruptor.MaxSequenceValue
const Mod = 1000000 * 10 // 1 million * N
const RingSize = 1024 * 128
const RingMask = RingSize - 1

var ringBuffer [RingSize]int64
