package disruptor

import "runtime"

type SingleWriter struct {
	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // all of the readers have advanced up to this sequence
	capacity int64
	previous int64
}

func NewSingleWriter(written *Cursor, upstream Barrier, capacity int64) *SingleWriter {
	assertPowerOfTwo(capacity)

	return &SingleWriter{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: InitialCursorSequenceValue,
	}
}

func assertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		// Wikipedia entry: http://bit.ly/1krhaSB
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}
}

func (this *SingleWriter) Reserve(count int64) int64 {
	this.previous += count

	for spin := int64(0); this.previous-this.capacity > this.upstream.Load(); spin++ {
		if spin&SpinMask == 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}
	}

	return this.previous
}

func (this *SingleWriter) Commit(_, upper int64) {
	this.written.Store(upper)
}

const SpinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
