package disruptor

import "sync/atomic"

type Sequence []uint64

func (this Sequence) atomicLoad() uint64 {
	return atomic.LoadUint64(&this[0])
}
func (this Sequence) store(value uint64) {
	this[0] = value
}
func (this Sequence) close() {
	this[0] = Uint64MaxValue
}

func NewSequence() Sequence {
	return Sequence(make([]uint64, FillCPUCacheLine))
}

// TODO: use build tags for i386, amd64, and ARM-v4,5,6,7,8 processors
// i386, ARM? = 32-byte cache line vs 64-byte cache line for amd64
const FillCPUCacheLine = 8
