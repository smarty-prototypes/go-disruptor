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

func (this *Reader) Receive() (int64, int64) {
	next := this.readerCursor.Load() + 1
	ready := this.upstreamBarrier.Load()

	if next <= ready {
		return next, ready - next
	} else if next <= this.writerCursor.Load() {
		return next, Gating
	} else {
		return next, Idling
	}
}
