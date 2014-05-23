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

func (this *SharedWriterBarrier) LoadBarrier(lower int64) int64 {
	shift, mask := this.shift, this.mask
	upper := this.reservation.Load()
	// fmt.Printf("\t\t\t\t\t[BARRIER] Next (lower): %d, Reservation (upper): %d\n", lower, upper)
	// fmt.Println("\t\t\t\t\t[BARRIER] Committed:", this.committed)
	for ; lower <= upper; lower++ {
		// fmt.Printf("\t\t\t\t\t[BARRIER] Inside Loop. Index: %d, Value: %d\n", sequence&mask, sequence>>shift)
		if this.committed[lower&mask] != int32(lower>>shift) {
			// fmt.Println("\t\t\t\t\t[BARRIER] Upstream Barrier:", sequence-1, this.committed)
			return lower - 1
		}
	}

	// fmt.Println("\t\t\t\t\t[BARRIER] Upstream Barrier (default):", lower)
	return upper
}
