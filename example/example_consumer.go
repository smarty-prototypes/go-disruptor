package main

import (
	"fmt"

	"github.com/smartystreets/go-disruptor"
)

func consume0(reader *disruptor.SimpleReader) {
	for {
		reader.Receive()
	}
}
func consume1(reader *disruptor.Reader) {
	// started := time.Now()

	fmt.Printf("\t\t\t\t\t[CONSUMER] Starting consumer...\n")
	for {
		sequence, remaining := reader.Receive()
		if remaining >= 0 {
			fmt.Printf("\t\t\t\t\t[CONSUMER] Received messages starting at sequence %d, with %d messages remaining\n", sequence, remaining)

			for remaining >= 0 {
				// if sequence%ReportingFrequency == 0 {
				// 	finished := time.Now()
				// 	fmt.Println(sequence, finished.Sub(started))
				// 	started = time.Now()
				// }

				message := ringBuffer[sequence&RingMask]
				fmt.Printf("\t\t\t\t\t[CONSUMER] Consuming sequence %d. Message Payload: %d\n", sequence, message)
				if sequence != message {
					alert := fmt.Sprintf("--------------\n\t\t\t\t\t[CONSUMER] ***Race Condition***::Sequence: %d, Message %d\n", sequence, message)
					fmt.Println(alert)
					panic(alert)
				}

				//ringBuffer[sequence&RingMask] = sequence % 2

				remaining--
				sequence++
			}

			fmt.Println("\t\t\t\t\t[CONSUMER] All messages consumed, committing up to sequence ", sequence-1)
			reader.Commit(sequence - 1)
		} else {
			if remaining == disruptor.Gating {
				fmt.Println("\t\t\t\t\t[CONSUMER] Consumer gating at sequence", sequence)
			} else if remaining == disruptor.Idling {
				fmt.Println("\t\t\t\t\t[CONSUMER] Consumer idling at sequence", sequence)
			}
			//time.Sleep(time.Millisecond * 10)
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
