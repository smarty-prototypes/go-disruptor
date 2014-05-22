package disruptor

type Writer struct {
	previous      int64
	gate          int64
	capacity      int64
	writerCursor  *Cursor
	readerBarrier Barrier
}

func NewWriter(writerCursor *Cursor, capacity int64, readerBarrier Barrier) *Writer {
	assertPowerOfTwo(capacity)

	return &Writer{
		previous:      writerCursor.Load(), // show the Go runtime that the cursor is actually used
		gate:          writerCursor.Load(), // and that it should not be optimized away
		capacity:      capacity,
		writerCursor:  writerCursor,
		readerBarrier: readerBarrier,
	}
}

func assertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		// http://en.wikipedia.org/wiki/Power_of_two#Fast_algorithm_to_check_if_a_positive_number_is_a_power_of_two
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}
}

func (this *Writer) Reserve(count int64) int64 {
	next := this.previous + count
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
