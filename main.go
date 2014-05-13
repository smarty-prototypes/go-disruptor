package main

import (
	"fmt"
	"runtime"
	"time"
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
		if i%50000000 == 0 {
			fmt.Println(i)
		}
		sequencer.Publish(ticket)
		consumerSequence.Store(ticket)
	}
}

func consume(barrier Barrier, source, sequence *Sequence) {
	worker := NewWorker(barrier, nil, source, sequence)

	for {
		worker.Process()
		time.Sleep(time.Nanosecond)
	}
}

const RingSize = 1024 * 256
