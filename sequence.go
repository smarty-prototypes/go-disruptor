package disruptor

import (
	"sync/atomic"
	"unsafe"
)

// atomicSequence is a cache-line-padded atomic int64 used to track sequence positions without false sharing.
type atomicSequence struct {
	_ [CacheLineBytes - unsafe.Sizeof(int64(0))]byte
	atomic.Int64
}

func newSequence() (this *atomicSequence) {
	for this = new(atomicSequence); uintptr(unsafe.Pointer(this))%CacheLineBytes != 0; this = new(atomicSequence) {
	}

	this.Store(defaultSequenceValue)
	return this
}

// newSequences allocates a contiguous, cache-line-aligned slice of *atomicSequence
func newSequences(count int) []*atomicSequence {
	var contiguous []atomicSequence
	for contiguous = make([]atomicSequence, count); uintptr(unsafe.Pointer(&contiguous[0]))%CacheLineBytes != 0; contiguous = make([]atomicSequence, count) {
		// not cache aligned, try again
	}

	// guaranteed cache alignment of underlying sequence values which are *contiguous*
	this := make([]*atomicSequence, count)
	for i := range this {
		this[i] = &contiguous[i]
		this[i].Store(defaultSequenceValue)
	}
	return this
}

const defaultSequenceValue = -1
