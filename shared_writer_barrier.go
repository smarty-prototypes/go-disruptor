package disruptor

import (
	"fmt"
	"math"
)

type SharedWriterBarrier struct {
	written   *Cursor
	committed []int32
	capacity  int64
	mask      int64
	shift     uint8
}

func NewSharedWriterBarrier(written *Cursor, capacity int64) *SharedWriterBarrier {
	fmt.Println("Shift: ", uint8(math.Log2(float64(capacity))))
	assertPowerOfTwo(capacity)

	return &SharedWriterBarrier{
		written:   written,
		committed: prepareCommitBuffer(capacity),
		capacity:  capacity,
		mask:      capacity - 1,
		shift:     uint8(math.Log2(float64(capacity))),
	}
}
func prepareCommitBuffer(capacity int64) []int32 {
	buffer := make([]int32, capacity)
	for i := range buffer {
		buffer[i] = int32(InitialSequenceValue)
	}
	return buffer
}

func (this *SharedWriterBarrier) Read(lower int64) int64 {
	shift, mask := this.shift, this.mask
	upper := this.written.Load()

	for sequence := lower; sequence <= upper; sequence++ {
		fmt.Printf("Reader found sequence %d (slot %d) populated with %d, expected %d\n", sequence, sequence&mask, this.committed[sequence&mask], int32(sequence>>shift))
		if this.committed[sequence&mask] != int32(sequence>>shift) {
			fmt.Printf("Reader cannot advance, awaiting a writer to commit reservation.\n")
			return lower - 1
		}
	}

	return upper
}
