package main

func (this *SingleProducerSequencer) Next(slotCount int64) int64 {
	nextValue := this.pad.Load()
	nextSequence := nextValue + slotCount
	wrap := nextSequence - this.ringSize
	cachedGate := this.pad[cachedGatePadIndex]

	if wrap > cachedGate || cachedGate > nextValue {
		minSequence := int64(0)
		for wrap > minSequence {
			minSequence = this.last.Load()
		}

		this.pad[cachedGatePadIndex] = minSequence
	}

	this.pad.Store(nextSequence)
	return nextSequence
}

func (this *SingleProducerSequencer) Publish(sequence int64) {
	// this.cursor.Store(sequence)
	this.cursor[0] = sequence
}

func NewSingleProducerSequencer(cursor *Sequence, ringSize int32, last Barrier) *SingleProducerSequencer {
	pad := NewSequence()
	pad[cachedGatePadIndex] = InitialSequenceValue

	return &SingleProducerSequencer{
		pad:      pad,
		cursor:   cursor,
		ringSize: int64(ringSize),
		last:     last,
	}
}

type SingleProducerSequencer struct {
	pad      *Sequence
	cursor   *Sequence
	ringSize int64
	last     Barrier
}

const cachedGatePadIndex = 1
