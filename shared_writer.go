package disruptor

import "sync/atomic"

type SharedWriter struct {
	capacity    int64
	gate        int64 // TODO: determine if this should be a *Cursor
	shift       uint8
	committed   []int32
	upstream    Barrier
	reservation *Cursor
}

func NewSharedWriter(shared *SharedWriterBarrier, upstream Barrier) *SharedWriter {
	return &SharedWriter{
		capacity:    shared.capacity,
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
		next := previous + count
		wrap := next - this.capacity

		if wrap > this.gate {
			min := this.upstream.Load()
			if wrap > min {
				return 0, Gating
			}

			this.gate = min // doesn't matter which write wins, BUT will most likely need to be a Cursor
		}

		if atomic.CompareAndSwapInt64(&this.reservation.value, previous, next) {
			return previous + 1, next
		}
	}
}

func (this *SharedWriter) Commit(lower, upper int64) {
	for mask := this.capacity - 1; lower <= upper; lower++ {
		this.committed[lower&mask] = int32(lower >> this.shift)
	}
}
