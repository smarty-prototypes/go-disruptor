package main

func (this *SingleProducerSequencer) Next(slots int64) int64 {
	current, gate := this.current, this.gate
	next := current + slots
	wrap := next - this.ringSize

	if wrap > gate || gate > current {
		min, last := int64(0), this.last
		for wrap > min {
			min = last.Load()
		}
		this.gate = min
	}

	this.current = next
	return next
}

func (this *SingleProducerSequencer) Publish(sequence int64) {
	this.cursor[0] = sequence
}

func NewSingleProducerSequencer(cursor *Sequence, ringSize int32, last Barrier) *SingleProducerSequencer {
	return &SingleProducerSequencer{
		current:  InitialSequenceValue,
		gate:     InitialSequenceValue,
		cursor:   cursor,
		ringSize: int64(ringSize),
		last:     last,
	}
}

type SingleProducerSequencer struct {
	current  int64
	gate     int64
	cursor   *Sequence
	ringSize int64
	last     Barrier
}
