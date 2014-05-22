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

	debug := make([]int64, 5)
	debug = debug[0:]

	for {
		sequence, remaining := reader.Receive()
		if remaining >= 0 {
			for remaining >= 0 {

				if sequence%ReportingFrequency == 0 {
					finished := time.Now()
					fmt.Println(sequence, finished.Sub(started))
					started = time.Now()
				}

				message := ringBuffer[sequence&RingMask]
				if sequence != message {
					for i := sequence - 5; i <= sequence; i++ {
						debug = append(debug, ringBuffer[i&RingMask])
					}

					alert := fmt.Sprintf("Race Condition::Sequence: %d, Message %d\n", sequence, message)
					fmt.Println(alert)
					fmt.Println("Partial Ring Buffer Snapshot:", debug)
					panic(alert)
				}

				//ringBuffer[sequence&RingMask] = sequence % 2

				remaining--
				sequence++
			}
			reader.Commit(sequence - 1)
		} else {
		}
	}
}

func consume2(reader *disruptor.Reader) {
	for {
		sequence, remaining := reader.Receive()
		if remaining >= 0 {
			for remaining >= 0 {
				message := ringBuffer[sequence&RingMask]
				if message != sequence%2 {
					alert := fmt.Sprintf("Race Condition (Layer 2)::Sequence: %d, Message %d\n", sequence, message)
					fmt.Print(alert)
					panic(alert)
				}

				remaining--
				sequence++
			}
			reader.Commit(sequence - 1)
		} else {
		}
	}
}
