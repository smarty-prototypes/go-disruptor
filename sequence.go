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

// newSequences allocates a slice of *atomicSequence in a contiguous space in memory
func newSequences(count int) []*atomicSequence {
	this := make([]*atomicSequence, count)
	for i := range this {
		this[i] = newSequence()
	}
	return this
}

const defaultSequenceValue = -1
