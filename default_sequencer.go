package disruptor

import (
	"context"
	"runtime"
)

type defaultSequencer struct {
	written  atomicSequence  // ring has been written up to this sequence
	upstream sequenceBarrier // all readers have advanced up to this sequence
	capacity int64
	current  int64
	gate     int64
}

func newSequencer(written atomicSequence, upstream sequenceBarrier, capacity int64) Sequencer {
	return &defaultSequencer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		current:  defaultSequenceValue,
		gate:     defaultSequenceValue,
	}
}

func (this *defaultSequencer) Reserve(ctx context.Context, count int64) int64 {
	if count <= 0 || count > this.capacity {
		return ErrReservationSize
	}

	this.current += count

	for spin := int64(0); this.current-this.capacity > this.gate; spin++ {
		if spin&spinMask == 0 {
			if ctx.Err() != nil {
				this.current -= count
				return ErrContextCanceled
			}

			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}

		this.gate = this.upstream.Load()
	}

	return this.current
}
func (this *defaultSequencer) Commit(_, upper int64) { this.written.Store(upper) }

const spinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
