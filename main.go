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

var ringBuffer [RingSize]int64

type TestHandler struct{}

func (this TestHandler) Consume(sequence, remaining int64) {
	if sequence%10000000 == 0 {
		fmt.Println(sequence)
	}
}
