package disruptor

import (
	"testing"
	"time"
)

func BenchmarkDisruptor(b *testing.B) {
	consumer, producer := NewSequence(), NewSequence()
	ring := make([]uint64, BufferSize)
	iterations := uint64(b.N)

	go func() {
		for current, maxAvailable := uint64(0), uint64(0); current < iterations; {
			for current >= maxAvailable {
				maxAvailable = consumer.AtomicLoad() + BufferSize
				time.Sleep(WaitStrategy)
			}

			ring[current&BufferMask] = current
			current++
			producer.Store(current)
		}
	}()

	for current, maxPublished := uint64(0), uint64(0); current < iterations; current++ {
		for current >= maxPublished {
			maxPublished = producer.AtomicLoad()
			time.Sleep(WaitStrategy)
		}

		message := ring[current&BufferMask]
		if message != current {
			panic("Out of sequence")
		}
		consumer.Store(current)
	}
}

const BufferSize = 1024 * 128
const BufferMask = BufferSize - 1
const WaitStrategy = time.Nanosecond
