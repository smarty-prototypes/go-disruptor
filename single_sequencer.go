package main

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

func (this *SingleProducerSequencer) Publish(sequence int64) {
	// TODO: put this one in a build tag as well
	this.cursor[SequencePayloadIndex] = sequence
	// this.cursor.Store(sequence)
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
