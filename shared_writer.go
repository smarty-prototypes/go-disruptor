package disruptor

import "sync/atomic"

type SharedWriter struct {
	capacity    int64
	mask        int64
	gate        int64 // TODO: determine if this should be a *Cursor
	shift       uint8
	committed   []int32
	upstream    Barrier
	reservation *Cursor
}

func NewSharedWriter(shared *SharedWriterBarrier, upstream Barrier) *SharedWriter {
	return &SharedWriter{
		capacity:    shared.capacity,
		mask:        shared.mask,
		gate:        InitialSequenceValue,
		shift:       shared.shift,
		committed:   shared.committed,
		upstream:    upstream,
		reservation: shared.reservation,
	}
}

func (this *SharedWriter) Reserve(count int64) (int64, int64) {
	for {
		previous := this.reservation.Load()
		upper := previous + count
		wrap := upper - this.capacity

		if wrap > this.gate {
			min := this.upstream.LoadBarrier(0)
			if wrap > min {
				return InitialSequenceValue, Gating
			}

			this.gate = min // doesn't matter which write wins, BUT will most likely need to be a Cursor
		}

		if atomic.CompareAndSwapInt64(&this.reservation.value, previous, upper) {
			return previous + 1, upper
		}
	}
}

func (this *SharedWriter) Commit(lower, upper int64) {
	for shift, mask := this.shift, this.mask; lower <= upper; lower++ {
		this.committed[lower&mask] = int32(lower >> shift)
	}
}
