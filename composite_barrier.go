package disruptor

type CompositeBarrier []*Sequence

func NewCompositeBarrier(sequences []*Sequence) CompositeBarrier { return sequences }

func (this CompositeBarrier) Load() int64 {
	var minimum = MaxCursorSequenceValue

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
