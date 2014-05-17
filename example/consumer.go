package main

import (
	"time"

	"github.com/smartystreets/go-disruptor"
)

func consume(writerBarrier *disruptor.Barrier, writerCursor, readerCursor *disruptor.Cursor) {
	reader := disruptor.NewReader(writerBarrier, &ConsumerHandler{time.Now()}, writerCursor, readerCursor)

	for {
		reader.Process()
	}
}
