package disruptor

const (
	Gating = -2
	Idling = -3
)

type Reader struct {
	upstreamBarrier Barrier
	writerCursor    *Cursor
	readerCursor    *Cursor
}

func NewReader(upstreamBarrier Barrier, writerCursor, readerCursor *Cursor) *Reader {
	return &Reader{
		upstreamBarrier: upstreamBarrier,
		writerCursor:    writerCursor,
		readerCursor:    readerCursor,
	}
}

// TODO: performance when current (or next?) sequence is received as a parameter to Receive
// instead of reading the cursor...
func (this *Reader) Receive() (int64, int64) {
	next := this.readerCursor.Load() + 1
	ready := this.upstreamBarrier.LoadBarrier(next)

	if next <= ready {
		return next, ready - next
	} else if gate := this.writerCursor.Load(); next <= gate {
		return next, Gating
	} else {
		return next, Idling
	}
}
