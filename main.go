package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)

	producerSequence := NewSequence()
	consumerSequence1 := NewSequence()
	// consumerSequence2 := NewSequence()

	producerBarrier := NewBarrier(producerSequence)
	consumerBarrier := NewBarrier(consumerSequence1) //, consumerSequence2)

	sequencer := NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)
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

func consume(barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, TestHandler{}, source, sequence)

	for {
		worker.Process()
	}
}

const MaxIterations = MaxSequenceValue
const Mod = 1000000 * 1 // 1 million * N
const RingSize = 1024
const RingMask = RingSize - 1

var ringBuffer [RingSize]int64

type TestHandler struct{}

func (this TestHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]

	if message != sequence {
		text := fmt.Sprintf("[Consumer] ERROR Sequence: %d, Message: %d\n", sequence, message)
		fmt.Printf(text)
		panic(text)
	}

	if sequence%Mod == 0 && sequence > 0 {
		// fmt.Printf("[Consumer] Sequence: %d, Message: %d\n", sequence, message)
	}
}
