package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/go-disruptor"
)

const Mod = 1000000 * 10 // 1 million * N

func consume(reader *disruptor.Reader) {
	started := time.Now()

	for {
		sequence, remaining := reader.Receive()
		if remaining >= 0 {
			for remaining >= 0 {

				if sequence%Mod == 0 {
					finished := time.Now()
					fmt.Println(sequence, finished.Sub(started))
					started = time.Now()
				}

				if sequence != ringBuffer[sequence&RingMask] {
					message := ringBuffer[sequence&RingMask]
					panic(fmt.Sprintf("Sequence: %d, Message %d\n", sequence, message))
				}

				remaining--
				sequence++
			}
			reader.Commit(sequence)
		}
	}
}
func easyConsume(reader *disruptor.EasyReader) {
	for {
		reader.Receive()
	}
}
