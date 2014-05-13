package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(2)

	producerSequence := NewSequence()
	consumerSequence := NewSequence()

	producerBarrier := NewBarrier(producerSequence)
	consumerBarrier := NewBarrier(consumerSequence)

	sequencer := NewMonoSequencer(producerSequence, RingSize, consumerBarrier)

	go consume(producerBarrier, producerSequence, consumerSequence)

	for i := int64(0); i < MaxSequenceValue; i++ {
		ticket := sequencer.Next(1)
		ringBuffer[ticket&RingMask] = ticket
		sequencer.Publish(ticket)
	}
}

func consume(barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, TestHandler{}, source, sequence)

	for {
		worker.Process()
	}
}

const RingSize = 1024
const RingMask = RingSize - 1
const Mod = 1000 * 1000 * 10 // 10 Million
var ringBuffer [RingSize]int64

type TestHandler struct{}

func (this TestHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]
	if message != sequence {
		panic(fmt.Sprintf("\n", message, sequence))
	} else {
		if sequence%Mod == 0 {
			fmt.Println("Current Sequence:", sequence)
		}
	}
}
