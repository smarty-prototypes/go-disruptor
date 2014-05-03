package disruptor

import (
	"testing"
	"time"
)

func BenchmarkDisruptor(b *testing.B) {
	iterations := uint64(b.N)

	cursor := NewSequence()
	handler := func(value uint64) {}
	worker1 := NewWorker(cursor, handler, WaitStrategy)
	worker2 := NewWorker(cursor, handler, WaitStrategy)
	workers := []Worker{worker1, worker2}

	go func() {
		for current, max := uint64(0), uint64(0); current < iterations; {
			for current >= max {
				max = minSequence(workers) + BufferSize
				//max = worker.sequence.atomicLoad() + BufferSize
				time.Sleep(WaitStrategy)
			}

			current++
			cursor.store(current)
		}

		cursor.close()
	}()

	go worker1.Process()
	worker2.Process()
}
func minSequence(workers []Worker) uint64 {
	min := Uint64MaxValue

	for i := 0; i < len(workers); i++ {
		seq := workers[i].sequence.atomicLoad()
		if seq < min {
			min = seq
		}
	}

	return min
}

const BufferSize = 1024 * 128
const BufferMask = BufferSize - 1
const WaitStrategy = time.Nanosecond
