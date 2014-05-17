package disruptor

type Writer struct {
	previous      int64
	gate          int64
	writerCursor  *Cursor
	ringSize      int64
	readerBarrier *Barrier
}

func NewWriter(writerCursor *Cursor, ringSize int32, readerBarrier *Barrier) *Writer {
	if !isPowerOfTwo(ringSize) {
		panic("The ring size must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}

	return &Writer{
		previous:      InitialCursorValue,
		gate:          InitialCursorValue,
		writerCursor:  writerCursor,
		ringSize:      int64(ringSize),
		readerBarrier: readerBarrier,
	}
}

func isPowerOfTwo(value int32) bool {
	return value > 0 && (value&(value-1)) == 0
}

// TODO: rename to Reserve
func (this *Writer) Next(items int64) int64 {
	next := this.previous + items
	wrap := next - this.ringSize

	if wrap > this.gate {
		min := this.readerBarrier.Load()
		for wrap > min || min < 0 {
			min = this.readerBarrier.Load()
		}

		this.gate = min
	}

	this.previous = next
	return next
}
