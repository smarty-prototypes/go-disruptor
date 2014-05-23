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

func (this *SharedWriter) Reserve(id int, count int64) int64 {
	for {
		previous := this.reservation.Load()
		next := previous + count
		wrap := next - this.capacity

		// fmt.Printf("[WRITER %d] Previous: %d, Next: %d, Wrap: %d\n", id, previous, next, wrap)

		if wrap > this.gate {
			// fmt.Printf("[WRITER %d] Previous gate: %d\n", id, this.gate)
			min := this.upstream.Load()
			if wrap > min {
				// fmt.Printf("[WRITER %d] New gate (waiting for consumers): %d\n", id, min)
				return Gating
			}

			// fmt.Printf("[WRITER %d] New gate found: %d\n", id, min)
			this.gate = min // doesn't matter which write wins, BUT will most likely need to be a Cursor
		}

		// fmt.Printf("[WRITER %d] Updating reservation. Previous: %d, Next: %d\n", id, previous, next)
		if atomic.CompareAndSwapInt64(&this.reservation.value, previous, next) {
			// fmt.Printf("[WRITER %d] Reservation updated\n", id)
			return next
			// } else {
			// 	fmt.Printf("[WRITER %d] Reservation rejected, retrying\n", id)
		}
	}
}

func (this *SharedWriter) Commit(sequence int64) {
	this.committed[sequence&this.mask] = int32(sequence >> this.shift)
}
