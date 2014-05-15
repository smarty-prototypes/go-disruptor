package disruptor

func (this *SingleProducerSequencer) Next(items int64) int64 {
	previous, gate := this.previous, this.gate
	next := previous + items
	wrap := next - this.ringSize

	if wrap > gate || gate > previous {
		barrier := this.barrier
		min := barrier.Load()
		for wrap > min || min < 0 {
			min = barrier.Load()
		}

		this.gate = min
	}

	this.previous = next
	return next
}

func NewSingleProducerSequencer(cursor *Sequence, ringSize int32, barrier Barrier) *SingleProducerSequencer {
	return &SingleProducerSequencer{
		previous: InitialSequenceValue,
		gate:     InitialSequenceValue,
		cursor:   cursor,
		ringSize: int64(ringSize),
		barrier:  barrier,
	}
}

type SingleProducerSequencer struct {
	previous int64
	gate     int64
	cursor   *Sequence
	ringSize int64
	barrier  Barrier
}
