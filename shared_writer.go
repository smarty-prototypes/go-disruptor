package disruptor

import (
	"math"
	"sync/atomic"
)

type SharedWriter struct {
	capacity    int64
	gate        int64 // TODO: this will most likely need to be a cursor
	shift       uint8
	committed   []int32
	upstream    Barrier
	reservation *Cursor
}

func NewSharedWriter(reservation *Cursor, capacity int64, upstream Barrier) *SharedWriter {
	assertPowerOfTwo(capacity)

	shift := uint8(math.Log2(float64(capacity)))
	buffer := initializeCommittedBuffer(capacity)

	return &SharedWriter{
		capacity:    capacity,
		gate:        InitialSequenceValue,
		shift:       shift,
		committed:   buffer,
		upstream:    upstream,
		reservation: reservation,
	}
}
func initializeCommittedBuffer(capacity int64) []int32 {
	buffer := make([]int32, capacity)
	for i := range buffer {
		buffer[i] = int32(InitialSequenceValue)
	}
	return buffer
}

func (this *SharedWriter) Reserve(count int64) (int64, int64) {
	if count <= 0 {
		panic("Reservation must be a positive number.")
	} else if count > this.capacity {
		panic("Reservation cannot exceed the capacity.")
	}

	for {
		previous := this.reservation.Load()
		next := previous + count
		wrap := next - this.capacity

		if wrap > this.gate {
			min := this.upstream.Load()
			if wrap > min {
				return 0, Gating
			}

			this.gate = min // doesn't matter which write wins, BUT will most likely need to be a Cursor
		}

		if atomic.CompareAndSwapInt64(&this.reservation.value, previous, next) {
			return previous + 1, next
		}
	}
}

func (this *SharedWriter) Commit(lower, upper int64) {
	for mask := this.capacity - 1; lower <= upper; lower++ {
		this.committed[lower&mask] = int32(lower >> this.shift)
	}
}

func (this *SharedWriter) Load() int64 {
	sequence := this.reservation.Load()

	for mask := this.capacity - 1; sequence >= 0; sequence-- {
		if this.committed[sequence&mask] == int32(sequence>>this.shift) {
			return sequence
		}
	}

	return sequence
}
