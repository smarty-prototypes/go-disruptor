package main

func (this *SingleProducerSequencer) Next(items int64) int64 {
	previous, gate := this.previous, this.gate
	next := previous + items
	wrap := next - this.ringSize
	// fmt.Printf("Producer:: Last claim: %d, Next: %d, Wrap: %d, Gate:%d\n", previous, next, wrap, gate)

	if wrap > gate || gate > previous {
		barrier := this.barrier
		min := barrier.Load(1)
		// fmt.Printf("Producer:: (a) Wrap: %d, Current Gate, %d, Proposed Gate:%d\n", wrap, gate, min)

		for wrap > min || min < 0 {
			// fmt.Printf("Producer:: (b) Wrap: %d, Current Gate, %d, Proposed Gate:%d\n", wrap, gate, min)
			min = barrier.Load(1)
			// if wrap <= min {
			// 	fmt.Printf("Producer:: Consumers have caught up to producer.\n")
			// }
		}

		// fmt.Printf("Producer:: (c) Wrap: %d, Current Gate, %d, Proposed Gate:%d\n", wrap, gate, min)
		this.gate = min
	}

	this.previous = next
	return next
}

func (this *SingleProducerSequencer) Publish(sequence int64) {
	this.cursor.Store(sequence)
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
