package disruptor

import (
	"math"
	"sync/atomic"
)

type multiSequencer struct {
	capacity  uint32          // 4B  — read every Reserve + Commit
	shift     uint8           // 1B  — read every Commit
	upper     atomicSequence  // 8B  — Load+CAS every Reserve
	gate      atomicSequence  // 8B  — read every Reserve (wrap check)
	committed []atomic.Int32  // 24B — read every Commit (slice header)
	upstream  sequenceBarrier // 16B — spin loop only
	waiter    WaitStrategy    // 16B — spin loop only
}

func (this *multiSequencer) Reserve(count int64) int64 {
	capacity := int64(this.capacity)
	if count <= 0 || count > capacity {
		return ErrReservationSize
	}

	// using atomic Add because it scales even with contention compared to CAS
	// this was at the cost of allowing Reserve to be canceled.
	var (
		upper = this.upper.Add(count) // claims the slot for the caller
		lower = upper - count
		wrap  = upper - capacity
		gate  = this.gate.Load()
	)

	// fast path
	if wrap <= gate && gate <= lower {
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

func (this *multiSequencerConfiguration) NewSequencer(upstream sequenceBarrier, waiting WaitStrategy) Sequencer {
	return &multiSequencer{
		upper:     this.written,
		gate:      newSequence(),
		capacity:  uint32(this.capacity),
		shift:     this.shift,
		committed: this.committed,
		upstream:  upstream,
		waiter:    waiting,
	}
}
