package disruptor

import (
	"context"
	"math"
	"sync/atomic"
)

type multiSequencer struct {
	capacity  uint32              // 4B  — read every Reserve + Commit
	spinMask  uint32              // 4B  — spin loop only
	shift     uint8               // 1B  — read every Commit (+7B padding)
	written   atomicSequence      // 8B  — Load+CAS every Reserve
	gate      atomicSequence      // 8B  — read every Reserve (wrap check)
	committed []atomic.Int32      // 24B — read every Commit (slice header)
	upstream  sequenceBarrier     // 16B — spin loop only
	waiter    ReserveWaitStrategy // 16B — spin loop only
}

func (this *multiSequencer) Reserve(ctx context.Context, count int64) int64 {
	capacity := int64(this.capacity)
	if count <= 0 || count > capacity {
		return ErrReservationSize
	}

	spinMask := int64(this.spinMask)
	for spin := int64(0); ; spin++ {
		previous := this.written.Load()
		upper := previous + count
		wrap := upper - capacity
		cachedGate := this.gate.Load()

		if wrap > cachedGate || cachedGate > previous {
			gate := this.upstream.Load(0)
			this.gate.Store(gate)

			for innerSpin := int64(0); wrap > gate; innerSpin++ {
				if innerSpin&spinMask == 0 && this.waiter.Wait(ctx) != nil {
					return ErrContextCanceled
				}

				gate = this.upstream.Load(0)
				this.gate.Store(gate)
			}
		}

		if this.written.CompareAndSwap(previous, upper) {
			return upper
		}

		if spin&spinMask == 0 && this.waiter.Wait(ctx) != nil {
			return ErrContextCanceled
		}
	}
}

func (this *multiSequencer) Commit(lower, upper int64) {
	for mask := int64(this.capacity) - 1; lower <= upper; lower++ {
		this.committed[lower&mask].Store(int32(lower >> this.shift))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type multiSequencerBarrier struct {
	written   atomicSequence
	committed []atomic.Int32 // Go uses address of the slice element as a pointer-receiver methods, so it operates in-place, not on a copy. Otherwise it would be a heap allocation.
	capacity  int64
	shift     uint8
}

func (this *multiSequencerBarrier) Load(lower int64) int64 {
	upper := this.written.Load()

	// walk up the slice when finding the next available slot
	for mask := this.capacity - 1; lower <= upper; lower++ {
		if this.committed[lower&mask].Load() != int32(lower>>this.shift) {
			return lower - 1
		}
	}

	return upper
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type multiSequencerConfiguration struct {
	written   atomicSequence
	committed []atomic.Int32
	capacity  int64
	shift     uint8
}

func newMultiSequencerConfiguration(written atomicSequence, capacity uint32) *multiSequencerConfiguration {
	committed := make([]atomic.Int32, capacity)
	for i := range committed {
		committed[i].Store(int32(defaultSequenceValue))
	}

	return &multiSequencerConfiguration{
		written:   written,
		committed: committed,
		capacity:  int64(capacity),
		shift:     uint8(math.Log2(float64(capacity))),
	}
}

func (this *multiSequencerConfiguration) NewBarrier() *multiSequencerBarrier {
	return &multiSequencerBarrier{
		written:   this.written,
		committed: this.committed,
		capacity:  this.capacity,
		shift:     this.shift,
	}
}

func (this *multiSequencerConfiguration) NewSequencer(upstream sequenceBarrier, waiting ReserveWaitStrategy) Sequencer {
	return &multiSequencer{
		written:   this.written,
		gate:      newSequence(),
		capacity:  uint32(this.capacity),
		spinMask:  uint32(waiting.SpinMask()),
		shift:     this.shift,
		committed: this.committed,
		upstream:  upstream,
		waiter:    waiting,
	}
}
