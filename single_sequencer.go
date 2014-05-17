package disruptor

type SingleProducerSequencer struct {
	previous int64
	gate     int64
	cursor   *Sequence
	ringSize int64
	barrier  Barrier
}

func NewSingleProducerSequencer(cursor *Sequence, ringSize int32, barrier Barrier) *SingleProducerSequencer {
	if !isPowerOfTwo(ringSize) {
		panic("The ring size must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.")
	}

	return &SingleProducerSequencer{
		previous: InitialSequenceValue,
		gate:     InitialSequenceValue,
		cursor:   cursor,
		ringSize: int64(ringSize),
		barrier:  barrier,
	}
}

func isPowerOfTwo(value int32) bool {
	return value > 0 && (value&(value-1)) == 0
}

func (this *SingleProducerSequencer) Next(items int64) int64 {
	next := this.previous + items
	wrap := next - this.ringSize

	if wrap > this.gate {
		min := this.barrier.Load()
		for wrap > min || min < 0 {
			min = this.barrier.Load()
		}

		this.gate = min
	}

	this.previous = next
	return next
}
