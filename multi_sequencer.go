package disruptor

import (
	"math"
	"sync/atomic"
)

// TODO: add padding around fields to prevent false sharing with CPU cache lines
type multiSequencer struct {
	upper     atomicSequence  // 8B  — atomic Add every Reserve
	gate      atomicSequence  // 8B  — read every Reserve (wrap check)
	committed []atomic.Int32  // 24B — read every Reserve + Commit (slice header; len is capacity)
	upstream  sequenceBarrier // 16B — spin loop only
	waiter    WaitStrategy    // 16B — spin loop only
	shift     uint8           // 1B  — read every Commit
}                                 // 80B total (73B + 7B tail padding) — spans two 64B cache lines

func newMultiSequencer(upper atomicSequence, committed []atomic.Int32, shift uint8, upstream sequenceBarrier, waiter WaitStrategy) *multiSequencer {
	return &multiSequencer{
		upper:     upper,
		gate:      newSequence(),
		shift:     shift,
		committed: committed,
		upstream:  upstream,
		waiter:    waiter,
	}
}

func (this *multiSequencer) Reserve(count int64) int64 {
	capacity := int64(len(this.committed))
	if count <= 0 || count > capacity {
		return ErrReservationSize
	}

	// using atomic Add because it scales even with contention compared to CAS
	// this was at the cost of allowing Reserve to be canceled.
	var (
		upper = this.upper.Add(count) // claims the slot for the caller
		wrap  = upper - capacity
		gate  = this.gate.Load()
	)

	// fast path
	if wrap <= gate && gate <= upper-count {
		return upper
	}

	// slow path
	for gate = this.upstream.Load(0); wrap > gate; gate = this.upstream.Load(0) {
		this.waiter.Reserve()
	}

	this.gate.Store(gate)
	return upper
}

func (this *multiSequencer) Commit(lower, upper int64) {
	for mask := int64(len(this.committed)) - 1; lower <= upper; lower++ {
		this.committed[lower&mask].Store(int32(lower >> this.shift))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type multiSequencerBarrier struct {
	committed []atomic.Int32  // 24B — walked every Load (loop body; len is capacity)
	written   atomicSequence  // 8B  — read once per Load (upper bound)
	shift     uint8           // 1B  — read every Load (shift computation)
}                                 // 33B total — fits in a single 64B cache line, no padding

func newMultiSequencerBarrier(written atomicSequence, committed []atomic.Int32, shift uint8) *multiSequencerBarrier {
	return &multiSequencerBarrier{written: written, committed: committed, shift: shift}
}

func (this *multiSequencerBarrier) Load(lower int64) int64 {
	upper := this.written.Load()

	for mask := int64(len(this.committed)) - 1; lower <= upper; lower++ {
		if this.committed[lower&mask].Load() != int32(lower>>this.shift) {
			return lower - 1
		}
	}

	return upper
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newCommittedBuffer(capacity uint32) ([]atomic.Int32, uint8) {
	committed := make([]atomic.Int32, capacity)
	for i := range committed {
		committed[i].Store(int32(defaultSequenceValue))
	}
	return committed, uint8(math.Log2(float64(capacity)))
}
