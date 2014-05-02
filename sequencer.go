package disruptor

type Sequencer struct {
	cursor Sequence
}

func (this Sequencer) Next() uint64 {
	return 0
}
func (this Sequencer) Publish(sequence uint64) {
}

func NewSequencer() Sequencer {
	return Sequencer{}
}

// consumers are their own thing, but they share the sequence with the sequencer...
// which facilitates the gating
// to signal a shutdown to a consumer, simply give a uint64.MaxValue to the consumer???

// or perhaps we do create some kind of payload which we reuse, e.g. some kind of "delivery"
// concept handed off to the consumer, e.g.
// type Delivery struct {
//		Sequence uint64
//		Remaining uint32
//		RingIndex uint32
//		Shutdown bool
// }

// or perhaps there's another method on the consumer called Shutdown()?

// Consumer vs Handler
