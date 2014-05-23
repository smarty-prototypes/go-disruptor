package disruptor

type CompositeBarrier struct {
	cursors []*Cursor
}

func NewCompositeBarrier(upstream ...*Cursor) *CompositeBarrier {
	cursors := make([]*Cursor, len(upstream))
	copy(cursors, upstream)
	return &CompositeBarrier{cursors}
}

func (this *CompositeBarrier) LoadBarrier(current int64) int64 {
	minimum := MaxSequenceValue
	for _, item := range this.cursors {
		sequence := item.Load()
		if sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
