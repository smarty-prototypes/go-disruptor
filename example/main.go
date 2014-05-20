package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const MaxConsumersPerGroup = 1
const MaxConsumerGroups = 1

func main() {
	runtime.GOMAXPROCS(MaxConsumerGroups*MaxConsumersPerGroup + 1)

	writerCursor := disruptor.NewCursor()
	writerBarrier := disruptor.NewBarrier(writerCursor)
	readerBarrier := startConsumerGroups(writerBarrier, writerCursor)
	writer := disruptor.NewWriter(writerCursor, RingSize, readerBarrier)
	publish(writer)
}

func startConsumerGroups(upstream disruptor.Barrier, writer *disruptor.Cursor) disruptor.Barrier {
	for i := 0; i < MaxConsumerGroups; i++ {
		upstream = startConsumerGroup(upstream, writer)
	}

	return upstream
}
func startConsumerGroup(upstreamBarrier disruptor.Barrier, writerCursor *disruptor.Cursor) disruptor.Barrier {
	readerCursors := []*disruptor.Cursor{}

	for i := 0; i < MaxConsumersPerGroup; i++ {
		readerCursor := disruptor.NewCursor()
		readerCursors = append(readerCursors, readerCursor)
		reader := disruptor.NewReader(upstreamBarrier, writerCursor, readerCursor)

		// wildly sporadic latency for single-item publish, e.g. 2 seconds, 65 ms, etc.
		// faster for 2-3+ items per publish
		go consume(reader)

		// constant time regardless of the number of items
		// go easyConsume(disruptor.NewEasyReader(reader, NewExampleConsumerHandler()))
	}

	return disruptor.NewBarrier(readerCursors...)
}
