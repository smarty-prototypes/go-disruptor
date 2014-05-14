package main

func (this *SingleProducerSequencer) Next(items int64) int64 {
	claimed, gate := this.claimed, this.gate
	next := claimed + items
	wrap := next - this.ringSize
	// fmt.Printf("Producer:: Last claim: %d, Next: %d, Wrap: %d, Gate:%d\n", claimed, next, wrap, gate)

	if wrap > gate {
		last := this.last
		min := last.Load(1)
		// fmt.Printf("Producer:: (a) Wrap: %d, Current Gate, %d, Proposed Gate:%d\n", wrap, gate, min)

		for wrap > min || min < 0 {
			min = last.Load(1)
			// fmt.Printf("Producer:: (b) Wrap: %d, Current Gate, %d, Proposed Gate:%d\n", wrap, gate, min)

			// if wrap <= min {
			// 	fmt.Printf("Producer:: Consumers have caught up to producer.\n")
			// }

		}

		// fmt.Printf("Producer:: (c) Wrap: %d, Current Gate, %d, Proposed Gate:%d\n", wrap, gate, min)
		this.gate = min
	}

	this.claimed = next
	return next
}

func (this *SingleProducerSequencer) Publish(sequence int64) {
	this.cursor[0] = sequence
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
