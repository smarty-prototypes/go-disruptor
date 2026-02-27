package disruptor

import (
	"context"
	"math"
	"sync/atomic"
)

type sharedSequencer struct {
	// cache line 1 — hot path (Reserve/Commit/Load)
	written   *atomicSequence // 8B  — atomic Add every Reserve; read every Load (upper bound of written)
	gate      *atomicSequence // 8B  — read every Reserve (wrap check)
	committed []atomic.Int32  // 24B — Store every Commit; scanned every Load (slice header; len is capacity)
	shift     uint8           // 1B  — read every Commit and Load
	_         [23]byte        // 23B — padding to 64B boundary

	// cache line 2 — slow path only
	upstream sequenceBarrier // 16B — spin loop only
	waiter   WaitStrategy    // 16B — spin loop only
	_        [32]byte        // 32B — tail padding
} // 128B total — fills two 64B cache lines

func newSharedSequencer(capacity uint32, upper *atomicSequence, waiter WaitStrategy) *sharedSequencer {
	committed := make([]atomic.Int32, capacity)
	for i := range committed {
		committed[i].Store(defaultSequenceValue)
	}

	return &sharedSequencer{
		written:   upper,
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

	var (
		upper      = this.written.Add(count) // claims the slot for the caller (not using CAS operation)
		lower      = upper - count
		wrap       = upper - capacity
		cachedGate = this.gate.Load()
	)

	// fast path
	if wrap <= cachedGate && cachedGate <= lower {
		return upper
	}

	// slow path
	for cachedGate = this.upstream.Load(0); wrap > cachedGate; cachedGate = this.upstream.Load(0) {
		this.waiter.Reserve()
	}

	this.gate.Store(cachedGate)
	return upper
}

func (this *sharedSequencer) TryReserve(ctx context.Context, count int64) int64 {
	if count <= 0 || count > int64(len(this.committed)) {
		return ErrReservationSize
	}

	for {
		lower := this.written.Load()
		upper := lower + count

		if this.hasAvailableCapacity(lower, count) && this.written.CompareAndSwap(lower, upper) {
			return upper // successfully claimed slot
		} else if this.waiter.TryReserve(ctx) != nil {
			return ErrContextCanceled
		}
	}
}
func (this *sharedSequencer) hasAvailableCapacity(lower, count int64) bool {
	var (
		upper      = lower + count
		wrap       = upper - int64(len(this.committed))
		cachedGate = this.gate.Load()
	)

	// fast path
	if wrap <= cachedGate && cachedGate <= lower {
		return true
	}

	// slow path
	gate := this.upstream.Load(0)
	this.gate.Store(gate)
	return wrap <= gate
}

func (this *sharedSequencer) Commit(lower, upper int64) {
	for mask := int64(len(this.committed)) - 1; lower <= upper; lower++ {
		this.committed[lower&mask].Store(int32(lower >> this.shift))
	}
}

func (this *sharedSequencer) Load(lower int64) int64 {
	upper := this.written.Load()

	for mask := int64(len(this.committed)) - 1; lower <= upper; lower++ {
		if this.committed[lower&mask].Load() != int32(lower>>this.shift) {
			return lower - 1
		}
	}

	return upper
}
