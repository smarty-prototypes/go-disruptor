package disruptor

type Barrier interface {
	Load() int64
}

func NewBarrier(upstream ...*Cursor) Barrier {
	if len(upstream) == 0 {
		panic("At least one upstream cursor is required.")
	} else if len(upstream) == 1 {
		return upstream[0]
	} else {
		return NewCompositeBarrier(upstream...)
	}
}
