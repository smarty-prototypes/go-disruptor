package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/go-disruptor"
)

const Mod = 1000000 * 10 // 1 million * N

func consume(writerBarrier disruptor.Barrier, writerCursor, readerCursor *disruptor.Cursor) {
	reader := disruptor.NewReader(writerBarrier, writerCursor, readerCursor)
	started := time.Now()

	for {
		sequence, remaining := reader.Receive()
		if remaining == disruptor.Idle {
		} else if remaining == disruptor.Gating {
		} else {
			for ; remaining >= 0; remaining-- {
				sequence++

				if sequence%Mod == 0 {
					finished := time.Now()
					fmt.Println(sequence, finished.Sub(started))
					started = time.Now()
				}

				if sequence != ringBuffer[sequence&RingMask] {
					message := ringBuffer[sequence&RingMask]
					panic(fmt.Sprintf("Sequence: %d, Message %d\n", sequence, message))
				}

			}

			reader.Commit(sequence)
		}
	}
}
