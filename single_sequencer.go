package main

func (this SingleProducerSequencer) Next(slotCount int64) int64 {
	nextValue := this.pad.Load()
	nextSequence := nextValue + slotCount
	wrapPoint := nextSequence - this.ringSize
	cachedGatingSequence := this.pad[cachedGatingSequencePadIndex]

	if wrapPoint > cachedGatingSequence || cachedGatingSequence > nextValue {
		minSequence := int64(0)
		for wrapPoint > minSequence {
			minSequence = this.last.Load()
		}

		this.pad[cachedGatingSequencePadIndex] = minSequence
	}

	this.pad.Store(nextSequence)
	return nextSequence
}

func (this SingleProducerSequencer) Publish(sequence int64) {
	this.cursor.Store(sequence)
}

func NewSingleProducerSequencer(cursor *Sequence, ringSize int32, last Barrier) SingleProducerSequencer {
	pad := NewSequence()
	pad[cachedGatingSequencePadIndex] = InitialSequenceValue

	return SingleProducerSequencer{
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

const cachedGatingSequencePadIndex = 1
