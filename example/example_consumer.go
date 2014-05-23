package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/go-disruptor"
)

func consume0(reader *disruptor.SimpleReader) {
	for {
		reader.Receive()
	}
}
func consume1(reader *disruptor.Reader) {
	started := time.Now()

	for {
		if lower, upper := reader.Receive(); upper >= 0 {
			for sequence := lower; sequence <= upper; sequence++ {
				if sequence%ReportingFrequency == 0 {
					finished := time.Now()
					fmt.Println(sequence, finished.Sub(started))
					started = time.Now()
				}

				message := ringBuffer[sequence&RingMask]
				if sequence != message {
					alert := fmt.Sprintf("***Race Condition***\tSequence: %d, Message %d\n", sequence, message)
					fmt.Println(alert)
					panic(alert)
				}

				ringBuffer[sequence&RingMask] = sequence % 2
			}

			reader.Commit(upper)
		}
	}
}

func consume2(reader *disruptor.Reader) {
	for {
		if lower, upper := reader.Receive(); upper >= 0 {
			for sequence := lower; sequence <= upper; sequence++ {
				message := ringBuffer[sequence&RingMask]
				if message != sequence%2 {
					alert := fmt.Sprintf("Race Condition (Layer 2)::Sequence: %d, Message %d\n", sequence, message)
					fmt.Print(alert)
					panic(alert)
				}
			}
			reader.Commit(upper)
		}
	}
}
