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
	// worker2 := NewWorker(cursor, handler, WaitStrategy)
	// worker3 := NewWorker(cursor, handler, WaitStrategy)
	// worker4 := NewWorker(cursor, handler, WaitStrategy)
	workers := []Worker{worker1 /*, worker2, worker3, worker4*/}

	// TODO: multi-producer sequencer
	// TODO: diamond gating workers, e.g. 1 runs, then 2&3 simultaneously after 1, then 4 after 2&3 are done.
	// TODO: DSL

	go func() {
		for current, max := uint64(0), uint64(0); current < iterations; {
			for current >= max {
				max = minSequence(workers) + BufferSize
				time.Sleep(WaitStrategy)
			}

			current++
			cursor.store(current)
		}

		cursor.close()
	}()

	// go worker4.Process()
	// go worker3.Process()
	// go worker2.Process()
	worker1.Process()
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
