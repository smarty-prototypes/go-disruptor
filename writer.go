package disruptor

type Writer struct {
	previous int64
	gate     int64
	capacity int64
	written  *Cursor
	upstream Barrier
}

func NewWriter(written *Cursor, capacity int64, upstream Barrier) *Writer {
	assertPowerOfTwo(capacity)

	return &Writer{
		previous: InitialSequenceValue,
		gate:     InitialSequenceValue,
		capacity: capacity,
		written:  written,
		upstream: upstream,
	}
}

func assertPowerOfTwo(value int64) {
	if value > 0 && (value&(value-1)) != 0 {
		// http://en.wikipedia.org/wiki/Power_of_two#Fast_algorithm_to_check_if_a_positive_number_is_a_power_of_two
		panic("The ring capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}
}

func (this *Writer) Reserve(count int64) (int64, int64) {
	upper := this.previous + count
	wrap := upper - this.capacity

	if wrap > this.gate {
		min := this.upstream.LoadBarrier(0)
		if wrap > min {
			return InitialSequenceValue, Gating
		}

		this.gate = min
	}

	this.previous = upper
	return upper - count + 1, upper
}
