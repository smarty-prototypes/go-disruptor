package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(3)

	producerSequence := NewSequence()
	consumerSequence1 := NewSequence()
	consumerSequence2 := NewSequence()

	producerBarrier := NewBarrier(producerSequence)
	consumerBarrier := NewBarrier(consumerSequence1, consumerSequence2)

	sequencer := NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)
	go consume(producerBarrier, producerSequence, consumerSequence1)
	go consume(producerBarrier, producerSequence, consumerSequence2)

	started := time.Now()
	for i := int64(0); i < 10; i++ {
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

	time.Sleep(time.Millisecond * 100)
}

func consume(barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, TestHandler{}, source, sequence)

	for {
		worker.Process()
	}
}

const Mod = 1000000 * 100 // 1 million * 10
const RingSize = 2
const RingMask = RingSize - 1

var ringBuffer [RingSize]int64

type TestHandler struct{}

func (this TestHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]
	if message != sequence {
		fmt.Printf("ERROR Consumer:: Sequence: %d, Message: %d\n", sequence, message)
		panic(fmt.Sprintf("Consumer:: Sequence: %d, Message: %d\n", sequence, message))
	}

	if sequence%Mod == 0 {
		//fmt.Printf("Consumer:: Sequence: %d, Message: %d\n", sequence, message)
	}
}
