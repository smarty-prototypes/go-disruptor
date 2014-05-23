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
	lower := this.readerCursor.Load() + 1
	upper := this.upstreamBarrier.LoadBarrier(lower)

	if lower <= upper {
		return lower, upper
	} else if gate := this.writerCursor.Load(); lower <= gate {
		return InitialSequenceValue, Gating
	} else {
		return InitialSequenceValue, Idling
	}
}
