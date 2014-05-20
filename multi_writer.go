package disruptor

import (
	"math"
	"sync/atomic"
)

type MultiWriter struct {
	capacity  int64
	gate      int64
	shift     uint8
	committed []int32
	upstream  Barrier
	claimed   *Cursor
}

func NewMultiWriter(claimed *Cursor, capacity int64, upstream Barrier) *MultiWriter {
	assertPowerOfTwo(capacity)

	shift := uint8(math.Log2(float64(capacity)))
	buffer := initializeCommitBuffer(capacity)

	return &MultiWriter{
		capacity:  capacity,
		gate:      InitialSequenceValue,
		shift:     shift,
		committed: buffer,
		upstream:  upstream,
		claimed:   claimed,
	}
}
func initializeCommitBuffer(capacity int64) []int32 {
	buffer := make([]int32, capacity)
	for i := int64(0); i < capacity; i++ {
		buffer[i] = -1
	}
	return buffer
}

func (this *MultiWriter) Reserve(count int64) (int64, int64) {
	if count <= 0 {
		panic("Reservation must be a positive number.")
	} else if count > this.capacity {
		panic("Reservation cannot exceed the capacity.")
	}

	for {
		previous := this.claimed.Load()
		next := previous + count
		wrap := next - this.capacity

		if wrap > this.gate {
			min := this.upstream.Load()
			if wrap > min {
				return 0, Gating
			}

			this.gate = min
		}

		if atomic.CompareAndSwapInt64(&this.claimed.value, previous, next) {
			return previous + 1, next
		}
	}
}

func (this *MultiWriter) Commit(lower, upper int64) {
	for mask := this.capacity - 1; lower <= upper; lower++ {
		this.committed[lower&mask] = int32(lower >> this.shift)
	}
}

func (this *MultiWriter) LoadBarrier(lower, upper int64) int64 {
	for mask := this.capacity - 1; lower <= upper; lower++ {
		if this.committed[lower&mask] < int32(lower>>this.shift) {
			return lower - 1
		}
	}

	return upper
}
