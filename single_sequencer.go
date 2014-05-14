package main

func (this *SingleProducerSequencer) Next(slots int64) int64 {
	current, cachedGate := this.current, this.cachedGate
	next := current + slots
	wrap := next - this.ringSize

	if wrap > cachedGate /*|| cachedGate > current*/ {
		min, last := int64(0), this.last
		for wrap > min {
			min = last.Load()
		}
		this.cachedGate = min
	}

	this.current = next
	return next
}

func (this *SingleProducerSequencer) Publish(sequence int64) {
	this.cursor[0] = sequence
}

func NewSingleProducerSequencer(cursor *Sequence, ringSize int32, last Barrier) *SingleProducerSequencer {
	return &SingleProducerSequencer{
		current:    InitialSequenceValue,
		cachedGate: InitialSequenceValue,
		cursor:     cursor,
		ringSize:   int64(ringSize),
		last:       last,
	}
}

type SingleProducerSequencer struct {
	current    int64
	cachedGate int64
	cursor     *Sequence
	ringSize   int64
	last       Barrier
}
