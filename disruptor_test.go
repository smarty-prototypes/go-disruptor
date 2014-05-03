package disruptor

import (
	"testing"
	"time"
)

func BenchmarkDisruptor(b *testing.B) {
	iterations := uint64(b.N)

	cursor := NewSequence()
	handler := func(value uint64) {}
	worker := NewWorker(cursor, handler, WaitStrategy)

	go func() {
		for current, max := uint64(0), uint64(0); current < iterations; {
			for current >= max {
				max = worker.sequence.atomicLoad() + BufferSize
				time.Sleep(WaitStrategy)
			}

			current++
			cursor.store(current)
		}

		cursor.close()
	}()

	worker.Process()
}

const BufferSize = 1024 * 128
const BufferMask = BufferSize - 1
const WaitStrategy = time.Nanosecond
