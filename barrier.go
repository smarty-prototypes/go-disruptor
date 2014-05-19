package disruptor

type Barrier interface {
	Load() int64
}

type MultiBarrier struct {
	cursors []*Cursor
}

func (this *MultiBarrier) Load() int64 {
	minimum := MaxSequenceValue

	for _, item := range this.cursors {
		cursor := item.Load()
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}

func NewBarrier(upstream ...*Cursor) Barrier {
	if len(upstream) == 0 {
		panic("At least one upstream cursor is required.")
	} else if len(upstream) == 1 {
		return upstream[0]
	} else {
		cursors := make([]*Cursor, len(upstream))
		copy(cursors, upstream)
		return &MultiBarrier{cursors}
	}
}
