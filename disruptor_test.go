package disruptor

import (
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkDisruptor(b *testing.B) {
	consumer, producer := NewSequence(), NewSequence()
	ring := make([]uint64, BufferSize)
	iterations := uint64(b.N)

	go func() {
		for current, max := uint64(0), uint64(0); current < iterations; {
			for current >= max {
				max = consumer.AtomicLoad() + BufferSize
				time.Sleep(WaitStrategy)
			}

			ring[current&BufferMask] = current
			current++
			producer.Store(current)
		}
	}()

	for current, max := uint64(0), uint64(0); current < iterations; current++ {
		for current >= max {
			max = producer.AtomicLoad()
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
const FillCPUCacheLine = 8
const WaitStrategy = time.Nanosecond

type Sequence []uint64

func NewSequence() Sequence {
	return Sequence(make([]uint64, FillCPUCacheLine))
}
func (this Sequence) AtomicLoad() uint64 {
	return atomic.LoadUint64(&this[0])
}
func (this Sequence) Store(value uint64) {
	this[0] = value
}
