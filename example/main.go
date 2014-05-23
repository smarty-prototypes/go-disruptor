package main

import (
	"runtime"

	"github.com/smartystreets/go-disruptor"
)

const (
	MaxConsumersPerGroup = 1
	MaxConsumerGroups    = 1
	MaxProducers         = 2
	ItemsToPublish       = 4
	ReportingFrequency   = 1000000 * 10 // 1 million * N
	RingSize             = 1024 * 16
	RingMask             = RingSize - 1
)

var ringBuffer [RingSize]int64

func main() {
	runtime.GOMAXPROCS(MaxConsumerGroups*MaxConsumersPerGroup + MaxProducers)

	written := disruptor.NewCursor()
	shared := disruptor.NewSharedWriterBarrier(written, RingSize)
	upstream := startConsumerGroups(shared, written)
	writer := disruptor.NewSharedWriter(shared, upstream)
	// writer := disruptor.NewWriter(written, upstream, RingSize)
	// startExclusiveProducer(writer)
	startSharedProducers(writer)
}

// func startExclusiveProducer(writer *disruptor.Writer) {
//publish(writer)
// }

func startSharedProducers(writer *disruptor.SharedWriter) {
	for i := 0; i < MaxProducers-1; i++ {
		go publish(writer)
	}

	publish(writer)
}

func startConsumerGroups(upstream disruptor.Barrier, written *disruptor.Cursor) disruptor.Barrier {
	for i := 0; i < MaxConsumerGroups; i++ {
		upstream = startConsumerGroup(i, upstream, written)
	}

	return upstream
}
func startConsumerGroup(group int, upstream disruptor.Barrier, written *disruptor.Cursor) disruptor.Barrier {
	cursors := []*disruptor.Cursor{}

	for i := 0; i < MaxConsumersPerGroup; i++ {
		read := disruptor.NewCursor()
		cursors = append(cursors, read)
		reader := disruptor.NewReader(read, written, upstream)

		// constant time regardless of the number of items
		// go consume0(disruptor.NewSimpleReader(reader, NewExampleConsumerHandler()))

		// TODO: wildly sporadic latency for single-item publish, e.g. 2 seconds, 65 ms, etc.
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

	return disruptor.NewCompositeBarrier(cursors...)
}
