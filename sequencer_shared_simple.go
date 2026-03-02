package disruptor

// simpleSharedSequencer is an alternative multi-writer Sequencer that uses ordered (serialized) commits instead of
// per-slot commit tracking. Commit busy-spins until all preceding sequences have been committed, then stores the upper
// sequence in a single atomic write. This makes Load O(1) (single atomic read) instead of O(N) (slot scan), but
// serializes commits which may cause contention with many more than 3-4 producers.
//
// Trade-offs vs sharedSequencer (per-slot commit):
//   - Commit: serialized waits (ordered) vs independent per-slot stores (per-slot)
//   - Load: single atomic read (ordered) vs scan from lower to upper (per-slot)
//   - Struct: 64B / 1 cache line (ordered) vs 128B / 2 cache lines
type simpleSharedSequencer struct {
	// cache line 1 — hot path (Reserve/Commit/Load)
	reservedSequence       *atomicSequence // 8B  — atomic Add every Reserve
	cachedConsumerSequence *atomicSequence // 8B  — read every Reserve (wrap check)
	committedSequence      *atomicSequence // 8B  — spin-read + Store every Commit; read every Load
	capacity               uint32          // 4B  — buffer capacity (power of 2)
	_                      [4]byte         // 4B  — tail padding
	consumerBarrier        sequenceBarrier // 16B — slow path
	waiter                 WaitStrategy    // 16B — slow path
} // 64B total — fills a single 64B cache line

func newSimpleSharedSequencer(capacity uint32, reservedSequence *atomicSequence, waiter WaitStrategy) *simpleSharedSequencer {
	return &simpleSharedSequencer{
		reservedSequence:       reservedSequence,
		cachedConsumerSequence: newSequence(),
		committedSequence:      newSequence(),
		capacity:               capacity,
		waiter:                 waiter,
	}
}

func (this *simpleSharedSequencer) Reserve(count uint32) int64 {
	if count == 0 || count > this.capacity {
		return ErrReservationSize
	}

	var (
		slots                    = int64(count)
		reservedSequence         = this.reservedSequence.Add(slots)
		previousReservedSequence = reservedSequence - slots
		minimumSequence          = reservedSequence - int64(this.capacity)
		consumerSequence         = this.cachedConsumerSequence.Load()
	)

	// fast path
	if minimumSequence <= consumerSequence && consumerSequence <= previousReservedSequence {
		return reservedSequence
	}

	// slow path
	for spin := int64(0); ; spin++ {
		consumerSequence = this.consumerBarrier.Load(0)
		if minimumSequence <= consumerSequence {
			break
		}
		this.waiter.Reserve(spin)
	}

	this.cachedConsumerSequence.Store(consumerSequence)
	return reservedSequence
}

func (this *simpleSharedSequencer) TryReserve(count uint32) int64 {
	if count == 0 || count > this.capacity {
		return ErrReservationSize
	}

	// fast path
	slots := int64(count)
	previousReservedSequence := this.reservedSequence.Load()
	if !this.hasAvailableCapacity(previousReservedSequence, slots) {
		return ErrCapacityUnavailable
	}

	// slow path
	if !this.reservedSequence.CompareAndSwap(previousReservedSequence, previousReservedSequence+slots) {
		return ErrCapacityUnavailable
	}

	return previousReservedSequence + slots
}
func (this *simpleSharedSequencer) hasAvailableCapacity(previousReservedSequence, count int64) bool {
	var (
		reservedSequence = previousReservedSequence + count
		minimumSequence  = reservedSequence - int64(this.capacity)
		consumerSequence = this.cachedConsumerSequence.Load()
	)

	if minimumSequence <= consumerSequence && consumerSequence <= previousReservedSequence {
		return true
	}

	consumerSequence = this.consumerBarrier.Load(0)
	this.cachedConsumerSequence.Store(consumerSequence)
	return minimumSequence <= consumerSequence
}

func (this *simpleSharedSequencer) Commit(lower, upper int64) {
	for this.committedSequence.Load() != lower-1 {
		// spin loop
	}
	this.committedSequence.Store(upper)
}

func (this *simpleSharedSequencer) Load(_ int64) int64 { return this.committedSequence.Load() }
