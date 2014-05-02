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
