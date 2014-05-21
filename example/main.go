package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const MaxConsumersPerGroup = 1
const MaxConsumerGroups = 2

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
		upstream = startConsumerGroup(i, upstream, writer)
	}

	return upstream
}
func startConsumerGroup(group int, upstreamBarrier disruptor.Barrier, writerCursor *disruptor.Cursor) disruptor.Barrier {
	readerCursors := []*disruptor.Cursor{}

	for i := 0; i < MaxConsumersPerGroup; i++ {
		readerCursor := disruptor.NewCursor()
		readerCursors = append(readerCursors, readerCursor)
		reader := disruptor.NewReader(upstreamBarrier, writerCursor, readerCursor)

		// constant time regardless of the number of items
		// go consume0(disruptor.NewSimpleReader(reader, NewExampleConsumerHandler()))

		// wildly sporadic latency for single-item publish, e.g. 2 seconds, 65 ms, etc.
		// faster for 2-3+ items per publish
		// go consume1(reader)

		if group == 0 {
			go consume1(reader)
		} else if group == 1 {
			go consume2(reader)
		}
	}

	return disruptor.NewBarrier(readerCursors...)
}
