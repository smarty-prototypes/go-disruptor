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
	go consume(1, producerBarrier, producerSequence, consumerSequence1)
	// go consume(2, producerBarrier, producerSequence, consumerSequence2)

	// started := time.Now()
	for i := int64(0); i < MaxIterations; i++ {
		// fmt.Printf("Producer:: Attempting to claim next sequence.\n")
		ticket := sequencer.Next(1)
		// fmt.Printf("Producer:: Claimed sequence: %d, Assigning slot...\n", ticket)
		ringBuffer[ticket&RingMask] = ticket
		// fmt.Printf("Producer:: Claimed sequence: %d, Publishing...\n", ticket)
		sequencer.Publish(ticket)
		// fmt.Printf("Producer:: Claimed sequence: %d, Published\n", ticket)
		// if ticket%Mod == 0 && ticket > 0 {
		// 	finished := time.Now()
		// 	elapsed := finished.Sub(started)
		// 	fmt.Println(ticket, elapsed)
		// 	started = time.Now()
		// }
	}

	time.Sleep(time.Nanosecond * 50)
	// fmt.Println("Graceful shutdown\n-------------------------------------------------\n")
}

func consume(name int, barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, TestHandler{}, source, sequence)

	for {
		// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Attempting to process messages.\n", name)
		if worker.Process(name)+1 > MaxIterations {
			break
		}
	}
}

const MaxIterations = MaxSequenceValue
const Mod = 1000000 * 1 // 1 million * N
const RingSize = 2
const RingMask = RingSize - 1

var ringBuffer [RingSize]int64

type TestHandler struct{}

func (this TestHandler) Consume(sequence, remaining int64) {
	message := ringBuffer[sequence&RingMask]
	// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Sequence: %d, Message: %d\n", sequence, message)

	if message != sequence {
		fmt.Printf("\t\t\t\t\t\t\t\t\tERROR Consumer:: Sequence: %d, Message: %d\n", sequence, message)
		panic(fmt.Sprintf("Consumer:: Sequence: %d, Message: %d\n", sequence, message))
	}

	if sequence%Mod == 0 && sequence > 0 {
		fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer:: Sequence: %d, Message: %d\n", sequence, message)
	}
}
