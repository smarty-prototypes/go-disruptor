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
	//consumerSequence2 := NewSequence()

	producerBarrier := NewBarrier(producerSequence)
	consumerBarrier := NewBarrier(consumerSequence1)
	// consumerBarrier := NewBarrier(consumerSequence1, consumerSequence2)

	sequencer := NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)
	go consume(1, producerBarrier, producerSequence, consumerSequence1)
	// go consume(2, producerBarrier, producerSequence, consumerSequence2)

	started := time.Now()
	for i := int64(0); i < MaxSequenceValue; i++ {
		ticket := sequencer.Next(1)
		//ringBuffer[ticket&RingMask] = ticket
		sequencer.Publish(ticket)
		if ticket%Mod == 0 {
			finished := time.Now()
			elapsed := finished.Sub(started)
			fmt.Println(ticket, elapsed)
			started = time.Now()
		}
	}
}

func consume(name int, barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, TestHandler{name}, source, sequence)

	for {
		worker.Process()
	}
}

const Mod = 1000000 * 100 // 1 million * 100
const RingSize = 1024 * 256
const RingMask = RingSize - 1

var ringBuffer [RingSize]int64

type TestHandler struct{ name int }

func (this TestHandler) Consume(sequence, remaining int64) {
	// message := ringBuffer[sequence&RingMask]
	// if message != sequence {
	// 	//fmt.Printf("\t\t\t\t\t\t\t\t\tERROR Consumer %d:: Sequence: %d, Message: %d\n", this.name, sequence, message)
	// 	//panic(fmt.Sprintf("Consumer %d:: Sequence: %d, Message: %d\n", this.name, sequence, message))
	// } else if sequence%Mod == 0 {
	// 	//fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Sequence: %d, Message: %d\n", this.name, sequence, message)
	// }
}
