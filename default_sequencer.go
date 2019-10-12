package disruptor

import "runtime"

type DefaultSequencer struct {
	written  *Sequence // the ring buffer has been written up to this sequence
	upstream Barrier   // all of the readers have advanced up to this sequence
	capacity int64
	previous int64
	gate     int64
}

func NewSequencer(written *Sequence, upstream Barrier, capacity int64) *DefaultSequencer {
	return &DefaultSequencer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: defaultSequenceValue,
		gate:     defaultSequenceValue,
	}
}

func (this *DefaultSequencer) Reserve(count int64) int64 {
	if count <= 0 {
		panic(ErrMinimumReservationSize)
	}

	this.previous += count
	for spin := int64(0); this.previous-this.capacity > this.gate; spin++ {
		if spin&SpinMask == 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}

		this.gate = this.upstream.Load()
	}
	return this.previous
}

func (this *DefaultSequencer) Commit(_, upper int64) { this.written.Store(upper) }

const SpinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
