package disruptor

import (
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkDisruptor(b *testing.B) {
	consumer, producer := Sequence{}, Sequence{}
	ring := make([]uint64, BufferSize)
	iterations := uint64(b.N)

	go func() {
		for current, max := uint64(0), uint64(0); current < iterations; current++ {
			for current >= max {
				max = atomic.LoadUint64(&consumer[0]) + BufferSize
				time.Sleep(WaitStrategy)
			}

			ring[current&BufferMask] = current
			producer[0] = current + 1
		}
	}()

	for current, max := uint64(0), uint64(0); current < iterations; current++ {
		for current >= max {
			max = atomic.LoadUint64(&producer[0])
			time.Sleep(WaitStrategy)
		}

		message := ring[current&BufferMask]
		if message != current {
			panic("Out of sequence")
		}
		consumer[0] = current
	}
}

const BufferSize = 1024 * 128
const BufferMask = BufferSize - 1
const FillCPUCacheLine = 8
const WaitStrategy = time.Nanosecond

type Sequence [FillCPUCacheLine]uint64
