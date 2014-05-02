package disruptor

type MultiSequencer struct {
	cursor Sequence
}

func (this MultiSequencer) Next() uint64 {
	return 0
}
func (this MultiSequencer) Publish(sequence uint64) {
}

func NewMultiSequencer() MultiSequencer {
	return MultiSequencer{}
}
