package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const MaxConsumers = 1

func main() {
	runtime.GOMAXPROCS(MaxConsumers + 2)

	writerCursor := disruptor.NewCursor()
	writerBarrier := disruptor.NewBarrier(writerCursor)

	// writerCursor.Store(disruptor.MaxSequenceValue)
	// startReaders(writerBarrier, writerCursor)
	// time.Sleep(time.Second * 10)

	readerCursors := startReaders(writerBarrier, writerCursor)
	readerBarrier := disruptor.NewBarrier(readerCursors...)

	writer := disruptor.NewWriter(writerCursor, RingSize, readerBarrier)
	publish(writer)
}
func startReaders(writerBarrier disruptor.Barrier, writerCursor *disruptor.Cursor) (readerCursors []*disruptor.Cursor) {
	for i := 0; i < MaxConsumers; i++ {
		readerCursor := disruptor.NewCursor()
		readerCursors = append(readerCursors, readerCursor)
		go consume(writerBarrier, writerCursor, readerCursor)
	}

	return readerCursors
}
