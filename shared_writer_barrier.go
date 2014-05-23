package disruptor

import "math"

type SharedWriterBarrier struct {
	capacity  int64
	mask      int64
	shift     uint8
	committed []int32
	written   *Cursor
}

func NewSharedWriterBarrier(written *Cursor, capacity int64) *SharedWriterBarrier {
	assertPowerOfTwo(capacity)

	return &SharedWriterBarrier{
		capacity:  capacity,
		mask:      capacity - 1,
		shift:     uint8(math.Log2(float64(capacity))),
		committed: prepareCommitBuffer(capacity),
		written:   written,
	}
}
func prepareCommitBuffer(capacity int64) []int32 {
	buffer := make([]int32, capacity)
	for i := range buffer {
		buffer[i] = int32(InitialSequenceValue)
	}
	return buffer
}

func (this *SharedWriterBarrier) LoadBarrier(lower int64) int64 {
	shift, mask := this.shift, this.mask
	upper := this.written.Load()

	for ; lower <= upper; lower++ {
		if this.committed[lower&mask] != int32(lower>>shift) {
			return lower - 1
		}
	}

	return upper
}
