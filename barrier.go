package disruptor

type Barrier []*Sequence

func NewBarrier(upstream ...*Sequence) Barrier {
	buffer := make([]*Sequence, len(upstream))
	copy(buffer, upstream)
	return Barrier(buffer)
}

func (this Barrier) Load() int64 {
	minimum := MaxSequenceValue

	for _, item := range this {
		cursor := item.Load()
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}
