package disruptor

type Barrier struct {
	single  bool
	cursors []*Sequence
}

func NewBarrier(upstream ...*Sequence) *Barrier {
	cursors := make([]*Sequence, len(upstream))
	copy(cursors, upstream)
	return &Barrier{
		single:  len(cursors) == 1,
		cursors: cursors,
	}
}

func (this *Barrier) Load() int64 {
	if this.single {
		return this.cursors[0].Load()
	}

	minimum := MaxSequenceValue

	for _, item := range this.cursors {
		cursor := item.Load()
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}
