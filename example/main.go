package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const MaxConsumers = 1

func main() {
	runtime.GOMAXPROCS(MaxConsumers + 1)

	writerCursor := disruptor.NewCursor()
	writerBarrier := disruptor.NewBarrier(writerCursor)

	readerCursors := startReaders(writerBarrier, writerCursor)
	readerBarrier := disruptor.NewBarrier(readerCursors...)

	writer := disruptor.NewWriter(writerCursor, RingSize, readerBarrier)
	publish(writer)
}
func startReaders(barrier disruptor.Barrier, sequence *disruptor.Cursor) (readerCursors []*disruptor.Cursor) {
	for i := 0; i < MaxConsumers; i++ {
		sequence := disruptor.NewCursor()
		readerCursors = append(readerCursors, sequence)
		go consume(barrier, sequence, sequence)
	}

	return readerCursors
}
