package disruptor

type DefaultSequencer struct {
	written  *Sequence // the ring buffer has been written up to this sequence
	upstream Barrier   // all of the readers have advanced up to this sequence
	capacity int64
	previous int64
}

func NewSequencer(written *Sequence, upstream Barrier, capacity int64) *DefaultSequencer {
	return &DefaultSequencer{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: defaultSequenceValue,
	}
}

func (this *DefaultSequencer) Reserve(count int64) int64 {
	if this.previous+count-this.capacity > this.upstream.Load() {
		return defaultSequenceValue // no room for the reservation
	}

	this.previous += count
	return this.previous
}

func (this *DefaultSequencer) Commit(_, upper int64) { this.written.Store(upper) }
