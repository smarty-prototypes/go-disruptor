package disruptor

import "runtime"

type SingleSequencer struct {
	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // all of the readers have advanced up to this sequence
	capacity int64
	previous int64
}

func NewSingleSequencer(written *Cursor, upstream Barrier, capacity int64) *SingleSequencer {
	return &SingleSequencer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: InitialCursorSequenceValue,
	}
}

func (this *SingleSequencer) Reserve(count int64) int64 {
	this.previous += count

	for spin := int64(0); this.previous-this.capacity > this.upstream.Load(); spin++ {
		if spin&spinMask == 0 {
			runtime.Gosched() // http://bit.ly/1xiDINZ
		}
	}

	return this.previous
}

func (this *SingleSequencer) Commit(_, upper int64) { this.written.Store(upper) }

const spinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
