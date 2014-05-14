package main

func (this *SingleProducerSequencer) Next(items int64) int64 {
	claimed, gate := this.claimed, this.gate
	next := claimed + items
	wrap := next - this.ringSize

	if wrap > gate {
		last := this.last
		min := last.Load()

		for wrap > min || min < 0 {
			min = last.Load()
		}

		this.gate = min
	}

	this.claimed = next
	return next
}

func (this *SingleProducerSequencer) Publish(sequence int64) {
	this.cursor[SequencePayloadIndex] = sequence
}

func NewSingleProducerSequencer(cursor *Sequence, ringSize int32, last Barrier) *SingleProducerSequencer {
	return &SingleProducerSequencer{
		claimed:  InitialSequenceValue,
		gate:     InitialSequenceValue,
		cursor:   cursor,
		ringSize: int64(ringSize),
		last:     last,
	}
}

type SingleProducerSequencer struct {
	claimed  int64
	gate     int64
	cursor   *Sequence
	ringSize int64
	last     Barrier
}
