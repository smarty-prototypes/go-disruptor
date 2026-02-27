package disruptor

type defaultSequencer struct {
	capacity uint32          // 4B  — read every Reserve
	_        [4]byte         // 4B  — explicit alignment padding
	upper    int64           // 8B  — read+write every Reserve
	gate     int64           // 8B  — read every Reserve (wrap check)
	written  *atomicSequence // 8B  — ring has been written up to this sequence
	upstream sequenceBarrier // 16B — all readers have advanced up to this sequence
	waiter   WaitStrategy    // 16B — spin loop only
}                            // 64B total — fills a single 64B cache line

func newSequencer(capacity int64, written *atomicSequence, upstream sequenceBarrier, waiter WaitStrategy) Sequencer {
	return &defaultSequencer{
		capacity: uint32(capacity),
		upper:    defaultSequenceValue,
		gate:     defaultSequenceValue,
		written:  written,
		upstream: upstream,
		waiter:   waiter,
	}
}

func (this *defaultSequencer) Reserve(count int64) int64 {
	capacity := int64(this.capacity)
	if count <= 0 || count > capacity {
		return ErrReservationSize
	}

	// fast path
	lower := this.upper
	this.upper += count
	wrap := this.upper - capacity
	if wrap <= this.gate && this.gate <= lower {
		return this.upper
	}

	// slow path
	for this.gate = this.upstream.Load(0); wrap > this.gate; this.gate = this.upstream.Load(0) {
		this.waiter.Reserve()
	}

	return this.upper
}
func (this *defaultSequencer) Commit(_, upper int64) { this.written.Store(upper) }
