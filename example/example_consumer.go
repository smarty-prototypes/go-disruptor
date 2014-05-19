package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/go-disruptor"
)

const Mod = 1000000 * 10 // 1 million * N

func consume(writerBarrier disruptor.Barrier, writerCursor, readerCursor *disruptor.Cursor) {
	reader := disruptor.NewReader(writerBarrier, &ConsumerHandler{time.Now()}, writerCursor, readerCursor)

	for {
		reader.Process()
	}
}
func (this *ConsumerHandler) Consume(sequence, remaining int64) {
	if sequence%Mod == 0 {
		finished := time.Now()
		fmt.Println(sequence, finished.Sub(this.started))
		this.started = time.Now()
	}

	if sequence != ringBuffer[sequence&RingMask] {
		message := ringBuffer[sequence&RingMask]
		panic(fmt.Sprintf("Race condition--Cursor: %d, Message: %d\n", sequence, message))
	}
}

type ConsumerHandler struct{ started time.Time }
