package main

import (
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(3)

	producerSequence := NewSequence()
	consumerSequence1 := NewSequence()
	consumerSequence2 := NewSequence()

	producerBarrier := NewBarrier(producerSequence)
	// consumerBarrier := NewBarrier(consumerSequence1)
	consumerBarrier := NewBarrier(consumerSequence1, consumerSequence2)

	sequencer := NewSingleProducerSequencer(producerSequence, RingSize, consumerBarrier)
	go consume(1, producerBarrier, producerSequence, consumerSequence1)
	go consume(2, producerBarrier, producerSequence, consumerSequence2)
	// time.Sleep(time.Millisecond * 10)

	// started := time.Now()
	for i := int64(0); i < MaxSequenceValue; i++ {
		// fmt.Printf("Producer:: Attempting to claim next sequence.\n")
		ticket := sequencer.Next(1)
		ringBuffer[ticket&RingMask] = ticket
		//time.Sleep(time.Nanosecond * 500)
		//runtime.Gosched()
		// fmt.Printf("Producer:: Claimed sequence: %d, Publishing...\n", ticket)
		sequencer.Publish(ticket)
		// fmt.Printf("Producer:: Claimed sequence: %d, Published\n", ticket)
		// if ticket%Mod == 0 {
		// 	finished := time.Now()
		// 	elapsed := finished.Sub(started)
		// 	fmt.Println(ticket, elapsed)
		// 	started = time.Now()
		// }
	}
}

func consume(name int, barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, TestHandler{name}, source, sequence)

	for {
		// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Attempting to process messages.\n", name)
		worker.Process(name)
	}
}

const Mod = 1000000 * 100 // 1 million * 10
const RingSize = 1024 * 256
const RingMask = RingSize - 1

// RingMask = RingSize - 1 // slightly faster than a mod operation
// 0&3 = 0
// 1&3 = 1
// 2&3 = 2
// 3&3 = 3
// 4&3 = 0

var ringBuffer [RingSize]int64

type TestHandler struct{ name int }

func (this TestHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]
	//fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Sequence: %d, Message: %d\n", this.name, sequence, message)
	if message != sequence {
		fmt.Printf("\t\t\t\t\t\t\t\t\tERROR Consumer %d:: Sequence: %d, Message: %d\n", this.name, sequence, message)
		panic(fmt.Sprintf("Consumer %d:: Sequence: %d, Message: %d\n", this.name, sequence, message))
	} else if sequence%Mod == 0 {
		fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Sequence: %d, Message: %d\n", this.name, sequence, message)
	}
}
