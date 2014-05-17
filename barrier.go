package disruptor

type Barrier func() int64

func NewBarrier(upstream ...*Cursor) Barrier {
	cursors := make([]*Cursor, len(upstream))
	copy(cursors, upstream)

	if len(upstream) == 0 {
		panic("At least one upstream cursor is required.")
	} else if len(upstream) == 1 {
		first := cursors[0]
		return func() int64 { return first.Load() }
	} else {
		return func() int64 {
			minimum := MaxCursorValue

			for _, item := range cursors {
				cursor := item.Load()
				if cursor < minimum {
					minimum = cursor
				}
			}

			return minimum
		}
	}
}
