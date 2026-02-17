package disruptor

import (
	"context"
	"math"
	"runtime"
	"sync/atomic"
)

type multiSequencer struct {
	written   atomicSequence
	gate      atomicSequence
	upstream  sequenceBarrier
	committed []atomic.Int32
	capacity  int64
	shift     uint8
}

func (this *multiSequencer) Reserve(ctx context.Context, count int64) int64 {
	if count <= 0 || count > this.capacity {
		return ErrReservationSize
	}

	// block until desired number of slots becomes available
	for spin := uint64(0); ; spin++ {
		previous := this.written.Load()
		upper := previous + count
		wrap := upper - this.capacity
		cachedGate := this.gate.Load()

		if wrap > cachedGate || cachedGate > previous {
			gate := this.upstream.Load(0)
			this.gate.Store(gate)

			for innerSpin := int64(0); wrap > gate; innerSpin++ {
				if innerSpin&spinMask == 0 {
					if ctx.Err() != nil {
						return ErrContextCanceled
					}

					runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
				}

				gate = this.upstream.Load(0)
				this.gate.Store(gate)
			}
		}

		if this.written.CompareAndSwap(previous, upper) {
			return upper
		}

		if spin&spinMask == 0 {
			if ctx.Err() != nil {
				return ErrContextCanceled
			}

			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}
	}
}

func (this *multiSequencer) Commit(lower, upper int64) {
	for mask := this.capacity - 1; lower <= upper; lower++ {
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

func (this *multiSequencerConfiguration) NewSequencer(upstream sequenceBarrier) Sequencer {
	return &multiSequencer{
		written:   this.written,
		gate:      newSequence(),
		upstream:  upstream,
		committed: this.committed,
		capacity:  this.capacity,
		shift:     this.shift,
	}
}
