package disruptor

const (
	Gating = -2
	Idle   = -3
)

type Reader struct {
	upstreamBarrier *Barrier
	callback        Consumer
	writerCursor    *Cursor
	readerCursor    *Cursor
}

func NewReader(upstreamBarrier *Barrier, callback Consumer, writerCursor, readerCursor *Cursor) *Reader {
	return &Reader{
		upstreamBarrier: upstreamBarrier,
		callback:        callback,
		writerCursor:    writerCursor,
		readerCursor:    readerCursor,
	}
}

func (this *Reader) Process() int64 {
	next := this.readerCursor.Load() + 1
	ready := this.upstreamBarrier.Load()

	if next <= ready {
		for next <= ready {
			this.callback.Consume(next, ready-next)
			next++
		}

		next--
		this.readerCursor.Store(next)
		return next
	} else if next <= this.writerCursor.Load() {
		return Gating
	} else {
		return Idle
	}
}
