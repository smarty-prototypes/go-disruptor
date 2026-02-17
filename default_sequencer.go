package disruptor

import "context"

type defaultSequencer struct {
	capacity uint32              // 4B  — read every Reserve
	spinMask uint32              // 4B  — spin loop only
	upper    int64               // 8B  — read+write every Reserve
	gate     int64               // 8B  — read every Reserve (wrap check)
	written  atomicSequence      // 8B  — ring has been written up to this sequence
	upstream sequenceBarrier     // 16B — all readers have advanced up to this sequence
	waiter   ReserveWaitStrategy // 16B — spin loop only
}

func newSequencer(capacity int64, written atomicSequence, upstream sequenceBarrier, waiter ReserveWaitStrategy) Sequencer {
	return &defaultSequencer{
		capacity: uint32(capacity),
		spinMask: uint32(waiter.SpinMask()),
		upper:    defaultSequenceValue,
		gate:     defaultSequenceValue,
		written:  written,
		upstream: upstream,
		waiter:   waiter,
	}
}

func (this *defaultSequencer) Reserve(ctx context.Context, count int64) int64 {
	capacity := int64(this.capacity)
	if count <= 0 || count > capacity {
		return ErrReservationSize
	}

	// fast path
	this.upper += count
	wrap := this.upper - capacity
	if wrap <= this.gate && this.gate <= (this.upper-count) {
		return this.upper
	}

	// slow path
	spinMask := int64(this.spinMask)
	for spin := int64(0); wrap > this.gate; spin++ {
		if spin&spinMask == 0 && this.waiter.Wait(ctx) != nil {
			this.upper = this.upper - count // undo reservation
			return ErrContextCanceled
		}

		this.gate = this.upstream.Load(0)
	}

	return this.upper
}
func (this *defaultSequencer) Commit(_, upper int64) { this.written.Store(upper) }
