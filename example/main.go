package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const MaxConsumers = 1

func main() {
	runtime.GOMAXPROCS(MaxConsumers + 3)

	writerCursor := disruptor.NewCursor()
	writerBarrier := disruptor.NewBarrier(writerCursor)

	readerCursors1 := startReaders(writerBarrier, writerCursor)
	readerBarrier1 := disruptor.NewBarrier(readerCursors1...)

	readerCursors2 := startReaders(readerBarrier1, writerCursor)
	readerBarrier2 := disruptor.NewBarrier(readerCursors2...)

	readerCursors3 := startReaders(readerBarrier2, writerCursor)
	readerBarrier3 := disruptor.NewBarrier(readerCursors3...)

	writer := disruptor.NewWriter(writerCursor, RingSize, readerBarrier3)
	publish(writer)
}
func startReaders(upstreamBarrier disruptor.Barrier, writerCursor *disruptor.Cursor) (readerCursors []*disruptor.Cursor) {
	for i := 0; i < MaxConsumers; i++ {
		readerCursor := disruptor.NewCursor()
		readerCursors = append(readerCursors, readerCursor)
		reader := disruptor.NewReader(upstreamBarrier, writerCursor, readerCursor)

		// wildly sporadic latency for single-item publish, e.g. 2 seconds, 65 ms, etc.
		// faster for 2-3+ items per publish
		go consume(reader)

		// constant time regardless of the number of items
		// go easyConsume(disruptor.NewEasyReader(reader, NewExampleConsumerHandler()))
	}

	return readerCursors
}
