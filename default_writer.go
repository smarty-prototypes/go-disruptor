package disruptor

import (
	"runtime"
	"sync/atomic"
)

type defaultWriter struct {
	written  *atomic.Int64   // ring has been written up to this sequence
	upstream sequenceBarrier // all readers have advanced up to this sequence
	capacity int64
	previous int64
}

func newWriter(written *atomic.Int64, upstream sequenceBarrier, capacity int64) Writer {
	return &defaultWriter{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: defaultCursorValue,
	}
}

func (this *defaultWriter) Reserve(count int64) int64 {
	if count <= 0 || count > this.capacity {
		return ErrReservationSize
	}

	var gate int64 = defaultCursorValue // TODO: this field may need to be stateful

	this.previous += count
	for spin := int64(0); this.previous-this.capacity > gate; spin++ {
		if spin&spinMask == 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}

		gate = this.upstream.Load()
	}
	return this.previous
}

func (this *defaultWriter) Commit(_, upper int64) { this.written.Store(upper) }

const spinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
