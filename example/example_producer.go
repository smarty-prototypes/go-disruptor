package main

import (
	"fmt"

	"github.com/smartystreets/go-disruptor"
)

func publish(id int, writer *disruptor.SharedWriter) {

	fmt.Printf("[PRODUCER %d] Starting producer...\n", id)

	for {
		if sequence := writer.Reserve(id, ItemsToPublish); sequence != disruptor.Gating {
			fmt.Printf("[PRODUCER %d] Writing %d to slot %d\n", id, sequence, sequence)
			ringBuffer[sequence&RingMask] = sequence
			fmt.Printf("[PRODUCER %d] Committing sequence %d\n", id, sequence)
			writer.Commit(sequence)
		} else {
			// fmt.Println("[PRODUCER] Gating")
			//time.Sleep(time.Millisecond)
		}
	}
}
