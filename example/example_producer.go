package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/go-disruptor"
)

func publish(id int, writer *disruptor.SharedWriter) {

	fmt.Printf("[PRODUCER %d] Starting producer...\n", id)

	for {
		// TODO: return lower, upper instead? or some kind of struct "Reservation"
		// upon which commit can be invoked?
		if sequence := writer.Reserve(id, ItemsToPublish); sequence != disruptor.Gating {
			// fmt.Printf("[PRODUCER %d] Writing %d to slot %d\n", id, sequence, sequence)
			ringBuffer[sequence&RingMask] = sequence
			// fmt.Printf("[PRODUCER %d] Committing from sequence %d\n", id, sequence)
			writer.Commit(sequence, sequence+ItemsToPublish-1)
		} else {
			// 	// fmt.Printf("[PRODUCER %d] Gating\n", id)
			time.Sleep(time.Millisecond * 10)
		}
	}
}
