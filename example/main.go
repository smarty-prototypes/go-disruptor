package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const (
	MaxConsumersPerGroup = 1
	MaxConsumerGroups    = 1
	MaxProducers         = 2
	ItemsToPublish       = 1
	ReportingFrequency   = 10000 //1000000 * 10 // 1 million * N
	RingSize             = 2
	RingMask             = RingSize - 1
)

var ringBuffer [RingSize]int64

func main() {
	runtime.GOMAXPROCS(MaxConsumerGroups*MaxConsumersPerGroup + MaxProducers)

	writerCursor := disruptor.NewCursor()
	writerBarrier := disruptor.NewSharedWriterBarrier(writerCursor, RingSize)
	readerBarrier := startConsumerGroups(writerBarrier, writerCursor)
	writer := disruptor.NewSharedWriter(writerBarrier, readerBarrier)

	startProducers(writer)
}
func startProducers(writer *disruptor.SharedWriter) {
	for i := 0; i < MaxProducers-1; i++ {
		go publish(i, writer)
	}

	publish(MaxProducers-1, writer)
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
		//go consume0(disruptor.NewSimpleReader(reader, NewExampleConsumerHandler()))

		// wildly sporadic latency for single-item publish, e.g. 2 seconds, 65 ms, etc.
		// faster for 2-3+ items per publish
		// go consume1(reader)

		if group == 0 {
			go consume1(reader)
		} else if group == 1 {
			go consume2(reader)
		} else {
			panic("only two consumer groups currently supported.")
		}
	}

	return disruptor.NewBarrier(readerCursors...)
}
