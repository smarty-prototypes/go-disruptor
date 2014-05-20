package disruptor

type Writer struct {
	previous      int64
	gate          int64
	capacity      int64
	writerCursor  *Cursor
	readerBarrier Barrier
}

func NewWriter(writerCursor *Cursor, capacity int64, readerBarrier Barrier) *Writer {
	if !isPowerOfTwo(capacity) {
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}

	return &Writer{
		previous:      writerCursor.Load(), // show the Go runtime that the cursor is actually used
		gate:          writerCursor.Load(), // and that it should not be optimized away
		capacity:      capacity,
		writerCursor:  writerCursor,
		readerBarrier: readerBarrier,
	}
}

func isPowerOfTwo(value int64) bool {
	return value > 0 && (value&(value-1)) == 0
}

func (this *Writer) Reserve(items int64) int64 {
	next := this.previous + items
	wrap := next - this.capacity

	if wrap > this.gate {
		min := this.readerBarrier.Load()
		if wrap > min {
			return Gating
		}

		this.gate = min
	}

	this.previous = next
	return next
}
