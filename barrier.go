package disruptor

type Barrier func() int64

func NewBarrier(upstream ...*Cursor) Barrier {
	cursors := make([]*Cursor, len(upstream))
	copy(cursors, upstream)

	if len(upstream) == 0 {
		panic("At least one upstream cursor is required.")
	} else if len(upstream) == 1 {
		single := cursors[0]
		return func() int64 { return single.Load() }
	} else {
		return func() int64 {
			minimum := MaxSequenceValue

			for _, item := range cursors {
				sequence := item.Load()
				if sequence < minimum {
					minimum = sequence
				}
			}

			return minimum
		}
	}
}
