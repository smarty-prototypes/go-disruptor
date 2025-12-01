package disruptor

import "runtime"

type defaultWriter struct {
	written  atomicSequence  // ring has been written up to this sequence
	upstream sequenceBarrier // all readers have advanced up to this sequence
	capacity int64
	previous int64
	gate     int64
}

func newWriter(written atomicSequence, upstream sequenceBarrier, capacity int64) Writer {
	return &defaultWriter{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: defaultSequenceValue,
		gate:     defaultSequenceValue,
	}
}

func (this *defaultWriter) Reserve(count int64) int64 {
	if count <= 0 || count > this.capacity {
		return ErrReservationSize
	}

	this.previous += count
	for spin := int64(0); this.previous-this.capacity > this.gate; spin++ {
		if spin&spinMask == 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}

		this.gate = this.upstream.Load()
	}
	return this.previous
}
func (this *defaultWriter) Commit(_, upper int64) { this.written.Store(upper) }

const spinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
