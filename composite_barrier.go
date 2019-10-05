package disruptor

type CompositeBarrier []*Cursor

func NewCompositeBarrier(cursors []*Cursor) CompositeBarrier {
	return CompositeBarrier(cursors)
}

func (this CompositeBarrier) Load() int64 {
	var minimum = MaxSequenceValue

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
