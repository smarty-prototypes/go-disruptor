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
		sequence, remaining := reader.Receive()
		if remaining >= 0 {
			for remaining >= 0 {

				if sequence%ReportingFrequency == 0 {
					finished := time.Now()
					fmt.Println(sequence, finished.Sub(started))
					started = time.Now()
				}

				if sequence != ringBuffer[sequence&RingMask] {
					message := ringBuffer[sequence&RingMask]
					alert := fmt.Sprintf("Race Condition::Sequence: %d, Message %d\n", sequence, message)
					fmt.Print(alert)
					panic(alert)
				}

				ringBuffer[sequence&RingMask] = sequence % 2

				remaining--
				sequence++
			}
			reader.Commit(sequence - 1)
		} else {
			//time.Sleep(time.Nanosecond)
		}
	}
}

func consume2(reader *disruptor.Reader) {
	for {
		sequence, remaining := reader.Receive()
		if remaining >= 0 {
			for remaining >= 0 {
				index := sequence & RingMask
				message := ringBuffer[index]
				if message != sequence%2 {
					panic("Race condition!")
				}

				remaining--
				sequence++
			}
			reader.Commit(sequence - 1)
		} else {
			//time.Sleep(time.Nanosecond)
		}
	}
}
