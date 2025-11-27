package disruptor

import (
	"runtime"
	"sync/atomic"
)

type defaultWriter struct {
	written  *atomic.Int64 // the ring buffer has been written up to this sequence
	upstream Barrier       // all of the readers have advanced up to this sequence
	capacity int64
	previous int64
}

func NewWriter(written *atomic.Int64, upstream Barrier, capacity int64) Writer {
	return &defaultWriter{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: defaultCursorValue,
	}
}

func (this *defaultWriter) Reserve(count int64) int64 {
	if count <= 0 {
		panic(ErrMinimumReservationSize)
	}

	this.previous += count
	for spin, gate := int64(0), int64(defaultCursorValue); this.previous-this.capacity > gate; spin++ {
		if spin&SpinMask == 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}

		gate = this.upstream.Load()
	}
	return this.previous
}

func (this *defaultWriter) Commit(_, upper int64) { this.written.Store(upper) }

const SpinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
