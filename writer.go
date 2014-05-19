package disruptor

type Writer struct {
	previous      int64
	gate          int64
	writerCursor  *Cursor
	ringSize      int64
	readerBarrier Barrier
}

func NewWriter(writerCursor *Cursor, ringSize int32, readerBarrier Barrier) *Writer {
	if !isPowerOfTwo(ringSize) {
		panic("The ring size must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}

	return &Writer{
		previous:      writerCursor.Load(), // show the Go runtime that the cursor is actually used
		gate:          writerCursor.Load(), // and that it should not be optimized away
		writerCursor:  writerCursor,
		ringSize:      int64(ringSize),
		readerBarrier: readerBarrier,
	}
}

func isPowerOfTwo(value int32) bool {
	return value > 0 && (value&(value-1)) == 0
}

func (this *Writer) Reserve(items int64) int64 {
	next := this.previous + items
	wrap := next - this.ringSize

	if wrap > this.gate {
		min := this.readerBarrier()
		if wrap > min {
			return Gating
		}

		this.gate = min
	}

	this.previous = next
	return next
}
