package disruptor

type CompositeBarrier struct {
	cursors []*Cursor
}

func NewCompositeBarrier(upstream []*Cursor) *CompositeBarrier {
	cursors := make([]*Cursor, len(upstream))
	copy(cursors, upstream)
	return &CompositeBarrier{cursors}
}

func (this *CompositeBarrier) Load() int64 {
	minimum := MaxSequenceValue

	for _, item := range this.cursors {
		cursor := item.Load()
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}
