package disruptor

import "math"

type SharedWriterBarrier struct {
	committed   []int32
	capacity    int64
	mask        int64
	shift       uint8
	reservation *Cursor
}

func NewSharedWriterBarrier(reservation *Cursor, capacity int64) *SharedWriterBarrier {
	assertPowerOfTwo(capacity)

	return &SharedWriterBarrier{
		committed:   prepareCommitBuffer(capacity),
		capacity:    capacity,
		mask:        capacity - 1,
		shift:       uint8(math.Log2(float64(capacity))),
		reservation: reservation,
	}
}
func prepareCommitBuffer(capacity int64) []int32 {
	buffer := make([]int32, capacity)
	for i := range buffer {
		buffer[i] = int32(InitialSequenceValue)
	}
	return buffer
}

func (this *SharedWriterBarrier) Load() int64 {

	for sequence := this.reservation.Load(); sequence >= 0; sequence-- {
		if this.committed[sequence&this.mask] == int32(sequence>>this.shift) {
			// fmt.Printf("\t\t\t\t\t[SHARED-WRITER-BARRIER] Barrier Sequence: %d\n", sequence)
			return sequence
		}
	}

	// fmt.Printf("\t\t\t\t\t[SHARED-WRITER-BARRIER] Barrier Sequence: -1\n")
	return InitialSequenceValue
}
