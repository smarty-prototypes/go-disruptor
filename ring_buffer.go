package disruptor

type RingBuffer []uint64

func NewRingBuffer(bufferSize int) RingBuffer {
	if !isPowerOfTwo(bufferSize) {
		panic("The buffer size must be a power of two.")
	}

	return RingBuffer(make([]uint64, bufferSize))
}
func isPowerOfTwo(value int) bool {
	return value != 0 && (value&(value-1)) == 0
}
