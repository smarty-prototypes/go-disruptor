package disruptor

import (
	"math"
	"sync/atomic"
)

type sharedSequencer struct {
	// cache line 1 — hot path (Reserve/Commit/Load)
	upper     *atomicSequence // 8B  — atomic Add every Reserve; read every Load (upper bound)
	gate      *atomicSequence // 8B  — read every Reserve (wrap check)
	committed []atomic.Int32  // 24B — Store every Commit; scanned every Load (slice header; len is capacity)
	shift     uint8           // 1B  — read every Commit and Load
	_         [23]byte        // 23B — padding to 64B boundary

	// cache line 2 — slow path only
	upstream sequenceBarrier // 16B — spin loop only
	waiter   WaitStrategy    // 16B — spin loop only
	_        [32]byte        // 32B — tail padding
}                            // 128B total — fills two 64B cache lines

func newSharedSequencer(capacity uint32, upper *atomicSequence, waiter WaitStrategy) *sharedSequencer {
	committed := make([]atomic.Int32, capacity)
	for i := range committed {
		committed[i].Store(int32(defaultSequenceValue))
	}
	return &sharedSequencer{
		upper:     upper,
		gate:      newSequence(),
		shift:     uint8(math.Log2(float64(capacity))),
		committed: committed,
		waiter:    waiter,
	}
}

func (this *sharedSequencer) Reserve(count int64) int64 {
	capacity := int64(len(this.committed))
	if count <= 0 || count > capacity {
		return ErrReservationSize
	}

	// using atomic Add because it scales even with contention compared to CAS
	// this was at the cost of NOT allowing Reserve to be canceled.
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

func (this *sharedSequencer) Commit(lower, upper int64) {
	for mask := int64(len(this.committed)) - 1; lower <= upper; lower++ {
		this.committed[lower&mask].Store(int32(lower >> this.shift))
	}
}

func (this *sharedSequencer) Load(lower int64) int64 {
	upper := this.upper.Load()

	for mask := int64(len(this.committed)) - 1; lower <= upper; lower++ {
		if this.committed[lower&mask].Load() != int32(lower>>this.shift) {
			return lower - 1
		}
	}

	return upper
}
