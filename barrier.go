package disruptor

func (this Barrier) Load() int64 {
	minimum := MaxSequenceValue

	for _, item := range this {
		cursor := item[0]
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}

func NewBarrier(upstream ...*Sequence) Barrier {
	buffer := make([]*Sequence, len(upstream))
	copy(buffer, upstream)
	return Barrier(buffer)
}

type Barrier []*Sequence
