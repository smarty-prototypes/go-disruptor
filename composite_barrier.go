package disruptor

type CompositeBarrier []*Sequence

func NewCompositeBarrier(sequences []*Sequence) CompositeBarrier { return sequences }

func (this CompositeBarrier) Load() int64 {
	var minimum = maxCursorSequenceValue

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}

const maxCursorSequenceValue int64 = (1 << 63) - 1
