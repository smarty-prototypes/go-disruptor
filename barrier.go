package disruptor

type Barrier struct {
	single  bool
	cursors []*Cursor
}

func NewBarrier(upstream ...*Cursor) *Barrier {
	cursors := make([]*Cursor, len(upstream))
	copy(cursors, upstream)

	// TODO: the "Load" function could be set here as a callback
	// such that the public "Load" points to one of two private load functions
	return &Barrier{
		single:  len(cursors) == 1,
		cursors: cursors,
	}
}

func (this *Barrier) Load() int64 {
	if this.single {
		return this.cursors[0].Load()
	}

	minimum := MaxCursorValue

	for _, item := range this.cursors {
		cursor := item.Load()
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}
